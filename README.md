# goscript

Run small, standalone Go scripts without creating a module or repository.

You write a single executable file with a shebang, and run it like a script.

---

## Usage

### Example

```go
#!/usr/bin/env goscript
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("args:", os.Args[1:])
}
```

```bash
chmod +x myscript
./myscript a b c
```

---

## Installation

```bash
go install zcag/goscript@latest
```
Make sure $GOPATH/bin (or $HOME/go/bin) is in your PATH.

---

## How it works

1. Script is executed via shebang (`#!/usr/bin/env goscript`)
2. `goscript` reads the script file and hashes its raw contents
3. Cache is checked under `~/.cache/goscript`
4. On cache miss:

   * shebang line is stripped (line count preserved)
   * script is written as `main.go` into a cache work directory
   * a binary is built and stored in the cache
5. The cached binary is executed

Subsequent runs reuse the cached binary.

---

## Requirements

* Scripts must use `package main`
* Scripts must define `func main()`

This is real Go code — no DSL, no auto-wrapping.

---

## Cache layout

```
~/.cache/goscript/
├── bin/
│   └── <hash>/app
└── work/
    └── <hash>/main.go
```

Cache key is currently based on the raw script contents.

---

## TODO
* Dependenct handling, so you can also include external dependencies
* include `go version` / `GOOS` / `GOARCH` in cache key
* support `goscript -c "fmt.Println('yeye')"`
* faster exec (`syscall.Exec`)


## Todo
* support importing packages
* incorporate go version into cache
