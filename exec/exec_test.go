package exec

import (
	"bytes"
	"testing"
)

func TestExec(t *testing.T) {
	buffer := &bytes.Buffer{}
	Exec(buffer, "echo Hello!")
	got := buffer.String()
	want := "Hello!\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
