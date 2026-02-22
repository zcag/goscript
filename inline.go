package main

import (
	"bytes"
	_ "embed"
	"regexp"
	"strings"
)

// Template markers
const (
	markerInlineBody = "// {{INLINE_BODY}}"
	markerPipeBefore = "// {{PIPE_BEFORE}}"
	markerPipeBody   = "// {{PIPE_BODY}}"
	markerPipeAfter  = "// {{PIPE_AFTER}}"
	markerBatchBody  = "// {{BATCH_BODY}}"
	markerFields     = "{{FIELDS_EXPR}}"
)

//go:embed script.go.template
var simpleTmpl []byte

//go:embed pipe.go.template
var pipeTmpl []byte

//go:embed batch.go.template
var batchTmpl []byte

// InlineMode describes which template to use.
type InlineMode int

const (
	ModeSimple InlineMode = iota // no magic vars
	ModePipe                     // x/line/i/f referenced → per-line loop
	ModeBatch                    // only lines referenced → pre-load all, no loop
)

// pipeVars are magic variables that trigger per-line loop mode.
var pipeVarNames = []string{"x", "line", "i", "idx", "index", "f", "fields"}

// batchVarNames trigger batch mode (all lines at once).
var batchVarNames = []string{"lines"}

// DetectInlineMode inspects the code snippet for magic variable references.
func DetectInlineMode(code string) InlineMode {
	for _, v := range pipeVarNames {
		if containsIdent(code, v) {
			return ModePipe
		}
	}
	for _, v := range batchVarNames {
		if containsIdent(code, v) {
			return ModeBatch
		}
	}
	return ModeSimple
}

// containsIdent checks whether s contains identifier name as a whole word.
func containsIdent(s, name string) bool {
	pattern := `\b` + regexp.QuoteMeta(name) + `\b`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// InlineToScript converts a code snippet to a complete Go source file.
// It auto-detects the appropriate template based on magic variable usage.
// parallel > 0 activates concurrent loop mode (pipe only).
func InlineToScript(code string, fieldSep string, parallel int) ([]byte, InlineMode, error) {
	mode := DetectInlineMode(code)
	code = wrapLastExpr(code)
	body := indentAsBlock(strings.TrimSpace(code), "\t")

	switch mode {
	case ModePipe:
		return buildPipeScript(body, fieldSep, parallel), mode, nil
	case ModeBatch:
		return buildBatchScript(body), mode, nil
	default:
		return buildSimpleScript(body), mode, nil
	}
}

func buildSimpleScript(body string) []byte {
	return bytes.Replace(simpleTmpl, []byte(markerInlineBody), []byte(body), 1)
}

func buildBatchScript(body string) []byte {
	return bytes.Replace(batchTmpl, []byte(markerBatchBody), []byte(body), 1)
}

func buildPipeScript(body string, fieldSep string, parallel int) []byte {
	_ = parallel // parallel handled separately via parallel template
	var fieldsExpr string
	if fieldSep == "" {
		fieldsExpr = "strings.Fields(x)"
	} else {
		fieldsExpr = `strings.Split(x, "` + fieldSep + `")`
	}

	out := pipeTmpl
	out = bytes.Replace(out, []byte(markerFields), []byte(fieldsExpr), 1)
	out = bytes.Replace(out, []byte(markerPipeBefore), []byte(""), 1)
	out = bytes.Replace(out, []byte(markerPipeBody), []byte(body), 1)
	out = bytes.Replace(out, []byte(markerPipeAfter), []byte(""), 1)
	return out
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
