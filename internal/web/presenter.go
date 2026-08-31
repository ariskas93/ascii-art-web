package web

import (
	"errors"
	"strings"
)

// Present converts renderer-controlled ANSI text into safe browser runs and
// ANSI-free download text. Package C implements it after Phase 0.
func Present(g GeneratedASCII) (lines []StyledLine, plain string, err error) {
	var plainBuilder strings.Builder
	var runBuilder strings.Builder
	var currentLine StyledLine
	colored := false

	flushRun := func() {
		if runBuilder.Len() == 0 {
			return
		}
		currentLine.Runs = append(currentLine.Runs, StyledRun{
			Text:    runBuilder.String(),
			Colored: colored,
		})
		runBuilder.Reset()
	}
	flushLine := func() {
		flushRun()
		lines = append(lines, currentLine)
		currentLine = StyledLine{}
	}

	text := g.ANSIText
	for i := 0; i < len(text); {
		ch := text[i]
		if ch == '\n' {
			plainBuilder.WriteByte(ch)

			flushLine()
			i++
			continue
		}
		if ch == '\033' {
			end := i + 1
			for end < len(text) && text[end] != 'm' {
				end++
			}
			if end >= len(text) {
				plain = plainBuilder.String()
				return lines, plain, errors.New("incomplete ANSI escape sequence")
			}
			sequence := text[i : end+1]
			flushRun()
			if sequence == "\033[0m" {
				colored = false
			} else if isSupportedANSIColor(sequence) {
				colored = true
			} else {
				plain = plainBuilder.String()
				return lines, plain, errors.New("unsupported ANSI escape sequence")
			}
			i = end + 1
			continue

		}
		plainBuilder.WriteByte(ch)
		runBuilder.WriteByte(ch)
		i++

	}
	if runBuilder.Len() > 0 || len(currentLine.Runs) > 0 {
		flushLine()
	}
	plain = plainBuilder.String()
	return
}

func isSupportedANSIColor(sequence string) bool {
	switch sequence {
	case "\033[30m", "\033[31m", "\033[32m", "\033[33m":
		return true
	case "\033[34m", "\033[35m", "\033[36m", "\033[37m":
		return true
	case "\033[38;5;208m":
		return true
	default:
		return false
	}
}
