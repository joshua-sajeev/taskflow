package exec

import (
	"fmt"
	"io"
	"os/exec"
)

func Exec(out io.Writer, command string) {

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Fprint(out, string(output))
}
