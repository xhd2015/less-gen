package gofile

import (
	"fmt"
	"strings"
)

type File struct {
	PkgName string

	Imports []*Import
	Decls   []Node
}

func (c *File) Import(imp *Import) {
	c.Imports = append(c.Imports, imp)
}

func (c *File) AddDecl(decl Node) {
	c.Decls = append(c.Decls, decl)
}

func (c *File) Format(ctx Context) string {
	var imps []string
	for _, imp := range c.Imports {
		imps = append(imps, imp.String())
	}
	var impStmt string
	if len(imps) > 0 {
		if len(imps) == 1 {
			impStmt = "import " + imps[0]
		} else {
			impStmt = fmt.Sprintf("import (\n%s\n)", strings.Join(imps, "\n"))
		}
	}

	decls := make([]string, 0, len(c.Decls))
	for _, decl := range c.Decls {
		decls = append(decls, decl.Format(ctx))
	}
	return fmt.Sprintf("package %s\n%s\n%s\n", c.PkgName, impStmt, strings.Join(decls, "\n"))
}
