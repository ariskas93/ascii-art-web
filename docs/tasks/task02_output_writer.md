# Task 02 - Output Writer

Owner: Programmer 2

## Goal

Add file output support inside `internal/output` without moving rendering or CLI parsing into this package.

## Recommended Functions

Use one shared writer function so terminal printing and file writing behave the same way:

```go
func WriteASCII(w io.Writer, lines []string) error
func PrintASCII(lines []string) error
func WriteASCIIFile(fileName string, lines []string) error
```

If the team wants a smaller first step, the minimum required function is:

```go
func WriteASCIIFile(fileName string, lines []string) error
```

## Responsibilities

`WriteASCIIFile` should:

- Create or overwrite the target file.
- Write every ASCII-art row exactly once.
- Add one newline after each row, matching terminal output.
- Preserve spaces and alignment.
- Preserve ANSI color codes if the renderer already added them.
- Return errors instead of printing or exiting.

`PrintASCII` should keep terminal output behavior, but it can call `WriteASCII(os.Stdout, lines)` internally.

## Tests

Add or update output tests for:

- `PrintASCII` still prints one newline per output row.
- `WriteASCIIFile` creates a file.
- Existing file contents are overwritten.
- Written file content exactly matches the ASCII rows with newlines.
- Empty output writes an empty file or no rows consistently with the shared writer behavior.
- File creation/write errors are returned.

Use `t.TempDir()` for file-writing tests.
