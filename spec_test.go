package main

import "testing"

func TestScripts(t *testing.T) {
	assertScript(t,
	 `#!/bin/env goscript
		package main
		import "fmt"
		func main() {
			fmt.Println("ok")
		}
	 `,
	 "ok",
	)

	assertScript(t,
   `#!/bin/env goscript
		package main
		import (
			"fmt"
			"os"
		)
		func main() {
			fmt.Println(os.Args[1:])
		}
		`,
		`^\[a b\]\n$`,
		"a",
		"b",
	)

	assertScriptError(t,
	 `#!/bin/env goscript
		package main
		func main() { SINTAX }
	 `,
	 "main.go:3.17: undefined: SINTAX",
	)

	assertScript(t,
	 `#!/usr/bin/env goscript
		package main

		import (
			"fmt"
			"github.com/pkg/errors"
		)

		func main() {
			err := errors.New("boomalaka")
			fmt.Println(err.Error())
		}
		`,
		"boomalaka",
	)
}
