package xgravity

import (
	"go/ast"
)

const (
	MethodGetEntity = "get"
)

type Entity struct {
	Name        string
	PackagePath string
	Endpoints   []Endpoint

	StructProperties []*ast.Field
}

type Endpoint struct {
	Path   string
	Method string
	Name   string

	Input  []*ast.Field
	Result []*ast.Field
}
