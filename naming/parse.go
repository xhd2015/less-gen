package naming

import "strings"

// prefix pkg path and basic name
// example:  a/b/c.D => a,b,c D
func SplitDotRef(ref string) ([]string, string) {
	dotIdx := strings.LastIndex(ref, ".")
	if dotIdx < 0 {
		return nil, ref
	}
	basicName := ref[dotIdx+1:]
	prefix := ref[:dotIdx]
	if prefix == "" {
		return nil, basicName
	}

	return strings.Split(prefix, "/"), basicName
}
