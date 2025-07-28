package strcase

import "strings"

func CamelToSnake(s string) string {
	list := SplitCamelCase(s)
	for i, e := range list {
		list[i] = strings.ToLower(e)
	}
	return strings.Join(list, "_")
}

func SnakeToCamel(s string) string {
	list := strings.Split(s, "_")
	for i, e := range list {
		list[i] = Capitalize(e)
	}
	return strings.Join(list, "")
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func Decapitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func SplitCamelCase(s string) []string {
	if s == "" {
		return nil
	}
	var list []string
	n := len(s)

	last := 0
	i := 0
	for i < n {
		j := findCamelCaseEnd(s, i)
		if j > last {
			list = append(list, s[last:j])
		}
		last = j
		i = j + 1
	}
	if last < n {
		list = append(list, s[last:])
	}
	return list
}

func isUpperCase(s string) bool {
	return strings.ToUpper(s) == s && !isDigit(s)
}

func isLowerCase(s string) bool {
	return strings.ToLower(s) == s && !isDigit(s)
}

func isDigit(s string) bool {
	if len(s) != 1 {
		return false
	}
	return s[0] >= '0' && s[0] <= '9'
}

//	Xxxx Yyyyy
//
// find boundary
// assumption: s[i-1] is upper case
func findCamelCaseEnd(s string, i int) int {
	j := nextLower(s, i)
	if j >= len(s) {
		return j
	}

	// none is upper case
	if j == i {
		return nextUpper(s, j+1)
	}

	// multiple contineous upper case

	// j lower, j-1: upper
	// j-1 = i

	// j-1: last upper
	return j - 1
}

func nextUpper(s string, i int) int {
	for i < len(s) && !isUpperCase(s[i:i+1]) {
		i++
	}
	return i
}

func nextLower(s string, i int) int {
	for i < len(s) && !isLowerCase(s[i:i+1]) {
		i++
	}
	return i
}
