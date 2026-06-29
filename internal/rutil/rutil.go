package rutil

import (
	"unicode/utf8"
)

func Tail(s string, start int) string {
	for start > 0 && s != "" {
		_, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		start--
	}
	return s
}

func Head(s string, end int) string {
	tail := Tail(s, end)
	return s[:len(s)-len(tail)]
}

func Body(s string, start, end int) string {
	s = Tail(s, start)
	return Head(s, end-start)
}

func Split(s string, col int) (string, string) {
	tail := Tail(s, col)
	head := s[:len(s)-len(tail)]
	return head, tail
}

func SplitBody(s string, start, end int) (string, string, string) {
	head, tail := Split(s, start)
	body, tail := Split(tail, end-start)
	return head, body, tail
}
