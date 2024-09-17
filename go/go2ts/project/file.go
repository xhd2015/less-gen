package project

import (
	"go/ast"
	"go/token"
)

type File struct {
	pkg  *Pkg
	file *ast.File
}

func (c *File) GoFile() *ast.File {
	return c.file
}

func (c *File) GetCommentFor(pos token.Pos) *ast.CommentGroup {
	fset := c.pkg.project.fset
	posLine := fset.Position(pos).Line
	for _, comment := range c.file.Comments {
		endLine := fset.Position(comment.End()).Line
		if endLine+1 == posLine {
			return comment
		}
	}
	return nil
}
