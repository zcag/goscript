package main

import (
	"os"
	"path/filepath"
)

type CacheKey string

type Resolved struct {
	Key    CacheKey
	Binary string
	WorkDir string
}

func LookupCache(key CacheKey) (*Resolved, bool, error) {
	bin := cacheBinPath(key)
	if _, err := os.Stat(bin); err == nil {
		return &Resolved{Key: key, Binary: bin, WorkDir: cacheWorkDir(key)}, true, nil
	} else if os.IsNotExist(err) {
		return nil, false, nil
	} else {
		return nil, false, err
	}
}

func cacheRoot() string {
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".cache")
}

func cacheWorkDir(key CacheKey) string {
	return filepath.Join(cacheRoot(), "goscript", "work", string(key))
}

func cacheBinPath(key CacheKey) string {
	return filepath.Join(cacheRoot(), "goscript", "bin", string(key), "app")
}
