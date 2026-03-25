package editor

import (
	"tea.kareha.org/lab/levi/internal/console"
)

const Esc rune = 0x1b

type Keyboard struct{}

func NewKeyboard() Keyboard {
	return Keyboard{}
}

func (kb *Keyboard) ReadRune() rune {
	return console.ReadRune()
}
