# Agent Log

This log summarizes the implementation and documentation work completed for the ASCII Art Color project.

## Initial State

- The project already rendered ASCII art from provided banner files.
- Argument mode supported text and optional banner selection.
- Interactive mode existed as an extra feature.
- Color support existed around preset names and earlier ANSI-style input experiments.

## Color Input Redesign

The user-facing color system was updated while keeping ANSI escape sequences as the internal terminal output mechanism.

Supported color inputs now are:

- preset names: `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `black`, `orange`
- full hexadecimal colors: `#RRGGBB`
- RGB comma form: `rgb,<red>,<green>,<blue>`
- HSL comma form: `hsl,<hue>,<saturation>%,<lightness>%`

RGB and HSL prefixes are case-insensitive. Preset names remain case-sensitive.

Unsupported user-facing inputs include:

- raw ANSI notation like `ansi:31`
- CSS-style `rgb(...)` and `hsl(...)`
- alpha forms such as `rgba(...)` and `hsla(...)`

## Renderer Changes

Updated `internal/renderer/color.go`:

- kept the function name `getANSIColor`
- changed it to return `(string, error)`
- added validation for `#RRGGBB`
- added RGB parsing and validation
- added HSL parsing and conversion to RGB
- added exact error messages:
  - `Invalid color name`
  - `Invalid hexadecimal color code`
  - `Invalid RGB color code`
  - `Invalid HSL color code`

Updated `internal/renderer/build.go`:

- color validation errors now come from `getANSIColor`
- `BuildASCII` propagates the exact renderer error
- ASCII-art coloring still happens during rendering by wrapping marked character block rows

## CLI Behavior

The CLI keeps color value validation out of the parser.

`internal/cli/args.go` is responsible for:

- accepting only exact `--color=<color>` flag spelling
- rejecting empty color values
- rejecting multiple color flags
- preserving non-empty color values for renderer validation
- supporting banner-compatible color forms:
  - `go run . [STRING] [BANNER] --color=<color>`
  - `go run . [STRING] [BANNER] --color=<color> [SUBSTRING]`

When `--color=<color>` is first, the default banner is `standard`.

When banner and color are used together, the banner must come before the color flag:

```bash
go run . "hello world" shadow --color=cyan world
```

## Interactive Color Support

Interactive mode now collects color settings after banner selection.

The interactive flow is:

1. Prompt for text.
2. Prompt for banner choice.
3. Prompt for optional color.
4. If a color is entered, prompt for optional substring.
5. Render and print the ASCII-art output.
6. Return to the text prompt.

The color prompt shown to the user is:

```text
Enter color, or press Enter for no color (/exit to quit)
Examples: red, #ff0000, rgb,255,0,0, hsl,0,100%,50%
Color:
```

Interactive color rules:

- empty color input means no color
- non-empty color input is passed to the renderer for validation
- empty substring input after a color means color the whole string
- non-empty substring input colors every matching occurrence
- substring input preserves spaces and only removes the line ending
- `/exit` works at the text, banner, color, and substring prompts
- invalid color values print the renderer error and return to the text prompt

The exported interactive loop reuses one buffered reader for `os.Stdin`, which keeps scripted interactive runs from losing buffered input between cycles.

## Output CLI Parser Handoff

Task 01 for output mode is implemented in the CLI package, with `main.go` intentionally left unchanged for the main integration owner.

Added `cli.Options`:

```go
type Options struct {
	Text       string
	BannerPath string
	Color      string
	Substring  string
	OutputFile string
}
```

Added output-aware entry points:

- `RunArgumentModeOptions(args []string) (Options, error)`
- `PickModeOptions(args []string) (Options, bool, error)`

The old tuple-returning functions still exist for compatibility:

- `RunArgumentMode(args []string) (string, string, string, string, error)`
- `PickMode(args []string) (string, string, string, string, bool, error)`

Output parsing rules:

- `--output=<fileName.txt>` is accepted only as the first argument after the program name
- malformed output flags return the exact output usage message
- `OutputFile == ""` means normal terminal output
- `OutputFile != ""` means the rendered ASCII-art rows should be written to that file by the integration layer
- output mode composes with existing banner and color forms
- ANSI color codes are preserved because output files should remain color-readable in the terminal

Interactive output collection was also added to the options-aware interactive path:

- after text, banner, color, and substring prompts, the CLI asks whether to write output to a file
- blank, `n`, and `no` keep `OutputFile == ""`
- `y` and `yes` prompt for an output file name
- empty output file names are rejected with `No valid output file name`
- invalid output choices are rejected with `Invalid output choice`
- `/exit` works at the output choice and output file name prompts
- the old tuple-based interactive path remains available for current `main.go` compatibility

Accepted parser examples:

```bash
go run . --output=banner.txt "hello"
go run . --output=banner.txt "hello" standard
go run . --output=banner.txt --color=red "hello"
go run . --output=banner.txt --color=red kit "a king kitten have kit"
go run . --output=banner.txt "hello world" shadow --color=cyan world
```

Invalid output formats return:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```

Important integration note:

- `output.WriteASCIIFile(fileName, lines)` already exists in `internal/output`
- `main.go` has not been updated to call `PickModeOptions` or `WriteASCIIFile`
- argument output and interactive output are not expected to write files until the main integration task is completed

## Test Coverage Added

Updated CLI tests to cover:

- hex color values
- RGB and HSL color values
- banner-compatible color forms
- empty color values
- color flags in invalid positions
- multiple color flags
- shell-split RGB/HSL values
- invalid banner with color mode
- interactive color collection, whole-string color, substring color, exit prompts, and EOF behavior
- output option parsing, malformed output flags, output plus banner, and output plus color
- interactive output choice, output file name collection, output prompt exits, invalid output choices, and empty output file retries

Updated renderer color tests to cover:

- valid preset names
- valid hex values, including uppercase hex
- valid RGB values with case-insensitive prefix
- valid HSL values with case-insensitive prefix
- invalid preset names
- invalid hex lengths and characters
- invalid RGB values and shapes
- invalid HSL values and shapes
- unsupported CSS-style color functions

Updated renderer build tests to cover:

- invalid color name errors
- invalid hex errors
- invalid RGB errors
- invalid HSL errors
- whole-string color rendering
- substring color rendering

Verification run:

```bash
go test ./...
```

## Documentation Updates

Updated project documentation to match the working program:

- `README.md`
- `docs/prd.md`
- `docs/golden_tests.md`
- `docs/audit_test.md`
- `docs/exercise.md`
- `docs/tasks/tasks01_cli.md`
- `docs/tasks/task02_banner.md`
- `docs/tasks/task03_renderer.md`
- `docs/tasks/task04_main.md`
- `docs/tasks/Task01_ai_log.md`
- `docs/workFlow.md`
- `internal/cli/test_cases.md`

Added:

- `to_study.md`
- `agent_log.md`

## Current Working State

The program supports:

```bash
go run . "hello"
go run . "hello" shadow
go run . --color=red "hello"
go run . --color=#ff0000 "hello"
go run . --color=RGB,255,0,0 "hello"
go run . --color=hsl,0,100%,50% "hello"
go run . --color=orange GuYs "HeY GuYs"
go run . "hello world" shadow --color=cyan world
```

Interactive mode supports the same renderer color values through prompts after banner selection.

The CLI parser is ready for output file integration, but output mode is not wired into `main.go` yet.

The project uses only standard Go packages and keeps responsibilities split across CLI, banner loading, renderer, output, and main integration packages.
