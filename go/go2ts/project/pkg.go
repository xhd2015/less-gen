package project

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Pkg struct {
	project *Project
	pkg     *packages.Package
}

func (c *Pkg) GoPkg() *packages.Package {
	return c.pkg
}

func (c *Pkg) Path() string {
	return c.pkg.PkgPath
}

func (c *Pkg) Types() *types.Package {
	return c.pkg.Types
}
func (c *Pkg) TypesInfo() *types.Info {
	return c.pkg.TypesInfo
}

func (c *Pkg) Lookup(name string) *Object {
	obj := c.pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return nil
	}
	return &Object{
		obj: obj,
	}
}

func (c *Pkg) GetContainingFile(pos token.Pos) *File {
	fset := c.project.fset
	fileName := fset.Position(pos).Filename

	for _, file := range c.pkg.Syntax {
		filePos := fset.Position(file.Pos())
		if filePos.Filename == fileName {
			return &File{
				pkg:  c,
				file: file,
			}
		}
	}
	return nil
}
