package goparse

import (
	"os"
	"path/filepath"
	"strings"
)

type Pkg struct {
	Files []*File
	Dir   string
}

func ParsePkg(dir string) (*Pkg, error) {
	files, err := ParseGoDir(dir)
	if err != nil {
		return nil, err
	}
	return &Pkg{
		Files: files,
		Dir:   dir,
	}, nil
}

func ParseGoDir(dir string) ([]*File, error) {
	subFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var goFiles []*File
	for _, file := range subFiles {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".go") {
			continue
		}
		goFile, err := Parse(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		goFiles = append(goFiles, goFile)
	}
	return goFiles, nil
}

func (c *Pkg) GetStruct(name string) *Struct {
	for _, goFile := range c.Files {
		goStruct := goFile.GetStruct(name)
		if goStruct != nil {
			return goStruct
		}
	}
	return nil
}
