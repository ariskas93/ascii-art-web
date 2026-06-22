# Workflow - ASCII Art Output

## Objective

Add support for writing the rendered ASCII-art result into a file with:

```bash
go run . --output=<fileName.txt> "hello"
go run . --output=<fileName.txt> "hello" standard
```

The program must still support:

```bash
go run . "hello"
```

If other optional projects are already implemented, such as color, the program should accept those correctly formatted options together with output.

Interactive mode should also collect output intent through prompts. When the user chooses file output, the CLI stores the file name in `OutputFile`; otherwise `OutputFile` stays empty and the result is printed to the terminal.

## Shared Contract

Use one parsed options value across the CLI and main integration:

```go
type Options struct {
	Text       string
	BannerPath string
	Color      string
	Substring  string
	OutputFile string
}
```

`OutputFile == ""` means print to stdout.

`OutputFile != ""` means write to the named file.

## Task Split

### Task 01 - CLI Output Option

Owner: Programmer 1

File: `docs/tasks/task01_cli_output.md`

Responsibilities:

- Parse `--output=<fileName.txt>`.
- In options-aware interactive mode, ask whether the user wants to write to a file or print to the terminal.
- If the interactive user chooses file output, ask for the output file name.
- Reject malformed output flags with the exact output usage message.
- Keep single-string mode working.
- Keep banner and correctly formatted color option compatibility.
- Return the shared options struct.

### Task 02 - Output Writer

Owner: Programmer 2

File: `docs/tasks/task02_output_writer.md`

Responsibilities:

- Add file writing in `internal/output`.
- Prefer shared writer logic for terminal and file output.
- Preserve spacing, newlines, and ANSI codes.
- Return errors instead of printing or exiting.
- Add output package tests.

### Task 03 - Main Integration

Owner: Programmer 3

File: `docs/tasks/task03_main_integration.md`

Responsibilities:

- Connect parsed `OutputFile` to the final output step.
- Print to terminal when no output file is set.
- Write to file when output file mode is set.
- Keep renderer and banner responsibilities unchanged.
- Run full tests and manual checks.

## Implementation Order

1. Agree on the `Options` struct and usage message constants.
2. Programmer 1 updates CLI parsing and CLI tests.
3. Programmer 2 adds output writer functions and output tests.
4. Programmer 3 integrates the final branch in `main.go`.
5. Run `go test ./...` and manual output checks.

## Exact Output Usage Message

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```
