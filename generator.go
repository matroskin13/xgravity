package xgravity

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
)

type Files map[string][]byte

func CreateApiPackage(currentPackage string, entities []Entity) (Files, error) {
	files := Files{}
	fset := token.NewFileSet()
	content := createEmptyFile("api")

	f, err := parser.ParseFile(fset, "api.go", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok && d.Tok == token.IMPORT {
			d.Specs = append(d.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: strconv.Quote(currentPackage),
				},
				Name: &ast.Ident{
					Name: "parent",
				},
			})
		}
	}

	for _, decl := range GenerateInterfaces(entities) {
		f.Decls = append(f.Decls, decl)
	}

	var buf bytes.Buffer

	printer.Fprint(&buf, fset, f)

	httpBytes, err := GetHTTPTemplate(currentPackage, entities)
	if err != nil {
		return nil, err
	}

	files["api.go"] = buf.Bytes()
	files["http.go"] = httpBytes

	return files, nil
}

func createEmptyFile(packageName string) string {
	return "package " + packageName + "\r\n\r\nimport()"
}

func GenerateInterfaces(entities []Entity) []*ast.GenDecl {
	var decls []*ast.GenDecl

	for _, entity := range entities {
		var methods []*ast.Field

		for _, endpoint := range entity.Endpoints {
			methods = append(methods, &ast.Field{
				Names: []*ast.Ident{
					{
						Name: endpoint.Name,
						Obj:  ast.NewObj(ast.Fun, endpoint.Name),
					},
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: endpoint.Input,
					},
					Results: &ast.FieldList{
						List: endpoint.Result,
					},
				},
			})
		}

		interfaceName := entity.Name + "Api"

		decl := &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{
						Name: interfaceName,
						Obj:  ast.NewObj(ast.Typ, interfaceName),
					},
					Type: &ast.InterfaceType{
						Methods: &ast.FieldList{
							List: methods,
						},
					},
				},
			},
		}

		decls = append(decls, decl)
	}

	return decls
}
