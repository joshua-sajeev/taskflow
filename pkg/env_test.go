package pkg

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		setEnv   bool
		expected string
	}{
		{
			name:     "returns env value when set",
			key:      "TEST_VAR",
			fallback: "default",
			envValue: "from_env",
			setEnv:   true,
			expected: "from_env",
		},
		{
			name:     "returns env value with special characters",
			key:      "SPECIAL_VAR",
			fallback: "default",
			envValue: "value:with:colons",
			setEnv:   true,
			expected: "value:with:colons",
		},
		{
			name:     "returns fallback when env not set",
			key:      "UNSET_VAR",
			fallback: "default_value",
			setEnv:   false,
			expected: "default_value",
		},
		{
			name:     "returns empty string from env when set to empty",
			key:      "EMPTY_VAR",
			fallback: "default",
			envValue: "",
			setEnv:   true,
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := GetEnv(tt.key, tt.fallback)
			if result != tt.expected {
				t.Errorf("GetEnv() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetEnv_MissingRequiredVariable(t *testing.T) {
	if os.Getenv("TEST_SUBPROCESS") == "1" {
		GetEnv("MISSING_REQUIRED_VAR", "")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestGetEnv_MissingRequiredVariable")
	cmd.Env = append(os.Environ(), "TEST_SUBPROCESS=1")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected process to exit with error, but it succeeded")
	}

	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.Success() {
		t.Errorf("expected non-zero exit status, got: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "missing required environment variable") {
		t.Errorf("expected error message to contain 'missing required environment variable', got: %s", outputStr)
	}
	if !strings.Contains(outputStr, "MISSING_REQUIRED_VAR") {
		t.Errorf("expected error message to contain 'MISSING_REQUIRED_VAR', got: %s", outputStr)
	}
}

func TestGetEnv_EmptyFallbackWithSetEnv(t *testing.T) {
	key := "SET_VAR_EMPTY_FALLBACK"
	value := "actual_value"

	os.Setenv(key, value)
	defer os.Unsetenv(key)

	result := GetEnv(key, "")
	if result != value {
		t.Errorf("GetEnv() = %q, want %q", result, value)
	}
}
