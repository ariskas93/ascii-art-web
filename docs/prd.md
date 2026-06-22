# PRD - ASCII-ART

---

## 1. Problem Statement

ASCII-ART is a Go CLI program that renders text as multi-line ASCII art using one of three provided banner files.

The required audit behavior is argument mode:

- Argument mode: `go run . [STRING] [BANNER]`

The program may also support an additional interactive mode:

- Interactive mode: `go run .`

The official ASCII-ART audit requirement is:

```bash
go run . [STRING] [BANNER]
```

The ASCII Art Color audit requirement is:

```bash
go run . --color=<color> [substring] [STRING]
```

The second argument represents the banner/template name.

The program must still support execution with a single `[STRING]` argument, which defaults to the standard banner.

Interactive mode remains an additional feature layered on top of the mandatory audit specification. It must not change or weaken the required argument-mode behavior.

The input text is mandatory in argument mode and collected by prompt in interactive mode. The banner is optional in argument mode and defaults to `standard`. In interactive mode the user selects a banner from a numbered menu, then may choose an optional color and optional substring to color.

Supported input:
- Letters (a-z, A-Z)
- Numbers (0-9)
- Spaces
- Special characters (based on banner support)
- Newline sequences (\n)
Constraints:
- Each character is represented by 8 lines
- Banner files follow a fixed format (no modification required)
- Banner files may use LF or CRLF line endings

Allowed packages:
- Only standard Go packages are allowed.

---

## 2. CLI Contract

The official audit behavior is argument mode:

```bash
go run . [STRING] [BANNER]
```

Interactive mode is an extended feature added on top of the mandatory specification.

### Interactive mode

Run:

```bash
go run .
```

The program operates in an interactive loop:

- The program starts and prompts the user to input a string or type /exit to terminate.
- Empty text is rejected with `No valid string input`.
- Whitespace-only text is accepted because spaces are renderable characters.
- After a valid string is entered, the user is prompted to select one of three banner files using numeric input:
  - `1` -> `banners/standard.txt`
  - `2` -> `banners/shadow.txt`
  - `3` -> `banners/thinkertoy.txt`
- Invalid banner choices are rejected with `Invalid banner choice` and the banner prompt repeats.
- After a valid banner is selected, the user is prompted to enter a color or press Enter for no color.
- The color prompt displays examples for preset, hexadecimal, RGB, and HSL color formats.

```text
Enter color, or press Enter for no color (/exit to quit)
Examples: red, #ff0000, rgb,255,0,0, hsl,0,100%,50%
Color:
```

- If the color prompt is empty, rendering continues without color.
- If a color is entered, the user is prompted to enter a substring or press Enter to color the whole string.
- The color prompt trims surrounding whitespace from the color value.
- The substring prompt preserves spaces inside the substring and only removes the line ending.
- Invalid color values are reported as rendering errors, and interactive mode returns to the string input prompt instead of terminating.
- The program generates and prints the corresponding ASCII-art output.
- After printing, the program returns to the string input prompt.
- This cycle repeats continuously until the user types /exit at any prompt, which terminates the program.

The repeat loop lives in `main.go`. The CLI package collects one valid text, banner, color, and substring selection per call.

### Argument mode

Run:

```bash
go run . "Hello"
go run . "Hello" standard
go run . "Hello" shadow
go run . "Hello\nThere" thinkertoy.txt
go run . --color=red "Hello"
go run . --color=red ll "Hello"
go run . --color=#ff0000 "Hello"
go run . --color=RGB,255,0,0 "Hello"
go run . --color=hsl,0,100%,50% "Hello"
go run . "Hello" shadow --color=cyan
```

Official format:

```bash
go run . [STRING] [BANNER]
```

The banner argument is optional.

If omitted, the program defaults to:

```text
banners/standard.txt
```

Argument rules:

- `args[1]` is the text to render.
- `args[2]` is optional and represents the banner/template name.
- If no banner is provided, the program uses `banners/standard.txt`.
- Accepted banner values are `standard`, `standard.txt`, `shadow`, `shadow.txt`, `thinkertoy`, and `thinkertoy.txt`.
- Numeric banner choices such as `1`, `2`, or `3` are only valid in interactive mode.
- The program must still work correctly with only one string argument.
- Argument mode renders once and exits.

Valid argument formats are:

```bash
go run . [STRING]
go run . [STRING] [BANNER]
go run . --color=<color> [STRING]
go run . --color=<color> [SUBSTRING] [STRING]
go run . [STRING] [BANNER] --color=<color>
go run . [STRING] [BANNER] --color=<color> [SUBSTRING]
```

Any invalid argument format must print:

```text
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

Invalid cases include:

- Too many arguments
- Missing string argument
- Unsupported banner name
- Numeric banner choices in argument mode, such as `1`, `2`, or `3`
- Any argument format outside the valid formats listed above

Unsupported banner names, including values such as `other`, `abc`, or numeric choices in argument mode, must print the official usage message.

If other optional ASCII-art projects are later implemented, the CLI may additionally support valid `[OPTION]` arguments while preserving backward compatibility with:

```bash
go run . [STRING] [BANNER]
```

Color option rules:

- The color flag must be formatted exactly as `--color=<color>`.
- Invalid color flag formats must print:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

- If no substring is provided, the whole input string is colored.
- If a substring is provided, every matching occurrence is colored.
- Substring matching is case-sensitive and is applied per normalized line.
- Color matching must not cross newline boundaries.
- Color mode defaults to `banners/standard.txt` unless a valid banner is provided by the extended `[STRING] [BANNER] --color=<color>` form.
- Supported named colors include `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, and `orange`.
- Hexadecimal color notation must use the full `#RRGGBB` form, for example `#ff0000` or `#FF0000`.
- RGB color notation must use the comma-only form `rgb,<red>,<green>,<blue>`, for example `rgb,255,0,0`.
- HSL color notation must use the comma-only form `hsl,<hue>,<saturation>%,<lightness>%`, for example `hsl,0,100%,50%`.
- RGB and HSL notation names are case-insensitive, so `RGB,255,0,0` and `HSL,0,100%,50%` are valid.
- RGB values must be integer values from `0` to `255`.
- HSL hue must be an integer from `0` to `360`; hue values greater than `360` are invalid.
- HSL saturation and lightness must be integer percentages from `0%` to `100%`.
- CSS-style function syntax such as `rgb(255, 0, 0)` or `hsl(0, 100%, 50%)` is not supported.
- Invalid hexadecimal color codes must print `Invalid hexadecimal color code`.
- Invalid RGB color codes must print `Invalid RGB color code`.
- Invalid HSL color codes must print `Invalid HSL color code`.
- Invalid color names must print `Invalid color name`.

---

## 3. ASCII-ART Functions

The program uses a mapping-based approach to handle banner files.
Each selected banner is loaded for the current render cycle and converted into a lookup structure that allows fast and simple access to the ASCII-art representation of each character.

### 3.1 Banner Mapping

Each banner file contains all printable ASCII characters (from ASCII 32 `' '` to ASCII 126 `'~'`).

Each character:
- Has a height of **8 lines**
- Is separated by **one empty line**
- Occupies **9 lines total**

The program parses the banner file and builds a map:

```go
map[rune][]string
```

This map associates each character with its corresponding 8-line ASCII-art representation.

**Example:**

```go
banner['A'] = []string{
    "line1",
    "line2",
    "line3",
    "line4",
    "line5",
    "line6",
    "line7",
    "line8",
}
```

This allows constant-time access to any character's ASCII-art block during rendering.

### 3.2 LoadBanner

```go
LoadBanner(path string) (map[rune][]string, error)
```

**Responsibility:**
- Read the selected banner file from disk
- Normalize CRLF line endings to LF
- Split its contents into lines
- Parse and map each ASCII character to its 8-line representation
- Return an error if the file cannot be read or is malformed

### 3.3 NormalizeInput

```go
NormalizeInput(input string) []string
```

**Responsibility:**
- Process the input string
- Split it into multiple lines using escaped `\n` sequences
- Return a slice of strings, where each element represents one line to render

**Example:**

```text
"Hello\nThere" -> ["Hello", "There"]
```

### 3.4 BuildASCII

```go
BuildASCII(lines []string, banner map[rune][]string, substring string, color string) ([]string, error)
```

**Responsibility:**
- Convert normalized input into ASCII-art output
- Render normally when `color == ""`
- Color the whole string when `color != ""` and `substring == ""`
- Color every matching substring occurrence when both `color` and `substring` are provided

**Process:**
- For each input line:
  - Detect the input character positions that should be colored
  - Create 8 output rows
  - For each character in the line:
    - Retrieve its ASCII-art block from the banner map
    - Wrap the character block rows with ANSI color codes when its input position is colored
    - Append each of the 8 rows horizontally
- Combine all processed lines into a final output slice

**Error handling:**
- If a character is not found in the banner map, return an error
- If a color name, hexadecimal code, RGB code, or HSL code is invalid, return an error

### 3.5 Color Helpers

```go
getANSIColor(color string) (string, error)
colorize(text string, colorCode string) string
coloredPositions(s string, substring string) map[int]bool
```

**Responsibility:**
- Convert supported color names, full `#RRGGBB` values, RGB values, and HSL values to ANSI escape codes
- Wrap colored ASCII-art segments with ANSI color and reset codes
- Find every input character position that belongs to a matching substring
- Return all positions when `substring == ""`

### 3.6 PrintASCII

```go
PrintASCII(output []string)
```

**Responsibility:**
- Print the ASCII-art output line by line to standard output
- Ensure correct formatting (each row printed on a new line)
- Do not print any extra text

### 3.7 Processing Flow

```text
cli.PickMode
-> LoadBanner
-> NormalizeInput
-> Detect colored positions
-> BuildASCII
-> PrintASCII
```

### 3.8 Design Rationale

- The selected banner is parsed once per render cycle
- The `map[rune][]string` structure provides fast lookup and clean rendering logic
- Responsibilities are clearly separated, making the code easier to test and maintain
- The rendering process is independent from file parsing and CLI parsing, reducing coupling between components
- Color is applied during rendering based on original input character positions, not by searching the rendered ASCII-art output
- `main.go` stays minimal and only coordinates package calls

### 3.9 Edge Cases

- Empty input → handled at CLI level
- Whitespace-only input → accepted because spaces are renderable
- Single `\n` → returns [""] (one empty line)
- Trailing `\n` → last empty element is removed
- Multiple `\n` → produce empty lines minus trailing one
- CRLF banner files → normalized by the banner loader
- Unsupported characters → return an error
- Invalid banner file → return an error during loading
- Missing substring in color mode → render uncolored ASCII-art output
- Empty substring in color mode → color the whole string

### 3.10 Banner Template Support

The renderer supports three official banner templates:

- `standard`
- `shadow`
- `thinkertoy`

Each template represents a different ASCII-art style while preserving the same printable ASCII mapping structure.

The selected banner name is resolved to its corresponding banner file path before loading.

The banner argument acts as a rendering template selector.

Example:

```text
standard    -> banners/standard.txt
shadow      -> banners/shadow.txt
thinkertoy  -> banners/thinkertoy.txt
```

The rendering pipeline remains independent from the selected template because all banners are converted into the same internal structure:

```go
map[rune][]string
```

This allows the renderer to operate identically regardless of banner style.

---

## 4. Project file structure

```text
ascii-art/
|-- main.go
|-- banners/
|   |-- standard.txt
|   |-- shadow.txt
|   `-- thinkertoy.txt
|-- internal/
|   |-- cli/
|   |   |-- args.go
|   |   |-- mode_picker.go
|   |   |-- interactive.go
|   |   |-- cli_test.go
|   |   |-- interactive_test.go
|   |   `-- test_cases.md
|   |-- banner/
|   |   |-- loader.go
|   |   `-- loader_test.go
|   |-- renderer/
|   |   |-- build.go
|   |   |-- color.go
|   |   |-- normalize.go
|   |   |-- substring.go
|   |   |-- build_test.go
|   |   |-- color_test.go
|   |   |-- substring_test.go
|   |   `-- normalize_test.go
|   `-- output/
|       |-- print.go
|       `-- print_test.go
|-- docs/
|   |-- audit_test.md
|   |-- exercise.md
|   |-- golden_tests.md
|   |-- prd.md
|   |-- workFlow.md
|   `-- tasks/
|       |-- Task01_ai_log.md
|       |-- tasks01_cli.md
|       |-- task02_banner.md
|       |-- task03_renderer.md
|       `-- task04_main.md
|-- go.mod
|-- agent_log.md
|-- to_study.md
`-- README.md
```
                
## 5. TASKS

### Task 01 - CLI

Package: internal/cli

Goal:
Collect user input, select a supported banner, and tell `main.go` whether the program should continue or exit.

Responsibilities:

- Choose interactive mode when no extra arguments are provided
- Choose argument mode when text is provided through command-line arguments
- Prompt for text, banner choice, optional color, and optional substring in interactive mode
- Validate text, optional banner, optional color, and optional substring in argument mode
- Validate the official CLI format: `go run . [STRING] [BANNER]`
- Validate the official color format: `go run . --color=<color> [substring] [STRING]`
- Print the official usage message for invalid formats
- Print the official color usage message for invalid color formats
- Support single-string execution with default standard banner
- Return text, banner path, color, substring, exit flag, and errors to `main.go`

Files:

internal/cli/args.go
internal/cli/mode_picker.go
internal/cli/interactive.go
internal/cli/cli_test.go
internal/cli/interactive_test.go


### Task 02 - Banner Loading and Mapping

Package: internal/banner

Goal:
Load the selected banner file and convert it into a character map.

Responsibilities:
- Read the banner file from disk
- Normalize CRLF line endings to LF
- Split the file into lines
- Map printable ASCII characters from 32 to 126
- Store each character as []string with 8 lines
- Return an error if the file is missing or malformed

Files:
internal/banner/loader.go
internal/banner/loader_test.go


### Task 03 - Renderer

Package: internal/renderer

Goal:
Normalize user input and build the final ASCII-art output.

Responsibilities:
- Convert escaped newline sequences \n into separate input lines
- Preserve empty lines caused by multiple newlines
- Create 8 output rows for each input line
- For each character, retrieve its banner from map[rune][]string
- Append character rows horizontally
- Resolve color names, hexadecimal colors, RGB values, and HSL values into ANSI escape codes
- Detect substring positions before rendering and color matching character blocks
- Handle empty input lines correctly
- Return an error for unsupported characters

Files:
internal/renderer/build.go
internal/renderer/build_test.go
internal/renderer/color.go
internal/renderer/color_test.go
internal/renderer/normalize.go
internal/renderer/normalize_test.go
internal/renderer/substring.go
internal/renderer/substring_test.go

### Task 04 - Output Printing

Package: internal/output

Goal:
Print the generated ASCII-art rows.

Responsibilities:
- Print each output row on a new line
- Avoid extra messages
- Keep printing separate from rendering logic

Files:
internal/output/print.go
internal/output/print_test.go


### Task 05 - Main Program Integration

File: main.go

Goal:
Connect all packages into one working program.

Responsibilities:
- Call cli.PickMode(os.Args)
- Load the selected banner
- Normalize input
- Build ASCII output
- Pass color and substring settings into the renderer
- Print the result
- Repeat only in interactive mode
- Exit after one render in argument mode

Files:
main.go


### Task 06 - Documentation

Goal:
Document how the program works and how to run it.

Responsibilities:
- Finalize README.md
- Finalize docs/prd.md
- Create task documentation files
- Include usage examples
- Include known limitations

Files:
README.md
docs/prd.md
docs/workFlow.md
docs/tasks/tasks01_cli.md
docs/tasks/task02_banner.md
docs/tasks/task03_renderer.md
docs/tasks/task04_main.md

Suggested Task Order
1. Banner Loading and Mapping
2. Input Normalization
3. ASCII Renderer
4. Output Printing
5. CLI
6. Main Program Integration
7. Documentation

## 6. Non-Goals

- The program will not support multiple independent strings through the command line.
- The program will not support custom banner files outside the three provided files:
    - standard.txt
    - shadow.txt
    - thinkertoy.txt
- The program will not support custom user-created templates outside the provided banner files.
- The program will not allow users to edit or modify banner files from inside the program.
- The program will not support arrow-key navigation for banner selection.
- The program will not support multiple color flags in the same command.
- The program will not support exporting output to a file.
- The program will not support alignment options such as center, left, right, or justify.
- The program will not support resizing, scaling, or changing the height of ASCII-art characters.
- The program will not support non-printable ASCII characters.
- The program will not support Unicode characters outside the printable ASCII range.
- The program will not run as a web application or GUI application.

---

## 7. Acceptance Criteria

### Functional Acceptance Criteria

- [x] The program starts with `go run .` in interactive mode.
- [x] The program supports `go run . [STRING] [BANNER]` in argument mode.
- [x] The program still supports a single `[STRING]` argument.
- [x] Interactive mode prompts the user to input a string.
- [x] Interactive mode prompts the user to select one of three banner files.
- [x] Interactive mode prompts the user to optionally enter a color after banner selection.
- [x] Interactive mode prompts the user to optionally enter a substring when a color is provided.
- [x] Interactive mode accepts an empty color prompt as uncolored rendering.
- [x] Interactive mode accepts an empty substring prompt as whole-string color rendering.
- [x] Argument mode accepts `standard`, `standard.txt`, `shadow`, `shadow.txt`, `thinkertoy`, and `thinkertoy.txt`.
- [x] Argument mode accepts `--color=<color> [STRING]` and colors the whole rendered string.
- [x] Argument mode accepts `--color=<color> [SUBSTRING] [STRING]` and colors every matching substring.
- [x] Argument mode accepts `[STRING] [BANNER] --color=<color> [SUBSTRING]` as the banner-compatible color form.
- [x] The selected banner acts as the rendering template.
- [x] The program supports `standard`, `shadow`, and `thinkertoy` templates.
- [x] The program prints the input string as ASCII art using the selected banner.
- [x] The program supports letters, numbers, spaces, and printable special characters from the banner.
- [x] The program handles escaped `\n` sequences inside the input string.
- [x] Interactive mode repeats the input flow after printing the output.
- [x] Interactive mode exits when the user types `/exit` at the text, banner, color, or substring prompt.
- [x] Argument mode exits after one render.

### Error Handling Acceptance Criteria

- [x] Empty interactive text is rejected with `No valid string input`.
- [x] Invalid interactive banner choice is rejected with `Invalid banner choice` and the prompt repeats.
- [x] Invalid interactive color values print the rendering error and return to the text input prompt.
- [x] Missing argument-mode text prints the usage message.
- [x] Empty argument-mode text prints the usage message.
- [x] Invalid argument-mode banner names print the usage message.
- [x] Invalid CLI formats print:

```text
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

- [x] Invalid color flag formats print:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

- [x] Too many arguments print the usage message.
- [x] Unsupported characters return a rendering error.
- [x] Invalid color names, hexadecimal color codes, RGB color codes, and HSL color codes return clear rendering errors.
- [x] Missing or corrupted banner files return an error.

### Technical Acceptance Criteria

- [x] The project uses only standard Go packages.
- [x] The project follows the agreed file structure.
- [x] Banner files are parsed into `map[rune][]string`.
- [x] CRLF banner files are normalized before parsing.
- [x] Rendering logic is separated from CLI logic.
- [x] Color matching and ANSI wrapping are testable inside the renderer package.
- [x] Printing logic is separated from rendering logic.
- [x] No ASCII-art generation logic exists inside `main.go`.
- [x] Unit tests exist for banner loading, input normalization, rendering, output printing, and CLI validation.
- [x] Unit tests exist for color lookup, color wrapping, substring matching, and color CLI validation.

### Output Correctness Criteria

Golden test expectations are documented in `docs/golden_tests.md`.

- [x] Official banners load with 95 mapped runes.
- [x] Empty normalized lines produce blank output rows.
- [x] Adjacent characters are composed horizontally row by row.
- [x] `PrintASCII` prints exactly one newline per output row.
- [x] Input `"\n"` with standard.txt produces one empty line
- [x] Input `"Hello\n"` with standard.txt prints only Hello
- [x] Color is applied to ASCII-art character blocks, not by replacing text after rendering.
- [x] ANSI reset codes are applied after colored output segments.

### Audit Test Cases

The implementation must pass the official audit-style commands below. Commands that pipe to `cat -e` must preserve all spaces and print one `$` marker per output line when viewed through `cat -e`.

Invalid format:

```bash
go run . "banana" standard abc
```

Expected output:

```text
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

Required rendering cases:

```bash
go run . "hello" standard | cat -e
go run . "hello world" shadow | cat -e
go run . "nice 2 meet you" thinkertoy | cat -e
go run . "you & me" standard | cat -e
go run . "123" shadow | cat -e
go run . "/(\")" thinkertoy | cat -e
go run . "ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow | cat -e
go run . "\"#$%&/()*+,-./" thinkertoy | cat -e
go run . "It's Working" thinkertoy | cat -e
```

Each command must render the input in the selected template with all characters aligned across exactly 8 output rows per normalized input line.

The program must also render correctly when given randomized upper/lowercase letters, numbers, printable special characters, or mixed printable input followed by `standard`, `shadow`, or `thinkertoy`.

Any unsupported banner name used as the second argument must print the official usage message.

Color audit cases:

```bash
go run . --color red "banana"
go run . --color=red "hello world"
go run . --color=#ff0000 "hello world"
go run . --color=RGB,255,0,0 "hello world"
go run . --color=hsl,120,100%,50% "hello world"
go run . --color=green "1 + 1 = 2"
go run . --color=yellow "(%&) ??"
go run . --color=orange GuYs "HeY GuYs"
go run . --color=blue B 'RGB()'
go run . "hello world" shadow --color=cyan world
```

The first command must print the exact color usage message. The other commands must render the selected text in color while preserving ASCII-art alignment. The banner-compatible form must use the selected banner and color only the requested substring when one is provided.

---

## 8. Architecture

The project follows a **Pipeline Architecture**.

This means the program processes data through a clear sequence of steps, where each step has one responsibility and passes its result to the next step.

```text
CLI mode selection
-> Option parsing
-> Banner loading and mapping
-> Input normalization
-> Colored position detection
-> ASCII rendering
-> Output printing
-> Repeat loop (interactive mode only)
```

### Why Pipeline Architecture

Pipeline architecture fits this project because the program has a fixed and predictable flow. The input must be collected first, then the banner and optional color settings must be selected, then the banner is loaded, then the string is normalized, rendered, and printed. Each stage depends on the output of the previous stage.

This keeps the code simple, readable, and easy to test.

### Benefits

- Clear separation of responsibilities
- Easier unit testing for each step
- Easier debugging because problems can be isolated to one stage
- Cleaner `main.go`, since it only coordinates the pipeline
- Easier future extension, such as adding file export or more banner styles

### Tradeoffs

- More files and packages than a single-file solution
- Slightly more structure than the minimum required
- Requires clear communication between packages

### Architecture Sketch

```text
+------------------------+
| cli.PickMode(os.Args) |
+------------------------+
            |
            v
+-------------------------------------+
| text + banner + color + substring   |
+-------------------------------------+
            |
            v
+------------------------+
| banner.LoadBanner      |
+------------------------+
            |
            v
+------------------------+
| renderer.NormalizeInput|
+------------------------+
            |
            v
+------------------------+
| renderer.BuildASCII    |
+------------------------+
            |
            v
+------------------------+
| output.PrintASCII      |
+------------------------+
            |
            v
+-------------------------------------+
| repeat loop (interactive mode only) |
+-------------------------------------+
```

---


## 9. Milestones

1. Banner Parsing Works
2. Rendering Engine Works
3. Output Printing Works
4. Interactive CLI Works
5. Full Integration Completed
6. Testing & Cleanup Completed

---

## 10. Risks / Resolved Notes

### Risks

- Incorrect banner parsing leading to broken output
- Incorrect handling of `\n` sequences
- Unsupported characters causing errors or crashes
- Output misalignment due to incorrect row building
- Infinite loop or improper CLI exit handling
- Missing or invalid banner files
- Tight coupling between CLI, rendering, and banner logic

### Resolved Notes

- Error message text is simple and tested where it is user-facing.
- Banner files are loaded once per selected render cycle.
- Trailing banner file newlines are handled by `LoadBanner`.
- User-facing color values support preset names, `#RRGGBB`, `rgb,<red>,<green>,<blue>`, and `hsl,<hue>,<saturation>%,<lightness>%`.
- RGB and HSL use comma-only syntax so audit commands can pass them as one `--color=` argument without shell-sensitive parentheses.
- Color is still rendered internally with ANSI escape sequences.

## 11. Integration Contract

These contracts define the boundaries between packages and are intended to minimize coupling between CLI, banner parsing, rendering, and output responsibilities.

Current function signatures:

### Banner

```go
LoadBanner(path string) (map[rune][]string, error)
```

### Renderer

```go
NormalizeInput(input string) []string
BuildASCII(lines []string, banner map[rune][]string, substring string, color string) ([]string, error)
getANSIColor(color string) (string, error)
colorize(text string, colorCode string) string
coloredPositions(s string, substring string) map[int]bool
```

### Output

```go
PrintASCII(output []string)
```

### CLI

```go
PickMode(args []string) (input string, bannerPath string, color string, substring string, exit bool, err error)
RunArgumentMode(args []string) (input string, bannerPath string, color string, substring string, err error)
RunInteractive() (input string, bannerPath string, color string, substring string, exit bool, err error)
```
