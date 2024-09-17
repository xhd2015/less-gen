package ts

import (
	"go/token"
	"go/types"

	"github.com/xhd2015/less-gen/go/load"
	"golang.org/x/tools/go/packages"
)

type Project struct {
	fset *token.FileSet
	pkgs []*packages.Package

	pkgMapping map[string]*Pkg
}

type Pkg struct {
	pkg *packages.Package
}

type Object struct {
	obj types.Object
}

type Struct struct {
}

func Load(args []string, opts *load.LoadOptions) (*Project, error) {
	fset, pkgs, err := load.LoadPackages(args, opts)
	if err != nil {
		return nil, err
	}
	pkgMapping := make(map[string]*Pkg, len(pkgs))
	for _, pkg := range pkgs {
		pkgMapping[pkg.PkgPath] = &Pkg{pkg: pkg}
	}
	return &Project{
		fset:       fset,
		pkgs:       pkgs,
		pkgMapping: pkgMapping,
	}, nil
}

func (c *Project) GetPkg(pkgPath string) *Pkg {
	return c.pkgMapping[pkgPath]
}

func (c *Pkg) GoPkg() *packages.Package {
	return c.pkg
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

func (c *Object) GoObject() types.Object {
	return c.obj
}

func (c *Object) AsStruct() types.Object {
	return c.obj
}
