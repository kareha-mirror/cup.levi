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
	Loc     Loc
	ViewLoc Loc
	VirtCol int // 0-based
	Pos     Pos

	Lines    []string
	Modified bool

	Path    string
	Stamp   Stamp
	NewFile bool
	CRLF    bool

	Marks   map[rune]Loc
	History History
	Depth   int
}

func (b *Buf) NumLines() int {
	return len(b.Lines)
}

// panics if row is out of range
func (b *Buf) Line(row int) string {
	// empty case
	if len(b.Lines) == 0 && row == 0 {
		return ""
	}

	return b.Lines[row]
}

// panics if row is out of range
func (b *Buf) SetLine(row int, line string) {
	// lazy init on empty case
	if len(b.Lines) == 0 && row == 0 {
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

func LineSep(crlf bool) string {
	if crlf {
		return "\r\n"
	} else {
		return "\n"
	}
}

func (b *Buf) Text(crlf bool) string {
	// empty case
	if len(b.Lines) == 0 {
		return ""
	}

	sep := LineSep(crlf)
	return strings.Join(b.Lines, sep) + sep
}

func (b *Buf) SetText(text string) {
	// empty case
	if text == "" {
		b.Lines = b.Lines[:0]
		return
	}

	// clip last newline if exists
	if text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
		b.CRLF = false
		if text != "" && text[len(text)-1] == '\r' {
			text = text[:len(text)-1]
			b.CRLF = true
		}
	} else if strings.Index(text, "\r\n") != -1 {
		b.CRLF = true
	}

	b.Lines = strings.Split(text, LineSep(b.CRLF))
}

func (b *Buf) Mark(r rune) {
	// lazy init
	if b.Marks == nil {
		b.Marks = map[rune]Loc{}
	}

	b.Marks[r] = b.Loc
}
