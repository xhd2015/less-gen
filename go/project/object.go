package project

import "go/types"

type Object struct {
	obj types.Object
}

func (c *Object) GoObject() types.Object {
	return c.obj
}
