# Golden Tests

These edge cases describe the expected behavior at the package boundaries. Keep the tests small and focused so each package proves only its own responsibility.

| Area | Edge case | Expected behavior | Test file |
| --- | --- | --- | --- |
| Banner loader | All official banners load successfully | `standard`, `shadow`, and `thinkertoy` each produce 95 mapped runes, and boundary runes `' '` and `'~'` each have 8 lines | `internal/banner/loader_test.go` |
| Banner loader | Windows CRLF banner content | Loaded banner lines must not contain hidden `\r` characters | `internal/banner/loader_test.go` |
| Banner loader | Missing character block | A banner with fewer than 95 blocks is invalid | `internal/banner/loader_test.go` |
| CLI argument mode | Whitespace-only text argument | Whitespace text is valid because spaces are renderable characters | `internal/cli/cli_test.go` |
| CLI argument mode | Escaped newline text argument | Text like `Hello\nThere` is preserved for the renderer to normalize | `internal/cli/cli_test.go` |
| CLI color mode | Whole-string color form | `--color=red "hello"` returns text `hello`, color `red`, empty substring, and `banners/standard.txt` | `internal/cli/cli_test.go` |
| CLI color mode | Whole-string hex color form | `--color=#ff0000 "hello"` returns text `hello`, color `#ff0000`, empty substring, and `banners/standard.txt` | `internal/cli/cli_test.go` |
| CLI color mode | Whole-string RGB color form | `--color=RGB,255,0,0 "hello"` returns text `hello`, color `RGB,255,0,0`, empty substring, and `banners/standard.txt` | `internal/cli/cli_test.go` |
| CLI color mode | Whole-string HSL color form | `--color=hsl,0,100%,50% "hello"` returns text `hello`, color `hsl,0,100%,50%`, empty substring, and `banners/standard.txt` | `internal/cli/cli_test.go` |
| CLI color mode | Substring color form | `--color=orange GuYs "HeY GuYs"` returns text `HeY GuYs`, color `orange`, substring `GuYs`, and `banners/standard.txt` | `internal/cli/cli_test.go` |
| CLI color mode | One-character substring | `--color=blue B "RGB()"` returns substring `B` and does not treat `RGB()` as a banner | `internal/cli/cli_test.go` |
| CLI color mode | Banner-compatible color form | `"hello world" shadow --color=cyan world` returns banner `banners/shadow.txt`, color `cyan`, substring `world` | `internal/cli/cli_test.go` |
| CLI color mode | Invalid separated color flag | `--color red "banana"` returns the exact color usage message | `internal/cli/cli_test.go` |
| CLI color mode | Empty color value | `--color= "banana"` returns the exact color usage message | `internal/cli/cli_test.go` |
| CLI color mode | Misspelled color flag | `-color=red` and `--colour=red` return the exact color usage message | `internal/cli/cli_test.go` |
| CLI color mode | Missing string after color flag | `--color=red` returns the exact color usage message | `internal/cli/cli_test.go` |
| CLI color mode | Multiple color flags | A command containing two `--color=<color>` flags is invalid | `internal/cli/cli_test.go` |
| CLI color mode | Order-flexible whole-string color | `"hello" --color=red` returns text `hello`, color `red`, and an empty substring | `internal/cli/cli_test.go` |
| CLI output mode | Order-flexible output flag | `"hello" --output=out.txt` returns text `hello` and output file `out.txt` | `internal/cli/cli_test.go` |
| CLI alignment mode | Valid alignment values | `--align=left`, `--align=center`, `--align=right`, and `--align=justify` return the selected alignment | `internal/cli/cli_test.go` |
| CLI alignment mode | Invalid alignment formats | Separated, empty, misspelled, uppercase, unsupported, and duplicate alignment flags return the exact alignment usage message | `internal/cli/cli_test.go` |
| CLI option mode | Unknown options | Unknown `--...` flags return the supported-options usage message | `internal/cli/cli_test.go` |
| CLI banner mode | Invalid banner value | Invalid banner names return `Invalid banner: <value>` followed by the general usage message | `internal/cli/cli_test.go` |
| CLI color mode | Shell-split RGB/HSL | Split color values such as `--color=rgb,255, 0, 0` are invalid argument formats | `internal/cli/cli_test.go` |
| CLI interactive mode | Windows-style user input | `\r\n` input is trimmed correctly at prompts | `internal/cli/interactive_test.go` |
| Renderer normalization | Leading escaped newline | `\nHello` becomes `["", "Hello"]` | `internal/renderer/normalize_test.go` |
| Renderer normalization | Trailing escaped newline | `Hello\n` becomes `["Hello"]` | `internal/renderer/normalize_test.go` |
| Renderer normalization | Real newline character | Real newline characters are not split by `NormalizeInput`; only escaped `\n` sequences are split | `internal/renderer/normalize_test.go` |
| Renderer color helpers | Supported preset colors | `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, and `orange` return non-empty ANSI codes | `internal/renderer/color_test.go` |
| Renderer color helpers | Supported hex notation | `#ff0000`, `#FF0000`, and other full `#RRGGBB` values return ANSI truecolor escape codes | `internal/renderer/color_test.go` |
| Renderer color helpers | Supported RGB notation | `rgb,255,0,0` and case variants such as `RGB,255,0,0` return ANSI truecolor escape codes | `internal/renderer/color_test.go` |
| Renderer color helpers | Supported HSL notation | `hsl,0,100%,50%` and case variants such as `HSL,0,100%,50%` return ANSI truecolor escape codes | `internal/renderer/color_test.go` |
| Renderer color helpers | Invalid color name | Unknown color names return `Invalid color name` | `internal/renderer/color_test.go` |
| Renderer color helpers | Invalid hex color | Short, long, or non-hex `#...` values return `Invalid hexadecimal color code` | `internal/renderer/color_test.go` |
| Renderer color helpers | Invalid RGB color | Out-of-range, non-integer, or wrong-shape `rgb...` values return `Invalid RGB color code` | `internal/renderer/color_test.go` |
| Renderer color helpers | Invalid HSL color | Out-of-range, missing percent, non-integer, or wrong-shape `hsl...` values return `Invalid HSL color code` | `internal/renderer/color_test.go` |
| Renderer color helpers | ANSI reset wrapping | `colorize("abc", "\033[31m")` returns color code + `abc` + `\033[0m` | `internal/renderer/color_test.go` |
| Renderer substring helpers | Empty substring | Empty substring marks every character position in the line | `internal/renderer/substring_test.go` |
| Renderer substring helpers | Missing substring | Missing substring returns an empty position map | `internal/renderer/substring_test.go` |
| Renderer substring helpers | Repeated matches | Multiple occurrences of the same substring all mark their character positions | `internal/renderer/substring_test.go` |
| Renderer substring helpers | Case-sensitive matching | `GuYs` matches `GuYs` but not `guys`, `GUYS`, or `Guys` | `internal/renderer/substring_test.go` |
| Renderer substring helpers | Two-letter substring | A two-character substring marks exactly those two positions for each match | `internal/renderer/substring_test.go` |
| Renderer build | Horizontal composition | Adjacent characters are joined row-by-row, not printed as separate blocks | `internal/renderer/build_test.go` |
| Renderer build | Empty input line between rendered lines | Empty normalized lines produce one blank output line and do not break following lines | `internal/renderer/build_test.go` |
| Renderer build color | Whole-string color | `BuildASCII(..., "", "#FF0000")` wraps every rendered character block row with truecolor red and reset codes | `internal/renderer/build_test.go` |
| Renderer build color | Substring color | `BuildASCII(..., "GuYs", "orange")` wraps only the `GuYs` character blocks | `internal/renderer/build_test.go` |
| Renderer build color | One-letter color | `BuildASCII(..., "B", "blue")` colors only the `B` in `RGB()` | `internal/renderer/build_test.go` |
| Renderer build color | Missing substring | If a provided substring is not found, output is valid uncolored ASCII art | `internal/renderer/build_test.go` |
| Renderer build color | Newline boundaries | A substring never matches across normalized input lines | `internal/renderer/build_test.go` |
| Renderer build color | Invalid color | An invalid name, hex code, RGB value, or HSL value returns an error before producing output | `internal/renderer/build_test.go` |
| Renderer build color | Spaces preserved | Coloring part of a string must not remove or trim spaces from rendered output | `internal/renderer/build_test.go` |
| Output printer | Exact stdout format | Each provided output row is printed with one newline, and an empty slice prints nothing | `internal/output/print_test.go` |
| Output alignment | ANSI-aware visible width | ANSI escape sequences do not count as visible width when padding aligned output | `internal/output/alignment_test.go` |
| Output alignment | Justify uses structured segments | Justify expands only spaces between rendered word segments, not spaces inside ASCII-art letters | `internal/output/alignment_test.go` |
| Output width | Linux non-terminal file | A regular file returns width `0` and detection status `false` | `internal/output/terminal_width_linux_test.go` |
| Output width | Live terminal width | Terminal output detects the current terminal width when output is produced instead of assuming a fixed width | manual terminal check |
| Output width | Redirected standard stream | Width detection tries standard output, standard input, and standard error before using fallbacks | `internal/output/*_test.go` |
| Output width | Detection fallback | If no live terminal is available, a valid `COLUMNS` value is used; otherwise width falls back to `80` | `internal/output/alignment_test.go` |
| File output alignment | Current terminal width | `--output=<file>` calculates alignment using the current launching terminal width, and writes the resulting padding as fixed spaces | `internal/output/*_test.go` plus manual file check |

---

## Color Golden Rules

Color output uses ANSI escape codes internally. User-facing color values support four notations:

- named colors, such as `red` or `orange`
- full hexadecimal colors in `#RRGGBB` form, such as `#ff0000`
- RGB colors in `rgb,<red>,<green>,<blue>` form, such as `rgb,255,0,0`
- HSL colors in `hsl,<hue>,<saturation>%,<lightness>%` form, such as `hsl,0,100%,50%`

Expected named color codes:

```text
black   -> \033[30m
red     -> \033[31m
green   -> \033[32m
yellow  -> \033[33m
blue    -> \033[34m
magenta -> \033[35m
cyan    -> \033[36m
white   -> \033[37m
orange  -> \033[38;5;208m
reset   -> \033[0m
```

Expected hexadecimal value notation:

```text
#ff0000 -> \033[38;2;255;0;0m
#FF0000 -> \033[38;2;255;0;0m
#00ff88 -> \033[38;2;0;255;136m
```

Expected RGB and HSL notation:

```text
rgb,255,0,0       -> \033[38;2;255;0;0m
RGB,0,255,0       -> \033[38;2;0;255;0m
hsl,0,100%,50%    -> \033[38;2;255;0;0m
HSL,120,100%,50%  -> \033[38;2;0;255;0m
```

Golden assertions for colored rendering:

- Color is applied to rendered ASCII-art character blocks, not to raw input before rendering.
- Every colored segment is followed by the reset code.
- Uncolored character blocks stay byte-for-byte identical to normal rendering.
- Output keeps the same visual 8-row structure per non-empty normalized line.
- Spaces and banner alignment are preserved.
- Substring matching is case-sensitive.
- Substring matching is repeated for every occurrence in a line.
- Substring matching does not cross escaped newline boundaries after normalization.
- A missing substring renders valid uncolored ASCII art.
- Multiple color flags in one command are not supported.

## Alignment Golden Rules

Alignment output will use structured renderer data.

Expected alignment values:

```text
left
center
right
justify
```

Golden assertions for aligned rendering:

- Left alignment is byte-for-byte equivalent to the flattened structured output.
- Center alignment adds left padding based on visible row width.
- Right alignment adds left padding based on visible row width.
- Rows wider than the configured width are printed unchanged.
- ANSI color sequences do not count toward visible width.
- Justify expands only spaces between word segments.
- Justify does not expand spaces inside ASCII-art letters.
- Justify leaves one-word lines left-aligned.
- Leading and trailing spaces are preserved.
- File output and terminal output use the same aligned rows.
- Terminal and file output use the current terminal width measured when output
  is produced.
- A saved text file keeps its written padding and does not dynamically realign
  when opened in a differently sized viewer.

## Audit Golden Outputs

These cases mirror the audit commands. Exact ASCII-art outputs are shown as they appear through `cat -e`, so those lines end with `$`. Color cases describe the ANSI and alignment assertions to verify.

### Invalid Argument Count

Command:

```bash
go run . "banana" standard abc
```

Expected output:

```text
Usage: go run . [STRING] [BANNER]

EX: go run . something standard
```

### Invalid Color Flag Format

Command:

```bash
go run . --color red "banana"
```

Expected output:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

### Empty Color Value

Command:

```bash
go run . --color= "banana"
```

Expected output:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

### Misspelled Color Flag

Commands:

```bash
go run . -color=red "banana"
go run . --colour=red "banana"
```

Expected output for each command:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

### Multiple Color Flags

Command:

```bash
go run . --color=red --color=blue "hello"
```

Expected behavior:

- The command is invalid because only one color flag is supported.
- It must not attempt to render `hello`.

### Invalid Color Notation

Commands:

```bash
go run . --color=#f00 "hello"
go run . --color=rgb,256,0,0 "hello"
go run . --color=hsl,361,100%,50% "hello"
```

Expected behavior:

- Short or malformed hexadecimal values return `Invalid hexadecimal color code`.
- Invalid RGB values return `Invalid RGB color code`.
- Invalid HSL values return `Invalid HSL color code`.
- Unsupported color names return `Invalid color name`.

### Whole String Red

Command:

```bash
go run . --color=red "hello world"
```

Expected behavior:

- Renders `hello world` with `banners/standard.txt`.
- Every rendered character block is wrapped with red `\033[31m` and reset `\033[0m`.
- Output remains aligned across 8 rows.

### Whole String Hex Value

Command:

```bash
go run . --color=#ff0000 "hello world"
```

Expected behavior:

- Renders `hello world` with `banners/standard.txt`.
- Every rendered character block is wrapped with `\033[38;2;255;0;0m` and reset `\033[0m`.
- CLI preserves the color value `#ff0000` exactly for renderer validation.

### Whole String Uppercase Hex Value

Command:

```bash
go run . --color=#FFAA00 "hello"
```

Expected behavior:

- Renders `hello` with `banners/standard.txt`.
- Every rendered character block is wrapped with `\033[38;2;255;170;0m` and reset `\033[0m`.
- Uppercase hex digits are accepted.

### Whole String RGB Value

Command:

```bash
go run . --color=RGB,255,0,0 "hello"
```

Expected behavior:

- Renders `hello` with `banners/standard.txt`.
- Every rendered character block is wrapped with `\033[38;2;255;0;0m` and reset `\033[0m`.
- The RGB notation prefix is case-insensitive.

### Whole String HSL Value

Command:

```bash
go run . --color=hsl,120,100%,50% "hello"
```

Expected behavior:

- Renders `hello` with `banners/standard.txt`.
- Every rendered character block is wrapped with `\033[38;2;0;255;0m` and reset `\033[0m`.
- The HSL notation prefix is case-insensitive.

### Whole String Green With Numbers

Command:

```bash
go run . --color=green "1 + 1 = 2"
```

Expected behavior:

- Renders numbers, spaces, and symbols using `banners/standard.txt`.
- Every rendered character block is wrapped with green `\033[32m` and reset `\033[0m`.
- Spaces remain visible as banner spacing; they are not trimmed.

### Whole String Yellow With Special Characters

Command:

```bash
go run . --color=yellow "(%&) ??"
```

Expected behavior:

- Renders all printable special characters using `banners/standard.txt`.
- Every rendered character block is wrapped with yellow `\033[33m` and reset `\033[0m`.
- Output remains 8 rows and aligned.

### Substring Orange

Command:

```bash
go run . --color=orange GuYs "HeY GuYs"
```

Expected behavior:

- Renders `HeY GuYs` with `banners/standard.txt`.
- Only the `G`, `u`, `Y`, and `s` character blocks in `GuYs` are wrapped with orange `\033[38;5;208m` and reset `\033[0m`.
- `HeY ` remains uncolored.
- Matching is case-sensitive.

### One Letter Blue

Command:

```bash
go run . --color=blue B 'RGB()'
```

Expected behavior:

- Renders `RGB()` with `banners/standard.txt`.
- Only the `B` character block is wrapped with blue `\033[34m` and reset `\033[0m`.
- `R`, `G`, `(`, and `)` remain uncolored.

### Two-Letter Substring

Command:

```bash
go run . --color=cyan ll "hello"
```

Expected behavior:

- Renders `hello` with `banners/standard.txt`.
- Only both `l` character blocks in the `ll` substring are wrapped with cyan `\033[36m` and reset `\033[0m`.
- `h`, `e`, and `o` remain uncolored.

### Missing Substring

Command:

```bash
go run . --color=red xyz "hello"
```

Expected behavior:

- Renders valid uncolored ASCII art for `hello`.
- Does not crash.
- Does not emit ANSI color codes because `xyz` is not present.

### Newline Boundary Color

Command:

```bash
go run . --color=magenta lohe "hello\nhello"
```

Expected behavior:

- Renders two normalized lines: `hello` and `hello`.
- Does not color across the escaped newline boundary.
- Since `lohe` does not appear inside either normalized line, output remains uncolored.

### Banner-Compatible Color Extension

Command:

```bash
go run . "hello world" shadow --color=cyan world
```

Expected behavior:

- Renders with `banners/shadow.txt`.
- Only the `world` character blocks are wrapped with cyan `\033[36m` and reset `\033[0m`.
- `hello ` remains uncolored.

### hello standard

Command:

```bash
go run . "hello" standard | cat -e
```

Expected output:

```text
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
```

### hello world shadow

Command:

```bash
go run . "hello world" shadow | cat -e
```

Expected output:

```text
                                                                                        $
_|                _| _|                                                     _|       _| $
_|_|_|     _|_|   _| _|   _|_|         _|      _|      _|   _|_|   _|  _|_| _|   _|_|_| $
_|    _| _|_|_|_| _| _| _|    _|       _|      _|      _| _|    _| _|_|     _| _|    _| $
_|    _| _|       _| _| _|    _|         _|  _|  _|  _|   _|    _| _|       _| _|    _| $
_|    _|   _|_|_| _| _|   _|_|             _|      _|       _|_|   _|       _|   _|_|_| $
                                                                                        $
                                                                                        $
```

### nice 2 meet you thinkertoy

Command:

```bash
go run . "nice 2 meet you" thinkertoy | cat -e
```

Expected output:

```text
                                                                       $
                       --                       o                      $
     o                o  o                      |                      $
o-o     o-o o-o         /        o-O-o o-o o-o -o-       o  o o-o o  o $
|  | | |    |-'        /         | | | |-' |-'  |        |  | | | |  | $
o  o |  o-o o-o       o--o       o o o o-o o-o  o        o--O o-o o--o $
                                                            |          $
                                                         o--o          $
```

### you & me standard

Command:

```bash
go run . "you & me" standard | cat -e
```

Expected output:

```text
                                                                $
                                ___                             $
 _   _    ___    _   _         ( _ )          _ __ ___     ___  $
| | | |  / _ \  | | | |        / _ \/\       | '_ ` _ \   / _ \ $
| |_| | | (_) | | |_| |       | (_>  <       | | | | | | |  __/ $
 \__, |  \___/   \__,_|        \___/\/       |_| |_| |_|  \___| $
 __/ /                                                          $
|___/                                                           $
```

### 123 shadow

Command:

```bash
go run . "123" shadow | cat -e
```

Expected output:

```text
                       $
  _|   _|_|   _|_|_|   $
_|_| _|    _|       _| $
  _|     _|     _|_|   $
  _|   _|           _| $
  _| _|_|_|_| _|_|_|   $
                       $
                       $
```

### Slash Parenthesis Quote thinkertoy

Command:

```bash
go run . "/(\")" thinkertoy | cat -e
```

Expected output:

```text
         o o    $
    o  / | | \  $
   /  o       o $
  o   |       | $
 /    o       o $
o      \     /  $
                $
                $
```

### Uppercase Alphabet shadow

Command:

```bash
go run . "ABCDEFGHIJKLMNOPQRSTUVWXYZ" shadow | cat -e
```

Expected output:

```text
                                                                                                                                                                                                                                                              $
  _|_|   _|_|_|     _|_|_| _|_|_|   _|_|_|_| _|_|_|_|   _|_|_| _|    _| _|_|_|       _| _|    _| _|       _|      _| _|      _|   _|_|   _|_|_|     _|_|     _|_|_|     _|_|_| _|_|_|_|_| _|    _| _|      _| _|          _| _|      _| _|      _| _|_|_|_|_| $
_|    _| _|    _| _|       _|    _| _|       _|       _|       _|    _|   _|         _| _|  _|   _|       _|_|  _|_| _|_|    _| _|    _| _|    _| _|    _|   _|    _| _|           _|     _|    _| _|      _| _|          _|   _|  _|     _|  _|         _|   $
_|_|_|_| _|_|_|   _|       _|    _| _|_|_|   _|_|_|   _|  _|_| _|_|_|_|   _|         _| _|_|     _|       _|  _|  _| _|  _|  _| _|    _| _|_|_|   _|  _|_|   _|_|_|     _|_|       _|     _|    _| _|      _| _|    _|    _|     _|         _|         _|     $
_|    _| _|    _| _|       _|    _| _|       _|       _|    _| _|    _|   _|   _|    _| _|  _|   _|       _|      _| _|    _|_| _|    _| _|       _|    _|   _|    _|       _|     _|     _|    _|   _|  _|     _|  _|  _|     _|  _|       _|       _|       $
_|    _| _|_|_|     _|_|_| _|_|_|   _|_|_|_| _|         _|_|_| _|    _| _|_|_|   _|_|   _|    _| _|_|_|_| _|      _| _|      _|   _|_|   _|         _|_|  _| _|    _| _|_|_|       _|       _|_|       _|         _|  _|     _|      _|     _|     _|_|_|_|_| $
                                                                                                                                                                                                                                                              $
                                                                                                                                                                                                                                                              $
```

### Special Characters thinkertoy

Command:

```bash
go run . "\"#$%&/()*+,-./" thinkertoy | cat -e
```

Expected output:

```text
o o         | |                                                  $
| |  | |   -O-O-      O          o  / \  o | o                 o $
    -O-O- o | |   o  /    o     /  o   o  \|/   |             /  $
     | |   -O-O-    /    /|    o   |   | --O-- -o-           o   $
    -O-O-   | | o  /  o o-O-  /    o   o  /|\   |    o-o    /    $
     | |   -O-O-  O       |  o      \ /  o | o     o     O o     $
            | |                                    |             $
                                                                 $
```

### It's Working thinkertoy

Command:

```bash
go run . "It's Working" thinkertoy | cat -e
```

Expected output:

```text
          o                                              $
o-O-o  o  |           o       o         o                $
  |    |              |       |         | /  o           $
  |   -o-   o-o       o   o   o o-o o-o OO     o-o  o--o $
  |    |     \         \ / \ /  | | |   | \  | |  | |  | $
o-O-o  o    o-o         o   o   o-o o   o  o | o  o o--O $
                                                       | $
                                                    o--o $
```
