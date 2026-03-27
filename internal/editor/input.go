package editor

import (
	"unicode/utf8"
)

type Input struct {
	head, tail string
	body       *RuneBuf
}

const maxBodyLen = 1024

func NewInput() *Input {
	return &Input{
		head: "",
		tail: "",
		body: new(RuneBuf),
	}
}

func (inp *Input) Reset() {
	inp.head = ""
	inp.tail = ""
	if inp.body.Len() > maxBodyLen {
		inp.body = new(RuneBuf)
	} else {
		inp.body.Reset()
	}
}

func (inp *Input) Init(line string, col int) {
	inp.Reset()
	rs := []rune(line)
	inp.head = string(rs[:col])
	if col < len(rs) {
		inp.tail = string(rs[col:])
	} else {
		inp.tail = ""
	}
}

func (inp *Input) WriteRune(r rune) {
	inp.body.WriteRune(r)
}

func (inp *Input) Line() string {
	return inp.head + inp.body.String() + inp.tail
}

func (inp *Input) Newline() []string {
	lines := []string{
		inp.head + inp.body.String(),
		inp.tail,
	}
	inp.head = ""
	inp.body.Reset()
	// tail is intentionally preserved
	return lines
}

func (inp *Input) Column() int {
	s := inp.head + inp.body.String()
	return utf8.RuneCountInString(s)
}

func (inp *Input) Backspace() bool {
	return inp.body.Backspace()
}
