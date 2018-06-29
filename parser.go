package xgravity

import (
	//"go/ast"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const (
	CommandGenerateRest = "//xgravity:gen-rest"
)

type StructAction struct {
	Name   string
	Values []string
}

func GetEntities(packagePath, filename string, b []byte) ([]Entity, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, b, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var entities []Entity
	var genErr error

	ast.Inspect(f, func(n ast.Node) bool {
		if gen, ok := n.(*ast.GenDecl); ok && gen.Tok == token.TYPE && gen.Doc != nil {
			action, err := GetStructAction(gen.Doc)
			if err != nil {
				genErr = err
				return true
			}

			entity := Entity{
				PackagePath: packagePath,
			}

			for _, spec := range gen.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := typeSpec.Type.(*ast.StructType); ok {
						entity.Name = typeSpec.Name.Name
						entity.StructProperties = st.Fields.List
					}
				}
			}

			switch action.Name {
			case CommandGenerateRest:
				endpoints, err := ParseRest(action, &entity)
				if err != nil {
					genErr = err
					return true
				}

				entity.Endpoints = append(entity.Endpoints, endpoints...)
			}

			entities = append(entities, entity)
		}

		return true
	})

	return entities, genErr
}

func GetStructAction(comments *ast.CommentGroup) (*StructAction, error) {
	action := StructAction{}

	for _, comment := range comments.List {
		if comment != nil {
			items := strings.Split(comment.Text, " ")

			switch items[0] {
			case CommandGenerateRest:
				action.Name = CommandGenerateRest

				if len(items) < 2 {
					return nil, errors.New("methods not found")
				}

				action.Values = strings.Split(items[1], ",")
			default:
				return nil, errors.New("command not support")
			}
		}
	}

	return &action, nil
}
