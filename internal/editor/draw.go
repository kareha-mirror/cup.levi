package editor

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"

	"tea.kareha.org/cup/levi/internal/buf"
	"tea.kareha.org/cup/levi/internal/rutil"
)

type ViewMeta struct {
	Loc buf.Loc
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
				r := rutil.RuneAt(line, col)
				if r == '\t' {
					b.Pos.X +=
						termi.TabWidth - (b.Pos.X % termi.TabWidth) - 1
				}
				if ed.mode == ModeInsert && b.Loc.Col > 0 {
					if r == '\t' {
						b.Pos.X++
					} else {
						b.Pos.X += termi.RuneWidth(r)
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
	view := []string{}
	viewMeta := []ViewMeta{}
	b := ed.Buf()
	numLines := max(b.NumLines(), 1)
	sb := strings.Builder{}

	y := 0
	first := true
	for row := viewLoc.Row; row < numLines+ed.inp.NumLines()-1; row++ {
		tail := row == b.Loc.Row && ed.mode == ModeInsert
		lines := termi.Wrap(ed.Line(row), ed.w, tail)

		col := 0
		for _, line := range lines {
			if real {
				sb.WriteString(termi.MoveCursor(0, y))
				if ed.colors != nil {
					if row == b.Loc.Row {
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
			loc := buf.Loc{col, row}
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
	view, viewMeta := ed.RenderBuffer(ed.Buf().ViewLoc)
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

	if ed.msg.IsSingle() {
		fmt.Print(ed.msg.view[0])
		ed.msg.Reset()
	} else if ed.msg.IsMulti() {
		numLines := len(ed.msg.view)
		y := ed.h - numLines - 1
		skip := 0
		if y < 0 {
			skip = -y
			y = 0
		}
		fmt.Print(termi.MoveCursor(0, y))
		fmt.Print(termi.SetInvert)
		fmt.Print("+=+=+=+=+=+=+=+")
		fmt.Print(termi.ResetInvert)
		fmt.Print(termi.ClearTail)
		y++
		for i := skip; i < len(ed.msg.view); i++ {
			fmt.Print(termi.MoveCursor(0, y))
			fmt.Printf("%s", ed.msg.view[i])
			y++
		}
		ed.msg.Reset()
		ed.redraw = true
	} else if ed.mode == ModePrompt {
		fmt.Print(termi.Render(":" + ed.prompt.String()))
	} else if ed.mode == ModeSearch {
		head := "/"
		if ed.searchs.backward {
			head = "?"
		}
		fmt.Print(termi.Render(
			fmt.Sprintf("%s%s", head, ed.searchs.pattern.String()),
		))
	} else if !ed.cfg.Silent {
		mode := ""
		switch ed.mode {
		case ModeCommand:
			mode = "levi"
		case ModeInsert:
			mode = "ins:"
		}

		seq := ed.parser.Cache
		cursor := ""
		if !ed.cmdOk {
			cursor = "_"
		}
		sep := ""
		code := ed.args.Code()
		if len(code) > 0 {
			sep = " : "
		}
		end := ""
		if ed.cmdOk {
			end = "."
		}

		fmt.Print(termi.Render(
			fmt.Sprintf("(%s)%s%s%s%s%s", mode, seq, cursor, sep, code, end),
		))
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
		line := termi.Render(":" + ed.prompt.String())
		rc := utf8.RuneCountInString(line)
		x := termi.StringWidth(line, rc)
		fmt.Print(termi.MoveCursor(x, ed.h-1))
	case ModeSearch:
		// "/" or "?"
		line := termi.Render("/" + ed.searchs.pattern.String())
		rc := utf8.RuneCountInString(line)
		x := termi.StringWidth(line, rc)
		fmt.Print(termi.MoveCursor(x, ed.h-1))
	default:
		p := ed.Buf().Pos
		fmt.Print(termi.MoveCursor(p.X, p.Y))
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
