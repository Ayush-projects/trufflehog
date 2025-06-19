package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	for _, arg := range os.Args[1:] {
		src, err := os.ReadFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", arg, err)
			continue
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, arg, src, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error %s: %v\n", arg, err)
			continue
		}
		var regexes []string
		var description string
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "MustCompile" && len(call.Args) > 0 {
						arg0 := call.Args[0]
						switch a := arg0.(type) {
						case *ast.BasicLit:
							regexes = append(regexes, strings.Trim(a.Value, "`\""))
						case *ast.BinaryExpr:
							if lit, ok := a.Y.(*ast.BasicLit); ok {
								regexes = append(regexes, strings.Trim(lit.Value, "`\""))
							}
						}
					}
				}
			}
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == "Description" && fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
					if len(fn.Body.List) > 0 {
						if ret, ok := fn.Body.List[0].(*ast.ReturnStmt); ok && len(ret.Results) > 0 {
							if lit, ok := ret.Results[0].(*ast.BasicLit); ok {
								description = strings.Trim(lit.Value, "`\"")
							}
						}
					}
				}
			}
			return true
		})
		if len(regexes) > 0 {
			fmt.Printf("%s|%s|%s\n", arg, description, strings.Join(regexes, ","))
		}
	}
}
