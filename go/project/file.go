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

func (c *File) Path() string {
	pos := c.pkg.project.fset.Position(c.file.Package)
	return pos.Filename
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

func (c *File) SearchNode(pos token.Pos) ast.Node {
	return searchNode(c.file, pos)
}

func searchNode(node ast.Node, pos token.Pos) ast.Node {
	var found ast.Node
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		if n.Pos() == pos {
			found = n
			return false
		}
		return true
	})
	return found
}
