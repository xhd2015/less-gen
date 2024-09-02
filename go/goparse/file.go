package goparse

import (
	"go/ast"
	"go/token"
)

type File struct {
	AST  *ast.File
	Code string
	Fset *token.FileSet
}

func (c *File) GetStruct(name string) *Struct {
	if c.AST == nil {
		return nil
	}
	for _, decl := range c.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}

		if len(genDecl.Specs) != 1 {
			continue
		}
		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		if typeSpec.Name == nil || typeSpec.Name.Name != name {
			continue
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		return &Struct{
			Name:       typeSpec.Name.Name,
			TypeSpec:   typeSpec,
			StructType: structType,
			File:       c,
		}
	}
	return nil
}

func (c *File) GetPkgName() string {
	if c.AST == nil || c.AST.Name == nil {
		return ""
	}
	return c.AST.Name.Name
}

func (c *File) GetNodeCode(node ast.Node) string {
	if c == nil || c.Fset == nil || node == nil {
		return ""
	}
	start := c.Fset.Position(node.Pos()).Offset
	end := c.Fset.Position(node.End()).Offset

	if start < 0 || start >= end || end > len(c.Code) {
		return ""
	}
	return c.Code[start:end]
}
