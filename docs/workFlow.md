# Workflow - ASCII Art Alignment

## Objective

Add alignment support to the existing ASCII-Art program:

```bash
go run . --align=left "hello"
go run . --align=center "hello" standard
go run . --align=right "hello" shadow
go run . --align=justify "hello world" thinkertoy
```

The program must keep supporting:

```bash
go run . "hello"
go run . "hello" shadow
go run . --color=red ll "hello"
go run . --output=out.txt "hello"
```

Alignment must work together with color, banner selection, escaped newlines, interactive mode, and output files.

## Shared Contract

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

`Align == "left"` means left alignment.

Accepted alignment values:

```text
left
center
right
justify
```

## Architecture Decision

The project will replace flat rendering with structured rendering.

Old renderer output:

```go
[]string
```

New renderer output:

```go
type RenderedASCII []RenderedLine
```

The renderer is still responsible only for converting input text into ASCII-art. It does not know about alignment or terminal width.

The output package receives the structured output and decides how to print or write it.

## Why This Matters

Left, center, and right alignment can work on flat rows.

Justify cannot safely work on flat rows because spaces inside ASCII-art letters look the same as spaces between words.

Therefore, the renderer must preserve segments:

```text
word("hello")
space(" ")
word("world")
```

Then the output package can expand only the real word gaps.

## Task Split

### Task 01 - CLI Alignment Option

Owner: Programmer 1

File: `docs/tasks/task01_cli_alignment.md`

Responsibilities:

- Add `Align` to `cli.Options`.
- Parse `--align=<left|center|right|justify>`.
- Keep flags order-flexible.
- Reject malformed and duplicate alignment flags.
- Add interactive alignment prompt.
- Update CLI tests.

### Task 02 - Structured Renderer

Owner: Programmer 2

File: `docs/tasks/task02_structured_renderer.md`

Responsibilities:

- Change `BuildASCII` to return structured output.
- Preserve word and space-run segments.
- Keep color rendering behavior.
- Keep newline behavior.
- Add renderer tests for segment structure and flattening.

### Task 03 - Output Alignment and Main Integration

Owner: Programmer 3

File: `docs/tasks/task03_output_alignment_integration.md`

Responsibilities:

- Add output functions for aligned writing.
- Implement left, center, right, and justify.
- Detect terminal width using standard packages.
- Ignore ANSI escape sequences when measuring visible width.
- Update `main.go` for structured output and alignment.
- Add output and integration tests.

## Implementation Order

1. Programmer 1 adds `Align` parsing and tests.
2. Programmer 2 replaces flat rendering with structured rendering.
3. Programmer 3 adapts output and `main.go` to structured rendering.
4. Run package tests while each task is in progress.
5. Run `go test ./...`.
6. Run manual checks for all alignments with standard, shadow, color, and file output.

## Manual Checks

```bash
go run . --align=left "hello" standard
go run . --align=center "hello" standard
go run . --align=right "hello" standard
go run . --align=justify "hello world" standard
go run . "hello world" shadow --align=center --color=cyan world
go run . --output=aligned.txt --align=right "hello" standard
```

## Exact Alignment Usage Message

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --align=<left|center|right|justify> "something" standard
```
