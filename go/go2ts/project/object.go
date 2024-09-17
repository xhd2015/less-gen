package project

import "go/types"

type Object struct {
	obj types.Object
}

func (c *Object) GoObject() types.Object {
	return c.obj
}

func (c *Object) AsStruct() types.Object {
	return c.obj
}
