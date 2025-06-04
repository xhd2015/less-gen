package load

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/packages"
)

type Packages struct {
	Fset           *token.FileSet
	Packages       []*packages.Package
	PackageMapping map[string]*packages.Package
}

func Load(dir string, patterns ...string) (*Packages, error) {
	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Fset: fset,
		Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedDeps | packages.NeedName | packages.NeedImports | packages.NeedTypesInfo,
		Dir:  dir,
	}, patterns...)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}
	packageMapping := make(map[string]*packages.Package, len(pkgs))
	for _, pkg := range pkgs {
		packageMapping[pkg.PkgPath] = pkg
	}
	return &Packages{
		Fset:           fset,
		Packages:       pkgs,
		PackageMapping: packageMapping,
	}, nil
}
