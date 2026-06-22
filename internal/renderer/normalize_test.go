package renderer

import (
	"reflect"
	"testing"
)

func TestNormalizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple string",
			input:    "Hello",
			expected: []string{"Hello"},
		},
		{
			name:     "single newline",
			input:    "Hello\\nThere",
			expected: []string{"Hello", "There"},
		},
		{
			name:     "double newline",
			input:    "Hello\\n\\nThere",
			expected: []string{"Hello", "", "There"},
		},
		{
			name:     "only newline",
			input:    "\\n",
			expected: []string{""},
		},
		{
			name:     "trailing newline",
			input:    "Hello\\n",
			expected: []string{"Hello"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "real newline is not split",
			input:    "Hello\nThere",
			expected: []string{"Hello\nThere"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeInput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("NormalizeInput(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
