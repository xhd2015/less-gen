package astinfo

import (
	"fmt"
	"go/token"
)

func FileLine(fset *token.FileSet, pos token.Pos) string {
	return fmt.Sprintf("%s:%d", fset.Position(pos).Filename, fset.Position(pos).Line)
}
