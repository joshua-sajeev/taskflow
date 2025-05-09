package exec_test

import (
	"bytes"
	"strings"
	"testing"

	"taskflow/exec"
)

func TestExec(t *testing.T) {
	tests := []struct {
		name         string
		command      string
		wantInOutput string // check if this appears in stdout or stderr
		wantError    bool
	}{
		{
			name:         "simple echo",
			command:      "echo Hello!",
			wantInOutput: "Hello!",
			wantError:    false,
		},
		{
			name:         "nonexistent command",
			command:      "fakecommand",
			wantInOutput: "not found", // system-dependent message
			wantError:    true,
		},
		{
			name:         "command with stderr",
			command:      "ls /does-not-exist",
			wantInOutput: "No such file", // or just "No such", for safety
			wantError:    true,
		},
		{
			name:         "empty command",
			command:      "",
			wantInOutput: "",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			err := exec.Exec(buf, tt.command)
			output := buf.String()

			if tt.wantError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantInOutput != "" && !strings.Contains(output, tt.wantInOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.wantInOutput, output)
			}
		})
	}
}
