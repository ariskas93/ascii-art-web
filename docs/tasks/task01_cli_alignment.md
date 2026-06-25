# Task 01 - CLI Alignment Option

Owner: Programmer 1

## Goal

Extend the existing options parser so alignment is part of the shared CLI configuration:

```bash
go run . --align=left "hello"
go run . --align=center "hello" shadow
go run . "hello" shadow --align=right
go run . --output=out.txt --color=red ll "hello" --align=justify
```

The parser must preserve existing support for:

```bash
go run . "hello"
go run . "hello" shadow
go run . --color=red "hello"
go run . --color=red ll "hello"
go run . --output=out.txt "hello"
```

## Shared Contract

Add alignment to the existing options struct:

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

`Align == ""` or `Align == "left"` means default left alignment.

Accepted values:

```text
left
center
right
justify
```

## Parsing Rules

- Options may appear before, between, or after positional arguments.
- Options must use `--name=value` form.
- Supported option names are `--color=`, `--output=`, and `--align=`.
- Invalid spellings such as `-align=right`, `--alignment=right`, `--align right`, and `--align=` must be rejected.
- Duplicate alignment flags must be rejected.
- Unknown `--...` flags must be rejected.
- Positional text and banner behavior must stay compatible with the current parser.
- A color substring remains positional and is recognized only when it appears immediately after `--color=<color>`.

## Usage Message

For invalid alignment option formats, use:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --align=<left|center|right|justify> something standard
```

## Interactive Mode

In options-aware interactive mode:

- After banner/color/substring selection, ask for alignment.
- Empty alignment input means left alignment.
- `/exit` exits from the alignment prompt.
- Invalid alignment values are rejected and the alignment prompt repeats.
- CLI stores the selected alignment in `Options.Align`; it does not perform rendering or output alignment.

## Tests

Add or update CLI tests for:

- Single string still defaults to left alignment.
- Valid `--align=left`, `--align=center`, `--align=right`, and `--align=justify`.
- Alignment flag before text, between text and banner, and after banner.
- Alignment combined with `--color=`.
- Alignment combined with `--output=`.
- Alignment combined with color substring.
- Empty alignment value.
- Separated alignment flag.
- Misspelled alignment flag.
- Unsupported alignment value.
- Duplicate alignment flags.
- Interactive empty alignment defaults to left.
- Interactive invalid alignment repeats the prompt.

## Done When

- CLI returns `Options.Align`.
- Existing color, output, banner, and interactive tests still pass.
- New alignment parser tests pass.
