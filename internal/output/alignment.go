package output

import (
	"os"
	"strconv"
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
