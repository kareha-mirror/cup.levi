package editor

import (
	"tea.kareha.org/lab/levi/internal/console"
)

const runeEscape rune = 0x1b
const runeEnter rune = '\r'

type key int

const (
	keyNormal = iota
	keyUp
	keyDown
	keyRight
	keyLeft
)

type keyboard struct {
	buf []rune
}

func newKeyboard() keyboard {
	return keyboard{
		buf: make([]rune, 0),
	}
}

func (kb *keyboard) readKey() (key, rune) {
	if len(kb.buf) > 0 {
		r := kb.buf[0]
		kb.buf = kb.buf[1:]
		return keyNormal, r
	}
	r := console.ReadRune()
	if r != runeEscape {
		return keyNormal, r
	}
	r2 := console.ReadRune()
	if r2 != '[' {
		kb.buf = append(kb.buf, r2)
		return keyNormal, r
	}
	r3 := console.ReadRune()
	switch r3 {
	case 'A':
		return keyUp, 0
	case 'B':
		return keyDown, 0
	case 'C':
		return keyRight, 0
	case 'D':
		return keyLeft, 0
	}
	kb.buf = append(kb.buf, r2)
	kb.buf = append(kb.buf, r3)
	return keyNormal, r
}
