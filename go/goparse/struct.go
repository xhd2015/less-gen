package goparse

import "go/ast"

type Struct struct {
	File       *File
	Name       string
	TypeSpec   *ast.TypeSpec
	StructType *ast.StructType
}

func (c *Struct) GetFields() []*Field {
	if c.StructType == nil {
		return nil
	}
	fields := c.StructType.Fields
	if fields == nil || fields.List == nil {
		return nil
	}
	var goFields []*Field
	for _, field := range fields.List {

		for _, name := range field.Names {
			var fieldName string
			if name != nil {
				// TODO: verify anonymouse field
				fieldName = name.Name
			}
			goFields = append(goFields, &Field{
				Name:     fieldName,
				Type:     field.Type,
				TypeStr:  c.File.GetNodeCode(field.Type),
				GoStruct: c,
			})
		}
	}
	return goFields
}
