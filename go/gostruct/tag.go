package gostruct

// ParseTag strips the backticks from the tag
func ParseTag(tag string) string {
	if tag == "" {
		return ""
	}
	if len(tag) <= 2 {
		return ""
	}
	if tag[0] != '`' || tag[len(tag)-1] != '`' {
		return ""
	}
	return tag[1 : len(tag)-1]
}
