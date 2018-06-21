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
	CommandGenerate = "//xgravity:gen"

	ActionGen = "gen"
)

type StructAction struct {
	Name    string
	Methods []string
}

func GetEntities(filename string, b []byte) ([]Entity, error) {
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

			switch action.Name {
			case ActionGen:
				entity := Entity{Methods: action.Methods}

				for _, spec := range gen.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if st, ok := typeSpec.Type.(*ast.StructType); ok {
							entity.Name = typeSpec.Name.Name
							entity.Properties = st.Fields.List
						}
					}
				}

				entities = append(entities, entity)
			}
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
			case CommandGenerate:
				action.Name = ActionGen

				if len(items) < 2 {
					return nil, errors.New("methods not found")
				}

				action.Methods = strings.Split(items[1], ",")
			default:
				return nil, errors.New("command not support")
			}
		}
	}

	return &action, nil
}
