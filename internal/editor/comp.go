package editor

func (ed *Editor) Usage(s string) bool {
	ed.Ring("Usage: %s", s)
	return true
}
