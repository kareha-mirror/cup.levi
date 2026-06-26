package rkind

func IsBlankLine(line string) bool {
	for _, r := range line {
		if !IsBlank(r) {
			return false
		}
	}
	return true
}

func TrimPrefixBlanks(s string) string {
	for i, r := range s {
		if !IsBlank(r) {
			return s[i:]
		}
	}
	return ""
}
