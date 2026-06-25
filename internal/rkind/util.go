package rkind

func IsBlankLine(line string) bool {
	for _, r := range line {
		if !IsBlank(r) {
			return false
		}
	}
	return true
}
