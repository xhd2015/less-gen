package gofile

type StructField struct {
	Field  *Field
	Nested []*StructField
}

func (c *StructField) IsNested() bool {
	return len(c.Nested) > 0
}
