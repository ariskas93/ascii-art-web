# ASCII Art Color Audit

## Functional

### Allowed Packages

Has the requirement for the allowed packages been respected?

Reminder for this project: **only standard Go packages are allowed**.

---

## Usage Error Handling

### Test

```bash
go run . --color red "banana"
```

### Expected Output

```text
Usage: go run . [OPTION] [STRING]

EX: go run . --color=<color> <substring to be colored> "something"
```

Does it display the same result as above?

---

## Basic Color Tests

### Test

```bash
go run . --color=red "hello world"
```

Does it display the expected result?

---

### Test

```bash
go run . --color=#ff0000 "hello world"
go run . --color=RGB,255,0,0 "hello world"
go run . --color=hsl,120,100%,50% "hello world"
```

Do the hexadecimal, RGB, and HSL forms display the expected colored output?

---

### Test

```bash
go run . --color=green "1 + 1 = 2"
```

Does it display the expected result?

---

### Test

```bash
go run . --color=yellow "(%&) ??"
```

Does it display the expected result?

---

## Substring Color Tests

### Test

Specify a substring to be colored from the second letter until the last letter.

Does it display the expected result, with the corresponding substring shown in the selected color?

---

### Test

Specify one letter to be colored, using the second letter.

Does it display the expected result, with the corresponding letter shown in the selected color?

---

### Test

Specify a substring of only two letters to be colored.

Does it display the expected result, with the corresponding substring shown in the selected color?

---

## Specific Case Tests

### Test

```bash
go run . --color=orange GuYs "HeY GuYs"
```

Expected behavior: color `GuYs`.

Does it display the expected result?

---

### Test

```bash
go run . --color=blue B 'RGB()'
```

Expected behavior: color only the `B`.

Does it display the expected result?

---

## Random Input Tests

### Test

Pass a random string with lowercase and uppercase letters, and a random color using the `--color=` flag.

Does it display the expected result?

---

### Test

Pass a random string with lowercase letters, numbers, and spaces, and a random color using the `--color=` flag.

Does it display the expected result?

---

### Test

Pass a random string with special characters and a random color using the `--color=` flag, specifying one letter to be colored.

Does it display the expected result?

---

### Test

Pass a random string with lowercase letters, uppercase letters, spaces, and numbers, with a random color using the `--color=` flag, specifying a substring to be colored.

Does it display the expected result?

---

## Auditor Decision

As an auditor, is this project up to every standard?

If not, why are you failing the project?

Possible reasons:

- Empty Work
- Incomplete Work
- Invalid Compilation
- Cheating
- Crashing
- Leaks

---

# General

## Usability

Is it easy and intuitive to specify the letter or letters to be colored?

---

## Multiple Colors

Can you use more than one color in the same string?

Expected answer: no. The current implementation supports one `--color=<color>` flag per command.

---

# Basic

## Performance

Does the project run quickly and effectively?

Check for:

- Excessive recursion
- Unnecessary data requests
- Inefficient processing

---

## Output Structure

Is the output of the program well structured?

Are the characters displayed correctly in line?

---

## Tests

Is there a test file for this code?

Are the tests checking each possible case?

---

## Code Quality

Does the code obey good practices?

---

# Social

Did you learn anything from this project?

Would you recommend or nominate this program as an example for the rest of the school?

---

# Bonus

## Multiple Color Notations

Try using different `--color` flag notations, such as:

```bash
--color=red
--color=#ff0000
--color=rgb,255,0,0
--color=hsl,0,100%,50%
```

Are preset names, hexadecimal, RGB, and HSL accepted as alternate ways to choose one foreground color?

Invalid notation checks:

```bash
--color=#f00
--color=rgb,256,0,0
--color=hsl,361,100%,50%
```

Do invalid hexadecimal, RGB, and HSL values return their specific error messages?
