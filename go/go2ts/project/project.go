package project

import (
	"go/token"

	"github.com/xhd2015/less-gen/go/load"
	"golang.org/x/tools/go/packages"
)

type Project struct {
	fset      *token.FileSet
	entryPkgs []*Pkg

	pkgMapping map[string]*Pkg
}

func Load(args []string, opts *load.LoadOptions) (*Project, error) {
	fset, pkgs, err := load.LoadPackages(args, opts)
	if err != nil {
		return nil, err
	}
	pkgMapping := make(map[string]*Pkg, len(pkgs))
	project := &Project{
		fset:       fset,
		pkgMapping: pkgMapping,
	}

	packages.Visit(pkgs, func(p *packages.Package) bool {
		pkgMapping[p.PkgPath] = &Pkg{project: project, pkg: p}
		return true
	}, nil)

	entryPkgs := make([]*Pkg, 0, len(pkgs))
	for _, pkg := range pkgs {
		entryPkgs = append(entryPkgs, pkgMapping[pkg.PkgPath])
	}
	project.entryPkgs = entryPkgs
	return project, nil
}

func (c *Project) GetEntryPackages() []*Pkg {
	return c.entryPkgs
}

func (c *Project) GetPkg(pkgPath string) *Pkg {
	return c.pkgMapping[pkgPath]
}
func (c *Project) Fset() *token.FileSet {
	return c.fset
}
