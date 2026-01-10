package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Resolved struct {
	Key    CacheKey
	Binary string
	WorkDir string
}

func main() {
	config := ParseArgs(os.Args)

	content, err := read(config.ScriptPath)
	if err != nil { fatal(err) }

	resolved, err := resolve(content)
	if err != nil { fatal(err) }

	RunAndExit(resolved.Binary, config.Args)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "goscript:", err)
	os.Exit(1)
}

func read(script string) ([]byte, error) {
	abs, err := filepath.Abs(script)
	if err != nil { return nil, err }

	raw, err := os.ReadFile(abs)
	if err != nil { return nil, err }

	return raw, nil
}

func resolve(raw []byte) (*Resolved, error) {
	key := HashContent(raw)

	r, hit, err := LookupCache(key);
	if err != nil { return nil, err }
	if hit { return r, nil }

	r, err = PrepareScript(key, raw)
	if err != nil { return nil, err }

	return r, nil
}
