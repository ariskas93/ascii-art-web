# ASCII-ART

ASCII-ART is a Go command-line program that renders text as multi-line ASCII art using predefined banner files.

The program supports:

- interactive mode
- argument mode
- three banner styles: `standard`, `shadow`, and `thinkertoy`
- foreground color using preset names, full hexadecimal color codes, RGB, or HSL
- escaped newline sequences written as `\n`
- printable ASCII characters supported by the banner files

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
5. Print the ASCII-art output.
6. Return to the text prompt.
7. Continue until you type `/exit`.

Banner choices in interactive mode:

```text
1 -> standard
2 -> shadow
3 -> thinkertoy
```

At the color prompt, press Enter to render without color. The prompt shows examples for each supported color format:

```text
Enter color, or press Enter for no color (/exit to quit)
Examples: red, #ff0000, rgb,255,0,0, hsl,0,100%,50%
Color:
```

If you enter a color and then press Enter at the substring prompt, the whole rendered string is colored. If you enter a substring, every matching occurrence is colored.

You can type `/exit` at the text, banner, color, or substring prompt.

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

Argument mode renders once and exits.

### Color

Color works in both argument mode and interactive mode. In interactive mode, enter the color when prompted after choosing a banner.

Color the whole rendered string with a preset color:

```bash
go run . --color=red "Hello"
```

Color the whole rendered string with a full hexadecimal color code:

```bash
go run . --color=#ff0000 "Hello"
```

Color the whole rendered string with RGB:

```bash
go run . --color=rgb,255,0,0 "Hello"
```

Color the whole rendered string with HSL:

```bash
go run . --color=hsl,0,100%,50% "Hello"
```

Color every matching substring:

```bash
go run . --color=RGB,255,0,0 kit "a king kitten have kit"
```

Use a custom banner and color together by placing the banner before the color flag:

```bash
go run . "hello world" shadow --color=cyan world
```

Supported preset names are:

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

RGB and HSL prefixes are case-insensitive, so `RGB,255,0,0` and `HSL,0,100%,50%` also work. RGB values must be integers from `0` to `255`. HSL hue must be an integer from `0` to `360`; saturation and lightness must be integer percentages from `0%` to `100%`. The program does not support `rgb(...)`, `hsl(...)`, `rgba`, or `hsla` syntax.

Invalid color values print:

- invalid color name: `Invalid color name`
- invalid hexadecimal color code: `Invalid hexadecimal color code`
- invalid RGB color code: `Invalid RGB color code`
- invalid HSL color code: `Invalid HSL color code`

## Newlines

Use escaped `\n` sequences inside the input string:

```bash
go run . "Hello\nThere" standard
```

This renders `Hello` and `There` as separate ASCII-art blocks.

## Examples

Default banner:

```bash
go run . "A"
```

Shadow banner:

```bash
go run . "Zone01" shadow
```

Thinkertoy banner:

```bash
go run . "Go!" thinkertoy
```

Multiple lines:

```bash
go run . "Hello\nWorld" standard
```

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
|-- go.mod
|-- agent_log.md
|-- to_study.md
`-- README.md
```

## Architecture

The program follows a small pipeline:

```text
CLI input
-> banner loading
-> input normalization
-> color position detection
-> ASCII rendering
-> output printing
```

Package responsibilities:

- `internal/cli`: chooses interactive or argument mode and validates user input
- `internal/banner`: loads banner files into `map[rune][]string`
- `internal/renderer`: normalizes input, validates color values, finds colored positions, and builds ASCII-art rows
- `internal/output`: prints final rows to standard output
- `main.go`: connects the packages and keeps the interactive loop minimal

## Banner Format

Each banner file contains printable ASCII characters from 32 to 126.

Each character block contains:

- 1 separator line
- 8 visible art lines
- 9 total lines

The loader maps every character into:

```go
map[rune][]string
```

Each rune maps to exactly 8 strings.

The loader also normalizes CRLF line endings, which keeps `thinkertoy.txt` rendering correctly across environments.

## Error Handling

Interactive mode:

- empty text is rejected with `No valid string input`
- invalid banner choices are rejected with `Invalid banner choice`
- empty color input means no color
- empty substring input after a color means the whole rendered string is colored
- invalid color values print the rendering error and return to the text prompt
- `/exit` exits cleanly

Argument mode:

- missing text returns an error
- empty text returns an error
- invalid banner names return an error
- too many arguments return an error

Rendering:

- unsupported characters return an error
- invalid color names return `Invalid color name`
- invalid hexadecimal color codes return `Invalid hexadecimal color code`
- invalid RGB color codes return `Invalid RGB color code`
- invalid HSL color codes return `Invalid HSL color code`
- invalid or missing banner files return an error

## Testing

Run all tests:

```bash
go test ./...
```

The tests cover:

- banner loading and invalid banner formats
- CRLF banner line endings
- CLI argument mode
- CLI interactive mode
- escaped newline normalization
- ASCII row building
- color name, hexadecimal, RGB, and HSL validation
- substring position matching for color rendering
- output printing

Additional edge cases are documented in:

```text
docs/golden_tests.md
internal/cli/test_cases.md
```

## Limitations

- Custom banner file paths are not accepted from user input.
- Numeric banner choices are only used in interactive mode.
- Unicode characters outside the printable ASCII banner range are not supported.
- Alignment, resizing, and file export are not implemented.
