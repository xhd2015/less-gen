package gofile

type BuiltinType string

type Type interface {
	Node
}

const (
	Int            BuiltinType = "int"
	Int64          BuiltinType = "int64"
	String         BuiltinType = "string"
	Boolean        BuiltinType = "boolean"
	InterfaceEmpty BuiltinType = "interface{}"
)

type Named struct {
	PkgPath string
	Name    string
}

func (c BuiltinType) Format(ctx Context) string {
	return string(c)
}

func (c *Named) Format(ctx Context) string {
	ref := ctx.GetPkgRef(c.PkgPath)
	if ref != "" {
		return ref + "." + c.Name
	}
	return c.Name
}

func (c BuiltinType) Equals(typ Type) bool {
	if typ, ok := typ.(BuiltinType); ok {
		return c == typ
	}
	return false
}
