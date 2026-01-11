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
	cfg := ParseArgs(os.Args)

	var content []byte
	var err error

	content, err = read(cfg)
	if err != nil { fatal(err) }

	resolved, err := resolve(content)
	if err != nil { fatal(err) }

	if (cfg.Action == ActionBuild) {
		panic("build output not implemented")
	} else if (cfg.Action == ActionMigrate) {
		panic("migration not implemented")
	}

	RunAndExit(resolved.Binary, cfg.Args)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "goscript:", err)
	os.Exit(1)
}

func read(cfg Config) ([]byte, error) {
	if (cfg.Input == InputInline) { return InlineToScript(cfg.InlineCode) }

	abs, err := filepath.Abs(cfg.ScriptPath)
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
