package editor

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"tea.kareha.org/cup/termi"
)

type Msg struct {
	messages []string
	rings    []string
	view     []string
}

func (msg *Msg) IsSingle() bool {
	return len(msg.view) == 1
}

func (msg *Msg) IsMulti() bool {
	return len(msg.view) > 1
}

func (msg *Msg) Reset() {
	msg.view = msg.view[:0]
}

func splitTextIntoLines(text string) []string {
	if len(text) < 1 {
		return nil
	}
	if text[len(text)-1] == '\n' {
		text = text[:len(text)-1]
	}
	return strings.Split(text, "\n")
}

func (msg *Msg) Message(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	lines := splitTextIntoLines(text)
	msg.messages = append(msg.messages, lines...)
}

func (msg *Msg) Ring(format string, a ...any) {
	text := fmt.Sprintf(format, a...)
	lines := splitTextIntoLines(text)
	msg.rings = append(msg.rings, lines...)
}

func (msg *Msg) Error(format string, a ...any) {
	msg.Ring("Error: "+format, a...)
}

func (ed *Editor) RenderMsg(ring bool) {
	var source *[]string
	if ring {
		source = &ed.msg.rings
	} else {
		source = &ed.msg.messages
	}
	sb := strings.Builder{}

	for _, s := range *source {
		lines := termi.Wrap(s, ed.w, false)
		for _, line := range lines {
			if ring {
				sb.WriteString(termi.SetInvert)
			}
			sb.WriteString(line)
			if ring {
				sb.WriteString(termi.ResetInvert)
			}
			rc := utf8.RuneCountInString(line)
			if termi.StringWidth(line, rc) < ed.w {
				sb.WriteString(termi.ClearTail)
			}
			ed.msg.view = append(ed.msg.view, sb.String())
			sb.Reset()
		}
	}

	*source = nil
}

func (ed *Editor) Message(format string, a ...any) {
	ed.msg.Message(format, a...)
	ed.RenderMsg(false)
}

func (ed *Editor) Ring(format string, a ...any) {
	ed.msg.Ring(format, a...)
	ed.RenderMsg(true)
}

func (ed *Editor) Error(format string, a ...any) {
	ed.msg.Error(format, a...)
	ed.RenderMsg(true)
}

func (ed *Editor) Notice(format string, a ...any) {
	if ed.cfg.Silent {
		return
	}
	ed.Message("("+format+")", a...)
}

func (ed *Editor) Unimplemented(name string) {
	ed.Ring("not implemented (" + name + ")")
}

func (ed *Editor) ShowFileInfo() {
	b := ed.Buf()
	path := b.Path
	if path == "" {
		path = "(memory)"
	}
	if b.NewFile {
		ed.Message("%s: new file: line 1", path)
		return
	}
	modified := "unmodified"
	if b.Modified {
		modified = "modified"
	}
	info := "empty file"
	numLines := b.NumLines()
	if numLines > 0 {
		info = fmt.Sprintf("line %d", b.Loc.Row+1)
	}
	ed.Message("%s: %s: %s", path, modified, info)
}

// show statistical info of inserted text
func (ed *Editor) ShowStatOfInserted() {
	numLines := len(ed.inserted)
	numBytes := numLines - 1
	numRunes := numLines - 1
	if ed.Buf().CRLF {
		numBytes *= 2
		numRunes *= 2
	}
	for _, line := range ed.inserted {
		numBytes += len(line)
		numRunes += utf8.RuneCountInString(line)
	}
	if numBytes > 0 {
		if numLines < 2 {
			ed.Notice(
				"%d bytes, %d runes inserted",
				numBytes, numRunes,
			)
		} else {
			ed.Notice(
				"%d lines, %d bytes, %d runes inserted",
				numLines, numBytes, numRunes,
			)
		}
	}
}
