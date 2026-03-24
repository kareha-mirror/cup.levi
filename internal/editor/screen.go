package editor

import (
	"tea.kareha.org/lab/levi/internal/console"
)

type Screen struct {
	w, h int
}

func NewScreen() Screen {
	w, h := console.Size()
	return Screen{
		w: w,
		h: h,
	}
}

func (scr *Screen) Size() (int, int) {
	return scr.w, scr.h
}
