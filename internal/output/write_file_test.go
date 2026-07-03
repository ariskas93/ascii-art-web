package output

import (
	"ascii-art/internal/renderer"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteASCIIFile_createsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")

	err := WriteASCIIFile(path, []string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created, but it does not exist")
	}
}

func TestWriteASCIIFile_contentMatchesLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")
	lines := []string{"foo", "bar"}

	err := WriteASCIIFile(path, lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	expected := "foo\nbar\n"
	if string(data) != expected {
		t.Errorf("got %q, want %q", string(data), expected)
	}
}

func TestWriteASCIIFile_overwritesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")

	_ = os.WriteFile(path, []byte("old content\n"), 0644)

	err := WriteASCIIFile(path, []string{"new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	expected := "new\n"
	if string(data) != expected {
		t.Errorf("got %q, want %q", string(data), expected)
	}
}

func TestWriteASCIIFile_emptyLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")

	err := WriteASCIIFile(path, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	if len(data) != 0 {
		t.Errorf("expected empty file, got %q", string(data))
	}
}

func TestWriteASCIIFile_invalidPath_returnsError(t *testing.T) {
	err := WriteASCIIFile("/nonexistent_dir/output.txt", []string{"hello"})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestWriteAlignedASCIIFile_writesAlignedRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aligned.txt")
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

	err := WriteAlignedASCIIFile(path, rendered, "right", 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	expected := "    A0\n" +
		"    A1\n" +
		"    A2\n" +
		"    A3\n" +
		"    A4\n" +
		"    A5\n" +
		"    A6\n" +
		"    A7\n"
	if string(data) != expected {
		t.Fatalf("got %q, want %q", string(data), expected)
	}
}
