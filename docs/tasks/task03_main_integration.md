# Task 03 - Main Integration

Owner: Programmer 3

## Goal

Connect the CLI output option to the final output step in `main.go`.

The renderer should not know whether output goes to the terminal or to a file.

## Flow

The final pipeline should stay:

```text
CLI args
-> parse options
-> load banner
-> normalize input
-> build ASCII-art output
-> print or write output
```

After rendering:

```go
if opts.OutputFile != "" {
	err = output.WriteASCIIFile(opts.OutputFile, ascii)
} else {
	err = output.PrintASCII(ascii)
}
```

## Responsibilities

- Update `main.go` to receive the new CLI options shape.
- Keep `go run . "hello"` working.
- Write to a file when `OutputFile` is set.
- Do not also print ASCII art to the terminal when output file mode is used.
- Print errors to stderr and exit non-zero.
- Interactive mode now asks whether to print or write to a file through `PickModeOptions`; respect `OutputFile` the same way as argument mode.
- Do not modify banner parsing or renderer logic unless required by the CLI contract.

## Manual Checks

Run:

```bash
go test ./...
go run . "hello"
go run . --output=banner.txt "hello"
go run . --output=banner.txt "hello" standard
go run . --output=banner.txt "Hello There!" shadow
```

Then check:

```bash
cat -e banner.txt
```

Invalid output formats should print exactly:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```
