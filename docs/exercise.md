# ASCII Art Color Project

## Description

This project extends the ASCII Art program by adding support for colored output.

The program must manipulate colors using the following flag format:

```bash
--color=<color> <substring to be colored>
```

- `--color` is the required flag
- `<color>` is the desired color
- `<substring to be colored>` is the part of the text that should be colored

This implementation supports four user-facing color notations:

- preset color names, such as `red`
- full hexadecimal colors, such as `#ff0000`
- RGB comma form, such as `rgb,255,0,0`
- HSL comma form, such as `hsl,0,100%,50%`

ANSI escape sequences are still used internally for terminal output. Raw ANSI input such as `ansi:31` is not supported.

The substring can contain:

- A single character
- A word
- Multiple characters

If no substring is specified, the entire string must be colored.

---

## Invalid Flag Format

The flag must follow **exactly** this format:

```bash
--color=<color>
```

Any invalid format must return the following usage message:

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

---

## Example

### Command

```bash
go run . --color=red kit "a king kitten have kit"
```

### Expected Behavior

The substring `kit` must be colored in:

- `kitten`
- `kit`

Other supported examples:

```bash
go run . --color=#ff0000 "hello"
go run . --color=RGB,255,0,0 "hello"
go run . --color=hsl,120,100%,50% "hello"
go run . "hello world" shadow --color=cyan world
```

---

## Compatibility With Other Options

If other optional ASCII-Art projects are implemented, the program must still support:

- Additional valid `[OPTION]`
- Additional valid `[BANNER]`

The program must also continue working with a single `[STRING]` argument.

Example:

```bash
go run . "hello"
```

---

# Requirements

- The project must be written in Go
- The code must follow good practices
- Unit tests are recommended

---

# Allowed Packages

Only standard Go packages are allowed.

---

# Learning Objectives

This project will help you learn about:

- The Go file system (`fs`) API
- Color converters
- Data manipulation
- Terminal display
