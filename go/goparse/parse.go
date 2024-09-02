package goparse

import (
	"github.com/xhd2015/xgo/support/goparse"
)

func Parse(file string) (*File, error) {
	code, astFile, fset, err := goparse.Parse(file)
	if err != nil {
		return nil, err
	}
	return &File{
		AST:  astFile,
		Code: string(code),
		Fset: fset,
	}, nil
}
