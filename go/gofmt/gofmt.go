package gofmt

import (
	"os"
	"path/filepath"

	"github.com/xhd2015/xgo/support/cmd"
)

func Format(goFile string) error {
	dir := filepath.Dir(goFile)
	file := filepath.Base(goFile)
	// and cd /tmp to avoid go get unwanted dependencies, make goimports just like a real tool.
	err := cmd.Dir(dir).Run("go", "run", "-v", "-mod=mod", "golang.org/x/tools/cmd/goimports@latest", "-w", file)
	if err != nil {
		return err
	}

	return cmd.Dir(dir).Run("gofmt", "-w", "-s", file)
}

func TryFormatCode(code string) string {
	fcode, err := FormatCode(code)
	if err != nil {
		return code
	}
	return fcode
}

func FormatCode(code string) (string, error) {
	dir, err := os.MkdirTemp("", "gofmt")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	codeFile := filepath.Join(dir, "code.go")
	err = os.WriteFile(codeFile, []byte(code), 0755)
	if err != nil {
		return "", err
	}
	err = Format(codeFile)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(codeFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
