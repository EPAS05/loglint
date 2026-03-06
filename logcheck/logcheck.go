package logcheck

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "logcheck",
	Doc:      "reports wrong log messages",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		switch fun := call.Fun.(type) {
		case *ast.SelectorExpr:
			if isLoggingFunction(pass, fun) {
				//pass.Reportf(call.Pos(), "found logging call")
				if len(call.Args) > 0 {
					analyzeMessage(pass, call, call.Args[0])
				}
			}
		}
	})

	return nil, nil
}

func isLoggingFunction(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}
	pkgPath := pkg.Path()
	switch pkgPath {
	case "log/slog", "go.uber.org/zap":
		return true
	}

	return false
}

func analyzeMessage(pass *analysis.Pass, call *ast.CallExpr, msgExpr ast.Expr) {
	msg, ok := extractMessageString(pass, msgExpr)
	if !ok {
		return
	}

	//fmt.Println(msg)

	if ok, diag := checkLowercaseStart(msg); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: diag,
		})
	}

	if ok, diag := checkEnglishOnly(msg); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: diag,
		})
	}

	if ok, diag := checkSpecialChars(msg); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: diag,
		})
	}

	if ok, diag := checkSensitiveData(msg); ok {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: diag,
		})
	}
}

func extractMessageString(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s := e.Value
			if len(s) >= 2 {
				s = s[1 : len(s)-1]
			}
			return s, true
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left, ok1 := extractMessageString(pass, e.X)
			right, ok2 := extractMessageString(pass, e.Y)
			if ok1 && ok2 {
				return left + right, true
			}
		}
	case *ast.Ident:
		obj := pass.TypesInfo.ObjectOf(e)
		if obj == nil {
			return "", false
		}
		con, ok := obj.(*types.Const)
		if !ok {
			return "", false
		}
		val := con.Val()
		if val == nil {
			return "", false
		}
		return constant.StringVal(val), true
	}
	return "", false
}

func checkLowercaseStart(s string) (bool, string) {
	trimmed := strings.TrimLeft(s, " \t")
	if trimmed == "" {
		return true, "empty log message"
	}
	for _, r := range trimmed {
		if unicode.IsLetter(r) {
			if !unicode.IsLower(r) {
				return true, "log message should start with a lowercase letter"
			}
			return false, ""
		}
	}
	return false, ""
}

func checkEnglishOnly(s string) (bool, string) {
	for _, r := range s {
		if r > 127 {
			return true, "log message should contain only English characters"
		}
	}
	return false, ""
}

func checkSpecialChars(s string) (bool, string) {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			continue
		}
		if unicode.IsSymbol(r) || unicode.IsPunct(r) {
			return true, "log message should not contain special characters or emojis"
		}
	}
	return false, ""
}

func checkSensitiveData(s string) (bool, string) {
	lower := strings.ToLower(s)
	sensitiveWords := []string{
		"password", "pass", "pwd", "secret", "token:", "key", "credential",
	}
	for _, word := range sensitiveWords {
		if strings.Contains(lower, word) {
			return true, "log message may contain sensitive data (word: " + word + ")"
		}
	}
	return false, ""
}
