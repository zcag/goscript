package main

import (
	"io"
	"os"
	"path/filepath"
)

// migrateScript copies the resolved workDir (go.mod + main.go) to dst,
// producing a standalone Go module the user can develop further.
func migrateScript(r *Resolved, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	for _, name := range []string{"go.mod", "go.sum", "main.go"} {
		src := filepath.Join(r.WorkDir, name)
		// go.sum may not exist for simple scripts
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}
		if err := copyFile(src, filepath.Join(dst, name)); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
