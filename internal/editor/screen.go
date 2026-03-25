package editor

import (
	"tea.kareha.org/lab/levi/internal/console"
)

type screen struct{}

func newScreen() screen {
	return screen{}
}

func (scr *screen) size() (int, int) {
	return console.Size()
}
