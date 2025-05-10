package cmd

import (
	"fmt"
	"os"
	"taskflow/exec"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Let's you run bash commands",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No command provided")
			return
		}

		for _, command := range args {
			err := exec.Exec(os.Stdout, command)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
