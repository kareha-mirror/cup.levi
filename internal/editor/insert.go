package editor

import (
	"unicode/utf8"

	"tea.kareha.org/lab/termi"
)

type Insert struct {
	head, tail string
	body       *RuneBuf
}

const maxBodyLen = 1024

func NewInsert() *Insert {
	return &Insert{
		head: "",
		tail: "",
		body: new(RuneBuf),
	}
}

func (ins *Insert) Reset() {
	ins.head = ""
	ins.tail = ""
	if ins.body.Len() > maxBodyLen {
		ins.body = new(RuneBuf)
	} else {
		ins.body.Reset()
	}
}

func (ins *Insert) Write(r rune) {
	ins.body.WriteRune(r)
}

func (ins *Insert) Enter(line string, col int) {
	rs := []rune(line)
	ins.head = string(rs[:col])
	if col < len(rs) {
		ins.tail = string(rs[col:])
	} else {
		ins.tail = ""
	}
}

func (ins *Insert) Line() string {
	return ins.head + ins.body.String() + ins.tail
}

func (ins *Insert) Newline() []string {
	lines := []string{
		ins.head + ins.body.String(),
		ins.tail,
	}
	ins.head = ""
	ins.body.Reset()
	// tail is intentionally preserved
	return lines
}

func (ins *Insert) Width() int {
	s := ins.head + ins.body.String()
	rc := utf8.RuneCountInString(s)
	return termi.StringWidth(s, rc)
}

func (ins *Insert) Backspace() bool {
	return ins.body.Backspace()
}
