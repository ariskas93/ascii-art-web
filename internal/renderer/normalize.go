package renderer

import "strings"

func NormalizeInput(input string) []string {
	if input == "" {
		return []string{}
	}
	lines := strings.Split(input, "\\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
