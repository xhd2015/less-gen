package goparse

import "go/ast"

type Field struct {
	Name    string
	Type    ast.Expr
	TypeStr string

	GoStruct *Struct
}
