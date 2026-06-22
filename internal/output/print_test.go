package output

import "testing"

func TestPrintASCII_returnsNoError(t *testing.T) {
	lines := []string{"  _  ", " | | "}

	err := PrintASCII(lines)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
