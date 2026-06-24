package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

type Input struct {
	head, tail string
	bodies     []termi.RuneBuf
}

func NewInput() *Input {
	return &Input{
		head:   "",
		tail:   "",
		bodies: []termi.RuneBuf{termi.RuneBuf{}},
	}
}

func (inp *Input) Reset() {
	inp.head = ""
	inp.tail = ""
	inp.bodies = []termi.RuneBuf{termi.RuneBuf{}}
}

func (inp *Input) Init(line string, col int, ai bool) {
	inp.Reset()
	rs := []rune(line)
	inp.head = string(rs[:col])
	if col < len(rs) {
		inp.tail = string(rs[col:])
	} else {
		inp.tail = ""
	}
	if ai && isBlankLine(inp.head) {
		inp.bodies[0].WriteString(inp.head)
		inp.head = ""
	}
}

func (inp *Input) body() *termi.RuneBuf {
	return &inp.bodies[len(inp.bodies)-1]
}

func (inp *Input) WriteRune(r rune) {
	inp.body().WriteRune(r)
}

func (inp *Input) LineLen() int {
	return len(inp.bodies)
}

func (inp *Input) Line(n int) string {
	if n < 0 || n >= len(inp.bodies) {
		panic("invalid line number")
	}
	if n == 0 {
		if len(inp.bodies) < 2 {
			return inp.head + inp.bodies[0].String() + inp.tail
		}
		return inp.head + inp.bodies[0].String()
	}
	if n == len(inp.bodies)-1 {
		return inp.bodies[n].String() + inp.tail
	}
	return inp.bodies[n].String()
}

func (inp *Input) Lines() []string {
	if len(inp.bodies) < 2 {
		return []string{inp.head + inp.bodies[0].String() + inp.tail}
	}
	lines := []string{inp.head + inp.bodies[0].String()}
	i := 1
	for i < len(inp.bodies)-1 {
		lines = append(lines, inp.bodies[i].String())
		i++
	}
	lines = append(lines, inp.bodies[i].String()+inp.tail)
	return lines
}

func (inp *Input) Inserted() []string {
	lines := []string{}
	for _, body := range inp.bodies {
		lines = append(lines, body.String())
	}
	return lines
}

func (inp *Input) Newline(ai bool) {
	indent := ""
	if ai {
		line := ""
		if len(inp.bodies) < 2 {
			line = inp.head + inp.bodies[0].String()
		} else {
			line = inp.bodies[len(inp.bodies)-1].String()
		}
		indent = getIndent(line)
	}
	b := termi.RuneBuf{}
	b.WriteString(indent)
	inp.bodies = append(inp.bodies, b)
	if ai {
		inp.tail = trimLeftBlanks(inp.tail)
	}
}

func (inp *Input) Column() int {
	var s string
	if len(inp.bodies) < 2 {
		s = inp.head + inp.body().String()
		return utf8.RuneCountInString(s)
	} else {
		s = inp.body().String()
	}
	return utf8.RuneCountInString(s)
}

func (inp *Input) Backspace() bool {
	if inp.body().RemoveTail() {
		return true
	}
	if len(inp.bodies) < 2 {
		return true
	}
	inp.bodies = inp.bodies[:len(inp.bodies)-1]
	return false
}

func (ed *Editor) InsertRune(r rune) {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	ed.inp.WriteRune(r)
	b := ed.Buf()
	b.Loc.Col = ed.inp.Column()
}

func (ed *Editor) Backspace() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	b := ed.Buf()
	if !ed.inp.Backspace() {
		b.Loc.Row--
	}
	b.Loc.Col = ed.inp.Column()
	// col and row are already confined
}

func (ed *Editor) InsertNewline() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	b := ed.Buf()
	ed.inp.Newline(ed.cfg.AutoIndent)
	b.Loc.Row++
	b.Loc.Col = ed.inp.Column()
	// col is already confined
	// XXX row is not confined
}
