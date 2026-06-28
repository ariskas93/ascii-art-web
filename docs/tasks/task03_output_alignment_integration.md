# Task 03 - Output Alignment and Main Integration

Owner: Programmer 3

## Goal

Use the structured renderer output to implement terminal/file alignment in `internal/output` and connect it in `main.go`.

The output package owns alignment. The renderer owns ASCII-art conversion. `main.go` only coordinates.

## Supported Alignments

```text
left
center
right
justify
```

Default alignment is left.

## Terminal Width

Use only Go standard packages.

Recommended first implementation:

- Read `COLUMNS` from the environment.
- If `COLUMNS` is missing, invalid, or less than 1, fall back to `80`.
- Keep width detection at the edge of the output package so unit tests can pass explicit widths.
- File output uses the current terminal width at the time the command runs; the program does not try to predict the width of a future text editor or viewer.

Suggested API shape:

```go
func WriteAlignedASCII(w io.Writer, rendered renderer.RenderedASCII, align string, width int) error
func PrintAlignedASCII(rendered renderer.RenderedASCII, align string) error
func WriteAlignedASCIIFile(fileName string, rendered renderer.RenderedASCII, align string, width int) error
```

Exact names can change, but tests should be able to call a width-explicit function.

## Alignment Rules

Visible width must ignore ANSI escape sequences.

Left:

- Print flattened rows unchanged.

Center:

- Add `floor((width - visibleRowWidth) / 2)` spaces before each row.
- If the row is wider than the width, add no padding.

Right:

- Add `width - visibleRowWidth` spaces before each row.
- If the row is wider than the width, add no padding.

Justify:

- Use `RenderedLine.Segments`, not raw flattened rows.
- Expand only gaps between word segments.
- Do not expand spaces inside ASCII-art letters.
- Lines with fewer than two word segments are printed left-aligned.
- Leading and trailing space segments are preserved but are not used as expandable gaps.
- If the rendered words are already wider than the width, print the line left-aligned.
- Distribute extra spaces across expandable gaps from left to right.

## Main Integration

Update `main.go` after Programmer 2 changes the renderer return type:

```text
cli.PickModeOptions
-> banner.LoadBanner
-> renderer.NormalizeInput
-> renderer.BuildASCII
-> output.PrintAlignedASCII or output.WriteAlignedASCIIFile
```

Respect:

- `Options.OutputFile`
- `Options.Align`
- interactive repeat behavior
- rendering errors returning to the prompt only in interactive mode

## Tests

Add or update output/integration tests for:

- Left alignment preserves flattened output.
- Center alignment pads rows correctly.
- Right alignment pads rows correctly.
- Rows wider than width are not padded.
- ANSI color codes do not count toward visible width.
- Justify expands only between word segments.
- Justify leaves one-word lines unchanged.
- Justify preserves leading/trailing spaces.
- File output writes aligned rows.
- `PrintAlignedASCII` uses detected width while `WriteAlignedASCII` can be tested with explicit width.

## Manual Checks

Run:

```bash
go test ./...
go run . --align=left "hello" standard
go run . --align=center "hello" standard
go run . --align=right "hello" standard
go run . --align=justify "hello world" standard
go run . "hello world" shadow --align=center --color=cyan world
go run . --output=aligned.txt --align=right "hello" standard
```

## Done When

- All four alignment modes work from structured renderer output.
- File output and terminal output share the same alignment behavior.
- Existing color, output, banner, and newline behavior still works.
