package xgravity

import (
	"go/ast"
	"net/http"
)

const (
	RestMethodGET = "get"
)

func ParseRest(action *StructAction, entity *Entity) ([]Endpoint, error) {
	var endpoints []Endpoint

	for _, value := range action.Values {
		switch value {
		case RestMethodGET:
			endpoint := Endpoint{
				Path:   entity.Name + "/:id",
				Method: http.MethodGet,
				Name:   "get" + entity.Name + "ById",
				Input: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "id"},
						},
						Type: &ast.Ident{
							Name: "int",
						},
					},
				},
				Result: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "parent" + "." + entity.Name,
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			}

			endpoints = append(endpoints, endpoint)
		}

	}

	return endpoints, nil
}
