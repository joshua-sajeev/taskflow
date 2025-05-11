package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"taskflow/exec"
	"taskflow/scheduler"
)

var (
	cronExpr string
	shellCmd string
	fileName string
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Run shell commands on a schedule using cron syntax",
	Long: `The 'scheduler' command lets you run any shell command at scheduled intervals using cron expressions.

Syntax:
  taskflow scheduler --cron "<cron_expression>" --cmd "<shell_command>"
  OR
  taskflow scheduler --every "<interval>" --cmd "<shell_command>"

Examples:
  - Run a command every minute using cron expression:
      taskflow scheduler --cron "* * * * *" --cmd "echo Hello"

  - Run a script every 5 minutes using cron expression:
      taskflow scheduler --cron "*/5 * * * *" --cmd "bash myscript.sh"

  - Run a command every Sunday at midnight using cron expression:
      taskflow scheduler --cron "0 0 * * 6" --cmd "echo 'Weekly task'"

  - Run a command every weekday at 9:00 AM using cron expression:
      taskflow scheduler --cron "0 9 * * 1-5" --cmd "echo 'Workday start'"

  - Run a command every 1 second using @every:
      taskflow scheduler --every "1s" --cmd "echo 'Every second'"

  - Run a command every 5 minutes using @every:
      taskflow scheduler --every "5m" --cmd "bash myscript.sh"

Cron Expression Format (based on robfig/cron):
  ┌───────────── minute (0 - 59)
  │ ┌───────────── hour (0 - 23)
  │ │ ┌───────────── day of month (1 - 31)
  │ │ │ ┌───────────── month (1 - 12)
  │ │ │ │ ┌───────────── day of week (0 - 6) (Sunday to Saturday)
  │ │ │ │ │
  │ │ │ │ │
  * * * * *

Note:
  - The command will keep running until you manually stop it (Ctrl+C).
  - Valid range for fields in cron expression: minute (0-59), hour (0-23), day of month (1-31), month (1-12), day of week (0-6).
  - The '@every' interval accepts human-readable durations like "1s" for seconds, "5m" for minutes, "1h" for hours, etc.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		hasFile := fileName != ""
		hasCron := cronExpr != ""
		hasCmd := shellCmd != ""

		// All empty
		if !hasFile && !(hasCron && hasCmd) {
			return fmt.Errorf("you must provide either --file or both --cron and --cmd")
		}

		// Mixed usage
		if hasFile && (hasCron || hasCmd) {
			return fmt.Errorf("--file cannot be used with --cron or --cmd")
		}

		// Using --file only
		if hasFile {
			fmt.Printf("Reading cron jobs from file: %s\n", fileName)
			scheduler.ScheduleFromFile(fileName)
			select {} // keep the scheduler running
		}

		// Using --cron and --cmd
		if hasCron && hasCmd {
			fmt.Printf("Scheduling: [%s] with cron [%s]\n", shellCmd, cronExpr)

			err := scheduler.Schedule(cronExpr, func() {
				err := exec.Exec(os.Stdout, shellCmd)
				if err != nil {
					fmt.Println("Error:", err)
				}
			})
			if err != nil {
				return fmt.Errorf("failed to schedule job: %w", err)
			}

			select {} // keep the scheduler running
		}

		// fallback, shouldn't be reached due to previous conditions
		return fmt.Errorf("unexpected flag combination")
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	// From CLI
	scheduleCmd.Flags().StringVar(&cronExpr, "cron", "", "Cron expression for scheduling")
	scheduleCmd.Flags().StringVar(&shellCmd, "cmd", "", "Shell command to execute")

	// From file
	scheduleCmd.Flags().StringVar(&fileName, "file", "", "File containing jobs in cron syntax")
}
