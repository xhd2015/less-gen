package main

import (
	"go/types"

	"github.com/xhd2015/less-gen/go/gofile"
)

func GetStructFields(structType *types.Struct) []*gofile.StructField {
	numField := structType.NumFields()
	var fields []*gofile.StructField
	fields = make([]*gofile.StructField, 0, numField)
	for i := 0; i < numField; i++ {
		field := structType.Field(i)
		tag := structType.Tag(i)

		var nested []*gofile.StructField
		// fmt.Printf("field: %T %v\n", field.Type(), field.Type())
		if field.Anonymous() {
			underlyingType := field.Type()
			if ptrType, ok := field.Type().Underlying().(*types.Pointer); ok {
				underlyingType = ptrType
			}

			subType, ok := underlyingType.(*types.Struct)
			if ok {
				nested = GetStructFields(subType)
			}
		}
		fields = append(fields, &gofile.StructField{
			Field:  buildField(field, tag),
			Nested: nested,
		})
	}
	return fields
}

func buildField(field *types.Var, tag string) *gofile.Field {
	name := field.Name()

	ft := field.Type()
	// fmt.Printf("Underlying: %T %v\n", ft.Underlying(), ft.Underlying())
	var t gofile.Type = gofile.InterfaceEmpty
	if bt, ok := ft.Underlying().(*types.Basic); ok {
		if bt.Kind() == types.Int64 {
			t = gofile.Int64
		}
	}
	if named, ok := ft.(*types.Named); ok {
		typObj := named.Obj()
		t = &gofile.Named{
			PkgPath: typObj.Pkg().Path(),
			Name:    typObj.Name(),
		}
	}

	return &gofile.Field{
		Name: name,
		Type: t,
		Tag:  tag,
	}
}
