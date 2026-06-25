# CLI Test Cases

## Interactive Mode

| Input sequence          | Expected                              |
| ----------------------- | ------------------------------------- |
| `""`                    | retry text prompt                     |
| `"   "`                 | accepted as text                      |
| `/exit`                 | `exit=true`                           |
| `hello` + `1`           | text=`hello`, standard path           |
| `hello` + `2`           | text=`hello`, shadow path             |
| `hello` + `3`           | text=`hello`, thinkertoy path         |
| `hello` + `9` + `2`     | retry banner only, then shadow path   |
| `hello` + `abc` + `1`   | retry banner only, then standard path |
| `hello` + `""` + `3`    | retry banner only, then thinkertoy    |
| `hello` + `/exit`       | `exit=true` at banner prompt          |
| `""` + `hello` + `1`    | retry text prompt, then standard path |
| `Hello\nThere` + `1`    | escaped newline text is preserved     |
| `hello` + ` 2 `         | trims banner choice whitespace        |
| `hello\r\n3\r\n`        | Windows-style input selects thinkertoy |
| EOF before text         | returns an error and `exit=true`      |
| EOF before banner       | returns an error and `exit=true`      |

## Argument Mode

| Args                                    | Expected                    |
| --------------------------------------- | --------------------------- |
| `["ascii-art"]`                         | missing text error          |
| `["ascii-art", "hello"]`                | standard path               |
| `["ascii-art", "   "]`                  | whitespace text is valid    |
| `["ascii-art", "Hello\\nThere"]`        | escaped newline is preserved |
| `["ascii-art", "hello", "standard"]`   | standard path               |
| `["ascii-art", "hello", "standard.txt"]` | standard path             |
| `["ascii-art", "hello", "shadow"]`     | shadow path                 |
| `["ascii-art", "hello", "shadow.txt"]` | shadow path                 |
| `["ascii-art", "hello", "thinkertoy"]` | thinkertoy path             |
| `["ascii-art", "hello", "thinkertoy.txt"]` | thinkertoy path       |
| `["ascii-art", "--color=red", "hello"]` | standard path, color=`red`, whole string |
| `["ascii-art", "--color=#ff0000", "hello"]` | standard path, color=`#ff0000`, whole string |
| `["ascii-art", "--color=#FF0000", "hello"]` | standard path, color=`#FF0000`, whole string |
| `["ascii-art", "--color=RGB,255,0,0", "hello"]` | standard path, color=`RGB,255,0,0`, whole string |
| `["ascii-art", "--color=hsl,0,100%,50%", "hello"]` | standard path, color=`hsl,0,100%,50%`, whole string |
| `["ascii-art", "--color=orange", "GuYs", "HeY GuYs"]` | standard path, color=`orange`, substring=`GuYs` |
| `["ascii-art", "--color=blue", "B", "RGB()"]` | standard path, color=`blue`, substring=`B` |
| `["ascii-art", "hello", "shadow", "--color=cyan"]` | shadow path, color=`cyan`, whole string |
| `["ascii-art", "hello world", "thinkertoy", "--color=green", "world"]` | thinkertoy path, color=`green`, substring=`world` |
| `["ascii-art", "hello", "--color=red"]` | standard path, color=`red`, whole string |
| `["ascii-art", "hello", "--color=red", "shadow"]` | shadow path, color=`red`, whole string |
| `["ascii-art", "--output=banner.txt", "hello"]` | standard path, output file=`banner.txt` |
| `["ascii-art", "hello", "--output=banner.txt"]` | standard path, output file=`banner.txt` |
| `["ascii-art", "hello world", "shadow", "--output=banner.txt", "--color=cyan", "world"]` | shadow path, color=`cyan`, substring=`world`, output file=`banner.txt` |
| `["ascii-art", ""]`                    | empty text error            |
| `["ascii-art", "hello", "unknown"]`    | invalid banner error        |
| `["ascii-art", "hello", "1"]`          | invalid banner error        |
| `["ascii-art", "hello", "standard", "extra"]` | too many args error |
| `["ascii-art", "--color", "red", "banana"]` | color usage error |
| `["ascii-art", "--color=", "banana"]` | color usage error |
| `["ascii-art", "--color=", "ban", "banana"]` | color usage error |
| `["ascii-art", "-color=red", "banana"]` | color usage error |
| `["ascii-art", "--colour=red", "banana"]` | color usage error |
| `["ascii-art", "--color=red"]` | color usage error |
| `["ascii-art", "--color=red", "--color=blue", "hello"]` | color usage error |
| `["ascii-art", "hello", "standard", "--color=red", "--color=blue"]` | color usage error |
| `["ascii-art", "hello", "unknown", "--color=red"]` | invalid banner error |
| `["ascii-art", "--color=rgb,255,", "0,", "0", "hello"]` | color usage error |
| `["ascii-art", "--color=hsl,0,", "100%,", "50%", "hello"]` | color usage error |

## Planned Alignment Argument Cases

| Args | Expected |
| --- | --- |
| `["ascii-art", "--align=left", "hello"]` | align=`left` |
| `["ascii-art", "--align=center", "hello"]` | align=`center` |
| `["ascii-art", "hello", "--align=right"]` | align=`right` |
| `["ascii-art", "hello world", "shadow", "--align=justify"]` | shadow path, align=`justify` |
| `["ascii-art", "--output=banner.txt", "--align=center", "hello"]` | output file=`banner.txt`, align=`center` |
| `["ascii-art", "--color=cyan", "world", "hello world", "--align=center"]` | color=`cyan`, substring=`world`, align=`center` |
| `["ascii-art", "--align", "center", "hello"]` | alignment usage error |
| `["ascii-art", "--align=", "hello"]` | alignment usage error |
| `["ascii-art", "-align=center", "hello"]` | alignment usage error |
| `["ascii-art", "--alignment=center", "hello"]` | alignment usage error |
| `["ascii-art", "--align=middle", "hello"]` | alignment usage error |
| `["ascii-art", "--align=left", "--align=right", "hello"]` | alignment usage error |

## Mode Picker

| Args                         | Expected                                  |
| ---------------------------- | ----------------------------------------- |
| `["ascii-art", "hello"]`     | argument mode with standard path          |
| `["ascii-art", "hello", "shadow"]` | argument mode with shadow path      |
| `["ascii-art"]` with `/exit` | interactive mode exits with `exit=true`   |

## Loop Ownership

`RunInteractiveOptions` returns one valid options value. The repeat loop lives in `main.go`, which calls `PickModeOptions(os.Args)` again after printing or writing output when the program is in interactive mode.
