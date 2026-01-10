# goscript

Run small, standalone Go scripts without creating a module or repository.
Supports importing **any** go module, no `go.mod` required.

You write a single executable file with a shebang, and run it like a script.

## Usage

### Example

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

```bash
chmod +x myscript
./myscript a b c
```

## Installation

```bash
go install zcag/goscript@latest
```
Make sure $GOPATH/bin (or $HOME/go/bin) is in your PATH.

## How it works

1. Script is executed via shebang (`#!/usr/bin/env goscript`)
2. `goscript` reads the script file and hashes its raw contents
3. Cache is checked under `~/.cache/goscript`
4. On cache miss:
   * shebang line is stripped (line count preserved)
   * script is written as `main.go` into a cache work directory
   * go.mod is generated with all external dependencies
   * Dependencies are installed
   * a binary is built and stored in the cache
5. The cached binary is executed

Subsequent runs reuse the cached binary.

## Requirements

* Scripts must use `package main`
* Scripts must define `func main()`

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

Cache key is currently based on the raw script contents.

## TODO
* find some solution to lsp and syntax highlighters freaking out from shebang
* support `goscript -c "fmt.Println('yeye')"`
* support `goscript --build .local/bin/mybin myscript.go`
* `goscript --migrate myscript.go myproj/` to convert script to proper go project
* cache clean command
* include `go version` / `GOOS` / `GOARCH` in cache key
* faster exec (`syscall.Exec`)
    * support/test passing stdout to ./myscript.go
