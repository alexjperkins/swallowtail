package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"strings"

	"go/parser"
	"go/token"
	"io/ioutil"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

var fset *token.FileSet

func main() {
	if len(os.Args) < 2 {
		panic("usage: go run tools/protoc-gen-typhon/rewrite-router.go service.foo")
	}

	serviceName := os.Args[1]
	routerPath := fmt.Sprintf("%s/handler/router.go", serviceName)
	protoPath := fmt.Sprintf("%s/proto", serviceName)

	contents, err := ioutil.ReadFile(routerPath)
	if err != nil {
		panic(err)
	}

	fset = token.NewFileSet()
	routerFile, err := parser.ParseFile(fset, routerPath, contents, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	pkgs, err := parser.ParseDir(fset, protoPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	var protoPackage *ast.Package
	for _, p := range pkgs {
		protoPackage = p
	}

	// Find all router.Register(method, path, handlers.Foo) calls in the proto package
	registerCalls := findGeneratedRoutes(protoPackage)

	// Find the Handlers{ Foo: handleFoo, ... }.Router() struct in router.go
	handlersLiteral := findHandlersLiteral(routerFile)
	if handlersLiteral == nil {
		fmt.Fprintf(os.Stderr, "error: could not find *Handlers{}.Router() in %s/handler/router.go\n", serviceName)
		os.Exit(1)
	}

	// Generate new router.METHOD(path, handleFoo) calls
	newRegisterCalls := generateRoutes(handlersLiteral, registerCalls)

	// Remove all comments. astutil.Apply() won't automatically preserve their
	// location for us :-(
	routerFile.Comments = []*ast.CommentGroup{}

	// Remove the proto import from router.go
	removeProtoImport(routerFile, serviceName)

	// Find an existing init() declaration, and add the new routes if it exists.
	hasExistingInit := addRoutesToInit(routerFile, newRegisterCalls)

	// Rewrite
	//     var Router = &fooproto.FooHandlers{}.Router()
	// to
	//     var Router = typhon.NewRouter()
	replaceGeneratedRouterMethodCall(routerFile, newRegisterCalls)

	// If there's no existing init(), add one just before the Service() declaration.
	if !hasExistingInit {
		addInitDecl(routerFile, newRegisterCalls)
	}

	buf := new(bytes.Buffer)
	if err := format.Node(buf, fset, routerFile); err != nil {
		panic(err)
	}

	ioutil.WriteFile(routerPath, buf.Bytes(), 0)
}

func removeProtoImport(routerFile *ast.File, serviceName string) {
	pre := func(c *astutil.Cursor) bool {
		if importSpec, ok := c.Node().(*ast.ImportSpec); ok {
			if strings.Contains(importSpec.Path.Value, fmt.Sprintf("%s/proto", serviceName)) {
				c.Delete()
				return false
			}
		}
		return true
	}

	astutil.Apply(routerFile, pre, nil)
}

func addRoutesToInit(routerFile *ast.File, routes []ast.Stmt) (hasExistingInit bool) {
	pre := func(c *astutil.Cursor) bool {
		if isFuncDecl(c.Parent(), "init") {
			if block, ok := c.Node().(*ast.BlockStmt); ok {
				hasExistingInit = true

				c.Replace(&ast.BlockStmt{
					Lbrace: block.Lbrace,
					Rbrace: block.Rbrace,
					List:   append(routes, block.List...),
				})
			}
		}
		return true
	}

	astutil.Apply(routerFile, pre, nil)
	return hasExistingInit
}

func addInitDecl(routerFile *ast.File, routes []ast.Stmt) {
	// Insert the init decl before the first function decl
	seenFunc := false
	newDecls := []ast.Decl{}
	for _, decl := range routerFile.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok && !seenFunc {
			newDecls = append(newDecls, &ast.FuncDecl{
				Name: ast.NewIdent("init"),
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: routes,
				},
			})
			seenFunc = true
		}
		newDecls = append(newDecls, decl)
	}
	routerFile.Decls = newDecls
}

func replaceGeneratedRouterMethodCall(
	routerFile *ast.File,
	routes []ast.Stmt,
) {
	pre := func(c *astutil.Cursor) bool {
		if isHandlersRouterMethodCall(c.Node()) {
			c.Replace(&ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("typhon"),
					Sel: ast.NewIdent("Router"),
				},
			})
			return false
		}
		return true
	}

	astutil.Apply(routerFile, pre, nil)
}

func findGeneratedRoutes(protoPackage *ast.Package) (calls []*ast.CallExpr) {
	inspectNode := func(n ast.Node) bool {
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if decl.Name.Name != "Router" {
			return true
		}

		lines := decl.Body.List

		for _, l := range lines {
			expr, ok := l.(*ast.ExprStmt)
			if !ok {
				continue
			}
			call, ok := expr.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if selector.Sel.Name != "Register" {
				continue
			}
			calls = append(calls, call)
		}
		return false
	}

	for _, f := range protoPackage.Files {
		ast.Inspect(f, inspectNode)
	}
	return calls
}

func findHandlersLiteral(routerFile ast.Node) (handlers *ast.CompositeLit) {
	inspectNode := func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if selector.Sel.Name != "Router" {
			return true

		}

		compLit, ok := selector.X.(*ast.CompositeLit)
		if !ok {
			return true
		}
		handlers = compLit
		return false
	}

	ast.Inspect(routerFile, inspectNode)
	return handlers
}

func generateRoutes(
	handlers *ast.CompositeLit,
	registerCalls []*ast.CallExpr,
) []ast.Stmt {
	handlerFuncsByKey := map[string]string{}

	for _, field := range handlers.Elts {
		keyValue := field.(*ast.KeyValueExpr)
		keyIdent := keyValue.Key.(*ast.Ident)
		funcName := keyValue.Value.(*ast.Ident)
		handlerFuncsByKey[keyIdent.Name] = funcName.Name

	}
	_ = handlerFuncsByKey

	routes := []ast.Stmt{}
	for _, call := range registerCalls {
		method := call.Args[0].(*ast.BasicLit)

		path := call.Args[1]

		selector := call.Args[2].(*ast.SelectorExpr)
		funcName := handlerFuncsByKey[selector.Sel.Name]

		r := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("Router"),
					Sel: ast.NewIdent(strings.Trim(method.Value, "\"")),
				},
				Args: []ast.Expr{
					path,
					ast.NewIdent(funcName),
				},
			},
		}
		routes = append(routes, r)
	}
	return routes
}

func isHandlersRouterMethodCall(n ast.Node) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return selector.Sel.Name == "Router"
}

func isFuncDecl(n ast.Node, name string) bool {
	decl, ok := n.(*ast.FuncDecl)
	if !ok {
		return false
	}

	return decl.Name.Name == name
}

func printNode(node ast.Node) {
	buf := new(bytes.Buffer)
	format.Node(buf, fset, node)
	fmt.Println(buf.String())
}
