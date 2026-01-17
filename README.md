# goscript

Run small, standalone Go scripts without scaffolding a module or repository. `goscript` lets you write a single executable file (with a shebang) and run it like a script while still using real Go and external dependencies.

## Why goscript

- **No `go.mod` required:** import any module and let goscript manage dependencies.
- **Fast repeat runs:** scripts are compiled and cached automatically.
- **Real Go:** no DSL, no wrappers, no magic runtime.

## Installation

```bash
go install zcag/goscript@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

## Quick start: shebang script

Create a file called `hello.gos`:

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

Make it executable and run it:

```bash
chmod +x hello.gos
./hello.gos a b c
```

**Output**

```
args: [a b c]
boom
```

### Build a binary from your script

```bash
goscript -o hello-bin hello.gos
./hello-bin 1 2 3
```

**Output**

```
args: [1 2 3]
boom
```

## Inline mode

Use `goscript -c` to run small snippets directly from the command line. Imports are resolved automatically and the compiled binary is cached for future runs.

```bash
goscript -c 'fmt.Println("inline ok")'
```

**Output**

```
inline ok
```

You can also compile inline snippets into a binary with `-o`:

```bash
goscript -o inline-bin -c 'fmt.Println("why not")'
./inline-bin
```

**Output**

```
why not
```

## How caching works

On the first run, goscript compiles your script and stores the result in a cache. Subsequent runs reuse that binary.

**Steps**

1. The script is executed via shebang (`#!/usr/bin/env goscript`).
2. goscript reads the script file and hashes its raw contents.
3. The cache is checked under `~/.cache/goscript`.
4. On a cache miss:
   - The shebang line is stripped.
   - Missing imports are added.
   - A `main.go` file is generated in a cache work directory.
   - A `go.mod` file is created with external dependencies.
   - Dependencies are installed.
   - A binary is built and stored in the cache.
5. The cached binary is executed.

## Requirements

- Scripts must use `package main`.
- Scripts must define `func main()`.

This is real Go code — no DSL, no auto-wrapping.

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

## Tips

- **Pin dependencies** by using a versioned import path (e.g. `github.com/foo/bar@v1.2.3`).
- **Keep scripts executable** with `chmod +x yourscript` so you can run them directly.
- **Use `-o`** when you want a standalone binary for distribution.

## Roadmap (ideas)

- Helper methods for common arg/pipe handling.
- Faster exec (`syscall.Exec`).
- `goscript --migrate myscript.go myproj/` to convert a script into a proper Go project.
