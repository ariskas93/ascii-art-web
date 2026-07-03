package output

import (
	"ascii-art/internal/renderer"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func visibleWidth(s string) int {
	width := 0
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) {
			if s[i] == '\033' && s[i+1] == '[' {
				for i < len(s) && s[i] != 'm' {
					i++
				}
				continue
			}
		}
		width++
	}
	return width
}

func terminalWidth() int {
	value := os.Getenv("COLUMNS")
	number, err := strconv.Atoi(value)
	if err != nil {
		return 80
	}
	if number < 1 {
		return 80
	}
	return number
}

func WriteAlignedASCII(w io.Writer, rendered renderer.RenderedASCII, align string, width int) error {
	if align == "" || align == "left" {
		lines := renderer.Flatten(rendered)
		return WriteASCII(w, lines)
	}
	if align == "center" {
		lines := renderer.Flatten(rendered)
		for i, row := range lines {
			rowWidth := visibleWidth(row)
			if rowWidth < width {
				padding := (width - rowWidth) / 2
				lines[i] = strings.Repeat(" ", padding) + row
			}
		}
		return WriteASCII(w, lines)
	}
	if align == "right" {
		lines := renderer.Flatten(rendered)
		for i, row := range lines {
			rowWidth := visibleWidth(row)
			if rowWidth < width {
				padding := width - rowWidth
				lines[i] = strings.Repeat(" ", padding) + row
			}
		}
		return WriteASCII(w, lines)
	}
	if align == "justify" {
		return writeJustifiedASCII(w, rendered, width)
	}
	return fmt.Errorf("unsupported alignment: %s", align)
}

func writeJustifiedASCII(w io.Writer, rendered renderer.RenderedASCII, width int) error {
	var output []string

	for _, line := range rendered {
		if len(line.Segments) == 0 {
			output = append(output, "")
			continue
		}

		var normalRows [8]string
		for _, segment := range line.Segments {
			for r := 0; r < 8; r++ {
				normalRows[r] += segment.Rows[r]
			}
		}

		gapCount := 0
		for i, segment := range line.Segments {
			if segment.Kind == renderer.SegmentSpace &&
				i > 0 &&
				i < len(line.Segments)-1 &&
				line.Segments[i-1].Kind == renderer.SegmentWord &&
				line.Segments[i+1].Kind == renderer.SegmentWord {
				gapCount++
			}
		}

		rowWidth := visibleWidth(normalRows[0])
		if gapCount == 0 || rowWidth >= width {
			for _, row := range normalRows {
				output = append(output, row)
			}
			continue
		}

		extraSpaces := width - rowWidth
		baseExtra := extraSpaces / gapCount
		remainder := extraSpaces % gapCount

		var justifiedRows [8]string
		gapIndex := 0

		for i, segment := range line.Segments {
			isGap := segment.Kind == renderer.SegmentSpace &&
				i > 0 &&
				i < len(line.Segments)-1 &&
				line.Segments[i-1].Kind == renderer.SegmentWord &&
				line.Segments[i+1].Kind == renderer.SegmentWord

			additionalSpaces := 0
			if isGap {
				additionalSpaces = baseExtra
				if gapIndex < remainder {
					additionalSpaces++
				}
				gapIndex++
			}

			for r := 0; r < 8; r++ {
				justifiedRows[r] += segment.Rows[r] + strings.Repeat(" ", additionalSpaces)
			}
		}

		for _, row := range justifiedRows {
			output = append(output, row)
		}
	}

	return WriteASCII(w, output)
}
