package rkind

func IsBlankLine(s string) bool {
	for _, r := range s {
		if !IsBlank(r) {
			return false
		}
	}
	return true
}
