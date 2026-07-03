// Package output owns terminal width detection, alignment, and writing rendered ASCII.
package output

import (
	"ascii-art/internal/renderer"
	"os"
)

// PrintASCII writes flat ASCII rows to standard output.
func PrintASCII(lines []string) error {
	return WriteASCII(os.Stdout, lines)
}

// PrintAlignedASCII writes rendered ASCII to standard output using the detected terminal width.
func PrintAlignedASCII(rendered renderer.RenderedASCII, align string) error {
	return WriteAlignedASCII(os.Stdout, rendered, align, terminalWidth())
}
