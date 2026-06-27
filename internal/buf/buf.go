package buf

import (
	"strings"
	"time"
)

type Loc struct {
	// order matters
	Col int // 0-based
	Row int // 0-based
}

type Pos struct {
	// order matters
	X int // 0-based
	Y int // 0-based
}

type Stamp struct {
	Time time.Time
	Size int64
}

type Buf struct {
	Loc      Loc
	ViewLoc  Loc
	VirtCol  int // 0-based
	Pos      Pos
	Lines    []string
	Path     string
	Modified bool
	Stamp    Stamp
	NewFile  bool
	Marks    map[rune]Loc
	Snapshot []string
}

func (b *Buf) NumLines() int {
	return len(b.Lines)
}

func (b *Buf) Line(row int) string {
	if len(b.Lines) < 1 {
		return ""
	}
	return b.Lines[row]
}

func (b *Buf) SetLine(row int, line string) {
	if len(b.Lines) < 1 {
		b.Lines = append(b.Lines, "")
	}
	b.Lines[row] = line
}

func (b *Buf) CurrentLine() string {
	return b.Line(b.Loc.Row)
}

func (b *Buf) SetCurrentLine(line string) {
	b.SetLine(b.Loc.Row, line)
}

func (b *Buf) Text() string {
	if len(b.Lines) < 1 {
		return ""
	}
	return strings.Join(b.Lines, "\n") + "\n"
}

func (b *Buf) SetText(text string) {
	if len(text) < 1 {
		b.Lines = b.Lines[:0]
	} else {
		// should also support CRLF or not?
		if text[len(text)-1] == '\n' {
			text = text[:len(text)-1]
		}
		b.Lines = strings.Split(text, "\n")
	}
}

func (b *Buf) Mark(r rune) {
	if b.Marks == nil {
		b.Marks = map[rune]Loc{}
	}
	b.Marks[r] = b.Loc
}
