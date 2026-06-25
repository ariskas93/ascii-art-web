# ASCII-Art Alignment Objective

We are extending the existing ASCII-Art project by implementing the **alignment** optional objective.

The current project already supports:

- argument mode
- interactive mode
- `standard`, `shadow`, and `thinkertoy` banners
- color rendering with `--color=<color>`
- output files with `--output=<file>`

The next project goal is to add:

```bash
go run . --align=left "hello"
go run . --align=center "hello" shadow
go run . "hello" shadow --align=right
go run . --align=justify "hello world"
```

## Current Architecture

The project uses a pipeline:

```text
CLI
-> Banner Loader
-> Normalize Input
-> Renderer
-> Output
```

Current packages:

- `internal/cli`
- `internal/banner`
- `internal/renderer`
- `internal/output`

## Final Design Decision

Alignment will not be implemented inside the renderer.

The renderer will still be responsible for converting text into ASCII-art, but it will no longer flatten each rendered line into plain `[]string` too early.

Instead, `BuildASCII` will return structured rendering data.

The output package will use that structured data to apply:

- left alignment
- center alignment
- right alignment
- justify alignment

This replaces the older idea where output would receive only rendered `[]string` rows.

## Why Structured Rendering

Justify alignment needs to expand only the spaces between words.

If the output package receives only final rendered rows, it cannot safely tell the difference between:

- spaces between words
- spaces inside ASCII-art letters
- spaces from the banner style itself

Trying to justify already-flattened rows would damage the shapes of the ASCII-art characters.

Structured rendering solves this by preserving word and space-run segments.

Example input:

```text
"hello   world"
```

Structured rendering:

```text
word("hello")
space("   ")
word("world")
```

Each segment contains its own 8 rendered ASCII-art rows.

## Proposed Renderer Contract

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
func BuildASCII(lines []string, banner map[rune][]string, substring string, color string) (RenderedASCII, error)
```

The function name can stay `BuildASCII`, but the return type changes from flat rows to structured output.

## CLI Rules

The CLI should parse:

```text
--align=left
--align=center
--align=right
--align=justify
```

Flags may appear in any order, as long as they use the `--flag=value` form.

Examples:

```bash
go run . --align=center "hello"
go run . "hello" --align=center
go run . "hello" shadow --align=right
go run . --output=out.txt --align=justify "hello world"
go run . --color=cyan world "hello world" --align=center
```

Invalid forms:

```bash
go run . --align center "hello"
go run . --align= "hello"
go run . -align=center "hello"
go run . --alignment=center "hello"
```

## Shared Options

The CLI should return a configuration object containing:

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

Default alignment is left.

## Terminal Width

Use only standard Go packages.

Recommended first version:

- read `COLUMNS`
- parse it as an integer
- fallback to `80` when missing or invalid
- keep width detection in the output package
- expose width-explicit functions for unit tests

## Alignment Algorithms

### Left

Print the flattened structured output as-is.

### Center

For each rendered row:

```text
left padding = (terminal width - visible row width) / 2
```

If the row is wider than the terminal width, print it unchanged.

### Right

For each rendered row:

```text
left padding = terminal width - visible row width
```

If the row is wider than the terminal width, print it unchanged.

### Justify

Justify must use structured segments.

Rules:

- Expand only gaps between word segments.
- Do not expand spaces inside rendered ASCII-art letters.
- Do not expand leading or trailing spaces.
- Lines with fewer than two word segments are printed left-aligned.
- If the rendered words already exceed terminal width, print left-aligned.
- Distribute extra gap width from left to right.

## Task Split

Task files:

```text
docs/tasks/task01_cli_alignment.md
docs/tasks/task02_structured_renderer.md
docs/tasks/task03_output_alignment_integration.md
```

Ownership:

- Programmer 1: CLI alignment option and `Options.Align`
- Programmer 2: structured renderer replacement
- Programmer 3: output alignment and main integration

## Edge Cases

- one-word justify lines
- empty normalized lines
- leading spaces
- trailing spaces
- multiple spaces between words
- ANSI color codes in aligned output
- rows wider than terminal width
- missing or invalid `COLUMNS`
- alignment combined with output files
- alignment combined with colored substrings

## Development Rule

Do not implement justify by searching or stretching spaces inside final rendered rows.

Justify must be based on structured word and space segments from the renderer.
