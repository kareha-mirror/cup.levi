package editor

import (
	"fmt"
)

type Msg struct {
	message string
	ring    string
}

func (msg *Msg) Message(format string, a ...any) {
	msg.message = fmt.Sprintf(format, a...)
}

func (msg *Msg) Ring(format string, a ...any) {
	msg.ring = fmt.Sprintf(format, a...)
}

func (msg *Msg) Error(format string, a ...any) {
	msg.Ring("Error: "+format, a...)
}

func (ed *Editor) Message(format string, a ...any) {
	ed.msg.Message(format, a...)
}

func (ed *Editor) Ring(format string, a ...any) {
	ed.msg.Ring(format, a...)
}

func (ed *Editor) Error(format string, a ...any) {
	ed.msg.Error(format, a...)
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

func (ed *Editor) InitialInfo() {
	b := ed.Buf()
	path := b.Path
	if path == "" {
		path = "(memory)"
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
