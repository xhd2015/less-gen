package semantic

import (
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Context struct {
	fset *token.FileSet
	pkgs []*packages.Package
}

func NewContext(fset *token.FileSet, pkgs []*packages.Package) *Context {
	return &Context{
		fset: fset,
		pkgs: pkgs,
	}
}

// the relation is

// 1. find all calls to r.Any("api", processor.Gen(...))
// 2. parse function Request and Response
// 3. build relation
func (c *Context) FindCalls() {

}
