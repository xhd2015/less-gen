package gofile

import (
	"fmt"
	"strings"
)

type Struct struct {
	Fields []*Field
}

func (c *Struct) AddField(field *Field) *Struct {
	c.Fields = append(c.Fields, field)
	return c
}

func (c *Struct) Format(ctx Context) string {
	fieldStrs := make([]string, 0, len(c.Fields))
	for _, field := range c.Fields {
		fieldStrs = append(fieldStrs, field.Format(ctx))

	}
	return fmt.Sprintf("struct{\n%s\n}", strings.Join(fieldStrs, "\n"))
}

type Field struct {
	Name string
	Type Type
	Tag  string
}

func (c *Field) Format(ctx Context) string {
	var tag string
	if c.Tag != "" {
		tag = "`" + c.Tag + "`"
	}
	return strings.Join([]string{c.Name, c.Type.Format(ctx), tag}, " ")
}
