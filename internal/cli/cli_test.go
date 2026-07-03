// Package cli tests argument parsing and usage-message contracts.
package cli

import (
	"fmt"
	"os"
	"testing"
)

func TestRunArgumentMode(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantText      string
		wantBanner    string
		wantColor     string
		wantSubstring string
		wantErr       string
	}{
		{
			name:       "text only defaults to standard",
			args:       []string{"ascii-art", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
		},
		{
			name:       "whitespace text is valid",
			args:       []string{"ascii-art", "   "},
			wantText:   "   ",
			wantBanner: "banners/standard.txt",
		},
		{
			name:       "escaped newline text is preserved",
			args:       []string{"ascii-art", "Hello\\nThere"},
			wantText:   "Hello\\nThere",
			wantBanner: "banners/standard.txt",
		},
		{
			name:       "standard without extension",
			args:       []string{"ascii-art", "hello", "standard"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
		},
		{
			name:       "standard with extension",
			args:       []string{"ascii-art", "hello", "standard.txt"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
		},
		{
			name:       "shadow without extension",
			args:       []string{"ascii-art", "hello", "shadow"},
			wantText:   "hello",
			wantBanner: "banners/shadow.txt",
		},
		{
			name:       "shadow with extension",
			args:       []string{"ascii-art", "hello", "shadow.txt"},
			wantText:   "hello",
			wantBanner: "banners/shadow.txt",
		},
		{
			name:       "thinkertoy without extension",
			args:       []string{"ascii-art", "hello", "thinkertoy"},
			wantText:   "hello",
			wantBanner: "banners/thinkertoy.txt",
		},
		{
			name:       "thinkertoy with extension",
			args:       []string{"ascii-art", "hello", "thinkertoy.txt"},
			wantText:   "hello",
			wantBanner: "banners/thinkertoy.txt",
		},
		{
			name:       "whole string named color",
			args:       []string{"ascii-art", "--color=red", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "red",
		},
		{
			name:       "whole string hex color",
			args:       []string{"ascii-art", "--color=#ff0000", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "#ff0000",
		},
		{
			name:       "whole string uppercase hex color",
			args:       []string{"ascii-art", "--color=#FF0000", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "#FF0000",
		},
		{
			name:       "whole string mixed case hex color",
			args:       []string{"ascii-art", "--color=#00Ff88", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "#00Ff88",
		},
		{
			name:       "whole string RGB color",
			args:       []string{"ascii-art", "--color=RGB,255,0,0", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "RGB,255,0,0",
		},
		{
			name:       "whole string HSL color",
			args:       []string{"ascii-art", "--color=hsl,0,100%,50%", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "hsl,0,100%,50%",
		},
		{
			name:          "substring named color",
			args:          []string{"ascii-art", "--color=orange", "GuYs", "HeY GuYs"},
			wantText:      "HeY GuYs",
			wantBanner:    "banners/standard.txt",
			wantColor:     "orange",
			wantSubstring: "GuYs",
		},
		{
			name:          "one character substring",
			args:          []string{"ascii-art", "--color=blue", "B", "RGB()"},
			wantText:      "RGB()",
			wantBanner:    "banners/standard.txt",
			wantColor:     "blue",
			wantSubstring: "B",
		},
		{
			name:       "banner compatible whole string color",
			args:       []string{"ascii-art", "hello", "shadow", "--color=cyan"},
			wantText:   "hello",
			wantBanner: "banners/shadow.txt",
			wantColor:  "cyan",
		},
		{
			name:          "banner compatible substring color",
			args:          []string{"ascii-art", "hello world", "thinkertoy", "--color=green", "world"},
			wantText:      "hello world",
			wantBanner:    "banners/thinkertoy.txt",
			wantColor:     "green",
			wantSubstring: "world",
		},
		{
			name:       "color flag after text colors whole string",
			args:       []string{"ascii-art", "hello", "--color=red"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "red",
		},
		{
			name:       "color flag between text and banner colors whole string",
			args:       []string{"ascii-art", "hello", "--color=red", "shadow"},
			wantText:   "hello",
			wantBanner: "banners/shadow.txt",
			wantColor:  "red",
		},
		{
			name:       "banner compatible hex color",
			args:       []string{"ascii-art", "hello", "standard", "--color=#ffaa00"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
			wantColor:  "#ffaa00",
		},
		{
			name:       "alignment flag is accepted in tuple mode",
			args:       []string{"ascii-art", "--align=center", "hello"},
			wantText:   "hello",
			wantBanner: "banners/standard.txt",
		},
		{
			name:    "missing text",
			args:    []string{"ascii-art"},
			wantErr: UsageMessage,
		},
		{
			name:    "empty args",
			args:    []string{},
			wantErr: UsageMessage,
		},
		{
			name:    "empty text",
			args:    []string{"ascii-art", ""},
			wantErr: UsageMessage,
		},
		{
			name:    "invalid banner name",
			args:    []string{"ascii-art", "hello", "unknown"},
			wantErr: invalidBannerMessage("unknown"),
		},
		{
			name:    "numeric banner is invalid in argument mode",
			args:    []string{"ascii-art", "hello", "1"},
			wantErr: invalidBannerMessage("1"),
		},
		{
			name:    "too many plain arguments",
			args:    []string{"ascii-art", "hello", "standard", "extra"},
			wantErr: UsageMessage,
		},
		{
			name:    "separated color flag is invalid",
			args:    []string{"ascii-art", "--color", "red", "banana"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "empty color value is invalid",
			args:    []string{"ascii-art", "--color=", "banana"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "empty color value with substring is invalid",
			args:    []string{"ascii-art", "--color=", "ban", "banana"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "single color flag without text is invalid",
			args:    []string{"ascii-art", "--color=red"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "separated color flag without value is invalid",
			args:    []string{"ascii-art", "--color", "hello"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "short color flag spelling is invalid",
			args:    []string{"ascii-art", "-color=red", "banana"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "British color flag spelling is invalid",
			args:    []string{"ascii-art", "--colour=red", "banana"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "multiple color flags are invalid",
			args:    []string{"ascii-art", "--color=red", "--color=blue", "hello"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "multiple color flags in banner compatible form are invalid",
			args:    []string{"ascii-art", "hello", "standard", "--color=red", "--color=blue"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "invalid color flag in banner compatible form",
			args:    []string{"ascii-art", "hello", "shadow", "--color="},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "missing text in banner compatible color form is invalid",
			args:    []string{"ascii-art", "", "shadow", "--color=red"},
			wantErr: UsageMessage,
		},
		{
			name:    "invalid banner in banner compatible color form is invalid",
			args:    []string{"ascii-art", "hello", "unknown", "--color=red"},
			wantErr: invalidBannerMessage("unknown"),
		},
		{
			name:    "color form with too many arguments is invalid",
			args:    []string{"ascii-art", "--color=red", "he", "hello", "extra"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "shell split RGB color is invalid",
			args:    []string{"ascii-art", "--color=rgb,255,", "0,", "0", "hello"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "shell split HSL color is invalid",
			args:    []string{"ascii-art", "--color=hsl,0,", "100%,", "50%", "hello"},
			wantErr: ColorUsageMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotBanner, gotColor, gotSubstring, err := RunArgumentMode(tt.args)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotText != tt.wantText {
				t.Fatalf("text = %q, want %q", gotText, tt.wantText)
			}
			if gotBanner != tt.wantBanner {
				t.Fatalf("banner = %q, want %q", gotBanner, tt.wantBanner)
			}
			if gotColor != tt.wantColor {
				t.Fatalf("color = %q, want %q", gotColor, tt.wantColor)
			}
			if gotSubstring != tt.wantSubstring {
				t.Fatalf("substring = %q, want %q", gotSubstring, tt.wantSubstring)
			}
		})
	}
}

func TestRunArgumentModeOptions(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Options
		wantErr string
	}{
		{
			name: "text only defaults output file to empty",
			args: []string{"ascii-art", "hello"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output with default banner",
			args: []string{"ascii-art", "--output=banner.txt", "hello"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output with explicit standard banner",
			args: []string{"ascii-art", "--output=banner.txt", "hello", "standard"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output with explicit shadow banner",
			args: []string{"ascii-art", "--output=banner.txt", "hello", "shadow"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/shadow.txt",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output with whole string color",
			args: []string{"ascii-art", "--output=banner.txt", "--color=red", "hello"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				Color:      "red",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output with substring color",
			args: []string{"ascii-art", "--output=banner.txt", "--color=red", "kit", "a king kitten have kit"},
			want: Options{
				Text:       "a king kitten have kit",
				BannerPath: "banners/standard.txt",
				Color:      "red",
				Substring:  "kit",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output with banner compatible substring color",
			args: []string{"ascii-art", "--output=banner.txt", "hello world", "shadow", "--color=cyan", "world"},
			want: Options{
				Text:       "hello world",
				BannerPath: "banners/shadow.txt",
				Color:      "cyan",
				Substring:  "world",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output flag after text is valid",
			args: []string{"ascii-art", "hello", "--output=banner.txt"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name: "output and color flags can appear after banner",
			args: []string{"ascii-art", "hello world", "shadow", "--output=banner.txt", "--color=cyan", "world"},
			want: Options{
				Text:       "hello world",
				BannerPath: "banners/shadow.txt",
				Color:      "cyan",
				Substring:  "world",
				OutputFile: "banner.txt",
				Align:      defaultAlign,
			},
		},
		{
			name:    "separated output flag is invalid",
			args:    []string{"ascii-art", "--output", "hello"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "empty output file is invalid",
			args:    []string{"ascii-art", "--output=", "hello"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "short output flag spelling is invalid",
			args:    []string{"ascii-art", "-output=banner.txt", "hello"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "plural output flag spelling is invalid",
			args:    []string{"ascii-art", "--outputs=banner.txt", "hello"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "split output file argument is invalid",
			args:    []string{"ascii-art", "--output", "banner.txt", "hello"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "output flag without text is invalid",
			args:    []string{"ascii-art", "--output=banner.txt"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "multiple output flags are invalid",
			args:    []string{"ascii-art", "--output=one.txt", "--output=two.txt", "hello"},
			wantErr: OutputUsageMessage,
		},
		{
			name:    "invalid color flag with output keeps color usage",
			args:    []string{"ascii-art", "--output=banner.txt", "--color", "red", "hello"},
			wantErr: ColorUsageMessage,
		},
		{
			name:    "invalid banner with output keeps general usage",
			args:    []string{"ascii-art", "--output=banner.txt", "hello", "unknown"},
			wantErr: invalidBannerMessage("unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunArgumentModeOptions(tt.args)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("options = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRunArgumentModeOptionsAlignment(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Options
		wantErr string
	}{
		{
			name: "left alignment",
			args: []string{"ascii-art", "--align=left", "hello"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				Align:      "left",
			},
		},
		{
			name: "center alignment before text",
			args: []string{"ascii-art", "--align=center", "hello"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				Align:      "center",
			},
		},
		{
			name: "right alignment after text",
			args: []string{"ascii-art", "hello", "--align=right"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				Align:      "right",
			},
		},
		{
			name: "justify alignment after banner",
			args: []string{"ascii-art", "hello world", "shadow", "--align=justify"},
			want: Options{
				Text:       "hello world",
				BannerPath: "banners/shadow.txt",
				Align:      "justify",
			},
		},
		{
			name: "alignment with output",
			args: []string{"ascii-art", "--output=banner.txt", "--align=center", "hello"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				OutputFile: "banner.txt",
				Align:      "center",
			},
		},
		{
			name: "alignment with color substring",
			args: []string{"ascii-art", "--color=cyan", "ll", "hello", "--align=center"},
			want: Options{
				Text:       "hello",
				BannerPath: "banners/standard.txt",
				Color:      "cyan",
				Substring:  "ll",
				Align:      "center",
			},
		},
		{
			name:    "separated alignment flag is invalid",
			args:    []string{"ascii-art", "--align", "center", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "empty alignment value is invalid",
			args:    []string{"ascii-art", "--align=", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "short alignment flag is invalid",
			args:    []string{"ascii-art", "-align=center", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "long alignment spelling is invalid",
			args:    []string{"ascii-art", "--alignment=center", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "unsupported alignment value is invalid",
			args:    []string{"ascii-art", "--align=middle", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "uppercase alignment value is invalid",
			args:    []string{"ascii-art", "--align=CENTER", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "duplicate alignment flags are invalid",
			args:    []string{"ascii-art", "--align=center", "--align=right", "hello"},
			wantErr: AlignUsageMessage,
		},
		{
			name:    "unknown option with value shows option usage",
			args:    []string{"ascii-art", "--banana=yes", "hello"},
			wantErr: OptionUsageMessage,
		},
		{
			name:    "unknown option without value shows option usage",
			args:    []string{"ascii-art", "--banana", "hello"},
			wantErr: OptionUsageMessage,
		},
		{
			name:    "valid alignment without text returns general usage",
			args:    []string{"ascii-art", "--align=center"},
			wantErr: UsageMessage,
		},
		{
			name:    "valid alignment with invalid banner returns banner error",
			args:    []string{"ascii-art", "--align=center", "hello", "unknown"},
			wantErr: invalidBannerMessage("unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunArgumentModeOptions(tt.args)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("options = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestPickModeUsesArgumentMode(t *testing.T) {
	text, bannerPath, color, substring, exit, err := PickMode([]string{"ascii-art", "--color=#ff0000", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if text != "hello" {
		t.Fatalf("text = %q, want %q", text, "hello")
	}
	if bannerPath != "banners/standard.txt" {
		t.Fatalf("bannerPath = %q, want %q", bannerPath, "banners/standard.txt")
	}
	if color != "#ff0000" {
		t.Fatalf("color = %q, want %q", color, "#ff0000")
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
}

func TestPickModeOptionsUsesArgumentMode(t *testing.T) {
	opts, exit, err := PickModeOptions([]string{"ascii-art", "--output=banner.txt", "--color=#ff0000", "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	want := Options{
		Text:       "hello",
		BannerPath: "banners/standard.txt",
		Color:      "#ff0000",
		OutputFile: "banner.txt",
		Align:      defaultAlign,
	}
	if opts != want {
		t.Fatalf("options = %+v, want %+v", opts, want)
	}
}

func TestPickModeUsesInteractiveMode(t *testing.T) {
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	t.Cleanup(func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	})

	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdin pipe: %v", err)
	}
	t.Cleanup(func() {
		stdinReader.Close()
	})

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open dev null: %v", err)
	}
	t.Cleanup(func() {
		devNull.Close()
	})

	os.Stdin = stdinReader
	os.Stdout = devNull

	go func() {
		fmt.Fprintln(stdinWriter, "/exit")
		stdinWriter.Close()
	}()

	text, bannerPath, color, substring, exit, err := PickMode([]string{"ascii-art"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if text != "" {
		t.Fatalf("text = %q, want empty", text)
	}
	if bannerPath != "" {
		t.Fatalf("bannerPath = %q, want empty", bannerPath)
	}
	if color != "" {
		t.Fatalf("color = %q, want empty", color)
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
}
