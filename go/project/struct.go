package project

import (
	"go/types"
)

type Struct struct {
	goStruct *types.Struct
}

type Field struct {
	goField *types.Var
}

func (c *Struct) GoStruct() *types.Struct {
	return c.goStruct
}

func (c *Struct) GetFields() []*Field {
	return getStructFields(c.goStruct)
}

func (c *Field) GoField() *types.Var {
	return c.goField
}

// func (c *Field) A(){
// c.goField.
// }
func (c *Field) GetNested() []*Field {
	field := c.goField
	if !field.Anonymous() {
		return nil
	}
	underlyingType := field.Type()
	if ptrType, ok := field.Type().Underlying().(*types.Pointer); ok {
		underlyingType = ptrType
	}

	subType, ok := underlyingType.(*types.Struct)
	if !ok {
		return nil
	}
	return getStructFields(subType)
}

func getStructFields(structType *types.Struct) []*Field {
	numField := structType.NumFields()
	var fields []*Field
	fields = make([]*Field, 0, numField)
	for i := 0; i < numField; i++ {
		field := structType.Field(i)

		fields = append(fields, &Field{
			goField: field,
		})
	}
	return fields
}
