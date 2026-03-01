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
	t.Run("-j parallel pipe", specParallel)
	t.Run("match() helper extracts group", specMatchHelper)
	t.Run("named() helper extracts named group", specNamedHelper)
	t.Run("-r filters non-matching lines", specRegexFilter)
	t.Run("-r exposes m capture groups", specRegexM)
	t.Run("-r exposes n named groups", specRegexN)
	t.Run("-r sub() replaces first match", specRegexSub)
	t.Run("-r suball() replaces all matches", specRegexSuball)
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
	assertCmdArgs(t, goscriptPath, []string{`fmt.Println("inline-ok")`}, "inline-ok")
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
		[]string{`fmt.Println(strings.ToUpper(x))`},
		"hello\nworld\n",
		`HELLO\nWORLD`,
	)
}

func specPipeModeIdx(t *testing.T) {
	assertCmdStdin(t,
		[]string{`fmt.Printf("%d %s\n", i, x)`},
		"a\nb\nc\n",
		`0 a`,
	)
}

func specPipeModeFields(t *testing.T) {
	assertCmdStdin(t,
		[]string{`fmt.Println(f[1])`},
		"foo bar baz\nalpha beta gamma\n",
		`bar\nbeta`,
	)
}

func specBatchModeLines(t *testing.T) {
	assertCmdStdin(t,
		[]string{`fmt.Println(len(lines))`},
		"a\nb\nc\n",
		`^3\n$`,
	)
}

func specParallel(t *testing.T) {
	// With -j 4, all lines should still be processed (order not guaranteed).
	// We check that all 3 results appear somewhere in the output.
	cmd := exec.Command(goscriptPath, `fmt.Println(strings.ToUpper(x))`, "-j", "4")
	cmd.Stdin = strings.NewReader("hello\nworld\nfoo\n")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("parallel failed: %v\n%s", err, out)
	}
	got := string(out)
	for _, want := range []string{"HELLO", "WORLD", "FOO"} {
		if !strings.Contains(got, want) {
			t.Fatalf("parallel output missing %q:\n%s", want, got)
		}
	}
}

func specFieldSep(t *testing.T) {
	assertCmdStdin(t,
		[]string{`f[1]`, "-F", ","},
		"a,b,c\n1,2,3\n",
		`b\n2`,
	)
}

func specExplain(t *testing.T) {
	cmd := exec.Command(goscriptPath, "--explain", `strings.ToUpper(x)`)
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
		[]string{`strings.ToUpper(x)`},
		"hello\nworld\n",
		`HELLO\nWORLD`,
	)
}

func specAutoPrintBatch(t *testing.T) {
	assertCmdStdin(t,
		[]string{`len(lines)`},
		"a\nb\nc\n",
		`^3\n$`,
	)
}

func specAutoPrintSimple(t *testing.T) {
	assertCmdArgs(t, goscriptPath,
		[]string{`"hello from auto-print"`},
		"hello from auto-print",
	)
}

func specAutoPrintSkipsExplicit(t *testing.T) {
	// Explicit fmt.Println — auto-print must NOT double-wrap.
	assertCmdArgs(t, goscriptPath,
		[]string{`fmt.Println("explicit")`},
		`^explicit\n$`,
	)
}

func specMatchHelper(t *testing.T) {
	// match() without -r; all input lines match so nil-check not needed
	assertCmdStdin(t,
		[]string{`match("(\\d+)", x)[1]`},
		"foo123\nbaz456\n",
		`^123\n456\n$`,
	)
}

func specNamedHelper(t *testing.T) {
	assertCmdStdin(t,
		[]string{`named("(?P<num>\\d+)", x)["num"]`},
		"foo123\nbaz456\n",
		`123\n456`,
	)
}

func specRegexFilter(t *testing.T) {
	// -r should skip lines that don't match
	assertCmdStdin(t,
		[]string{`x`, "-r", `\d+`},
		"foo\nbar123\nbaz\nqux456\n",
		`^bar123\nqux456\n$`,
	)
}

func specRegexM(t *testing.T) {
	// m[1] is first capture group
	assertCmdStdin(t,
		[]string{`m[1]`, "-r", `(\d+)`},
		"foo123\nbar\nbaz456\n",
		`^123\n456\n$`,
	)
}

func specRegexN(t *testing.T) {
	// named capture groups via n
	assertCmdStdin(t,
		[]string{`n["num"]`, "-r", `(?P<num>\d+)`},
		"foo123\nbar\nbaz456\n",
		`^123\n456\n$`,
	)
}

func specRegexSub(t *testing.T) {
	// sub replaces first match using -r pattern
	assertCmdStdin(t,
		[]string{`sub("[$1]", x)`, "-r", `(\d+)`},
		"a1b2\nc3\n",
		`^a\[1\]b2\nc\[3\]\n$`,
	)
}

func specRegexSuball(t *testing.T) {
	// suball replaces all matches using -r pattern
	assertCmdStdin(t,
		[]string{`suball("[$1]", x)`, "-r", `(\d+)`},
		"a1b2\nc3\n",
		`^a\[1\]b\[2\]\nc\[3\]\n$`,
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
