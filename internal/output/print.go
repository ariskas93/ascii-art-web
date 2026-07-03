package output

import (
	"ascii-art/internal/renderer"
	"os"
)

func PrintASCII(lines []string) error {
	return WriteASCII(os.Stdout, lines)
}

func PrintAlignedASCII(rendered renderer.RenderedASCII, align string) error {
	return WriteAlignedASCII(os.Stdout, rendered, align, terminalWidth())
}
