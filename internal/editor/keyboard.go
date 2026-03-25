package editor

import (
	"tea.kareha.org/lab/levi/internal/console"
)

const Escape rune = 0x1b
const Enter rune = '\r'

type Keyboard struct{}

func NewKeyboard() Keyboard {
	return Keyboard{}
}

func (kb *Keyboard) ReadRune() rune {
	return console.ReadRune()
}
