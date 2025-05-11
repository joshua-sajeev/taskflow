package scheduler

import (
	"bufio"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"strings"
	"taskflow/exec"
)

var CronScheduler = cron.New(
	cron.WithParser(
		cron.NewParser(
			cron.SecondOptional |
				cron.Minute |
				cron.Hour |
				cron.Dom |
				cron.Month |
				cron.Dow |
				cron.Descriptor,
		),
	),
)

func Schedule(cronString string, job func()) error {

	_, err := CronScheduler.AddFunc(cronString, func() { go job() })
	if err != nil {
		return err
	}
	CronScheduler.Start()

	return nil
}

func ScheduleFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening cron file: %v", err)
	}
	defer file.Close()

	lines := make(chan string)
	readerr := make(chan error)
	done := make(chan bool)

	// Reader goroutine
	go func() {
		defer close(lines)
		defer close(readerr)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		readerr <- scanner.Err()
	}()

	// Scheduler goroutine
	go func() {
		for line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "|", 2)
			if len(parts) != 2 {
				log.Printf("Invalid format, expected '<cron> | <cmd>': %s", line)
				continue
			}

			spec := strings.TrimSpace(parts[0])
			cmd := strings.TrimSpace(parts[1])

			fmt.Printf("Scheduling: [%s] with cron [%s]\n", cmd, spec)

			err := Schedule(spec, func() {
				if err := exec.Exec(os.Stdout, cmd); err != nil {
					fmt.Println("Error executing job:", err)
				}
			})
			if err != nil {
				log.Printf("Failed to schedule job: %v", err)
			}
		}
		done <- true
	}()

	// Wait for reading and scheduling
	if err := <-readerr; err != nil {
		return fmt.Errorf("error reading cron file: %w", err)
	}
	<-done

	select {} // Keep running
}
