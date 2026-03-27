package editor

type Combuf struct {
	buf   []rune
	cache string
}

const maxCombufLen = 256

func NewCombuf() *Combuf {
	return &Combuf{
		buf:   make([]rune, 0),
		cache: "",
	}
}

func (cb *Combuf) String() string {
	return string(cb.buf)
}

func (cb *Combuf) InsertRune(r rune) {
	cb.buf = append(cb.buf, r)
	cb.cache = cb.String()
}

func (cb *Combuf) Clear() {
	if len(cb.buf) > maxCombufLen {
		cb.buf = make([]rune, 0)
	} else {
		cb.buf = cb.buf[:0]
	}
}

func (cb *Combuf) Cache() string {
	return cb.cache
}

func (cb *Combuf) ClearAll() {
	cb.Clear()
	cb.cache = ""
}
