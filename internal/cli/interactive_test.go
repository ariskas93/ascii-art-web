package cli

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRunInteractiveValidBannerChoices(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantBanner string
	}{
		{
			name:       "standard",
			input:      "hello\n1\n\n",
			wantBanner: "banners/standard.txt",
		},
		{
			name:       "shadow",
			input:      "hello\n2\n\n",
			wantBanner: "banners/shadow.txt",
		},
		{
			name:       "thinkertoy",
			input:      "hello\n3\n\n",
			wantBanner: "banners/thinkertoy.txt",
		},
		{
			name:       "trims banner whitespace",
			input:      "hello\n 2 \n\n",
			wantBanner: "banners/shadow.txt",
		},
		{
			name:       "windows line endings",
			input:      "hello\r\n3\r\n\r\n",
			wantBanner: "banners/thinkertoy.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			gotText, gotBanner, gotColor, gotSubstring, exit, err := runInteractive(strings.NewReader(tt.input), &out)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if exit {
				t.Fatal("exit = true, want false")
			}
			if gotText != "hello" {
				t.Fatalf("text = %q, want %q", gotText, "hello")
			}
			if gotBanner != tt.wantBanner {
				t.Fatalf("banner = %q, want %q", gotBanner, tt.wantBanner)
			}
			if gotColor != "" {
				t.Fatalf("color = %q, want empty", gotColor)
			}
			if gotSubstring != "" {
				t.Fatalf("substring = %q, want empty", gotSubstring)
			}
			if !strings.Contains(out.String(), "Choose the number of the desired banner:") {
				t.Fatalf("expected banner menu in output, got %q", out.String())
			}
		})
	}
}

func TestRunInteractiveRetriesEmptyText(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("\nhello\n1\n\n"), &out)
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
	if color != "" {
		t.Fatalf("color = %q, want empty", color)
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
	if !strings.Contains(out.String(), "No valid string input") {
		t.Fatalf("expected empty text warning in output, got %q", out.String())
	}
}

func TestRunInteractiveAcceptsWhitespaceText(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("   \n1\n\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if text != "   " {
		t.Fatalf("text = %q, want three spaces", text)
	}
	if bannerPath != "banners/standard.txt" {
		t.Fatalf("bannerPath = %q, want %q", bannerPath, "banners/standard.txt")
	}
	if color != "" {
		t.Fatalf("color = %q, want empty", color)
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
}

func TestRunInteractiveRetriesInvalidBannerChoices(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n9\nabc\n\n2\n\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if text != "hello" {
		t.Fatalf("text = %q, want %q", text, "hello")
	}
	if bannerPath != "banners/shadow.txt" {
		t.Fatalf("bannerPath = %q, want %q", bannerPath, "banners/shadow.txt")
	}
	if color != "" {
		t.Fatalf("color = %q, want empty", color)
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}

	const wantWarnings = 3
	if gotWarnings := strings.Count(out.String(), "Invalid banner choice"); gotWarnings != wantWarnings {
		t.Fatalf("invalid banner warnings = %d, want %d; output %q", gotWarnings, wantWarnings, out.String())
	}
}

func TestRunInteractiveExitAtTextPrompt(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("/exit\n"), &out)
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

func TestRunInteractiveExitAtBannerPrompt(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n/exit\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if text != "" {
		t.Fatalf("text = %q, want empty when exiting from banner prompt", text)
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

func TestRunInteractivePreservesEscapedNewlineText(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("Hello\\nThere\n1\n\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if text != "Hello\\nThere" {
		t.Fatalf("text = %q, want %q", text, "Hello\\nThere")
	}
	if bannerPath != "banners/standard.txt" {
		t.Fatalf("bannerPath = %q, want %q", bannerPath, "banners/standard.txt")
	}
	if color != "" {
		t.Fatalf("color = %q, want empty", color)
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
}

func TestRunInteractiveColorWholeString(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n1\nred\n\n"), &out)
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
	if color != "red" {
		t.Fatalf("color = %q, want %q", color, "red")
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
	if !strings.Contains(out.String(), "Enter color") {
		t.Fatalf("expected color prompt in output, got %q", out.String())
	}
	if !strings.Contains(out.String(), "Examples: red, #ff0000, rgb,255,0,0, hsl,0,100%,50%") {
		t.Fatalf("expected color examples in output, got %q", out.String())
	}
	if !strings.Contains(out.String(), "Color: ") {
		t.Fatalf("expected color input label in output, got %q", out.String())
	}
}

func TestRunInteractiveTrimsColorWhitespace(t *testing.T) {
	var out bytes.Buffer

	_, _, color, substring, exit, err := runInteractive(strings.NewReader("hello\n1\n red \n\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if color != "red" {
		t.Fatalf("color = %q, want %q", color, "red")
	}
	if substring != "" {
		t.Fatalf("substring = %q, want empty", substring)
	}
}

func TestRunInteractiveColorSubstring(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("a king kitten have kit\n2\ncyan\nkit\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if text != "a king kitten have kit" {
		t.Fatalf("text = %q, want %q", text, "a king kitten have kit")
	}
	if bannerPath != "banners/shadow.txt" {
		t.Fatalf("bannerPath = %q, want %q", bannerPath, "banners/shadow.txt")
	}
	if color != "cyan" {
		t.Fatalf("color = %q, want %q", color, "cyan")
	}
	if substring != "kit" {
		t.Fatalf("substring = %q, want %q", substring, "kit")
	}
	if !strings.Contains(out.String(), "Enter substring to color") {
		t.Fatalf("expected substring prompt in output, got %q", out.String())
	}
}

func TestRunInteractivePreservesSubstringSpaces(t *testing.T) {
	var out bytes.Buffer

	_, _, color, substring, exit, err := runInteractive(strings.NewReader("hello world\n1\nyellow\n world \n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if color != "yellow" {
		t.Fatalf("color = %q, want %q", color, "yellow")
	}
	if substring != " world " {
		t.Fatalf("substring = %q, want %q", substring, " world ")
	}
}

func TestRunInteractiveOptionsDefaultPrintsToTerminal(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\n\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}

	want := Options{
		Text:       "hello",
		BannerPath: "banners/standard.txt",
	}
	if opts != want {
		t.Fatalf("options = %+v, want %+v", opts, want)
	}
	if !strings.Contains(out.String(), "Write output to a file?") {
		t.Fatalf("expected output prompt in output, got %q", out.String())
	}
}

func TestRunInteractiveOptionsNoOutputFile(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n2\n\nno\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}

	want := Options{
		Text:       "hello",
		BannerPath: "banners/shadow.txt",
	}
	if opts != want {
		t.Fatalf("options = %+v, want %+v", opts, want)
	}
}

func TestRunInteractiveOptionsOutputFile(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\ny\nbanner.txt\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}

	want := Options{
		Text:       "hello",
		BannerPath: "banners/standard.txt",
		OutputFile: "banner.txt",
	}
	if opts != want {
		t.Fatalf("options = %+v, want %+v", opts, want)
	}
	if !strings.Contains(out.String(), "Enter output file name") {
		t.Fatalf("expected output file prompt in output, got %q", out.String())
	}
}

func TestRunInteractiveOptionsColorSubstringOutputFile(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello world\n2\ncyan\nworld\ny\nbanner.txt\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}

	want := Options{
		Text:       "hello world",
		BannerPath: "banners/shadow.txt",
		Color:      "cyan",
		Substring:  "world",
		OutputFile: "banner.txt",
	}
	if opts != want {
		t.Fatalf("options = %+v, want %+v", opts, want)
	}
}

func TestRunInteractiveOptionsRetriesInvalidOutputChoice(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\nmaybe\nn\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if opts.OutputFile != "" {
		t.Fatalf("OutputFile = %q, want empty", opts.OutputFile)
	}
	if !strings.Contains(out.String(), "Invalid output choice") {
		t.Fatalf("expected invalid output warning in output, got %q", out.String())
	}
}

func TestRunInteractiveOptionsRetriesEmptyOutputFile(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\ny\n\nbanner.txt\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exit {
		t.Fatal("exit = true, want false")
	}
	if opts.OutputFile != "banner.txt" {
		t.Fatalf("OutputFile = %q, want %q", opts.OutputFile, "banner.txt")
	}
	if !strings.Contains(out.String(), "No valid output file name") {
		t.Fatalf("expected empty output file warning in output, got %q", out.String())
	}
}

func TestRunInteractiveOptionsExitAtOutputPrompt(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\n/exit\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if opts != (Options{}) {
		t.Fatalf("options = %+v, want empty", opts)
	}
}

func TestRunInteractiveOptionsExitAtOutputFilePrompt(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\ny\n/exit\n"), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if opts != (Options{}) {
		t.Fatalf("options = %+v, want empty", opts)
	}
}

func TestRunInteractiveOptionsReturnsErrorOnEOFBeforeOutputChoice(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\n"), &out)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("err = %v, want EOF", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if opts != (Options{}) {
		t.Fatalf("options = %+v, want empty", opts)
	}
}

func TestRunInteractiveOptionsReturnsErrorOnEOFBeforeOutputFile(t *testing.T) {
	var out bytes.Buffer

	opts, exit, err := runInteractiveOptions(strings.NewReader("hello\n1\n\ny\n"), &out)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("err = %v, want EOF", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if opts != (Options{}) {
		t.Fatalf("options = %+v, want empty", opts)
	}
}

func TestRunInteractiveExitAtColorPrompt(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n1\n/exit\n"), &out)
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

func TestRunInteractiveExitAtSubstringPrompt(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n1\nred\n/exit\n"), &out)
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

func TestRunInteractiveReturnsErrorOnEOFBeforeText(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader(""), &out)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("err = %v, want EOF", err)
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

func TestRunInteractiveReturnsErrorOnEOFBeforeBanner(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n"), &out)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("err = %v, want EOF", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if text != "" {
		t.Fatalf("text = %q, want empty after banner EOF", text)
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

func TestRunInteractiveReturnsErrorOnEOFBeforeColor(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n1\n"), &out)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("err = %v, want EOF", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if text != "" {
		t.Fatalf("text = %q, want empty after color EOF", text)
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

func TestRunInteractiveReturnsErrorOnEOFBeforeSubstring(t *testing.T) {
	var out bytes.Buffer

	text, bannerPath, color, substring, exit, err := runInteractive(strings.NewReader("hello\n1\nred\n"), &out)
	if err == nil {
		t.Fatal("expected EOF error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("err = %v, want EOF", err)
	}
	if !exit {
		t.Fatal("exit = false, want true")
	}
	if text != "" {
		t.Fatalf("text = %q, want empty after substring EOF", text)
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
