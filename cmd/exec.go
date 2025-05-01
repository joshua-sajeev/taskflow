package cmd

import (
	"os"
	"taskflow/exec"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Let's you run bash command",
	Run: func(cmd *cobra.Command, args []string) {
		for _, command := range args {
			exec.Exec(os.Stdout, command)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
