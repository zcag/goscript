package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Resolved struct {
	Key     CacheKey
	Binary  string
	WorkDir string
}

func main() {
	cfg := ParseArgs(os.Args)

	if cfg.Action == ActionClearCache {
		clearCache()
		return
	}

	content, _, err := readContent(cfg)
	if err != nil {
		fatal(err)
	}

	if cfg.Action == ActionExplain {
		fmt.Print(string(content))
		return
	}

	resolved, err := resolve(content, cfg.NoCache)
	if err != nil {
		fatal(err)
	}

	switch cfg.Action {
	case ActionBuild:
		if err := copyBinary(resolved.Binary, cfg.OutputPath); err != nil {
			fatal(err)
		}
		fmt.Printf("Compiled into %s\n", cfg.OutputPath)
	case ActionMigrate:
		if err := migrateScript(resolved, cfg.MigrateDir); err != nil {
			fatal(err)
		}
		fmt.Printf("Migrated to %s\n", cfg.MigrateDir)
	case ActionRun:
		RunAndExit(resolved.Binary, cfg.Args)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "goscript:", err)
	os.Exit(1)
}

// readContent reads/generates the Go source for the given config.
// Returns the source bytes and the detected InlineMode (for informational use).
func readContent(cfg Config) ([]byte, InlineMode, error) {
	if cfg.Input == InputInline {
		src, mode, err := InlineToScript(cfg.InlineCode, cfg.FieldSep, cfg.Parallel)
		return src, mode, err
	}

	abs, err := filepath.Abs(cfg.ScriptPath)
	if err != nil {
		return nil, ModeSimple, err
	}
	raw, err := os.ReadFile(abs)
	if err != nil {
		return nil, ModeSimple, err
	}
	return raw, ModeSimple, nil
}

func resolve(raw []byte, noCache bool) (*Resolved, error) {
	key := HashContent(raw)

	if !noCache {
		r, hit, err := LookupCache(key)
		if err != nil {
			return nil, err
		}
		if hit {
			return r, nil
		}
	}

	return PrepareScript(key, raw)
}

func copyBinary(src string, dst string) error {
	var perms = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, perms, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func clearCache() {
	dir := filepath.Join(cacheRoot(), "goscript")
	if err := os.RemoveAll(dir); err != nil {
		fatal(err)
	}
	fmt.Println("Cache cleared.")
}
