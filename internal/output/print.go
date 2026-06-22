package output

import "os"

func PrintASCII(lines []string) error {
	return WriteASCII(os.Stdout, lines)
}
