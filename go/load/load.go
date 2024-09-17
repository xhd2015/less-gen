package load

import (
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

type LoadOptions struct {
	Dir        string
	ForTest    bool
	BuildFlags []string // see FlagBuilder
	LoadMode   []packages.LoadMode
}

func LoadPackages(args []string, opts *LoadOptions) (*token.FileSet, []*packages.Package, error) {
	fset := token.NewFileSet()
	dir := opts.Dir

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, nil, err
	}
	var loadMode packages.LoadMode
	if len(opts.LoadMode) > 0 {
		for _, m := range opts.LoadMode {
			loadMode |= m
		}
	} else {
		// all
		loadMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedDeps | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedModule
	}

	cfg := &packages.Config{
		Dir:        absDir,
		Mode:       loadMode,
		Fset:       fset,
		Tests:      opts.ForTest,
		BuildFlags: opts.BuildFlags,
	}
	pkgs, err := packages.Load(cfg, args...)
	if err != nil {
		return nil, nil, err
	}
	return fset, pkgs, nil
}
