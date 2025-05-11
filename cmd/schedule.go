package cmd

import (
	"fmt"
	"log"
	"os"
	"taskflow/exec"
	"taskflow/scheduler"

	"github.com/spf13/cobra"
)

var (
	cronExpr string
	shellCmd string
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Run shell commands on a schedule using cron syntax",
	Long: `The 'scheduler' command lets you run any shell command at scheduled intervals using cron expressions.

Syntax:
  taskflow scheduler --cron "<cron_expression>" --cmd "<shell_command>"

Examples:
  - Run a command every minute:
      taskflow scheduler --cron "* * * * *" --cmd "echo Hello"

  - Run a script every 5 minutes:
      taskflow scheduler --cron "*/5 * * * *" --cmd "bash myscript.sh"

Cron Expression Format:
  ┌───────────── minute (0 - 59)
  │ ┌───────────── hour (0 - 23)
  │ │ ┌───────────── day of month (1 - 31)
  │ │ │ ┌───────────── month (1 - 12)
  │ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
  │ │ │ │ │
  │ │ │ │ │
  * * * * *

Note:
  The command will keep running until you manually stop it (Ctrl+C).
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scheduler called")

		if cronExpr == "" || shellCmd == "" {
			log.Fatal("Both --cron and --cmd flags are required")
		}

		fmt.Printf("Scheduling: [%s] with cron [%s]\n", shellCmd, cronExpr)

		err := scheduler.Schedule(cronExpr, func() {
			err := exec.Exec(os.Stdout, shellCmd)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		})
		if err != nil {
			log.Fatalf("Failed to schedule job: %v", err)
		}
		if err != nil {
			log.Fatalf("Failed to schedule job: %v", err)
		}

		select {}

	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.Flags().StringVar(&cronExpr, "cron", "", "Cron expression for scheduling")
	scheduleCmd.Flags().StringVar(&shellCmd, "cmd", "", "Shell command to execute")
	scheduleCmd.MarkFlagRequired("cron")
	scheduleCmd.MarkFlagRequired("cmd")
}
