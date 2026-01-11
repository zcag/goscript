package main

import "testing"

func TestScripts(t *testing.T) {
	t.Run("basic args", specBasicArgs)
	t.Run("compile error thrown", specCompileError)
	t.Run("dependency loads", specDeps)
	t.Run("inline runs", specInline)
	t.Run("outputs binary", specOutput)
	t.Run("migrates script to dir", specOutput)
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

	assertCmdError(t, scr, "main.go:3.17: undefined: SINTAX")
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

func specInline(t *testing.T) {
	assertCmdErrorArgs(t, goscriptPath, []string{"-c", `fmt.Println("inline-ok")`}, "not implemented")
	// assertCmdArgs(t, goscriptPath, []string{"-c", `fmt.Println("inline-ok")`}, "inline-ok")
}

func specOutput(t *testing.T) {
	var scr = prepScript(t,
   `#!/bin/env goscript
		package main
		import "fmt"
		func main() { fmt.Println("output-ok") }`)

	assertCmdErrorArgs(t, goscriptPath, []string{scr, "-o", testOutBinPath}, "not implemented")

	// assertCmdArgs(t, goscriptPath, []string{scr, "-o", testOutBinPath}, "Compiled into")
	// assertCmd(t, testOutBinPath, "output-ok")
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

	assertCmdErrorArgs(t, goscriptPath, []string{scr, "-m", testMigrateDir}, "not implemented")

	// assertDirExists(t, testMigrateDir)
	// assertFileContent(t, testMigrateDir + "/go.mod", "github.com/package/errors")
	// assertFileContent(t, testMigrateDir + "/main.go", "boomalaka")
}
