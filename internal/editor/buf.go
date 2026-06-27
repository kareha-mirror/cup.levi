package editor

import (
	"fmt"
	"os"
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
)

func (ed *Editor) NewBuf() {
	b := &buf.Buf{NewFile: true}
	if ed.bufIdx < len(ed.bufs) {
		ed.bufs[ed.bufIdx] = b
	} else {
		ed.bufs = append(ed.bufs, b)
		ed.bufIdx = len(ed.bufs) - 1
	}
	ed.redraw = true
}

func (ed *Editor) Buf() *buf.Buf {
	return ed.bufs[ed.bufIdx]
}

func (ed *Editor) NumBufs() int {
	return len(ed.bufs)
}

func (ed *Editor) CheckQuit() {
	if len(ed.bufs) < 1 {
		ed.alive = false
	}
}

func (ed *Editor) Close(force bool) {
	if !force && ed.Buf().Modified {
		ed.Ring(
			"File modified since last complete write;" +
				" write or use ! to override.",
		)
		return
	}
	bufs := append([]*buf.Buf{}, ed.bufs[:ed.bufIdx]...)
	if ed.bufIdx+1 < len(ed.bufs) {
		bufs = append(bufs, ed.bufs[ed.bufIdx+1:]...)
	}
	ed.bufs = bufs
	n := len(ed.bufs)
	if ed.bufIdx >= n {
		ed.bufIdx = max(n-1, 0)
	}
}

func (ed *Editor) Load(path string, force bool) error {
	if !force && ed.Buf().Modified {
		ed.Ring(
			"File modified since last complete write;" +
				" write or use ! to override.",
		)
		return fmt.Errorf("file modified")
	}
	ed.NewBuf()
	b := ed.Buf()
	b.Path = path
	if path == "" {
		return nil
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	stamp := buf.Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	b.SetText(string(data))
	b.Stamp = stamp
	b.NewFile = false
	b.Modified = false
	return nil
}

func (ed *Editor) SaveAs(path string, force bool) error {
	if path == "" {
		ed.Ring("No filename specified")
		return fmt.Errorf("no filename specified")
	}
	info, err := os.Stat(path)
	newFile := ""
	stamp := buf.Stamp{}
	if err != nil {
		newFile = " new file:"
	} else {
		stamp = buf.Stamp{
			Time: info.ModTime(),
			Size: info.Size(),
		}
	}
	b := ed.Buf()
	if !force && path == b.Path && stamp != b.Stamp {
		ed.Ring(
			"%s: file modified more recently than this copy;"+
				" use ! to override.",
			path,
		)
		return fmt.Errorf("file modified more recently")
	}

	text := b.Text()
	err = os.WriteFile(path, []byte(text), 0666)
	if err != nil {
		return err
	}
	info, err = os.Stat(path)
	if err != nil {
		return err
	}
	stamp = buf.Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	ed.Message(
		"%s:%s %d lines, %d bytes, %d runes.",
		path, newFile, b.NumLines(), len(text), utf8.RuneCountInString(text),
	)

	if b.Path == "" {
		b.Path = path
	}
	if path == b.Path {
		b.Stamp = stamp
	}
	b.NewFile = false
	b.Modified = false
	return nil
}

func (ed *Editor) Save(force bool) error {
	return ed.SaveAs(ed.Buf().Path, force)
}
