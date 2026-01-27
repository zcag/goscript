# 🚀 goscript

> Run standalone Go scripts instantly—no modules, no setup, just code.

**goscript** lets you write and execute Go code as easily as Python or Bash scripts. Import **any** Go module without creating a `go.mod` file. Perfect for quick automation, one-off tasks, and portable scripts.

```bash
#!/usr/bin/env goscript
package main

import "github.com/pkg/errors"  // Auto-installed!

func main() {
    fmt.Println("Hello, world!")  // Auto-imported!
}
```

---

## ✨ Features

- 🔥 **Zero Configuration** — No `go.mod` required
- 📦 **Auto-Import** — Missing stdlib packages automatically added
- ⚡ **Fast Execution** — Built binaries are cached for instant re-runs
- 🎯 **Shebang Support** — Run scripts directly with `./myscript.go`
- 💻 **Inline Mode** — Execute Go code directly from command line with `-c`
- 🔨 **Binary Compilation** — Create standalone executables with `-o`
- 🌐 **Any Module** — Import third-party packages automatically

---

## 📦 Installation

```bash
go install github.com/zcag/goscript@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your PATH:

```bash
export PATH="$HOME/go/bin:$PATH"
```

---

## 🎯 Quick Start

### Basic Script Execution

Create a file `hello.go`:

```go
#!/usr/bin/env goscript
package main

func main() {
    fmt.Println("Hello from goscript!")
}
```

Make it executable and run:

```bash
chmod +x hello.go
./hello.go
```

**Output:**
```
Hello from goscript!
```

✨ Notice how `fmt` was **automatically imported**—no import statement needed!

---

## 📚 Usage Examples

### Example 1: Command-Line Arguments

**Script:** `args.go`
```go
#!/usr/bin/env goscript
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ./args.go <name> [age]")
        return
    }

    name := os.Args[1]
    fmt.Printf("Hello, %s!\n", name)

    if len(os.Args) > 2 {
        fmt.Printf("Age: %s years old\n", os.Args[2])
    }
}
```

**Usage:**
```bash
./args.go Alice
```

**Output:**
```
Hello, Alice!
```

```bash
./args.go Bob 25
```

**Output:**
```
Hello, Bob!
Age: 25 years old
```

---

### Example 2: External Dependencies

**Script:** `fetch.go`
```go
#!/usr/bin/env goscript
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, err := http.Get("https://api.github.com/zen")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println("GitHub Zen:", string(body))
}
```

**Usage:**
```bash
./fetch.go
```

**Output:**
```
GitHub Zen: Design for failure.
```

The first run will download dependencies and build the binary. Subsequent runs are instant!

---

### Example 3: Third-Party Packages

**Script:** `errors.go`
```go
#!/usr/bin/env goscript
package main

import (
    "fmt"
    "github.com/pkg/errors"
)

func riskyOperation() error {
    return errors.New("something went wrong")
}

func main() {
    err := riskyOperation()
    if err != nil {
        wrapped := errors.Wrap(err, "failed to execute operation")
        fmt.Println(wrapped)
    }
}
```

**Usage:**
```bash
./errors.go
```

**Output:**
```
failed to execute operation: something went wrong
```

goscript automatically:
1. Detects `github.com/pkg/errors` is needed
2. Creates a temporary `go.mod`
3. Downloads the dependency
4. Builds and caches the binary

---

### Example 4: JSON Processing

**Script:** `json-parse.go`
```go
#!/usr/bin/env goscript
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    data := `{"name":"Alice","age":30}`

    var person Person
    json.Unmarshal([]byte(data), &person)

    fmt.Printf("%s is %d years old\n", person.Name, person.Age)
}
```

**Usage:**
```bash
./json-parse.go
```

**Output:**
```
Alice is 30 years old
```

---

### Example 5: File Operations

**Script:** `file-read.go`
```go
#!/usr/bin/env goscript
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ./file-read.go <filename>")
        return
    }

    file, err := os.Open(os.Args[1])
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 1
    for scanner.Scan() {
        fmt.Printf("%3d: %s\n", lineNum, scanner.Text())
        lineNum++
    }
}
```

**Usage:**
```bash
./file-read.go README.md | head -5
```

**Output:**
```
  1: # 🚀 goscript
  2:
  3: > Run standalone Go scripts instantly—no modules, no setup, just code.
  4:
  5: **goscript** lets you write and execute Go code as easily as Python or Bash scripts.
```

---

## 💻 Inline Mode

Execute Go code directly from the command line with `-c`:

### Simple Print

```bash
goscript -c 'fmt.Println("Hello from inline!")'
```

**Output:**
```
Hello from inline!
```

### Using Variables

```bash
goscript -c 'x := 42; fmt.Printf("The answer is %d\n", x)'
```

**Output:**
```
The answer is 42
```

### Math Operations

```bash
goscript -c 'import "math"; fmt.Println(math.Sqrt(16))'
```

**Output:**
```
4
```

### Current Time

```bash
goscript -c 'import "time"; fmt.Println(time.Now().Format("2006-01-02 15:04:05"))'
```

**Output:**
```
2026-01-16 14:30:22
```

### Helper Functions

Inline mode includes a `die()` helper function for error handling:

```bash
goscript -c 'data, err := os.ReadFile("nonexistent.txt"); die(err); fmt.Println(string(data))'
```

**Output:**
```
error: open nonexistent.txt: no such file or directory
```

---

## 🔨 Building Binaries

Compile your script into a standalone executable:

### From Script File

```bash
goscript -o myapp hello.go
./myapp
```

**Output:**
```
Compiled into myapp
Hello from goscript!
```

The resulting binary is completely standalone and can be distributed without goscript installed.

### From Inline Code

```bash
goscript -o greet -c 'fmt.Println("Greetings, human!")'
./greet
```

**Output:**
```
Compiled into greet
Greetings, human!
```

### Real-World Example: HTTP Server

**Script:** `server.go`
```go
#!/usr/bin/env goscript
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server running on :8080")
    http.ListenAndServe(":8080", nil)
}
```

**Build and run:**
```bash
goscript -o webserver server.go
./webserver
```

**Output:**
```
Compiled into webserver
Server running on :8080
```

Now visit `http://localhost:8080/World` and see "Hello, World!"

---

## ⚡ Performance & Caching

goscript caches compiled binaries for lightning-fast re-execution:

**First run:**
```bash
time ./hello.go
Hello from goscript!

real    0m2.341s  # Includes build time
```

**Second run:**
```bash
time ./hello.go
Hello from goscript!

real    0m0.012s  # Instant! Uses cached binary
```

### Cache Structure

```
~/.cache/goscript/
├── bin/
│   └── <hash>/app          # Compiled binaries
└── work/
    └── <hash>/
        ├── main.go         # Generated script
        └── go.mod          # Generated dependencies
```

The cache key is based on the script's content. Modifying the script invalidates the cache and triggers a rebuild.

### Clear Cache

```bash
rm -rf ~/.cache/goscript
```

---

## 🎓 How It Works

1. **Script Execution**
   - Script is executed via shebang (`#!/usr/bin/env goscript`)

2. **Content Hashing**
   - Script content is hashed to create a unique cache key

3. **Cache Lookup**
   - Check if a compiled binary exists in `~/.cache/goscript/bin/<hash>/`

4. **On Cache Miss:**
   - Strip shebang line
   - Auto-import missing packages (like `fmt`)
   - Generate `main.go` and `go.mod` in work directory
   - Run `go build` to create binary
   - Store binary in cache

5. **Execute**
   - Run the cached binary with provided arguments

---

## 📋 Requirements

- **Go 1.16+** installed
- Scripts must use `package main`
- Scripts must define `func main()`

This is **real Go code**—no DSL, no magic, no auto-wrapping into functions. Everything is standard Go.

---

## 🔧 Advanced Usage

### Combining Options

Build a binary from a script and run it:
```bash
goscript -o mybin myscript.go && ./mybin arg1 arg2
```

### Debugging

To see the generated code, check the cache work directory:

```bash
# Run your script
./myscript.go

# Find the work directory
ls -la ~/.cache/goscript/work/

# View generated files
cat ~/.cache/goscript/work/<hash>/main.go
cat ~/.cache/goscript/work/<hash>/go.mod
```

### Pipeline Usage

goscript works great in pipelines:

```bash
cat data.txt | goscript -c 'scanner := bufio.NewScanner(os.Stdin); for scanner.Scan() { fmt.Println("Line:", scanner.Text()) }'
```

---

## 🆚 Comparison

### vs. `go run`

| Feature | goscript | go run |
|---------|----------|--------|
| Requires go.mod | ❌ No | ✅ Yes |
| Auto-imports stdlib | ✅ Yes | ❌ No |
| Shebang support | ✅ Yes | ❌ No |
| Inline execution | ✅ Yes | ❌ No |
| Caches binaries | ✅ Yes | ⚠️ Partial |
| Auto-install deps | ✅ Yes | ❌ No |

### vs. Python/Bash

- ✅ **Compiled performance** — 10-100x faster than interpreted languages
- ✅ **Type safety** — Catch errors at compile time
- ✅ **Rich stdlib** — net/http, encoding/json, crypto, and more built-in
- ✅ **Concurrency** — Goroutines make parallel tasks trivial
- ✅ **Cross-platform** — Single binary runs anywhere

---

## 🎯 Use Cases

### Quick Automation Scripts
```bash
#!/usr/bin/env goscript
package main

import "os/exec"

func main() {
    cmd := exec.Command("git", "status")
    cmd.Stdout = os.Stdout
    cmd.Run()
}
```

### API Testing
```bash
goscript -c 'resp, _ := http.Get("https://httpbin.org/json"); body, _ := io.ReadAll(resp.Body); fmt.Println(string(body))'
```

### System Administration
```go
#!/usr/bin/env goscript
package main

import "os/exec"

func main() {
    // Check disk usage
    cmd := exec.Command("df", "-h")
    cmd.Stdout = os.Stdout
    cmd.Run()
}
```

### Data Processing
```bash
goscript -c 'data := []int{1,2,3,4,5}; sum := 0; for _, n := range data { sum += n }; fmt.Println(sum)'
```

---

## 🤝 Contributing

Contributions are welcome! Check out the [issues](https://github.com/zcag/goscript/issues) or submit a PR.

### Development

```bash
# Clone the repo
git clone https://github.com/zcag/goscript.git
cd goscript

# Build
go build

# Run tests
go test ./...
```

---

## 📝 Roadmap

- [ ] Helper methods for inline mode (pipe/arg handling)
- [ ] Auto-import for common non-stdlib packages
- [ ] Faster execution via `syscall.Exec`
- [ ] `goscript --migrate` to convert script to full Go project
- [ ] Environment variable support for configuration
- [ ] Better error messages and diagnostics

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## ⭐ Star History

If you find goscript useful, give it a star on [GitHub](https://github.com/zcag/goscript)!

---

**Made with ❤️ for Gophers who want to script**

*Happy scripting! 🚀*
