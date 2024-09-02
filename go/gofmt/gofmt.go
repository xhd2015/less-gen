package gofmt

import (
	"path/filepath"

	"github.com/xhd2015/xgo/support/cmd"
)

func Format(goFile string) error {
	dir := filepath.Dir(goFile)
	file := filepath.Base(goFile)
	// and cd /tmp to avoid go get unwanted dependencies, make goimports just like a real tool.
	err := cmd.Dir(dir).Run("go", "run", "-v", "-mod=mod", "golang.org/x/tools/cmd/goimports", "-w", file)
	if err != nil {
		return err
	}

	return cmd.Dir(dir).Run("gofmt", "-w", "-s", file)
}
