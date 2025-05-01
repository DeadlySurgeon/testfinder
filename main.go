package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	emitSubTests  = flag.Bool("subtests", true, "Emit sub tests")
	emitFullTests = flag.Bool("fulltests", true, "Emit full tests")
	ignoreRunLit  = flag.Bool("ignorelit", false, "Ignores runs by literals")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "You must specify a path")
		os.Exit(1)
	}

	if err := run(flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(path string) (err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		return err
	}

	dir := filepath.Dir(path)
	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedFiles,
		Dir:   dir,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		fileIndex := -1
		for i, goFile := range pkg.GoFiles {
			if len(goFile) >= len(path) && goFile[len(goFile)-len(path):] == path {
				fileIndex = i
				break
			}
		}
		if fileIndex == -1 {
			continue
		}

		out, err := locations(pkg.Syntax[fileIndex], pkg)
		if err != nil {
			return err
		}

		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		return e.Encode(out)
	}

	return errors.New("unable to load test file")
}

type testInfo struct {
	TestFunc string `json:"testFunc"`
	SubTest  string `json:"subTest"`
	Line     int    `json:"line"`
}

func locations(file *ast.File, pkg *packages.Package) ([]testInfo, error) {
	fSet := pkg.Fset
	out := make([]testInfo, 0)
	var err error
	posToNode := map[token.Pos]ast.Node{}
	ast.Inspect(file, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		if n == nil {
			return true
		}
		posToNode[n.Pos()] = n

		// If we don't have a function, or a body, or the function name has the
		// Test prefix, skip it.
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "Test") {
			return true
		}
		testFunc := fn.Name.Name

		// Not a test function.
		if len(fn.Type.Params.List) != 1 || len(fn.Type.Params.List[0].Names) != 1 {
			return true
		}

		// Check that the parameter is a `t *testing.T`
		if !checkT(fn.Type.Params.List[0].Names[0]) {
			return true
		}

		if *emitFullTests {
			out = append(out, testInfo{
				TestFunc: testFunc,
				SubTest:  "",
				Line:     fSet.Position(fn.Pos()).Line,
			})
		}

		// Skip sub tests hunting if we do not need to emit it.
		if !*emitSubTests {
			return true
		}

		// Look inside the test function body
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			posToNode[n.Pos()] = n
			if err != nil {
				return true
			}
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) < 1 {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Run" {
				return true
			}

			// Check that the receiver is *testing.T
			if !checkT(sel.X) {
				return true
			}

			if len(call.Args) == 0 {
				return true
			}

			switch v := call.Args[0].(type) {
			case *ast.BasicLit:
				if !*ignoreRunLit {
					out = append(out, testInfo{
						TestFunc: testFunc,
						SubTest:  blStr(v),
						Line:     pkg.Fset.Position(sel.Pos()).Line,
					})
				}
			case *ast.Ident:
				// This is where we picked out a map value.
				obj := pkg.TypesInfo.Uses[v]
				if obj == nil {
					return true
				}

				noRN, ok := posToNode[obj.Parent().Pos()].(*ast.RangeStmt)
				if !ok {
					return false
				}

				out = processRangeX(out, call, testFunc, noRN.X, fSet)
			case *ast.SelectorExpr:
				// Should be t.(<-name) to get t?
				recvIdent, ok := v.X.(*ast.Ident)
				if !ok {
					return false
				}

				obj := pkg.TypesInfo.Uses[recvIdent]
				if obj == nil {
					return true
				}
				noRN, ok := posToNode[obj.Parent().Pos()].(*ast.RangeStmt)
				if !ok {
					return false
				}
				out = processRangeX(out, call, testFunc, noRN.X, fSet)
			default:
			}

			return false
		})
		return false
	})

	return out, nil
}

func processRangeX(out []testInfo, call *ast.CallExpr, testFunc string, f ast.Node, fset *token.FileSet) []testInfo {
	switch xv := f.(type) {
	case *ast.Ident:
		stmt, ok := xv.Obj.Decl.(*ast.AssignStmt)
		if !ok {
			return out
		}
		if len(stmt.Rhs) == 0 {
			return out
		}
		return processRangeX(out, call, testFunc, stmt.Rhs[0], fset)
	case *ast.CompositeLit:
		switch xv.Type.(type) {
		case *ast.ArrayType:
			for _, elt := range xv.Elts {
				switch v := elt.(type) {
				case *ast.CompositeLit:
					for _, field := range v.Elts {
						kv, ok := field.(*ast.KeyValueExpr)
						if !ok {
							continue
						}
						keyLit, valueLit := identStr(kv.Key), blStr(kv.Value)
						if keyLit == "" || valueLit == "" {
							continue
						}

						// Whatever man.
						tRunSelName := call.Args[0].(*ast.SelectorExpr).Sel.Name

						if keyLit != tRunSelName {
							continue
						}

						out = append(out, testInfo{
							TestFunc: testFunc,
							SubTest:  valueLit,
							Line:     fset.Position(kv.Key.Pos()).Line,
						})
					}
				}
			}

		case *ast.MapType:
			for _, elt := range xv.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if key := blStr(kv.Key); key != "" {
						out = append(out, testInfo{
							TestFunc: testFunc,
							SubTest:  key,
							Line:     fset.Position(kv.Key.Pos()).Line,
						})
					}
				}
			}
		}
	}
	return out
}

func blStr(n ast.Node) string {
	lit, ok := n.(*ast.BasicLit)
	if !ok {
		return ""
	}
	if lit.Kind != token.STRING {
		return ""
	}

	// Ignore error.
	v, _ := strconv.Unquote(lit.Value)
	return v
}

func identStr(n ast.Node) string {
	lit, ok := n.(*ast.Ident)
	if !ok {
		return ""
	}

	return lit.Name
}

func checkT(n ast.Node) bool {
	if n == nil {
		return false
	}
	recvIdent := n.(*ast.Ident)
	objF, ok := recvIdent.Obj.Decl.(*ast.Field)
	if !ok {
		return false
	}
	starEx, ok := objF.Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	starSel, ok := starEx.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	_ = starSel
	xIdent, ok := starSel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return xIdent.Name == "testing" && starSel.Sel.Name == "T"
}
