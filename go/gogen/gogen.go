package gogen

import (
	"go/ast"
	"go/token"
	"strings"
)

// example:
//
//	commands := []string{"go run github.com/xhd2015/ormx/cmd/ormx@latest gen", "go run github.com/xhd2015/ormx/cmd/ormx gen", "ormx gen"}
func FindGoGenerate(comments []*ast.CommentGroup, commands []string) []*ast.Comment {
	var found []*ast.Comment
	for _, comment := range comments {
		for _, line := range comment.List {
			suffix, ok := strings.CutPrefix(line.Text, "//go:generate ")
			if !ok {
				continue
			}
			for _, prefix := range commands {
				if strings.HasPrefix(suffix, prefix) {
					if isSpaceOrEnd(suffix[len(prefix):]) {
						found = append(found, line)
						break
					}
				}
			}
		}
	}
	return found
}

func isSpaceOrEnd(s string) bool {
	if s == "" {
		return true
	}
	return !token.IsIdentifier("_" + string(s[0]))
}
