package project

import (
	"fmt"
	"go/ast"
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

func (c *Pkg) Get(name string) (*Object, error) {
	obj := c.Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("%s not found", name)
	}
	return obj, nil
}

// GetStructType lookup the given name, validate
// it's a TypeName object, and get it's type
func (c *Pkg) GetStructType(name string) (*Struct, error) {
	obj, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	goObj := obj.GoObject()
	namedType, ok := goObj.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("want *types.TypeName, %s has %T", name, goObj)
	}
	ut := namedType.Type().Underlying()
	structT, ok := ut.(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("want *types.Struct, %s underlying has %T", name, ut)
	}
	return &Struct{
		goStruct: structT,
	}, nil
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

func (c *Pkg) SearchNode(pos token.Pos) ast.Node {
	for _, file := range c.pkg.Syntax {
		n := searchNode(file, pos)
		if n != nil {
			return n
		}
	}
	return nil
}
