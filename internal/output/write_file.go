package output

import "os"

func WriteASCIIFile(fileName string, lines []string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteASCII(f, lines)
}
