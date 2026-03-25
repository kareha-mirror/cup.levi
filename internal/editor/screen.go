package editor

import (
	"tea.kareha.org/lab/levi/internal/console"
)

type screen struct {
	w, h int
}

func newScreen() screen {
	w, h := console.Size()
	return screen{
		w: w,
		h: h,
	}
}

func (scr *screen) size() (int, int) {
	return scr.w, scr.h
}
