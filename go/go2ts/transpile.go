package go2ts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/less-gen/go/go2ts/basic"
)

func TranspileFile(file string) (string, error) {
	res, err := basic.LoadAndTranslate([]string{file}, &basic.Options{})
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", fmt.Errorf("no packages found: %s", file)
	}
	if len(res) > 1 {
		return "", fmt.Errorf("multiple packages found: %s", file)
	}
	return res[0].Code, nil
}

func TranspileCode(goCode string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "go2ts-transpile")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)
	absMainFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(absMainFile, []byte(goCode), 0755)
	if err != nil {
		return "", err
	}

	return TranspileFile(absMainFile)
}
