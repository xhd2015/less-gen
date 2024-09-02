package fs

import (
	"os"
	"path/filepath"
	"strings"
)

type FileOption func(file string) error

// NOTE: files' key always use "/" as separator
func MkdirWriteFiles(rootDir string, files map[string]string, opts ...FileOption) error {
	for file, code := range files {
		fullFile := filepath.Join(rootDir, filepath.Join(strings.Split(file, "/")...))
		err := MkdirWriteFile(fullFile, code)
		if err != nil {
			return err
		}
		for _, opt := range opts {
			err := opt(fullFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MkdirWriteFile(file string, code string) error {
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return os.WriteFile(file, []byte(code), 0755)
}
