# PRD - ASCII-ART Alignment

---

## 1. Problem Statement

ASCII-ART is a Go CLI program that renders text as multi-line ASCII art using one of three banner files.

The project already supports:

- argument mode
- interactive mode
- banner selection
- color rendering
- output file writing

This project iteration adds the alignment optional objective:

```bash
go run . --align=<left|center|right|justify> [STRING] [BANNER]
```

The program must preserve compatibility with the original audit behavior:

```bash
go run . [STRING] [BANNER]
```

and the color behavior:

```bash
go run . --color=<color> [substring] [STRING]
```

Allowed packages:

- Only standard Go packages are allowed.

---

## 2. Supported Features

The program supports:

- `standard`, `shadow`, and `thinkertoy` banners
- single-string rendering with default `standard` banner
- escaped newline sequences written as `\n`
- printable ASCII characters supported by the banner files
- foreground color with named colors, `#RRGGBB`, RGB, and HSL notation
- optional file output with `--output=<file>`
- optional alignment with `--align=<left|center|right|justify>`

The program must not support external dependencies or custom user-supplied banner paths.

---

## 3. CLI Contract

### Argument Mode

Valid base forms:

```bash
go run . "Hello"
go run . "Hello" standard
go run . "Hello" shadow
go run . "Hello\nThere" thinkertoy.txt
```

Valid option examples:

```bash
go run . --color=red "Hello"
go run . --color=red ll "Hello"
go run . --output=out.txt "Hello"
go run . --align=center "Hello"
go run . "Hello" shadow --align=right
go run . --output=out.txt --align=justify "Hello world"
go run . "Hello world" shadow --color=cyan world --align=center
```

Options may appear before, between, or after positional arguments, as long as they use `--flag=value` form.

Supported option flags:

```text
--color=<color>
--output=<file>
--align=<left|center|right|justify>
```

Invalid option forms:

```bash
go run . --color red "hello"
go run . --output out.txt "hello"
go run . --align center "hello"
go run . --align= "hello"
go run . -align=center "hello"
go run . --alignment=center "hello"
go run . --align=CENTER "hello"
```

Alignment values are case-sensitive and must be lowercase.

Valid:

```text
--align=left
--align=center
--align=right
--align=justify
```

Invalid:

```text
--align=Left
--align=CENTER
--align=RIGHT
--align=Justify
```

### Positional Arguments

After removing supported flags:

```text
[STRING]
[STRING] [BANNER]
```

are valid base forms.

Accepted banner values:

```text
standard
standard.txt
shadow
shadow.txt
thinkertoy
thinkertoy.txt
```

If only one positional argument exists, it is always text, even when it equals a banner name.

Example:

```bash
go run . "shadow"
```

renders the word `shadow`.

To use a banner, provide it as the second positional argument:

```bash
go run . "shadow" standard
```

### Color Substring Rule

The project does not add a `--substring` flag.

Substring coloring stays compatible with the exercise:

```bash
go run . --color=red kit "a king kitten"
```

A color substring is positional and is recognized only when it appears immediately after the `--color=<color>` flag.

### Usage Messages

General invalid argument format:

```text
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

Invalid color option format:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

Invalid output option format:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```

Invalid alignment option format:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --align=<left|center|right|justify> "something" standard
```

Unknown option format:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

OPTIONS:
--color=<color>
--output=<fileName.txt>
--align=<left|center|right|justify>

EX:
go run . --color=red "hello"
go run . --output=out.txt "hello" standard
go run . --align=center "hello" standard
```

Invalid banner format:

```text
Invalid banner: <banner>

Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

---

## 4. Interactive Mode

Run:

```bash
go run .
```

The program loops until `/exit`.

Interactive mode collects:

1. Text
2. Banner selection
3. Optional color
4. Optional substring when color is provided
5. Optional output file
6. Optional alignment

Interactive defaults:

- empty color means no color
- empty substring after a color means color the whole string
- empty output choice means print to terminal
- empty alignment means left alignment

Invalid interactive choices should repeat the relevant prompt instead of crashing.

---

## 5. Shared Options

The CLI returns one options value:

```go
type Options struct {
	Text       string
	BannerPath string
	Color      string
	Substring  string
	OutputFile string
	Align      string
}
```

Meanings:

- `BannerPath` defaults to `banners/standard.txt`
- `Color == ""` means no color
- `Substring == ""` with a color means color the whole string
- `OutputFile == ""` means print to terminal
- `Align == "left"` is the default alignment

---

## 6. Rendering Contract

### Banner Loading

```go
LoadBanner(path string) (map[rune][]string, error)
```

Responsibilities:

- read the selected banner file
- normalize CRLF line endings
- map printable ASCII characters from 32 to 126
- return an error for missing or malformed files

### Input Normalization

```go
NormalizeInput(input string) []string
```

Examples:

```text
"Hello\nThere" -> ["Hello", "There"]
"Hello\n\nThere" -> ["Hello", "", "There"]
"\nHello" -> ["", "Hello"]
```

### Structured ASCII Rendering

The old flat renderer returned:

```go
[]string
```

The alignment version replaces that with structured rendering:

```go
type SegmentKind int

const (
	SegmentWord SegmentKind = iota
	SegmentSpace
)

type RenderedSegment struct {
	Kind SegmentKind
	Text string
	Rows []string
}

type RenderedLine struct {
	Segments []RenderedSegment
}

type RenderedASCII []RenderedLine
```

Preferred final renderer API:

```go
BuildASCII(lines []string, banner map[rune][]string, substring string, color string) (RenderedASCII, error)
```

Renderer responsibilities:

- split normalized lines into word and space-run segments
- render each segment using the banner map
- preserve leading, trailing, and repeated spaces
- preserve color support
- preserve newline behavior
- return errors for unsupported characters
- remain independent from output alignment and terminal width

### Color Helpers

```go
getANSIColor(color string) (string, error)
colorize(text string, colorCode string) string
coloredPositions(s string, substring string) map[int]bool
```

Color matching is applied per normalized line and must not cross newline boundaries.

---

## 7. Output and Alignment Contract

The output package owns printing, file writing, terminal width detection, and alignment.

Alignment modes:

```text
left
center
right
justify
```

Terminal width detection:

- use only standard Go packages
- detect the current terminal width each time output is produced
- try standard output, standard input, and standard error because individual
  streams may be redirected
- use Linux standard-library system calls for live detection on Linux and WSL
- keep a non-Linux implementation that reports live detection as unavailable,
  allowing the fallback chain to continue
- use a valid `COLUMNS` value when live terminal detection fails
- fallback to `80` when live detection and `COLUMNS` are unavailable or invalid
- unit-test alignment through width-explicit functions

Terminal output and `--output=<file>` use the same current-terminal width
policy. When writing aligned output to a file, alignment is calculated using
the width of the terminal in which the command is running. The resulting
spaces are fixed in the file. The program does not try to detect or predict
the width of the editor, terminal, or viewer that opens the file later.

Visible width must ignore ANSI escape sequences.

### Left

Print flattened rows unchanged.

### Center

Pad each row with:

```text
(width - visibleRowWidth) / 2
```

Rows wider than width receive no padding.

### Right

Pad each row with:

```text
width - visibleRowWidth
```

Rows wider than width receive no padding.

### Justify

Justify must use structured segments.

Rules:

- expand only gaps between word segments
- never expand spaces inside ASCII-art letters
- do not use leading or trailing spaces as expandable gaps
- one-word lines are left-aligned
- rows wider than width are left-aligned
- distribute extra gap width from left to right

---

## 8. Processing Flow

```text
cli.PickModeOptions
-> banner.LoadBanner
-> renderer.NormalizeInput
-> renderer.BuildASCII
-> output.PrintAlignedASCII or output.WriteAlignedASCIIFile
```

`main.go` is the integration layer only.

---

## 9. Task Split

Task files:

```text
docs/tasks/task01_cli_alignment.md
docs/tasks/task02_structured_renderer.md
docs/tasks/task03_output_alignment_integration.md
```

Ownership:

- Programmer 1: CLI alignment option and interactive alignment prompt
- Programmer 2: structured renderer replacement
- Programmer 3: output alignment and main integration

---

## 10. Acceptance Criteria

### CLI

- [ ] `--align=left`, `--align=center`, `--align=right`, and `--align=justify` are accepted.
- [ ] Alignment flags can appear in any position.
- [ ] Invalid alignment flags print the exact alignment usage message.
- [ ] Duplicate alignment flags are rejected.
- [ ] Existing color and output flags still work.
- [ ] Single-string mode still works.

### Renderer

- [ ] `BuildASCII` returns structured output.
- [ ] Word and space segments are preserved.
- [ ] Existing visual output can be reproduced by flattening structured output.
- [ ] Color rendering still works.
- [ ] Escaped newlines still work.
- [ ] Unsupported characters return an error.

### Output

- [ ] Left alignment matches existing output.
- [ ] Center alignment pads correctly.
- [ ] Right alignment pads correctly.
- [ ] Justify expands only spaces between word segments.
- [ ] One-word justify lines are not stretched.
- [ ] ANSI escape codes do not count toward visible width.
- [ ] File output receives the same aligned rows as terminal output.

### Integration

- [ ] `go test ./...` passes.
- [ ] Alignment works with all three banners.
- [ ] Alignment works with color.
- [ ] Alignment works with output files.
- [ ] Interactive mode can select alignment.

---

## 11. Non-Goals

- No custom banner paths.
- No custom banner editing.
- No multiple color flags.
- No multiple independent input strings.
- No Unicode rendering outside printable ASCII banner support.
- No web or GUI mode.
- No resizing or scaling of ASCII-art character height.
- No justify implementation based on mutating already-flattened rendered rows.

---

## 12. Risks

- Justify could distort ASCII letters if implemented on flat rows.
- ANSI color codes can break width calculations if counted as visible characters.
- Terminal width may be unavailable in some environments.
- Parser changes can accidentally break original audit commands.
- Structured rendering can duplicate old renderer logic if not designed carefully.

Mitigation:

- Use structured word and space segments.
- Test visible-width calculations with ANSI strings.
- Fallback to width `80`.
- Keep existing CLI tests and add alignment cases.
- Add flattening tests to compare structured rendering to previous output.
