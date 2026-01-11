# goscript

Run small, standalone Go scripts without creating a module or repository.
Supports importing **any** go module, no `go.mod` required.

You write a single executable file with a shebang, and run it like a script.

## Usage

### Shebang Example

```go
#!/usr/bin/env goscript
package main

import (
	"fmt" // You don't have to define stdlibs
	"os"
    "github.com/pkg/errors" // These deps are automatically installed
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

You can create a binary from this singular goscript
```bash
goscript -o mybin myscript
./mybin a b c
```

### Inline Example

You can also run goscript directly with `goscript -c` to run go from the param directly.
It automagically handles packages/imports for you.

This feature also benefits from cached binaries, so it actually runs compiled binaries after running it at least once

```bash
goscript -c 'fmt.Println("yeye")'
````

This can be combined with `-o` to compile a binary from your inline go command.
Not sure why you would do this, but it's possible.

```bash
goscript -o mybin -c 'fmt.Println("why")'
./mybin
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
   * shebang line is stripped
   * any missing imports are added to the script
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
* add some helper methods to inline, to ease common pipe/arg handling and printing
* supoort auto importing non standart libraries
    * maybe have a big map of common libs
* faster exec (`syscall.Exec`)
    * support/test passing stdout to ./myscript.go
* `goscript --migrate myscript.go myproj/` to convert script to proper go project
