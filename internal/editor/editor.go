package editor

import (
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/lab/levi/internal/console"
	"tea.kareha.org/lab/levi/internal/util"
)

type mode int

const (
	modeCommand = iota
	modeInsert
)

type Editor struct {
	scr        *screen
	kb         *keyboard
	col, row   int
	x, y       int
	vrow       int
	lines      []string
	head, tail string
	insert     *strings.Builder
	mode       mode
	path       string
}

func Init(args []string) *Editor {
	var path string
	var lines []string
	if len(args) > 1 {
		path = args[1]
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		if len(data) > 0 {
			if data[len(data)-1] == '\n' {
				data = data[:len(data)-1]
			}
			lines = strings.Split(string(data), "\n")
		}
	}
	if len(lines) < 1 {
		lines = make([]string, 1)
	}

	console.Raw()

	scr := newScreen()
	kb := newKeyboard()

	return &Editor{
		scr:    &scr,
		kb:     &kb,
		col:    0,
		row:    0,
		x:      0,
		y:      0,
		vrow:   0,
		lines:  lines,
		head:   "",
		tail:   "",
		insert: new(strings.Builder),
		mode:   modeCommand,
		path:   path,
	}
}

func (ed *Editor) Finish() {
	console.Clear()
	console.HomeCursor()
	console.Cooked()
	console.ShowCursor()

	if ed.path != "" {
		text := strings.Join(ed.lines, "\n") + "\n"
		err := ioutil.WriteFile(ed.path, []byte(text), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func (ed *Editor) runeCount() int {
	return utf8.RuneCountInString(ed.lines[ed.row])
}

func (ed *Editor) lineHeight(line string) int {
	w, _ := ed.scr.size()
	rc := utf8.RuneCountInString(line)
	width := util.StringWidth(line, rc)
	return 1 + max(width-1, 0)/w
}

func (ed *Editor) drawBuffer() {
	_, h := ed.scr.size()

	y := 0
	for i := ed.vrow; i < len(ed.lines); i++ {
		var line string
		if ed.mode == modeInsert && i == ed.row {
			line = ed.head + ed.insert.String() + ed.tail
		} else {
			line = ed.lines[i]
		}

		console.MoveCursor(0, y)
		util.Print(line)

		y += ed.lineHeight(line)
		if y >= h-1 {
			break
		}
	}

	for ; y < h-1; y++ {
		console.MoveCursor(0, y)
		util.Print("~")
	}
}

func (ed *Editor) drawStatus() {
	_, h := ed.scr.size()

	console.MoveCursor(0, h-1)
	switch ed.mode {
	case modeCommand:
		util.Print("c")
	case modeInsert:
		util.Print("i")
	}
}

func (ed *Editor) updateCursor() {
	w, h := ed.scr.size()

	var dy int
	switch ed.mode {
	case modeCommand:
		ed.row = min(max(ed.row, 0), max(len(ed.lines)-1, 0))
		len := ed.runeCount()
		ed.col = min(ed.col, max(len-1, 0))

		// XXX approximation
		width := util.StringWidth(ed.lines[ed.row], ed.col)
		ed.x = width % w
		dy = width / w
	case modeInsert:
		// XXX approximation
		width := util.StringWidth(ed.head+ed.insert.String(), ed.col)
		ed.x = width % w
		dy = width / w
	}

	if ed.row < ed.vrow {
		ed.vrow = ed.row
	}

	y := 0
	for i := ed.vrow; i < ed.row; i++ {
		y += ed.lineHeight(ed.lines[i])
	}
	ed.y = y + dy

	for ed.y >= h-1 {
		ed.vrow++

		y := 0
		for i := ed.vrow; i < ed.row; i++ {
			y += ed.lineHeight(ed.lines[i])
		}
		ed.y = y + dy
	}
}

func (ed *Editor) repaint() {
	console.HideCursor()

	console.Clear()
	console.HomeCursor()

	ed.updateCursor()

	ed.drawBuffer()
	ed.drawStatus()

	console.MoveCursor(ed.x, ed.y)

	console.ShowCursor()
}

func (ed *Editor) exitInsert() {
	ed.lines[ed.row] = ed.head + ed.insert.String() + ed.tail
	ed.tail = ""
	ed.insert.Reset()
	ed.mode = modeCommand
	ed.moveLeft(1)
}

func (ed *Editor) insertNewline() {
	before := make([]string, 0, len(ed.lines)+1)
	before = append(before, ed.lines[:ed.row]...)
	var after []string
	if ed.row+1 < len(ed.lines) {
		after = ed.lines[ed.row+1:]
	} else {
		after = []string{}
	}
	newLines := []string{
		ed.head + ed.insert.String(),
		ed.tail,
	}
	ed.lines = append(append(before, newLines...), after...)

	ed.row++

	ed.col = 0
	ed.head = ""
	ed.insert.Reset()
}

func (ed *Editor) deleteBefore() {
	if ed.insert.Len() < 1 {
		return // TODO ring
	}
	insert := ed.insert.String()
	_, size := utf8.DecodeLastRuneInString(insert)
	insert = insert[:len(insert)-size]
	ed.insert.Reset()
	ed.insert.WriteString(insert)
	ed.col--
}

func (ed *Editor) insertRune(r rune) {
	ed.insert.WriteRune(r)
	ed.col++
}

func (ed *Editor) Main() {
	for {
		ed.repaint()

		k, r := ed.kb.readKey()
		switch ed.mode {
		case modeCommand:
			switch k {
			case keyNormal:
				switch r {
				case 'q':
					return
				case 'i':
					ed.enterInsert()
				case 'a':
					ed.enterInsertAfter()
				case 'h':
					ed.moveLeft(1)
				case 'l':
					ed.moveRight(1)
				case 'j':
					ed.moveDown(1)
				case 'k':
					ed.moveUp(1)
				}
			case keyUp:
				ed.moveUp(1)
			case keyDown:
				ed.moveDown(1)
			case keyRight:
				ed.moveRight(1)
			case keyLeft:
				ed.moveLeft(1)
			default:
				// TODO ring
			}
		case modeInsert:
			switch k {
			case keyNormal:
				switch r {
				case runeEscape:
					ed.exitInsert()
				case runeEnter:
					ed.insertNewline()
				case runeBackspace:
					ed.deleteBefore()
				case runeDelete:
					ed.deleteBefore()
				default:
					ed.insertRune(r)
				}
			case keyUp:
				ed.exitInsert()
				ed.moveUp(1)
			case keyDown:
				ed.exitInsert()
				ed.moveDown(1)
			case keyRight:
				ed.exitInsert()
				ed.moveRight(1)
			case keyLeft:
				ed.exitInsert()
				ed.moveLeft(1)
			default:
				// TODO ring
			}
		}
	}
}
