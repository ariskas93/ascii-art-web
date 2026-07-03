package output

import (
	"ascii-art/internal/renderer"
	"io"
	"os"
	"testing"
)

func TestPrintASCII_returnsNoError(t *testing.T) {
	lines := []string{"  _  ", " | | "}

	err := PrintASCII(lines)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPrintAlignedASCII_usesDetectedTerminalWidth(t *testing.T) {
	t.Setenv("COLUMNS", "6")
	rendered := renderer.RenderedASCII{
		{
			Segments: []renderer.RenderedSegment{
				{
					Kind: renderer.SegmentWord,
					Text: "A",
					Rows: []string{"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7"},
				},
			},
		},
	}

	got, err := captureStdout(t, func() error {
		return PrintAlignedASCII(rendered, "center")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "  A0\n" +
		"  A1\n" +
		"  A2\n" +
		"  A3\n" +
		"  A4\n" +
		"  A5\n" +
		"  A6\n" +
		"  A7\n"
	if got != expected {
		t.Fatalf("got %q, want %q", got, expected)
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not create stdout pipe: %v", err)
	}

	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	fnErr := fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("could not close stdout writer: %v", err)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("could not read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("could not close stdout reader: %v", err)
	}

	os.Stdout = oldStdout

	return string(data), fnErr
}
