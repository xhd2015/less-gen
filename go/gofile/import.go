package gofile

import (
	"strconv"
	"strings"
)

type Import struct {
	Name string
	Path string
}

func (c *Import) String() string {
	return strings.Join([]string{c.Name, strconv.Quote(c.Path)}, " ")
}
