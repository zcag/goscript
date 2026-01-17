# goscript

Run tiny, single-file Go programs like scripts—no module setup required. `goscript` reads your file (or inline code), resolves missing imports, builds a cached binary, and runs it. External dependencies are fetched automatically.

## Highlights

- **Zero setup:** no `go.mod` needed in your script folder.
- **Works with any Go module:** just import it and go.
- **Fast repeats:** compiled binaries are cached by script content.
- **Inline or file-based:** run a shebang script or one-liners.

---

## Installation

```bash
go install zcag/goscript@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is on your `PATH`.

---

## Quick start

### 1) Run a script file (shebang)

Create a file named `hello.gs`:

```go
#!/usr/bin/env goscript
package main

import (
    "fmt"
    "os"
    "github.com/pkg/errors"
)

func main() {
    fmt.Println("args:", os.Args[1:])
    fmt.Println(errors.New("boom"))
}
```

Run it:

```bash
chmod +x hello.gs
./hello.gs a b c
```

**Output:**

```
args: [a b c]
boom
```

### 2) Run inline code

```bash
goscript -c 'fmt.Println("hello from inline")'
```

**Output:**

```
hello from inline
```

### 3) Build a reusable binary

```bash
goscript -o hello-bin hello.gs
./hello-bin a b c
```

**Output:**

```
args: [a b c]
boom
```

---

## Usage

```text
goscript [flags] <script> [args...]
```

### Flags

- `-c, --code` – Inline Go code.
- `-o, --out` – Build output path (compile without running).
- `-m, --mig` – Migrate target dir (not implemented).

### Notes

- You must provide **exactly one** of:
  - a script path, or
  - inline code via `-c`.
- Scripts must be valid Go with `package main` and `func main()`.

---

## More examples

### Read stdin

```bash
echo "hello" | goscript -c 'data, _ := io.ReadAll(os.Stdin); fmt.Printf("stdin=%s", data)'
```

**Output:**

```
stdin=hello
```

### Use an external module

```bash
goscript -c 'fmt.Println(errors.New("wrapped"))'
```

**Output:**

```
wrapped
```

---

## How it works

1. `goscript` reads your script or inline code.
2. The raw contents are hashed and looked up in the cache.
3. On cache miss, `goscript`:
   - strips a shebang line (if present)
   - adds missing imports
   - generates a temporary `go.mod`
   - downloads dependencies
   - builds a binary into the cache
4. The cached binary is executed (or copied out with `-o`).

Subsequent runs reuse the cached binary, so they’re fast.

---

## Cache layout

```
~/.cache/goscript/
├── bin/
│   └── <hash>/app
└── work/
    └── <hash>
        ├── main.go
        └── go.mod
```

The cache key is based on the raw script contents.

---

## Requirements

- Scripts must use `package main`.
- Scripts must define `func main()`.

This is real Go code—no DSL, no auto-wrapping.

---

## Roadmap

- Add helper utilities for inline mode (args/stdin helpers).
- Support auto-importing common non-stdlib libraries.
- Faster exec path (`syscall.Exec`).
- `goscript --migrate myscript.go myproj/` to convert a script to a project.
