package output

import (
	"ascii-art/internal/renderer"
	"os"
)

func WriteASCIIFile(fileName string, lines []string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteASCII(f, lines)
}

func WriteAlignedASCIIFile(fileName string, rendered renderer.RenderedASCII, align string, width int) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteAlignedASCII(f, rendered, align, width)
}

func WriteAlignedASCIIFileWithTerminalWidth(fileName string, rendered renderer.RenderedASCII, align string) error {
	width := terminalWidth()
	return WriteAlignedASCIIFile(fileName, rendered, align, width)
}
