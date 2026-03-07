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

type Config struct {
	EnableLowercase bool `mapstructure:"enable_lowercase"`
	EnableEnglish bool `mapstructure:"enable_english"`
	EnableSpecial bool `mapstructure:"enable_special"`
	EnableSensitive bool `mapstructure:"enable_sensitive"`
	SensitiveWords []string `mapstructure:"sensitive_words"`
}

func DefaultConfig() *Config {
	return &Config{
		EnableLowercase: true,
		EnableEnglish:   true,
		EnableSpecial:   true,
		EnableSensitive: true,
		SensitiveWords: []string{
			"password", "pass", "pwd", "secret", "token:", "key", "credential",
		},
	}
}

var Analyzer = NewAnalyzer(DefaultConfig())

func NewAnalyzer(cfg *Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "logcheck",
		Doc:      "reports wrong log messages",
		Run:      runWithConfig(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func runWithConfig(cfg *Config) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		nodeFilter := []ast.Node{
			(*ast.CallExpr)(nil),
		}

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			call := n.(*ast.CallExpr)

			switch fun := call.Fun.(type) {
			case *ast.SelectorExpr:
				if isLoggingFunction(pass, fun) {
					if len(call.Args) > 0 {
						analyzeMessageWithConfig(pass, call, call.Args[0], cfg)
					}
				}
			}
		})

		return nil, nil
	}
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

func analyzeMessageWithConfig(pass *analysis.Pass, call *ast.CallExpr, msgExpr ast.Expr, cfg *Config) {
	msg, ok := extractMessageString(pass, msgExpr)
	if !ok {
		return
	}

	if cfg.EnableLowercase {
		if ok, diag := checkLowercaseStart(msg); ok {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: diag,
			})
		}
	}

	if cfg.EnableEnglish {
		if ok, diag := checkEnglishOnly(msg); ok {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: diag,
			})
		}
	}

	if cfg.EnableSpecial {
		if ok, diag := checkSpecialChars(msg); ok {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: diag,
			})
		}
	}

	if cfg.EnableSensitive {
		if ok, diag := checkSensitiveDataWithConfig(msg, cfg.SensitiveWords); ok {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: diag,
			})
		}
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

func checkSensitiveDataWithConfig(s string, words []string) (bool, string) {
	lower := strings.ToLower(s)
	for _, word := range words {
		if strings.Contains(lower, word) {
			return true, "log message may contain sensitive data (word: " + word + ")"
		}
	}
	return false, ""
}
