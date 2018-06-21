package xgravity

import (
	"go/ast"
)

const (
	MethodGetEntity = "get"
)

type Entity struct {
	Name    string
	Methods []string

	Properties []*ast.Field
}
