# Task 02 - Structured Renderer

Owner: Programmer 2

## Goal

Replace the old flat `BuildASCII` rendering model with structured rendering so output alignment, especially justify, can work without damaging ASCII-art characters.

The renderer still must not know about terminal width, file output, or alignment policy.

## Design

Introduce renderer-owned structures that preserve the relationship between original text segments and rendered ASCII rows:

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

Each segment has exactly 8 rows for non-empty rendered content. Empty normalized input lines should still be represented clearly, so the output package can print the correct blank line.

## Responsibilities

- Replace `BuildASCII` with a structured renderer contract.
- Split each normalized input line into word and space-run segments.
- Preserve user-provided spaces as `SegmentSpace` segments.
- Render each segment from the banner map.
- Preserve color behavior by applying ANSI codes to character blocks based on original input character positions.
- Preserve escaped newline behavior from `NormalizeInput`.
- Return errors for unsupported characters.
- Return existing color validation errors unchanged.
- Keep renderer independent from alignment and terminal width.

## Suggested API

Preferred final API:

```go
func BuildASCII(lines []string, banner map[rune][]string, substring string, color string) (RenderedASCII, error)
```

The name can stay `BuildASCII`, but the return type changes from `[]string` to structured output.

If useful for tests and compatibility inside the renderer package, add a helper:

```go
func Flatten(rendered RenderedASCII) []string
```

`Flatten` should join all segment rows in natural left-aligned order.

## Segment Rules

Input:

```text
"hello   world"
```

Segments:

```text
word("hello")
space("   ")
word("world")
```

Input:

```text
"   hello"
```

Segments:

```text
space("   ")
word("hello")
```

Input:

```text
"hello   "
```

Segments:

```text
word("hello")
space("   ")
```

## Tests

Add or update renderer tests for:

- Structured output contains word and space segments.
- Flat rendering through `Flatten` matches the previous `[]string` output.
- Multiple spaces are preserved as one space-run segment.
- Leading and trailing spaces are preserved.
- Empty normalized lines are preserved.
- Whole-string color still colors every character block.
- Substring color still colors only matching character blocks.
- Missing substring still renders uncolored output.
- Substring matching does not cross normalized line boundaries.
- Unsupported characters return an error.
- Invalid color values return existing error messages.

## Done When

- Renderer returns structured output.
- Existing visual ASCII output can still be reproduced by flattening.
- No alignment logic exists in the renderer package.
