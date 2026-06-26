package editor

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buf"
)

type ViewMeta struct {
	Loc buf.Loc
}

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
	b := ed.Buf()
	if b.Loc.Row-b.ViewLoc.Row >= ed.h-1 {
		b.ViewLoc.Row = b.Loc.Row - (ed.h - 2)
		b.ViewLoc.Col = 0
	}

	current := ed.Line(b.Loc.Row)
	col := b.Loc.Col
	if ed.mode == ModeInsert && b.Loc.Col > 0 {
		col--
	}
	b.Pos.X = 0
	dy := 0
	first := true
	colSum := 0
	if current != "" {
		lines := termi.Wrap(current, ed.w, ed.mode == ModeInsert)
		for _, line := range lines {
			rc := utf8.RuneCountInString(line)
			if first && colSum < b.ViewLoc.Col {
				colSum += rc
				col -= rc
				continue
			}
			first = false
			if col < rc {
				b.Pos.X = termi.StringWidth(line, col)
				r := runeAt(line, col)
				if r == '\t' {
					b.Pos.X += termi.TabWidth - (b.Pos.X % termi.TabWidth) - 1
				}
				if ed.mode == ModeInsert && b.Loc.Col > 0 {
					if termi.IsWide(r) || termi.IsEmoji(r) {
						b.Pos.X += 2
					} else {
						b.Pos.X++
					}
					if b.Pos.X > ed.w {
						if r == '\t' {
							b.Pos.X = termi.TabWidth
						} else {
							b.Pos.X = 2
						}
						dy++
					} else if b.Pos.X == ed.w {
						b.Pos.X = 0
						dy++
					}
				}
				break
			}
			col -= rc
			dy++
		}
	}

	if b.Loc.Row < b.ViewLoc.Row {
		b.ViewLoc.Row = b.Loc.Row
		b.ViewLoc.Col = 0
	}

	y := 0
	first = true
	colSum = 0
	for i := b.ViewLoc.Row; i < b.Loc.Row; i++ {
		lines := termi.Wrap(ed.Line(i), ed.w, false)
		for _, line := range lines {
			rc := utf8.RuneCountInString(line)
			if first && colSum < b.ViewLoc.Col {
				colSum += rc
				continue
			}
			first = false
			y++
		}
	}
	b.Pos.Y = y + dy

	first = true
	colSum = 0
	for b.Pos.Y >= ed.h-1 {
		b.ViewLoc.Row++
		b.ViewLoc.Col = 0

		y := 0
		for i := b.ViewLoc.Row; i < b.Loc.Row; i++ {
			lines := termi.Wrap(ed.Line(i), ed.w, false)
			for _, line := range lines {
				rc := utf8.RuneCountInString(line)
				if first && colSum < b.ViewLoc.Col {
					colSum += rc
					continue
				}
				first = false
				y++
			}
		}
		b.Pos.Y = y + dy
	}
}

func (ed *Editor) renderBuffer(
	viewLoc buf.Loc, real bool,
) ([]string, []ViewMeta) {
	b := ed.Buf()
	view := []string{}
	viewMeta := []ViewMeta{}
	numLines := max(b.NumLines(), 1)
	sb := strings.Builder{}

	y := 0
	first := true
	for i := viewLoc.Row; i < numLines+ed.inp.LineLen()-1; i++ {
		tail := i == b.Loc.Row && ed.mode == ModeInsert
		lines := termi.Wrap(ed.Line(i), ed.w, tail)

		col := 0
		for _, line := range lines {
			if real {
				sb.WriteString(termi.MoveCursor(0, y))
				if ed.colors != nil {
					if i == b.Loc.Row {
						sb.WriteString(ed.colors.Current.Seq())
					} else {
						sb.WriteString(ed.colors.Buffer.Seq())
					}
				}
				sb.WriteString(termi.Render(line))
			}
			rc := utf8.RuneCountInString(line)
			if first && viewLoc.Col > col {
				col += rc
				continue
			}
			first = false
			if real {
				if termi.StringWidth(line, rc) < ed.w {
					sb.WriteString(termi.ClearTail)
				}
				if ed.colors != nil {
					sb.WriteString(termi.ResetAttr)
				}
				view = append(view, sb.String())
				sb.Reset()
			}
			loc := buf.Loc{col, i}
			viewMeta = append(viewMeta, ViewMeta{loc})
			col += rc

			y++
			if y >= ed.h-1 {
				break
			}
		}

		if y >= ed.h-1 {
			break
		}
	}

	if real {
		for ; y < ed.h-1; y++ {
			sb.WriteString(termi.MoveCursor(0, y))
			if ed.colors != nil {
				sb.WriteString(ed.colors.Border.Seq())
			}
			sb.WriteString(termi.Render("~"))
			sb.WriteString(termi.ClearTail)
			if ed.colors != nil {
				sb.WriteString(termi.ResetAttr)
			}
			view = append(view, sb.String())
			sb.Reset()
		}
	}

	return view, viewMeta
}

func (ed *Editor) RenderBuffer(viewLoc buf.Loc) ([]string, []ViewMeta) {
	return ed.renderBuffer(viewLoc, true)
}

func (ed *Editor) RenderMeta(loc buf.Loc) []ViewMeta {
	_, viewMeta := ed.renderBuffer(loc, false)
	return viewMeta
}

func (ed *Editor) DrawBuffer() {
	b := ed.Buf()
	view, viewMeta := ed.RenderBuffer(b.ViewLoc)
	for i, line := range view {
		if i < len(ed.view) && line == ed.view[i] {
			continue
		}
		fmt.Print(line)
	}
	ed.view = view
	ed.viewMeta = viewMeta
}

func (ed *Editor) DrawStatus() {
	fmt.Print(termi.MoveCursor(0, ed.h-1))
	if ed.colors != nil {
		fmt.Print(ed.colors.Status.Seq())
	}

	if ed.msg.ring != "" {
		fmt.Print(termi.SetInvert)
		fmt.Print(ed.msg.ring)
		fmt.Print(termi.ResetInvert)
		ed.msg.ring = ""
	} else if ed.msg.message != "" {
		fmt.Print(ed.msg.message)
		ed.msg.message = ""
	} else if ed.mode == ModePrompt {
		fmt.Printf(":%s", ed.prompt.String())
	} else if ed.mode == ModeSearch {
		head := "/"
		if ed.search.backward {
			head = "?"
		}
		fmt.Printf("%s%s", head, ed.search.pattern.String())
	} else {
		mode := ""
		switch ed.mode {
		case ModeCommand:
			mode = "command"
		case ModeInsert:
			mode = "insert"
		}
		fmt.Printf("(%s)%s", mode, ed.parser.Cache())
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
	switch ed.mode {
	case ModePrompt:
		line := ":" + ed.prompt.String()
		rc := utf8.RuneCountInString(line)
		x := termi.StringWidth(line, rc)
		fmt.Print(termi.MoveCursor(x, ed.h-1))
	case ModeSearch:
		line := "/" + ed.search.pattern.String() // "/" or "?"
		rc := utf8.RuneCountInString(line)
		x := termi.StringWidth(line, rc)
		fmt.Print(termi.MoveCursor(x, ed.h-1))
	default:
		b := ed.Buf()
		fmt.Print(termi.MoveCursor(b.Pos.X, b.Pos.Y))
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
