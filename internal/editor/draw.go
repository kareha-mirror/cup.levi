package editor

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

func runeAt(s string, i int) rune {
	for _, r := range s {
		if i == 0 {
			return r
		}
		i--
	}
	return utf8.RuneError
}

func (ed *Editor) UpdateCursor() {
	b := ed.Buffer()
	current := ed.CurrentLine()
	col := b.col
	if ed.mode == ModeInsert && b.col > 0 {
		col--
	}
	b.x = 0
	dy := 0
	if current != "" {
		lines := termi.Wrap(current, ed.w, ed.mode == ModeInsert)
		for _, line := range lines {
			rc := utf8.RuneCountInString(line)
			if col < rc {
				b.x = termi.StringWidth(line, col)
				r := runeAt(line, col)
				if r == '\t' {
					b.x += termi.TabWidth - (b.x % termi.TabWidth) - 1
				}
				if ed.mode == ModeInsert && b.col > 0 {
					if termi.IsWide(r) || termi.IsEmoji(r) {
						b.x += 2
					} else {
						b.x++
					}
					if b.x > ed.w {
						if r == '\t' {
							b.x = termi.TabWidth
						} else {
							b.x = 2
						}
						dy++
					} else if b.x == ed.w {
						b.x = 0
						dy++
					}
				}
				break
			}
			col -= rc
			dy++
		}
	}

	if b.row < b.vrow {
		b.vrow = b.row
	}

	y := 0
	for i := b.vrow; i < b.row; i++ {
		lines := termi.Wrap(ed.Line(i), ed.w, false)
		y += len(lines)
	}
	b.y = y + dy

	for b.y >= ed.h-1 {
		b.vrow++

		y := 0
		for i := b.vrow; i < b.row; i++ {
			lines := termi.Wrap(ed.Line(i), ed.w, false)
			y += len(lines)
		}
		b.y = y + dy
	}
}

func (ed *Editor) DrawBuffer() {
	b := ed.Buffer()
	view := []string{}
	linesLen := max(len(b.lines), 1)
	sb := strings.Builder{}

	y := 0
	for i := b.vrow; i < linesLen+ed.inp.LineLen()-1; i++ {
		tail := i == b.row && ed.mode == ModeInsert
		lines := termi.Wrap(ed.Line(i), ed.w, tail)

		for _, line := range lines {
			sb.WriteString(termi.MoveCursor(0, y))
			if ed.colors != nil {
				if i == b.row {
					sb.WriteString(ed.colors.CurrentFg.Fg())
					sb.WriteString(ed.colors.CurrentBg.Bg())
				} else {
					sb.WriteString(ed.colors.TextFg.Fg())
					sb.WriteString(ed.colors.TextBg.Bg())
				}
			}
			sb.WriteString(termi.Render(line))
			rc := utf8.RuneCountInString(line)
			if termi.StringWidth(line, rc) < ed.w {
				sb.WriteString(termi.ClearTail)
			}
			if ed.colors != nil {
				sb.WriteString(termi.ResetAttr)
			}
			view = append(view, sb.String())
			sb.Reset()

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
		sb.WriteString(termi.MoveCursor(0, y))
		if ed.colors != nil {
			sb.WriteString(ed.colors.TextFg.Fg())
			sb.WriteString(ed.colors.TextBg.Bg())
		}
		sb.WriteString(termi.Render("~"))
		sb.WriteString(termi.ClearTail)
		if ed.colors != nil {
			sb.WriteString(termi.ResetAttr)
		}
		view = append(view, sb.String())
		sb.Reset()
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
	if ed.colors != nil {
		fmt.Print(ed.colors.StatusFg.Fg())
		fmt.Print(ed.colors.StatusBg.Bg())
	}

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
			m = "command"
		case ModeInsert:
			m = "insert"
		case ModeSearch:
			m = "search"
		default:
			m = "invalid"
		}

		fmt.Printf(
			"[%s] %s", ed.parser.Cache(), m,
		)
	}
	fmt.Print(termi.ClearTail)

	fmt.Print(termi.MoveCursor(ed.w-2, ed.h-1))
	if ed.esc {
		fmt.Print(" *")
	} else {
		fmt.Print(" .")
	}

	if ed.colors != nil {
		fmt.Print(termi.ResetAttr)
	}
}

func (ed *Editor) PlaceCursor() {
	if ed.mode == ModePrompt {
		line := ":" + ed.prompt.String()
		rc := utf8.RuneCountInString(line)
		x := termi.StringWidth(line, rc)
		fmt.Print(termi.MoveCursor(x, ed.h-1))
	} else {
		b := ed.Buffer()
		fmt.Print(termi.MoveCursor(b.x, b.y))
	}
}

func (ed *Editor) Draw() {
	ed.w, ed.h = termi.Size()
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
	ed.PlaceCursor()

	fmt.Print(termi.ShowCursor)
}
