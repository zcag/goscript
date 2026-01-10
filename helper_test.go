package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

var goscriptPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "goscript-test-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	goscriptPath = filepath.Join(tmp, "goscript")

	// Build bin
	if out, err := exec.Command("go", "build", "-o", goscriptPath, ".").CombinedOutput(); err != nil {
		panic(string(out))
	}

	// So shebang env finds the test binary first
	os.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	os.Exit(m.Run())
}

// Creates a tmp executable from given multiline string, asserts no error and stdout matches want (regex)
func assertScript(t *testing.T, src, want string, args ...string) {
	t.Helper()
	assertSpec(t, src, want, false, args...)
}

func assertScriptError(t *testing.T, src, want string) {
	t.Helper()
	assertSpec(t, src, want, true)
}

func assertSpec(t *testing.T, src, want string, wantErr bool, args ...string) {
	t.Helper()

	script := writeSpec(t, src)
	out, err := exec.Command(script, args...).CombinedOutput()

	if wantErr != (err != nil) {
		t.Fatalf("wantErr=%v gotErr=%v out=%s", wantErr, err != nil, out)
	}

	got := string(out)
	if !regexp.MustCompile(want).MatchString(got) {
		t.Fatalf("got %q want /%s/", got, want)
	}
}

func writeSpec(t *testing.T, src string) string {
	t.Helper()

	script := filepath.Join(t.TempDir(), "spec")
	if err := os.WriteFile(script, []byte(src), 0o755); err != nil {
		t.Fatal(err)
	}
	return script
}
