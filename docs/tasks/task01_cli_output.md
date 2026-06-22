# Task 01 - CLI Output Option

Owner: Programmer 1

## Goal

Extend argument parsing so the program accepts the output option:

```bash
go run . --output=<fileName.txt> "hello"
go run . --output=<fileName.txt> "hello" standard
```

The parser must still accept the existing single-string mode:

```bash
go run . "hello"
```

If color support is already implemented, the parser should keep accepting correctly formatted color options and banners together with output.

The options-aware interactive path should also collect output intent after the existing text, banner, color, and substring prompts.

## Shared Data Contract

Prefer replacing long return value lists with one options struct:

```go
type Options struct {
	Text       string
	BannerPath string
	Color      string
	Substring  string
	OutputFile string
}
```

`OutputFile == ""` means print to the terminal.

`OutputFile != ""` means write the rendered ASCII art to that file.

## Usage Message

For invalid output option formats, return exactly:

```text
Usage: go run . [OPTION] [STRING] [BANNER]

EX: go run . --output=<fileName.txt> something standard
```

## Accepted Output Cases

```bash
go run . --output=banner.txt "hello"
go run . --output=banner.txt "hello" standard
go run . --output=banner.txt "hello" shadow
go run . --output=banner.txt "hello" thinkertoy
```

## Invalid Output Cases

These must return the output usage message:

```bash
go run . --output "hello"
go run . --output= "hello"
go run . -output=banner.txt "hello"
go run . --outputs=banner.txt "hello"
go run . --output banner.txt "hello"
```

## Color Compatibility

If color is implemented, keep the existing color usage message for invalid color flags and use the output usage message only when the output flag is malformed.

Recommended supported combinations:

```bash
go run . --output=banner.txt --color=red "hello"
go run . --output=banner.txt --color=red kit "a king kitten have kit"
go run . --output=banner.txt "hello" standard --color=red
go run . --output=banner.txt "hello world" shadow --color=cyan world
```

Do not place rendering logic in the CLI package. The CLI only validates arguments and returns parsed options.

## Interactive Output Compatibility

In options-aware interactive mode:

- Ask whether the user wants to write output to a file or print to the terminal.
- Empty input, `n`, and `no` mean print to terminal and leave `OutputFile == ""`.
- `y` and `yes` mean ask for an output file name.
- Empty output file names are rejected and the file-name prompt repeats.
- `/exit` exits at the output choice and output file name prompts.
- The CLI only stores the selected file name in `OutputFile`; `main.go` decides whether to print or write.

## Tests

Add or update CLI tests for:

- Single string still works.
- Output flag with default banner.
- Output flag with explicit banner.
- Output flag stores only the file name in `OutputFile`.
- Invalid output flag formats return the exact output usage message.
- Existing color parsing still works.
- Output plus correctly formatted color options works if color support is present.
- Interactive output choice returns an empty or populated `OutputFile`.
- Interactive invalid output choices and empty output file names are handled.
