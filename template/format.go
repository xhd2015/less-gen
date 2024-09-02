package template

import "strings"

func Format(template string, vars map[string]string) string {
	template = RemoveComment(template)
	for k, v := range vars {
		template = strings.ReplaceAll(template, k, v)
	}
	return template
}

func RemoveComment(tpl string) string {
	lines := strings.Split(tpl, "\n")

	var inComment bool
	j := 0
	for _, line := range lines {
		noSpaceLine := strings.TrimSpace(line)
		if inComment {
			if noSpaceLine == "// </TEMPLATE>" {
				inComment = false
			}
			continue
		}

		if noSpaceLine == "// <TEMPLATE>" {
			inComment = true
			continue
		}
		lines[j] = line
		j++
	}

	return strings.Join(lines[:j], "\n")
}
