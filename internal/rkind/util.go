package rkind

import (
	"strings"
)

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

func IndentOf(line string) string {
	for i, r := range line {
		if !IsBlank(r) {
			return line[:i]
		}
	}
	return line
}

func Escape(s string) string {
	b := strings.Builder{}
	for _, r := range s {
		if r < 0x20 {
			b.WriteRune('^')
			b.WriteRune(r + '@')
		} else if r == 0x7f {
			b.WriteRune('^')
			b.WriteRune('?')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
