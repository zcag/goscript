package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"

	"golang.org/x/tools/imports"
)

const inlineMarker = "// {{INLINE_BODY}}"

//go:embed script.go.template
var templateSrc []byte

func InlineToScript(code string) ([]byte, error) {
	body := indentAsBlock(strings.TrimSpace(code), "\t")

	out := bytes.Replace(
		templateSrc,
		[]byte(inlineMarker),
		[]byte(body),
		1,
	)

	out, err := HandleImports(out)
	return out, err
}

func HandleImports(src []byte) ([]byte, error) {
	out, err := imports.Process("main.go", src, nil)
	if err != nil { return nil, fmt.Errorf("goimports: %w", err) }
	return out, nil
}

func indentAsBlock(s, prefix string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			lines[i] = ""
		} else {
			lines[i] = prefix + ln
		}
	}
	return strings.Join(lines, "\n")
}
