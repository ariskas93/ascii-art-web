// Package renderer converts normalized input text into structured ASCII-art rows.
package renderer

import "strings"

// NormalizeInput splits escaped newline sequences while preserving real newline characters.
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
