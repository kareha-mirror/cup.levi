package buf

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
	Context Loc

	snaps snaps
	Depth int

	stored string
}

func New(crlf bool, depth int) *Buf {
	return &Buf{
		NewFile: true,
		CRLF:    crlf,
		Depth:   depth,
	}
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
		b.Lines = []string{""}
	}

	b.Lines[row] = line
}

func (b *Buf) CurrentLine() string {
	return b.Line(b.Loc.Row)
}

func (b *Buf) SetCurrentLine(line string) {
	b.SetLine(b.Loc.Row, line)
	// empty case
	if b.NumLines() == 1 && b.Line(0) == "" {
		b.Lines = nil
	}
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

func HasCRLF(text string) bool {
	return strings.Contains(text, "\r\n")
}

func TextToLines(text string) []string {
	// empty case
	if text == "" {
		return nil
	}

	crlf := HasCRLF(text)
	sep := LineSep(crlf)
	if strings.HasSuffix(text, sep) {
		text = text[:len(text)-len(sep)]
	}

	return strings.Split(text, sep)
}

func (b *Buf) SetText(text string) {
	b.CRLF = HasCRLF(text)
	b.Lines = TextToLines(text)
}

func (b *Buf) Mark(r rune) {
	// lazy init
	if b.Marks == nil {
		b.Marks = map[rune]Loc{}
	}

	b.Marks[r] = b.Loc
}
