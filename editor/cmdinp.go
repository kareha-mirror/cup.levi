package editor

type CmdInput struct {
	buf  []rune
	last []rune
}

func (in *CmdInput) String() string {
	return string(in.buf)
}

func (in *CmdInput) LastString() string {
	return string(in.last)
}

func (in *CmdInput) WriteRune(r rune) {
	in.buf = append(in.buf, r)
	in.last = append([]rune{}, in.buf...)
}

func (in *CmdInput) Backspace() bool {
	if len(in.buf) < 1 {
		return false
	}
	in.buf = in.buf[:len(in.buf)-1]
	in.last = append([]rune{}, in.buf...)
	return true
}

func (in *CmdInput) Reset() {
	in.buf = in.buf[:0]
}

func (in *CmdInput) ResetAll() {
	in.Reset()
	in.last = in.last[:0]
}
