package editor

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

func (ed *Editor) DrawBuffer() {
	view := []string{}

	linesLen := max(len(ed.lines), 1)

	y := 0
	for i := ed.vrow; i < linesLen+ed.inp.LineLen()-1; i++ {
		tail := i == ed.row && ed.mode == ModeInsert
		lines := termi.Wrap(ed.Line(i), ed.w, tail)

		for _, line := range lines {
			b := strings.Builder{}
			b.WriteString(termi.MoveCursor(0, y))
			b.WriteString(termi.Render(line))
			b.WriteString(termi.ClearTail)
			view = append(view, b.String())

			y++
			if y >= ed.h-1 {
				break
			}
		}

		if y >= ed.h-1 {
			break
		}
	}

	for ; y < ed.h-1; y++ {
		b := strings.Builder{}
		b.WriteString(termi.MoveCursor(0, y))
		b.WriteString(termi.Render("~"))
		b.WriteString(termi.ClearTail)
		view = append(view, b.String())
	}

	for i, line := range view {
		if i < len(ed.view) && line == ed.view[i] {
			continue
		}
		fmt.Print(line)
	}
	ed.view = view
}

func (ed *Editor) DrawStatus() {
	fmt.Print(termi.MoveCursor(0, ed.h-1))

	if ed.ring != "" {
		fmt.Print(termi.SetInvert)
		fmt.Print(ed.ring)
		fmt.Print(termi.ResetInvert)
		ed.ring = ""
	} else if ed.message != "" {
		fmt.Print(ed.message)
		ed.message = ""
	} else if ed.mode == ModePrompt {
		fmt.Printf(":%s", ed.prompt.String())
	} else {
		var m string
		switch ed.mode {
		case ModeCommand:
			m = "vi command"
		case ModeInsert:
			m = "vi insert"
		case ModeSearch:
			m = "vi search"
		default:
			panic("invalid mode")
		}

		fmt.Printf(
			"[%s] %s %d,%d %s",
			ed.parser.Cache(), m, ed.row+1, ed.col+1, ed.path,
		)
	}
	fmt.Print(termi.ClearTail)

	fmt.Print(termi.MoveCursor(ed.w-2, ed.h-1))
	if ed.esc {
		fmt.Print(" *")
	} else {
		fmt.Print(" .")
	}

	fmt.Print(termi.MoveCursor(ed.x, ed.y))
}

func runeAt(s string, i int) rune {
	runes := []rune(s)
	if i < len(runes) {
		return runes[i]
	}
	return utf8.RuneError
}

func (ed *Editor) UpdateCursor() {
	lines := termi.Wrap(ed.CurrentLine(), ed.w, ed.mode == ModeInsert)
	col := ed.col
	if ed.mode == ModeInsert && ed.col > 0 {
		col--
	}
	ed.x = 0
	dy := 0
	if lines[0] != "" {
		for i := 0; i < len(lines); i++ {
			rc := utf8.RuneCountInString(lines[i])
			if col < rc {
				ed.x = termi.StringWidth(lines[i], col)
				if runeAt(lines[i], col) == '\t' {
					ed.x += termi.TabWidth - (ed.x % termi.TabWidth) - 1
				}
				if ed.mode == ModeInsert && ed.col > 0 {
					r := runeAt(lines[i], col)
					if termi.IsWide(r) || termi.IsEmoji(r) {
						ed.x += 2
					} else {
						ed.x++
					}
					if ed.x == ed.w {
						ed.x = 0
						dy++
					} else if ed.x > ed.w {
						ed.x = 2
						dy++
					}
				}
				break
			}
			col -= rc
			dy++
		}
	}

	if ed.row < ed.vrow {
		ed.vrow = ed.row
	}

	y := 0
	for i := ed.vrow; i < ed.row; i++ {
		lines := termi.Wrap(ed.Line(i), ed.w, false)
		y += len(lines)
	}
	ed.y = y + dy

	for ed.y >= ed.h-1 {
		ed.vrow++

		y := 0
		for i := ed.vrow; i < ed.row; i++ {
			lines := termi.Wrap(ed.Line(i), ed.w, false)
			y += len(lines)
		}
		ed.y = y + dy
	}
}

func (ed *Editor) Draw() {
	w, h := termi.Size()
	ed.w = w
	ed.h = h

	fmt.Print(termi.HideCursor)

	if ed.redraw {
		ed.view = []string{}
		fmt.Print(termi.Clear)
		ed.redraw = false
	}
	fmt.Print(termi.HomeCursor)

	ed.UpdateCursor()

	ed.DrawBuffer()
	ed.DrawStatus()

	//fmt.Print(termi.MoveCursor(ed.x, ed.y)) // already in ed.DrawStatus()

	fmt.Print(termi.ShowCursor)
}
