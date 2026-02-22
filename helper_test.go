package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var goscriptPath string

var testOutBinPath string
var testMigrateDir string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "goscript-test-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	goscriptPath = filepath.Join(tmp, "goscript")

	testOutBinPath = filepath.Join(tmp, "compiled")
	testMigrateDir = filepath.Join(tmp, "migrated")

	// Build bin
	if out, err := exec.Command("go", "build", "-o", goscriptPath, ".").CombinedOutput(); err != nil {
		panic(string(out))
	}

	// So shebang env finds the test binary first
	os.Setenv("PATH", tmp+string(os.PathListSeparator)+os.Getenv("PATH"))

	os.Exit(m.Run())
}

func assertCmd(t *testing.T, script string, want string) {
	t.Helper()
	_assertCmd(t, script, []string{}, want, false)
}

func assertCmdArgs(t *testing.T, script string, args []string, want string) {
	t.Helper()
	_assertCmd(t, script, args, want, false)
}

func assertCmdError(t *testing.T, script string, want string) {
	t.Helper()
	_assertCmd(t, script, []string{}, want, true)
}

func assertCmdErrorArgs(t *testing.T, script string, args []string, want string) {
	t.Helper()
	_assertCmd(t, script, args, want, true)
}

func _assertCmd(t *testing.T, script string, args []string, want string, wantErr bool) {
	t.Helper()

	out, err := exec.Command(script, args...).CombinedOutput()

	if wantErr != (err != nil) {
		t.Fatalf(
			"wantErr=%v gotErr=%v out=%s\ncommand: %s %v",
			wantErr,
			err != nil,
			out,
			script,
			args,
		)
	}

	got := string(out)
	if !regexp.MustCompile(want).MatchString(got) {
		t.Fatalf("got %q want /%s/\ncommand: %s %v", got, want, script, args)
	}
}

// assertCmdStdin runs goscript with given args, feeds stdinData, and matches output.
func assertCmdStdin(t *testing.T, args []string, stdinData string, want string) {
	t.Helper()
	cmd := exec.Command(goscriptPath, args...)
	cmd.Stdin = bytes.NewBufferString(stdinData)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command error: %v\noutput: %s\ncommand: %s %v", err, out, goscriptPath, args)
	}
	got := string(out)
	if !regexp.MustCompile(want).MatchString(got) {
		t.Fatalf("got %q want /%s/\ncommand: %s %v", got, want, goscriptPath, args)
	}
}

func assertDirExists(t *testing.T, dir string) {
	t.Helper()
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected dir %s to exist: %v", dir, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", dir)
	}
}

func assertFileContains(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if !strings.Contains(string(data), want) {
		t.Fatalf("file %s does not contain %q\ncontents:\n%s", path, want, data)
	}
}

func prepScript(t *testing.T, src string) string {
	t.Helper()

	script := filepath.Join(t.TempDir(), "spec")
	if err := os.WriteFile(script, []byte(src), 0o755); err != nil {
		t.Fatal(err)
	}
	return script
}
