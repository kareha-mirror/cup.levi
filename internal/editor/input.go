package editor

import (
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi/rbuf"
	"tea.kareha.org/cup/termi/rkind"
	"tea.kareha.org/cup/termi/rutil"
)

type Input struct {
	head, tail string
	bodies     []rbuf.RuneBuf
	offset     int
}

func (inp *Input) Reset() {
	inp.head = ""
	inp.tail = ""
	inp.bodies = []rbuf.RuneBuf{rbuf.RuneBuf{}}
	inp.offset = 0
}

func (inp *Input) Init(line string, col int, ai bool) {
	inp.Reset()
	inp.head, inp.tail = rutil.Split(line, col)
	if ai && rkind.IsBlankLine(inp.head) {
		inp.bodies[0].WriteString(inp.head)
		inp.offset = len(inp.head)
		inp.head = ""
	}
}

func (inp *Input) body() *rbuf.RuneBuf {
	return &inp.bodies[len(inp.bodies)-1]
}

func (inp *Input) WriteRune(r rune) {
	inp.body().WriteRune(r)
}

func (inp *Input) WriteString(s string) {
	inp.body().WriteString(s)
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

func IsSpaces(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r != ' ' {
			return false
		}
	}
	return true
}

func (inp *Input) Newline(indent string, ai bool) {
	body := rbuf.RuneBuf{}
	body.WriteString(indent)
	inp.bodies = append(inp.bodies, body)
	if ai {
		inp.tail = rkind.TrimPrefixBlanks(inp.tail)
	}
}

func (inp *Input) Column() int {
	if len(inp.bodies) < 2 {
		return utf8.RuneCountInString(inp.head) + inp.body().RuneCount()
	} else {
		return inp.body().RuneCount()
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
		ed.Error("invalid state")
		return
	}
	ed.inp.WriteRune(r)
	ed.Buf().Loc.Col = ed.inp.Column()
}

func isIndent(s string) bool {
	for _, r := range s {
		if r != '\t' && r != ' ' {
			return false
		}
	}
	return true
}

func (ed *Editor) InputBackspace() {
	if ed.mode != ModeInsert {
		ed.Error("invalid state")
		return
	}

	outdented := false
	if ed.cfg.DetectIndent && isIndent(ed.inp.body().String()) {
		if len(ed.inp.body().String()) >= len(ed.Buf().Indent) {
			for range ed.Buf().Indent {
				ed.inp.Backspace()
			}
			outdented = true
		}
	}

	b := ed.Buf()
	if !outdented && !ed.inp.Backspace() {
		b.Loc.Row--
	}
	b.Loc.Col = ed.inp.Column()
	// col and row are already confined
}

func (ed *Editor) InputNewline() {
	if ed.mode != ModeInsert {
		ed.Error("invalid state")
		return
	}

	indent := ""
	if ed.cfg.AutoIndent {
		if len(ed.inp.bodies) < 2 {
			indent = rkind.IndentOf(ed.inp.head + ed.inp.bodies[0].String())
		} else {
			indent = rkind.IndentOf(
				ed.inp.bodies[len(ed.inp.bodies)-1].String(),
			)
		}
	}
	b := ed.Buf()
	if ed.cfg.DetectIndent && !b.IndentDetected {
		if strings.Contains(indent, "\t") {
			b.Indent = "\t"
			b.IndentDetected = true
		} else if IsSpaces(indent) && len(indent) >= 2 {
			b.Indent = indent
			b.IndentDetected = true
		}
	}

	ed.inp.Newline(indent, ed.cfg.AutoIndent)
	b.Loc.Row++
	b.Loc.Col = ed.inp.Column()
	// col is already confined
	// row is confined in insert mode
}

func (ed *Editor) InputTab() {
	if ed.mode != ModeInsert {
		ed.Error("invalid state")
		return
	}
	if ed.cfg.DetectIndent && isIndent(ed.inp.body().String()) {
		ed.inp.WriteString(ed.Buf().Indent)
	} else {
		ed.inp.WriteRune('\t')
	}
	ed.Buf().Loc.Col = ed.inp.Column()
}
