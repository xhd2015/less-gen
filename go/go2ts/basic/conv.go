package basic

import (
	"fmt"
	"path/filepath"

	load "github.com/xhd2015/less-gen/go/load/legacy"
)

type Translate struct {
	PkgPath string
	Code    string
}

type Options struct {
	Dir        string
	ForTest    bool
	BuildFlags []string
}

func LoadAndTranslate(args []string, opts *Options) ([]*Translate, error) {
	var absDir string
	var buildFlags []string
	if opts != nil {
		if opts.Dir != "" {
			var err error
			absDir, err = filepath.Abs(opts.Dir)
			if err != nil {
				return nil, err
			}
		}
		buildFlags = opts.BuildFlags
	}

	fset, pkgs, err := load.LoadPackages(args, &load.LoadOptions{
		Dir:        absDir,
		ForTest:    false,
		BuildFlags: buildFlags,
	})
	if err != nil {
		return nil, fmt.Errorf("loading packages err: %v", err)
	}

	results := make([]*Translate, 0, len(pkgs))
	for _, pkg := range pkgs {
		code := processPkg(fset, pkg)

		results = append(results, &Translate{
			PkgPath: pkg.PkgPath,
			Code:    code,
		})
	}
	return results, nil
}
