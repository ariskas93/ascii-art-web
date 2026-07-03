// Package output owns terminal width detection, alignment, and writing rendered ASCII.
package output

import (
	"ascii-art/internal/renderer"
	"os"
)

// WriteASCIIFile writes flat ASCII rows to fileName, replacing any existing file.
func WriteASCIIFile(fileName string, lines []string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteASCII(f, lines)
}

// WriteAlignedASCIIFile writes aligned rendered ASCII to fileName using an explicit width.
func WriteAlignedASCIIFile(fileName string, rendered renderer.RenderedASCII, align string, width int) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteAlignedASCII(f, rendered, align, width)
}

// WriteAlignedASCIIFileWithTerminalWidth writes aligned rendered ASCII using the detected terminal width.
func WriteAlignedASCIIFileWithTerminalWidth(fileName string, rendered renderer.RenderedASCII, align string) error {
	width := terminalWidth()
	return WriteAlignedASCIIFile(fileName, rendered, align, width)
}
