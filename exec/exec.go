package exec

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func Exec(out io.Writer, command string) error {

	if command == "" {
		return errors.New("empty command")
	}
	cmd := exec.Command("sh", "-c", command)

	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command failed with exit code %d: %w", exitErr.ExitCode(), err)
		}
		return fmt.Errorf("command execution error: %w", err)
	}

	return nil
}
