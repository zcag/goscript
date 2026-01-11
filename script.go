package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func HashContent(b []byte) CacheKey {
	sum := sha256.Sum256(b)
	return CacheKey(hex.EncodeToString(sum[:])[:24])
}

func PrepareScript(key CacheKey, raw []byte) (*Resolved, error) {
	goSrc := StripShebang(raw)

	workDir := cacheWorkDir(key)
	err := os.MkdirAll(workDir, 0o755)
	if err != nil { return nil, err }

	// go.mod
	modPath := filepath.Join(workDir, "go.mod")
	mod := []byte(
		"module goscript/" + string(key) + "\n\n" +
		"go " + goVersionLine() + "\n",
	)
	err = os.WriteFile(modPath, mod, 0o644)
	if err != nil { return nil, err }

	// main.go
	goPath := filepath.Join(workDir, "main.go")
	if err := os.WriteFile(goPath, goSrc, 0o644); err != nil {
		return nil, err
	}

	bin := cacheBinPath(key)
	err = os.MkdirAll(filepath.Dir(bin), 0o755)
	if err != nil { return nil, err }

	err = RunQuiet("go", []string{"mod", "tidy"}, workDir);
	if err != nil { return nil, err }

	tmp := bin + ".tmp"
	err = RunQuiet(
		"go",
		[]string{"build", "-trimpath", "-o", tmp, goPath},
		workDir,
	);
	if err != nil { os.Remove(tmp); return nil, err }

	err = os.Rename(tmp, bin);
	if err != nil { os.Remove(tmp); return nil, err }

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

func goVersionLine() string {
	out, _ := exec.Command("go", "env", "GOVERSION").Output()
	// "go1.22.1" → "1.22"
	v := strings.TrimPrefix(strings.TrimSpace(string(out)), "go")
	parts := strings.Split(v, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return v
}
