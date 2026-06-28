# ASCII-ART

ASCII-ART is a Go command-line program that renders text as multi-line ASCII art using predefined banner files.

The current implementation supports:

- argument mode
- interactive mode
- three banner styles: `standard`, `shadow`, and `thinkertoy`
- foreground color using preset names, full hexadecimal color codes, RGB, or HSL
- escaped newline sequences written as `\n`
- output files with `--output=<file>`
- printable ASCII characters supported by the banner files

The active project objective is to add alignment. The CLI alignment option is part of Task 01, while the renderer and output packages still need Task 02 and Task 03 to make alignment affect the final visual output:

```bash
go run . --align=left "Hello"
go run . --align=center "Hello" shadow
go run . --align=right "Hello" standard
go run . --align=justify "Hello world" thinkertoy
```

Alignment design and task ownership are documented in:

```text
docs/prd.md
docs/exercise.md
docs/workFlow.md
docs/tasks/
```

## Requirements

- Go 1.22.2 or newer
- The provided banner files in the `banners/` directory

The project uses only Go standard library packages.

## How To Run

### Interactive Mode

Run the program without extra arguments:

```bash
go run .
```

The program will:

1. Ask for text.
2. Ask you to choose a banner.
3. Ask for an optional color.
4. If a color is provided, ask for an optional substring to color.
5. Ask whether to write output to a file.
6. Print or write the ASCII-art output.
7. Return to the text prompt.
8. Continue until you type `/exit`.

Banner choices in interactive mode:

```text
1 -> standard
2 -> shadow
3 -> thinkertoy
```

### Argument Mode

Run the program with text:

```bash
go run . "Hello"
```

This uses the standard banner by default.

Run with a selected banner:

```bash
go run . "Hello" shadow
go run . "Hello" thinkertoy
go run . "Hello" standard.txt
```

Accepted banner arguments:

- `standard`
- `standard.txt`
- `shadow`
- `shadow.txt`
- `thinkertoy`
- `thinkertoy.txt`

If only one positional argument is provided, it is treated as text even if it equals a banner name.

```bash
go run . "shadow"
```

renders the word `shadow`.

Argument mode renders once and exits.

## Order-Flexible Options

Supported options use `--flag=value` form:

```text
--color=<color>
--output=<file>
```

Options can appear before, between, or after positional arguments:

```bash
go run . --output=out.txt "Hello"
go run . "Hello" --output=out.txt
go run . "Hello" shadow --color=cyan
go run . "hello world" shadow --output=out.txt --color=cyan world
```

Separated option values are invalid:

```bash
go run . --color red "Hello"
go run . --output out.txt "Hello"
```

## Color

Color works in both argument mode and interactive mode.

Color the whole rendered string:

```bash
go run . --color=red "Hello"
go run . --color=#ff0000 "Hello"
go run . --color=rgb,255,0,0 "Hello"
go run . --color=hsl,0,100%,50% "Hello"
```

Color every matching substring:

```bash
go run . --color=RGB,255,0,0 kit "a king kitten have kit"
```

Use a custom banner and color together:

```bash
go run . "hello world" shadow --color=cyan world
```

Supported preset names:

- `black`
- `red`
- `green`
- `yellow`
- `blue`
- `magenta`
- `cyan`
- `white`
- `orange`

Supported color value formats:

- preset color name: `red`
- hexadecimal: `#RRGGBB`, such as `#ff0000`
- RGB: `rgb,<red>,<green>,<blue>`, such as `rgb,255,0,0`
- HSL: `hsl,<hue>,<saturation>%,<lightness>%`, such as `hsl,0,100%,50%`

RGB and HSL prefixes are case-insensitive.

Invalid color values print:

- invalid color name: `Invalid color name`
- invalid hexadecimal color code: `Invalid hexadecimal color code`
- invalid RGB color code: `Invalid RGB color code`
- invalid HSL color code: `Invalid HSL color code`

## Output Files

Write rendered ASCII art to a file:

```bash
go run . --output=banner.txt "Hello"
go run . "Hello" shadow --output=banner.txt
go run . --output=banner.txt --color=red "Hello"
```

When `--output=<file>` is used, the program writes to the file instead of printing the ASCII art to the terminal.

ANSI color codes are preserved in file output.

## Alignment Objective

Alignment is the active implementation target. The CLI can parse alignment options once Task 01 is complete; visual alignment depends on the structured renderer and output work in Task 02 and Task 03.

Planned syntax:

```bash
go run . --align=left "Hello"
go run . --align=center "Hello"
go run . --align=right "Hello"
go run . --align=justify "Hello world"
```

Planned values:

```text
left
center
right
justify
```

The agreed design replaces flat renderer output with structured rendering. The renderer will preserve word and space-run segments so the output package can implement justify without stretching spaces inside ASCII-art letters.

Task split:

- Programmer 1: CLI alignment option
- Programmer 2: structured renderer
- Programmer 3: output alignment and main integration

Task files:

```text
docs/tasks/task01_cli_alignment.md
docs/tasks/task02_structured_renderer.md
docs/tasks/task03_output_alignment_integration.md
```

## Newlines

Use escaped `\n` sequences inside the input string:

```bash
go run . "Hello\nThere" standard
```

This renders `Hello` and `There` as separate ASCII-art blocks.

## Project Structure

```text
ascii-art/
|-- main.go
|-- banners/
|   |-- standard.txt
|   |-- shadow.txt
|   `-- thinkertoy.txt
|-- internal/
|   |-- banner/
|   |   |-- loader.go
|   |   `-- loader_test.go
|   |-- cli/
|   |   |-- args.go
|   |   |-- interactive.go
|   |   |-- mode_picker.go
|   |   |-- cli_test.go
|   |   `-- interactive_test.go
|   |-- output/
|   |   |-- print.go
|   |   |-- write.go
|   |   |-- write_file.go
|   |   `-- print_test.go
|   `-- renderer/
|       |-- build.go
|       |-- color.go
|       |-- normalize.go
|       |-- substring.go
|       |-- build_test.go
|       |-- color_test.go
|       |-- substring_test.go
|       `-- normalize_test.go
|-- docs/
|   |-- audit_test.md
|   |-- exercise.md
|   |-- golden_tests.md
|   |-- prd.md
|   |-- workFlow.md
|   `-- tasks/
|       |-- task01_cli_alignment.md
|       |-- task02_structured_renderer.md
|       `-- task03_output_alignment_integration.md
|-- go.mod
`-- README.md
```

## Architecture

Current pipeline:

```text
CLI input
-> banner loading
-> input normalization
-> ASCII rendering
-> output printing or file writing
```

Alignment target pipeline:

```text
CLI input
-> banner loading
-> input normalization
-> structured ASCII rendering
-> aligned terminal/file output
```

Package responsibilities:

- `internal/cli`: chooses interactive or argument mode and validates user input
- `internal/banner`: loads banner files into `map[rune][]string`
- `internal/renderer`: normalizes input, validates color values, finds colored positions, and builds ASCII-art
- `internal/output`: prints, writes files, and will own alignment
- `main.go`: connects the packages

## Testing

Run all tests:

```bash
go test ./...
```

Current tests cover:

- banner loading and invalid banner formats
- CRLF banner line endings
- CLI argument mode
- CLI interactive mode
- order-flexible color and output options
- escaped newline normalization
- ASCII row building
- color name, hexadecimal, RGB, and HSL validation
- substring position matching for color rendering
- output printing and file writing

Additional edge cases are documented in:

```text
docs/golden_tests.md
internal/cli/test_cases.md
```

## Limitations

- Custom banner file paths are not accepted from user input.
- Numeric banner choices are only used in interactive mode.
- Unicode characters outside the printable ASCII banner range are not supported.
- Alignment is designed but not yet implemented.
- Resizing and scaling are not implemented.
