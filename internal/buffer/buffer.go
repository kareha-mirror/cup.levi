package buffer

import (
	"strings"
	"time"
)

type Loc struct {
	Col int // 0-based
	Row int // 0-based
}

type Pos struct {
	X int // 0-based
	Y int // 0-based
}

type Stamp struct {
	Time time.Time
	Size int64
}

type Buffer struct {
	Loc      Loc
	ViewLoc  Loc
	VirtCol  int // 0-based
	Pos      Pos
	Lines    []string
	Path     string
	Modified bool
	Stamp    Stamp
	Marks    map[rune]Loc
}

func (b *Buffer) NumLines() int {
	return len(b.Lines)
}

func (b *Buffer) Line(row int) string {
	if len(b.Lines) < 1 {
		return ""
	}
	return b.Lines[row]
}

func (b *Buffer) SetLine(row int, line string) {
	if len(b.Lines) < 1 {
		b.Lines = []string{""}
	}
	b.Lines[row] = line
}

func (b *Buffer) CurrentLine() string {
	if len(b.Lines) < 1 {
		return ""
	}
	return b.Lines[b.Loc.Row]
}

func (b *Buffer) SetCurrentLine(line string) {
	b.SetLine(b.Loc.Row, line)
}

func (b *Buffer) Text() string {
	if len(b.Lines) < 1 {
		return ""
	}
	return strings.Join(b.Lines, "\n") + "\n"
}

func (b *Buffer) SetText(text string) {
	if len(text) < 1 {
		b.Lines = []string{}
	} else {
		// TODO should also support CRLF or not?
		if text[len(text)-1] == '\n' {
			text = text[:len(text)-1]
		}
		b.Lines = strings.Split(text, "\n")
	}
}

func (b *Buffer) Mark(r rune) {
	if b.Marks == nil {
		b.Marks = map[rune]Loc{}
	}
	b.Marks[r] = b.Loc
}
