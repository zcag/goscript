# goscript

Run small, standalone Go scripts without creating a module or repository.
Supports importing **any** Go module automatically — no `go.mod` required.

---

## Script mode

Write a single executable file with a shebang and run it like a script.

```go
#!/usr/bin/env goscript
package main

import (
    "fmt"
    "os"
    "github.com/pkg/errors" // external deps installed automatically
)

func main() {
    fmt.Println("args:", os.Args[1:])
    fmt.Println(errors.New("boom"))
}
```

```bash
chmod +x myscript
./myscript a b c

# Compile to a standalone binary
goscript -o mybin myscript
./mybin a b c

# Graduate script to a proper Go module
goscript -m myproject/ myscript
```

Scripts must use `package main` and define `func main()` — it's real Go, not a DSL.

---

## Inline mode (`-c`)

Run Go snippets without a file. Auto-imports are handled; compiled binaries are cached.

```bash
goscript -c 'fmt.Println("hello")'
```

### Auto-print

The last expression in your snippet is automatically printed — no `fmt.Println` needed:

```bash
goscript -c '"hello " + "world"'     # prints: hello world
goscript -c 'os.Getenv("HOME")'      # prints: /home/you
goscript -c 'time.Now().Year()'      # prints: 2025
```

Explicit `fmt.Print*` calls are detected and left untouched (no double-printing).

---

## Pipe mode (automatic)

When your snippet references any of these **magic variables**, pipe mode is activated
automatically — no flags required:

| Variable | Type | Value |
|----------|------|-------|
| `x` / `line` | `string` | current input line |
| `i` / `idx` | `int` | line index (0-based) |
| `f` / `fields` | `[]string` | whitespace-split fields of `x` |
| `lines` | `[]string` | all lines (read before the loop) |

### Per-line processing

Reference `x` or `line` → loops over stdin line by line:

```bash
ls | goscript -c 'strings.ToUpper(x)'
ps aux | goscript -c 'f[0]'                    # first field of each line
cat data.csv | goscript -c 'f[1]' -F ','       # CSV second column
cat file.txt | goscript -c 'fmt.Sprintf("%3d  %s", i, x)'  # add line numbers
```

### Batch processing

Reference `lines` → reads all stdin first, then runs your code once:

```bash
wc -l < file.txt                               # shell way
cat file.txt | goscript -c 'len(lines)'        # Go way, auto-printed

# Sort lines
cat file.txt | goscript -c 'sort.Strings(lines); strings.Join(lines, "\n")'

# Count unique words
cat file.txt | goscript -c '
    counts := map[string]int{}
    for _, l := range lines { counts[l]++ }
    pp(counts)
'
```

### Field separator (`-F`)

`-F sep` changes how `f`/`fields` is split. Default is whitespace (like `strings.Fields`).

```bash
cat data.csv | goscript -c 'f[2]' -F ','       # third CSV column
cat /etc/passwd | goscript -c 'f[0]' -F ':'    # usernames
```

### Parallel processing (`-j N`)

`-j N` runs the loop body across N goroutines (output order not guaranteed):

```bash
cat urls.txt | goscript -c 'fetch(x)' -j 20

# Square numbers in parallel
seq 100 | goscript -c 'atoi(x) * atoi(x)' -j 8
```

---

## Helpers

Every inline/pipe/batch script gets these functions for free:

| Function | Description |
|----------|-------------|
| `die(err)` | Print error and exit 1. No-op if err is nil. |
| `must(v, err)` | Generic unwrap: `n := must(strconv.Atoi(s))` |
| `atoi(s)` | Parse int, exit on failure. |
| `atof(s)` | Parse float64, exit on failure. |
| `pp(v) string` | Format any value as indented JSON. Auto-print friendly. |
| `trim(s)` | `strings.TrimSpace(s)` shortcut. |
| `splitlines(s)` | Split a multi-line string into `[]string`. |

Examples:

```bash
seq 5 | goscript -c 'atoi(x) * atoi(x)'        # squares
echo '{"a":1}' | goscript -c 'pp(map[string]int{"x": 42})'
```

---

## Inspect generated code (`--explain`)

Print the Go source that would be compiled, without actually running it:

```bash
ls | goscript --explain -c 'strings.ToUpper(x)'
goscript --explain -c 'pp(os.Environ())'
```

Useful for understanding what's generated or debugging unexpected output.

---

## Cache control

```bash
goscript --no-cache -c 'fmt.Println("fresh compile")'   # skip cache, always recompile
goscript --clear-cache                                    # wipe entire cache
```

Cache lives at `~/.cache/goscript/` (respects `$XDG_CACHE_HOME`):

```
~/.cache/goscript/
├── bin/<hash>/app      ← compiled binary
└── work/<hash>/
    ├── main.go
    └── go.mod
```

The cache key is the SHA-256 of the generated source, so any change to your snippet
produces a fresh compile.

---

## Installation

```bash
go install github.com/zcag/goscript@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

---

## How it works

1. Source is read (from file or `-c` snippet)
2. For inline mode: magic variables are detected → appropriate template is chosen
3. Auto-print: if the last statement is a bare expression, it's wrapped in `fmt.Println`
4. Source is hashed → cache is checked
5. On miss: imports are resolved (via `goimports`), `go.mod` is generated, binary is compiled
6. Compiled binary is executed (or printed with `--explain`, or saved with `-o`)

Subsequent runs with the same source reuse the cached binary instantly.
