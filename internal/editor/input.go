package editor

import (
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/rkind"
)

type Input struct {
	head, tail string
	bodies     []termi.RuneBuf
	offset     int
}

func (inp *Input) Reset() {
	inp.head = ""
	inp.tail = ""
	inp.bodies = []termi.RuneBuf{termi.RuneBuf{}}
	inp.offset = 0
}

func (inp *Input) Init(line string, col int, ai bool) {
	inp.Reset()
	rs := []rune(line)
	inp.head = string(rs[:col])
	if col < len(rs) {
		inp.tail = string(rs[col:])
	}
	if ai && rkind.IsBlankLine(inp.head) {
		inp.bodies[0].WriteString(inp.head)
		inp.offset = len(inp.head)
		inp.head = ""
	}
}

func (inp *Input) body() *termi.RuneBuf {
	return &inp.bodies[len(inp.bodies)-1]
}

func (inp *Input) WriteRune(r rune) {
	inp.body().WriteRune(r)
}

func (inp *Input) NumLines() int {
	if inp.bodies == nil {
		return 1
	}
	return len(inp.bodies)
}

func (inp *Input) Line(row int) string {
	if row < 0 || row >= len(inp.bodies) {
		panic("invalid line number")
	}
	if row == 0 {
		if len(inp.bodies) < 2 {
			return inp.head + inp.bodies[0].String() + inp.tail
		}
		return inp.head + inp.bodies[0].String()
	}
	if row == len(inp.bodies)-1 {
		return inp.bodies[row].String() + inp.tail
	}
	return inp.bodies[row].String()
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
	first := inp.bodies[0].String()
	if inp.offset < len(first) {
		first = first[inp.offset:]
	}
	lines := append([]string{}, first)
	for i := 1; i < len(inp.bodies); i++ {
		lines = append(lines, inp.bodies[i].String())
	}
	return lines
}

func (inp *Input) Newline(ai bool) {
	indent := ""
	if ai {
		if len(inp.bodies) < 2 {
			indent = rkind.IndentOf(inp.head + inp.bodies[0].String())
		} else {
			indent = rkind.IndentOf(inp.bodies[len(inp.bodies)-1].String())
		}
	}
	body := termi.RuneBuf{}
	body.WriteString(indent)
	inp.bodies = append(inp.bodies, body)
	if ai {
		inp.tail = rkind.TrimPrefixBlanks(inp.tail)
	}
}

func (inp *Input) Column() int {
	if len(inp.bodies) < 2 {
		return utf8.RuneCountInString(inp.head) + inp.body().Len()
	} else {
		return inp.body().Len()
	}
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

//
// For Editor
//

func (ed *Editor) InputWriteRune(r rune) {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	ed.inp.WriteRune(r)
	b := ed.Buf()
	b.Loc.Col = ed.inp.Column()
}

func (ed *Editor) InputBackspace() {
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

func (ed *Editor) InputNewline() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	b := ed.Buf()
	ed.inp.Newline(ed.cfg.AutoIndent)
	b.Loc.Row++
	b.Loc.Col = ed.inp.Column()
	// col is already confined
	// row is confined in insert mode
}
