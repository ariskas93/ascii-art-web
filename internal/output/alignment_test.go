package output

import (
	"testing"
)

func TestVisibleWidth(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "plain text",
			input: "hello",
			want:  5,
		},
		{
			name:  "leading spaces",
			input: "  hi",
			want:  4,
		},
		{
			name:  "named ANSI color",
			input: "\033[31mhello\033[0m",
			want:  5,
		},
		{
			name:  "truecolor ANSI color",
			input: "\033[38;2;255;0;0mhello\033[0m",
			want:  5,
		},
		{
			name:  "colored leading spaces",
			input: "\033[36m  hi\033[0m",
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := visibleWidth(tt.input)
			if got != tt.want {
				t.Fatalf("visibleWidth(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
