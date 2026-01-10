package main

import "testing"

func TestScripts(t *testing.T) {
	t.Run("basic args", specBasicArgs)
	t.Run("compile error thrown", specCompileError)
	t.Run("dependency loads", specDeps)
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
