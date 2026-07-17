# Project Handoff: ASCII Art Justify

Last audited: 2026-07-17  
Repository branch at audit time: `main` tracking `origin/main`  
Audited commit: `f054079` (`cli integrated`)  
Go module: `ascii-art`  
Declared and verified Go version: `1.22.2`

## 1. Purpose of This Document

This document is the operational handoff for the next AI agent that will extend this program for a later school project. It describes the checked-in implementation, not only the intended assignment design.

The distinction matters here: the repository accepts alignment options at the CLI, but does not visually align output yet. Several planning documents describe future code as though it were the target end state, and one older agent log contains integration notes that are now stale.

When this handoff conflicts with a planning statement, verify the current code and tests before changing anything.

## 2. Executive Snapshot

The project is a Go CLI that renders printable ASCII text with one of three fixed banner files. It currently has working support for:

- argument mode and a looping interactive mode;
- `standard`, `shadow`, and `thinkertoy` banners;
- named, hexadecimal, RGB, and HSL foreground colors;
- coloring all text or every occurrence of a selected substring;
- escaped `\n` input normalization;
- terminal output or file output through `--output=<file>`;
- parsing and validating `--align=left|center|right|justify`.

Alignment rendering is **not implemented**. `Options.Align` is collected, but `main.go` never passes it to the renderer or output package. Verified commands using `left`, `center`, and `justify` produced byte-for-byte identical output.

Current alignment task status:

| Work item | Status | Evidence |
| --- | --- | --- |
| Task 01: CLI alignment parsing and prompt | Complete | `Options.Align`, parser validation, interactive prompt, and CLI tests exist |
| Task 02: structured renderer | Not implemented | `BuildASCII` still returns flat `[]string`; no segment types exist |
| Task 03: aligned output | Not implemented | no width detection, ANSI-aware width helper, or alignment functions exist |
| Existing output-file integration | Complete | `main.go` calls `WriteASCIIFile`; this predates/overlaps Task 03 |

## 3. Verified Baseline

At the time of this audit:

```text
go test ./...
?    ascii-art                    [no test files]
ok   ascii-art/internal/banner
ok   ascii-art/internal/cli
ok   ascii-art/internal/output
ok   ascii-art/internal/renderer
```

The worktree was clean before this handoff file was added. The tracked `ascii-art` file is a statically linked Linux x86-64 executable with debug information. Treat it as a generated/stale artifact; use `go run .` or rebuild from source rather than assuming the binary matches the current code.

The project uses only Go standard-library imports and has no third-party module requirements.

## 4. Source-of-Truth Order

Use this priority when documents disagree:

1. The current school assignment and any instructions supplied to the next agent.
2. Current tests and externally supplied audit requirements.
3. Current implementation.
4. `docs/prd.md` and the three alignment task documents for the agreed future design.
5. `README.md`, `docs/golden_tests.md`, and other explanatory documents.
6. `agent_log.md`, which describes an earlier color/output integration stage and contains stale claims.

There is no tracked `AGENTS.md`; `.gitignore` explicitly ignores it. The current agent session was given external `AGENTS.md` instructions centered on the color assignment. A future agent must read whatever active instructions it receives and reconcile them with this broader repository before editing.

Notable documentation discrepancies:

- The directory/project objective is alignment/justify, while the externally supplied agent instructions describe the preceding color extension.
- `agent_log.md` says `main.go` is not integrated with options-aware output. That is no longer true.
- Planning documents describe a future structured renderer. The types and APIs shown there do not exist yet.
- Some older documents focus on the color audit and do not represent the alignment completion state.

## 5. Repository Map

```text
.
├── main.go                         integration loop
├── go.mod                         module ascii-art; Go 1.22.2
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt              CRLF file; loader normalizes it
├── internal/
│   ├── banner/
│   │   └── loader.go               validates and maps the 95 banner glyphs
│   ├── cli/
│   │   ├── args.go                 option/operand parser and exact usage text
│   │   ├── interactive.go          interactive prompts and option collection
│   │   └── mode_picker.go          argument-vs-interactive dispatch
│   ├── renderer/
│   │   ├── normalize.go            splits literal backslash-n sequences
│   │   ├── build.go                current flat eight-row renderer
│   │   ├── color.go                color parsing and ANSI conversion
│   │   └── substring.go            marks input positions to color
│   └── output/
│       ├── write.go                writes rows to an io.Writer
│       ├── print.go                stdout adapter
│       └── write_file.go           truncating file adapter
├── docs/
│   ├── prd.md                      most complete alignment specification
│   ├── exercise.md                 architecture/design summary
│   ├── workFlow.md                 task ordering and shared contract
│   ├── golden_tests.md             existing and planned edge cases
│   ├── audit_test.md               older color-focused audit checklist
│   └── tasks/                       three alignment task briefs
├── README.md                        user-facing current behavior and limitation
├── agent_log.md                     historical; partly stale
└── ascii-art                        tracked compiled executable
```

Tests live beside each package. There are no `main` package integration tests.

## 6. Runtime Pipeline

Current flow:

```text
os.Args
  -> cli.PickModeOptions
  -> cli.Options
  -> banner.LoadBanner
  -> renderer.NormalizeInput
  -> renderer.BuildASCII (flat []string)
  -> output.PrintASCII or output.WriteASCIIFile
```

`main.go` determines interactive mode only from `len(os.Args) == 1`. It loops after successful interactive renders and also retries after renderer errors such as an invalid color. Argument mode renders once and exits.

CLI or banner/output errors print to stderr and exit with status 1. A renderer error returns to the initial text prompt in interactive mode, but exits status 1 in argument mode.

Important gap:

```text
cli.Options.Align -> currently unused after parsing
```

The intended alignment flow replaces flat renderer output with structured segments and gives alignment/terminal-width ownership to `internal/output`.

## 7. Current CLI Contract

### Argument mode

Representative working commands:

```bash
go run . "Hello"
go run . "Hello" shadow
go run . --color=red "Hello"
go run . --color=red ll "Hello"
go run . "hello world" shadow --color=cyan world
go run . --output=out.txt "Hello"
go run . --align=center "Hello"
```

Accepted banner operands are exactly:

```text
standard, standard.txt, shadow, shadow.txt, thinkertoy, thinkertoy.txt
```

With one positional operand, it is always text. Thus `go run . "shadow"` renders the word rather than selecting the banner.

Options are collected independently of their location and must use `--name=value`. Duplicate color, output, or alignment flags are rejected. Unknown `--...` flags use the general option usage message.

The color substring is not a flag. It is recognized as the positional operand immediately after `--color=<color>` when the remaining operands otherwise require a substring interpretation. This creates an ambiguity when values can also form valid `[STRING] [BANNER]`; preserve existing test behavior unless a new specification explicitly changes it.

`RunArgumentModeOptions` is the primary parser used by `main`. Legacy tuple APIs remain:

```go
RunArgumentMode(args) (text, bannerPath, color, substring string, err error)
PickMode(args) (text, bannerPath, color, substring string, exit bool, err error)
```

The legacy path does not return output or alignment settings and rejects output mode. Avoid building new features on it.

### Exact option validation

- Color flag: `--color=<non-empty value>`.
- Output flag: `--output=<non-empty file name>`.
- Alignment flag: lowercase `--align=left`, `center`, `right`, or `justify`.
- Separated values such as `--align center` are invalid.
- Misspellings such as `-align=`, `--alignment=`, `-color=`, and `--colour=` are recognized as malformed forms and return their option-specific usage messages.
- Empty text `""` is invalid; whitespace-only text is valid.

Exact usage strings are constants in `internal/cli/args.go`. Do not casually reword them because tests and school audits may compare them literally.

### Interactive mode

Running `go run .` repeatedly asks for:

1. non-empty text or `/exit`;
2. banner number `1`, `2`, or `3`;
3. optional color;
4. optional substring if a color was supplied;
5. alignment, defaulting to `left` on an empty response;
6. whether to write to a file and, if yes, its name.

The shared buffered stdin reader is intentional. Recreating it on every loop can lose already-buffered scripted input.

## 8. Renderer and Data Contracts

### Banner loading

```go
func LoadBanner(path string) (map[rune][]string, error)
```

The loader:

- reads a fixed path selected by the CLI;
- converts CRLF to LF;
- requires exactly 95 blocks for ASCII 32 through 126;
- assumes nine physical lines per block: an empty separator then eight glyph rows;
- returns `invalid banner format` for the wrong total structure.

Do not modify the banner files. The loader's CRLF normalization is required for `thinkertoy.txt`.

### Input normalization

```go
func NormalizeInput(input string) []string
```

Only the two-character sequence backslash + `n` is split. A real newline rune is not split and later becomes an unsupported character.

Current tested behavior:

```text
"Hello\\nThere"   -> ["Hello", "There"]
"Hello\\n\\nThere" -> ["Hello", "", "There"]
"\\nHello"         -> ["", "Hello"]
"Hello\\n"         -> ["Hello"]
"\\n"              -> [""]
""                  -> []
```

The removal of one trailing escaped-newline entry is a tested compatibility rule, even though other assignment variants may expect a final empty line. If the next assignment changes this rule, update tests and documentation deliberately.

### Current flat rendering

```go
func BuildASCII(
    lines []string,
    banner map[rune][]string,
    substring string,
    color string,
) ([]string, error)
```

Each non-empty normalized line becomes eight output rows. Each empty normalized line becomes one blank output row. Character glyphs are concatenated horizontally. Any rune absent from the banner produces `unsupported character`.

The flat `[]string` return type is the main obstacle to safe justify alignment: banner glyphs themselves contain spaces, so final rows cannot distinguish inter-word spaces from glyph-internal spaces.

### Color behavior

Supported named colors:

```text
black red green yellow blue magenta cyan white orange
```

Also supported:

```text
#RRGGBB
rgb,R,G,B
hsl,H,S%,L%
```

RGB/HSL prefixes are case-insensitive; named colors are case-sensitive. RGB components are integers `0..255`; HSL accepts integer hue `0..360` and integer percentages `0..100`. Hex must contain exactly six hexadecimal digits.

Color is applied while building ASCII art. The renderer first marks source-character positions, then wraps each selected character's glyph row in an ANSI prefix and `\033[0m`. Do not implement coloring by searching the final rendered rows.

Substring matching is case-sensitive, finds overlapping and repeated matches, and is performed per normalized line. An empty substring colors the whole line; a missing substring leaves output uncolored. These helpers use byte-oriented positions, which is compatible with the only supported input range (printable ASCII) but should be revisited if Unicode is ever added.

### Output behavior

`WriteASCII` writes every row plus one newline. An empty row slice writes nothing.

`WriteASCIIFile` uses `os.Create`, so it creates or truncates the requested path. It does not validate an extension or create parent directories. ANSI codes are preserved in file output.

Banner paths are relative to the process working directory. Running the compiled binary outside the repository root will fail to find `banners/...` unless this architecture is changed.

## 9. Alignment Design That Must Be Preserved

The documents agree on a structured renderer owned by `internal/renderer`:

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

The preferred `BuildASCII` return changes from `[]string` to `RenderedASCII`, optionally with a `Flatten` compatibility helper.

This design is important, not cosmetic. Justification must expand only real gaps between word segments. It must never scan a completed ASCII-art row and stretch arbitrary spaces, because those spaces are part of glyph shapes.

Ownership boundaries:

- `internal/renderer`: segment source lines and render segments; no terminal width or alignment policy.
- `internal/output`: flatten rows, measure visible widths, detect width, and apply left/center/right/justify.
- `main.go`: connect options, rendering, and output without implementing algorithms.
- `internal/cli`: validate and store alignment only.

Planned width rules:

- read `COLUMNS` using the standard library;
- fall back to 80 when missing, invalid, or below 1;
- expose a width-explicit writer for deterministic tests;
- ignore ANSI escape sequences when measuring visible width;
- add no padding when a row is wider than the target width.

Planned justify rules:

- expand only gaps between word segments;
- preserve leading/trailing space segments without treating them as expandable gaps;
- leave lines with fewer than two words left-aligned;
- leave over-width lines left-aligned;
- distribute remainder spaces from left to right.

Read `docs/tasks/task02_structured_renderer.md` before changing the renderer and `docs/tasks/task03_output_alignment_integration.md` before changing output. The design may still need adjustment for the next school project's specification, but do not discard the segment distinction accidentally.

## 10. Tests and Verification

Primary command:

```bash
go test ./...
```

Current package coverage includes:

- official, missing, malformed, and CRLF banner loading;
- CLI modes, exact error messages, option combinations, duplicates, and interactive retries/exits;
- normalization of escaped newlines;
- flat rendering, empty lines, unsupported characters, and horizontal glyph composition;
- color syntax/conversion and renderer error propagation;
- repeated, overlapping, missing, case-sensitive substring matches;
- terminal row writing and truncating file output.

Major missing tests:

- no alignment output tests;
- no structured renderer tests;
- no ANSI-aware visible-width tests;
- no `COLUMNS` fallback tests;
- no `main.go` integration tests.

After alignment work, run at minimum:

```bash
go test ./...
go run . --align=left "hello" standard
go run . --align=center "hello" standard
go run . --align=right "hello" standard
go run . --align=justify "hello world" standard
go run . "hello world" shadow --align=center --color=cyan world
go run . --output=aligned.txt --align=right "hello" standard
```

For visual width tests, set a deterministic `COLUMNS` value or call a width-explicit function. Compare visible widths after removing ANSI control sequences, not raw byte lengths.

## 11. Fragile Areas and Known Risks

- **Alignment is silently ignored.** The CLI accepts it successfully, which can make manual smoke tests look complete unless leading spaces are compared.
- **Flat rendering cannot safely support justify.** Preserve source word/space structure before alignment.
- **ANSI codes distort `len(row)`.** Width calculations need an ANSI-aware visible-width function.
- **Exact messages are contractual.** CLI and color errors have literal test expectations.
- **Newline semantics are unusual.** Only literal `\\n` is recognized, and one trailing escaped newline is dropped.
- **Empty normalized lines use one row, not eight.** Preserve this existing visual convention unless the assignment explicitly changes it.
- **Banner layout parsing is positional.** Changing blank separator handling can shift every glyph.
- **Relative banner paths couple execution to the repository root.** Be aware when adding integration tests or packaging.
- **Output files are destructive by design.** `os.Create` overwrites existing files.
- **Color and substring parsing has positional ambiguity.** Add tests before refactoring CLI operand resolution.
- **Legacy CLI APIs coexist with `Options`.** Prefer the options-aware APIs; removing legacy functions may break tests or downstream school stages.
- **Interactive prompt order is tested.** Alignment currently occurs before the output-file question.
- **Tracked compiled binary may drift.** Source and tests are authoritative.

## 12. Safe Extension Points

For the documented alignment objective, the safest sequence is:

1. Introduce renderer segment types and tests while retaining a flattening path that proves old left-aligned output is unchanged.
2. Adapt renderer color-position tracking so segment boundaries do not reset or corrupt original character positions.
3. Add output helpers with an explicit width and ANSI-aware visible measurement.
4. Implement left, center, and right before justify.
5. Implement justify from structured word/space segments.
6. Change `main.go` only after renderer/output contracts compile and have package tests.
7. Verify terminal and file output share the same alignment path.
8. Update README/status documents and either refresh or clearly archive `agent_log.md`.

For a different next school project, keep the same package boundary principle: parse configuration in CLI, transform input in renderer, handle destination-specific behavior in output, and keep `main.go` as orchestration.

## 13. Non-Goals in the Current Design

- third-party packages;
- user-supplied custom banner paths;
- Unicode outside printable ASCII 32..126;
- multiple independent input strings;
- multiple color flags;
- GUI/web output;
- glyph scaling or resizing;
- justify by modifying spaces in already-flattened rows.

Do not broaden these without an explicit next-assignment requirement.

## 14. First Actions for the Next Agent

1. Read the new assignment and active agent instructions completely.
2. Run `git status --short --branch` and preserve unrelated user changes.
3. Run `go test ./...` before editing.
4. Read `main.go`, `internal/cli/args.go`, and the package being changed.
5. Read the relevant task document; treat proposed types as design, not existing code.
6. Confirm whether the next assignment wants the current escaped-newline behavior or a different one.
7. Add tests that demonstrate the new behavior before or alongside implementation.
8. Keep the standard-library-only rule and exact existing CLI contracts unless explicitly superseded.
9. Do not trust the tracked binary or stale status prose over source/tests.

## 15. Handoff Summary

This is a healthy, tested base for a further ASCII-art assignment. Color, banner selection, interactive mode, and output files are integrated. Alignment parsing is complete but alignment rendering is absent. The next agent's most important responsibility is to preserve existing behavior while introducing the structured representation required for safe justification—or, if the next school project has a different goal, to use the established CLI/banner/renderer/output boundaries without inheriting stale assumptions from the planning documents.
