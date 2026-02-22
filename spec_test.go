package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestScripts(t *testing.T) {
	t.Run("basic args", specBasicArgs)
	t.Run("compile error thrown", specCompileError)
	t.Run("dependency loads", specDeps)
	t.Run("auto imports missing libs", specImports)
	t.Run("inline runs", specInline)
	t.Run("outputs binary", specOutput)
	t.Run("migrates script to dir", specMigrate)
	t.Run("pipe mode: x var per line", specPipeModeX)
	t.Run("pipe mode: i index var", specPipeModeIdx)
	t.Run("pipe mode: f fields var", specPipeModeFields)
	t.Run("batch mode: lines var", specBatchModeLines)
	t.Run("auto-print: pipe expr", specAutoPrintPipe)
	t.Run("auto-print: batch expr", specAutoPrintBatch)
	t.Run("auto-print: simple inline", specAutoPrintSimple)
	t.Run("auto-print: skips explicit print", specAutoPrintSkipsExplicit)
	t.Run("--explain prints source", specExplain)
	t.Run("-F field separator", specFieldSep)
}

func specBasicArgs(t *testing.T) {
	var scr = prepScript(t,
   `#!/bin/env goscript
		package main
		import (
			"fmt"
			"os"
		)
		func main() {
			fmt.Println(os.Args[1:])
		}`)

	assertCmdArgs(t, scr, []string{"a", "b"}, `^\[a b\]\n$`)
}

func specCompileError(t *testing.T) {
	var scr = prepScript(t,
	 `#!/bin/env goscript
		package main
		func main() { SINTAX }`)

	assertCmdError(t, scr, "undefined: SINTAX")
}

func specDeps(t *testing.T) {
	var scr = prepScript(t,
	 `#!/usr/bin/env goscript
		package main

		import (
			"fmt"
			"github.com/pkg/errors"
		)

		func main() {
			err := errors.New("boomalaka")
			fmt.Println(err.Error())
		}`)

	assertCmd(t, scr, "boomalaka")
}

func specImports(t *testing.T) {
	var scr = prepScript(t,
	 `#!/usr/bin/env goscript
		package main
		func main() { fmt.Println("auto-import works") }`)

	assertCmd(t, scr, "auto-import works")
}

func specInline(t *testing.T) {
	assertCmdArgs(t, goscriptPath, []string{"-c", `fmt.Println("inline-ok")`}, "inline-ok")
}

func specOutput(t *testing.T) {
	var scr = prepScript(t,
   `#!/bin/env goscript
		package main
		import "fmt"
		func main() { fmt.Println("output-ok") }`)

	assertCmdArgs(
		t,
		goscriptPath,
		[]string{scr, "-o", testOutBinPath},
		"Compiled into",
	)
	assertCmd(t, testOutBinPath, "output-ok")
}

func specPipeModeX(t *testing.T) {
	assertCmdStdin(t,
		[]string{"-c", `fmt.Println(strings.ToUpper(x))`},
		"hello\nworld\n",
		`HELLO\nWORLD`,
	)
}

func specPipeModeIdx(t *testing.T) {
	assertCmdStdin(t,
		[]string{"-c", `fmt.Printf("%d %s\n", i, x)`},
		"a\nb\nc\n",
		`0 a`,
	)
}

func specPipeModeFields(t *testing.T) {
	assertCmdStdin(t,
		[]string{"-c", `fmt.Println(f[1])`},
		"foo bar baz\nalpha beta gamma\n",
		`bar\nbeta`,
	)
}

func specBatchModeLines(t *testing.T) {
	assertCmdStdin(t,
		[]string{"-c", `fmt.Println(len(lines))`},
		"a\nb\nc\n",
		`^3\n$`,
	)
}

func specFieldSep(t *testing.T) {
	assertCmdStdin(t,
		[]string{"-c", `f[1]`, "-F", ","},
		"a,b,c\n1,2,3\n",
		`b\n2`,
	)
}

func specExplain(t *testing.T) {
	cmd := exec.Command(goscriptPath, "--explain", "-c", `strings.ToUpper(x)`)
	cmd.Stdin = strings.NewReader("")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("explain failed: %v\n%s", err, out)
	}
	src := string(out)
	if !strings.Contains(src, "package main") {
		t.Fatalf("explain output missing 'package main':\n%s", src)
	}
	if !strings.Contains(src, "strings.ToUpper(x)") {
		t.Fatalf("explain output missing user code:\n%s", src)
	}
}

func specAutoPrintPipe(t *testing.T) {
	// No fmt.Println — auto-print should kick in.
	assertCmdStdin(t,
		[]string{"-c", `strings.ToUpper(x)`},
		"hello\nworld\n",
		`HELLO\nWORLD`,
	)
}

func specAutoPrintBatch(t *testing.T) {
	assertCmdStdin(t,
		[]string{"-c", `len(lines)`},
		"a\nb\nc\n",
		`^3\n$`,
	)
}

func specAutoPrintSimple(t *testing.T) {
	assertCmdArgs(t, goscriptPath,
		[]string{"-c", `"hello from auto-print"`},
		"hello from auto-print",
	)
}

func specAutoPrintSkipsExplicit(t *testing.T) {
	// Explicit fmt.Println — auto-print must NOT double-wrap.
	assertCmdArgs(t, goscriptPath,
		[]string{"-c", `fmt.Println("explicit")`},
		`^explicit\n$`,
	)
}

func specMigrate(t *testing.T) {
	var scr = prepScript(t,
   `#!/bin/env goscript
		package main
		import (
			"fmt"
			"github.com/pkg/errors"
		)
		func main() {
			err := errors.New("boomalaka")
			fmt.Println(err.Error())
		}`)

	assertCmdArgs(t, goscriptPath, []string{scr, "-m", testMigrateDir}, "Migrated to")

	assertDirExists(t, testMigrateDir)
	assertFileContains(t, testMigrateDir+"/go.mod", "github.com/pkg/errors")
	assertFileContains(t, testMigrateDir+"/main.go", "boomalaka")
}
