package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: goscript <script.go> [args...]")
		os.Exit(2)
	}

	script := os.Args[1]
	args := os.Args[2:]

	abs, err := filepath.Abs(script)
	if err != nil {
		fatal(err)
	}

	raw, err := os.ReadFile(abs)
	if err != nil {
		fatal(err)
	}

	key := HashContent(raw)

	if r, ok, err := LookupCache(key); err != nil {
		fatal(err)
	} else if ok {
		RunAndExit(r.Binary, args)
	}

	r, err := PrepareScript(key, raw)
	if err != nil {
		fatal(err)
	}

	RunAndExit(r.Binary, args)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "goscript:", err)
	os.Exit(1)
}
