package editor

import (
	"os"
	"unicode/utf8"

	"tea.kareha.org/cup/levi/internal/buf"
)

// Creates new buffer and place it last of buffer list.
func (ed *Editor) NewBuf() {
	b := buf.New(ed.cfg.CRLF, ed.cfg.Depth)

	if ed.bufIdx < len(ed.bufs) {
		ed.bufs[ed.bufIdx] = b
	} else {
		if ed.bufIdx != len(ed.bufs) {
			ed.Error("Invalid buffer index")
		}
		ed.bufs = append(ed.bufs, b)
		ed.bufIdx = len(ed.bufs) - 1
	}

	ed.redraw = true
}

// Returns current buffer.
func (ed *Editor) Buf() *buf.Buf {
	return ed.bufs[ed.bufIdx]
}

// Returns number of buffers.
func (ed *Editor) NumBufs() int {
	return len(ed.bufs)
}

// Activates quit flag if buffer list is empty.
func (ed *Editor) CheckQuit() {
	if len(ed.bufs) < 1 {
		ed.alive = false
	}
}

// Closes current buffer and remove from buffer list.
func (ed *Editor) Close(force bool) bool {
	if !force && ed.Buf().Modified {
		ed.Ring(
			"File modified since last complete write;" +
				" write or use ! to override.",
		)
		return false
	}

	bufs := append([]*buf.Buf{}, ed.bufs[:ed.bufIdx]...)
	if ed.bufIdx+1 < len(ed.bufs) {
		bufs = append(bufs, ed.bufs[ed.bufIdx+1:]...)
	}
	ed.bufs = bufs

	if ed.lastBufIdx >= ed.bufIdx {
		ed.lastBufIdx = max(ed.lastBufIdx-1, 0)
	}

	numBufs := len(ed.bufs)
	if ed.bufIdx >= numBufs {
		ed.bufIdx = max(numBufs-1, 0)
	}

	return true
}

// Loads text from file to current buffer.
func (ed *Editor) Load(path string, force bool) bool {
	if !force && ed.Buf().Modified {
		ed.Ring(
			"File modified since last complete write;" +
				" write or use ! to override.",
		)
		return false
	}

	ed.NewBuf()
	b := ed.Buf()
	b.Path = path
	if path == "" {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	stamp := buf.Stamp{
		Time: info.ModTime(),
		Size: info.Size(),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		ed.Error("%v", err)
		return false
	}
	b.SetText(string(data))
	b.Stamp = stamp
	b.NewFile = false
	b.Marks = nil
	b.Modified = false
	b.StoreLine()
	return true
}

// Creates new buffer and loads text from file to it.
func (ed *Editor) Open(path string) bool {
	ed.lastBufIdx = ed.bufIdx

	ed.bufIdx = ed.NumBufs()
	return ed.Load(path, true)
}

// Saves current buffer to named file.
func (ed *Editor) SaveAs(path string, force bool) bool {
	if path == "" {
		ed.Ring("No filename specified")
		return false
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
		return false
	}

	text := b.Text(b.CRLF)
	err = os.WriteFile(path, []byte(text), 0666)
	if err != nil {
		ed.Error("%v", err)
		return false
	}
	info, err = os.Stat(path)
	if err != nil {
		ed.Error("%v", err)
		return false
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
	return true
}

// Saves current buffer to original file.
func (ed *Editor) Save(force bool) bool {
	return ed.SaveAs(ed.Buf().Path, force)
}
