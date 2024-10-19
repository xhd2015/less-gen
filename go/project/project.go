package project

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"

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

func (c *Project) GetOnlyEntryPackage() (*Pkg, error) {
	if len(c.entryPkgs) == 0 {
		return nil, fmt.Errorf("no packages loaded")
	}
	if len(c.entryPkgs) > 1 {
		pkgs := make([]string, 0, len(c.entryPkgs))
		for _, entryPkg := range c.entryPkgs {
			pkgs = append(pkgs, entryPkg.Path())
		}
		return nil, fmt.Errorf("multiple pakcages loaded: %s", strings.Join(pkgs, ","))
	}
	return c.entryPkgs[0], nil
}

func (c *Project) GetPkg(pkgPath string) *Pkg {
	return c.pkgMapping[pkgPath]
}

func (c *Project) Fset() *token.FileSet {
	return c.fset
}

func (c *Project) GetCode(node ast.Node) (string, error) {
	start := c.fset.Position(node.Pos())
	if !start.IsValid() {
		return "", fmt.Errorf("not found")
	}
	end := c.fset.Position(node.End())
	if !end.IsValid() {
		return "", fmt.Errorf("not found")
	}
	file, err := os.ReadFile(start.Filename)
	if err != nil {
		return "", err
	}
	if start.Offset >= len(file) || end.Offset > len(file) {
		return "", fmt.Errorf("invalid node")
	}
	return string(file[start.Offset:end.Offset]), nil
}
