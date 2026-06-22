package renderer

import "errors"

func BuildASCII(lines []string, banner map[rune][]string, substring string, color string) ([]string, error) {
	var colorCode string
	if color != "" {
		var err error
		colorCode, err = getANSIColor(color)
		if err != nil {
			return nil, err
		}
	}

	output := []string{}

	for _, n := range lines {
		if n == "" {
			output = append(output, "")
			continue
		}

		var positions map[int]bool
		if colorCode != "" {
			positions = coloredPositions(n, substring)
		}

		runes := []rune(n)
		var rows [8]string
		for i, ch := range runes {
			charLines, ok := banner[ch]
			if !ok {
				return nil, errors.New("unsupported character")
			}
			for r := 0; r < 8; r++ {
				if colorCode != "" && positions[i] {
					rows[r] += colorize(charLines[r], colorCode)
				} else {
					rows[r] += charLines[r]
				}
			}
		}
		for _, row := range rows {
			output = append(output, row)
		}
	}
	return output, nil
}
