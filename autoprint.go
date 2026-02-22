package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// wrapLastExpr checks the user's code snippet: if the last statement is a bare
// expression (not an assignment, not a loop, not an explicit print call), it
// wraps it in fmt.Println(...) so the result is automatically printed.
//
// This mirrors pyp's "auto-print" behaviour: the last expression is output.
// If the user already calls fmt.Print*, the snippet is left untouched.
func wrapLastExpr(code string) string {
	// We wrap the snippet in a synthetic function to parse it as statements.
	const preamble = "package main\nfunc _(){\n"
	src := preamble + code + "\n}"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		// Unparseable — leave as-is; compile errors will surface later.
		return code
	}

	if len(f.Decls) == 0 {
		return code
	}
	fn, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok || fn.Body == nil || len(fn.Body.List) == 0 {
		return code
	}

	last := fn.Body.List[len(fn.Body.List)-1]
	exprStmt, ok := last.(*ast.ExprStmt)
	if !ok {
		return code
	}

	// Skip if the expression is already a fmt.Print* / println / print call.
	if isPrintCall(exprStmt.X) {
		return code
	}

	// Extract the expression text from the original code using byte offsets.
	// Offset in the parsed file includes the preamble, so subtract it.
	preambleLen := len(preamble)
	start := fset.Position(exprStmt.Pos()).Offset - preambleLen
	end := fset.Position(exprStmt.End()).Offset - preambleLen

	if start < 0 || end > len(code) || start >= end {
		return code
	}

	exprText := strings.TrimSpace(code[start:end])
	return code[:start] + "fmt.Println(" + exprText + ")" + code[end:]
}

// isPrintCall returns true if expr is a call to fmt.Print*, log.Print*, or
// the builtin println/print.
func isPrintCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		// fmt.Println, fmt.Printf, log.Println, etc.
		pkg, ok := fn.X.(*ast.Ident)
		if !ok {
			return false
		}
		sel := fn.Sel.Name
		pkgName := pkg.Name
		if (pkgName == "fmt" || pkgName == "log") && strings.HasPrefix(sel, "Print") {
			return true
		}
		// Also skip fmt.Fprint* (explicit writer)
		if pkgName == "fmt" && strings.HasPrefix(sel, "Fprint") {
			return true
		}
	case *ast.Ident:
		// builtin println / print
		if fn.Name == "println" || fn.Name == "print" {
			return true
		}
	}
	return false
}
