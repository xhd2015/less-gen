package gofile

import "strings"

type Context interface {
	GetPkgRef(pkgPath string) string
}

type Node interface {
	Format(ctx Context) string
}

func NewContext() Context {
	return &context{}
}

type context struct {
}

func (c *context) GetPkgRef(pkgPath string) string {
	last := pkgPath
	idx := strings.LastIndex(pkgPath, "/")
	if idx >= 0 {
		last = pkgPath[idx+1:]
	}
	return last
}
