package main

import (
	"fmt"
	"os"
	"io"
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
		err = copyBinary(resolved.Binary, cfg.OutputPath);
		if err != nil { fatal(err) }
		fmt.Printf("Compiled into %s\n", cfg.OutputPath)
	} else if (cfg.Action == ActionMigrate) {
		panic("migration not implemented")
	} else if (cfg.Action == ActionRun) {
		RunAndExit(resolved.Binary, cfg.Args)
	} else {
		panic("No action. Shouldn't be possible to reach this state")
	}

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

func copyBinary(src string, dst string) error {
		var perms = os.O_CREATE|os.O_WRONLY|os.O_TRUNC

		in, err := os.Open(src); if err != nil { return err }
		out, err := os.OpenFile( dst, perms, 0o755); if err != nil { return err }
		_, err = io.Copy(out, in); if err != nil { return err }

		return nil
}
