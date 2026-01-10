package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunAndExit(bin string, args []string) {
	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		fmt.Fprintln(os.Stderr, "goscript:", err)
		os.Exit(1)
	}
	os.Exit(0)
}
