package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
)

func HashContent(b []byte) CacheKey {
	sum := sha256.Sum256(b)
	return CacheKey(hex.EncodeToString(sum[:])[:24])
}

func PrepareScript(key CacheKey, raw []byte) (*Resolved, error) {
	goSrc := StripShebang(raw)

	workDir := cacheWorkDir(key)
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return nil, err
	}

	goPath := filepath.Join(workDir, "main.go")
	if err := os.WriteFile(goPath, goSrc, 0o644); err != nil {
		return nil, err
	}

	bin := cacheBinPath(key)
	if err := os.MkdirAll(filepath.Dir(bin), 0o755); err != nil {
		return nil, err
	}

	// If already built, reuse.
	if _, err := os.Stat(bin); err == nil {
		return &Resolved{Key: key, Binary: bin, WorkDir: workDir}, nil
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	tmp := bin + ".tmp"
	cmd := exec.Command("go", "build", "-trimpath", "-o", tmp, goPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmp)
		return nil, err
	}
	if err := os.Rename(tmp, bin); err != nil {
		_ = os.Remove(tmp)
		return nil, err
	}

	return &Resolved{Key: key, Binary: bin, WorkDir: workDir}, nil
}

func StripShebang(b []byte) []byte {
	if len(b) >= 2 && b[0] == '#' && b[1] == '!' {
		for i := 0; i < len(b); i++ {
			if b[i] == '\n' {
				return b[i:]
			}
		}
		return nil
	}
	return b
}
