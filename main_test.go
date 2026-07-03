// Package main contains end-to-end CLI wiring tests.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainInvalidSeparatedAlignFlag(t *testing.T) {
	stdout, stderr, err := runMain(t, nil, "--align", "right", "something", "standard")
	if err == nil {
		t.Fatalf("expected command to fail")
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "Usage: go run . [OPTION] [STRING] [BANNER]") {
		t.Fatalf("stderr = %q, want alignment usage message", stderr)
	}
	if !strings.Contains(stderr, "--align=<left|center|right|justify>") {
		t.Fatalf("stderr = %q, want alignment usage message", stderr)
	}
}

func TestMainPrintsRightAlignedASCII(t *testing.T) {
	stdout, stderr, err := runMain(t, map[string]string{"COLUMNS": "80"}, "--align=right", "left", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	rows := outputRows(stdout)
	if len(rows) != 8 {
		t.Fatalf("got %d output rows, want 8", len(rows))
	}
	for i, row := range rows {
		if len(row) != 80 {
			t.Fatalf("row %d length = %d, want 80", i, len(row))
		}
	}
}

func TestMainWritesAlignedASCIIFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aligned.txt")

	stdout, stderr, err := runMain(t, map[string]string{"COLUMNS": "80"}, "--output="+path, "--align=right", "hello", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr: %s", err, stderr)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	rows := outputRows(string(data))
	if len(rows) != 8 {
		t.Fatalf("got %d output rows, want 8", len(rows))
	}
	for i, row := range rows {
		if len(row) != 80 {
			t.Fatalf("row %d length = %d, want 80", i, len(row))
		}
	}
}

// runMain executes the CLI in a subprocess so main.go wiring is tested end to end.
func runMain(t *testing.T, env map[string]string, args ...string) (string, string, error) {
	t.Helper()

	cmdArgs := append([]string{"run", "."}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = os.Environ()
	for key, value := range env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// outputRows splits command output into rendered rows without counting the final newline.
func outputRows(output string) []string {
	output = strings.TrimSuffix(output, "\n")
	if output == "" {
		return nil
	}
	return strings.Split(output, "\n")
}
