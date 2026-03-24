package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

func enableRawMode() (*term.State, error) {
	return term.MakeRaw(int(os.Stdin.Fd()))
}

func disableRawMode(state *term.State) {
	term.Restore(int(os.Stdin.Fd()), state)
}

func clearScreen() {
	fmt.Print("\x1b[2J")
}

func goHome() {
	fmt.Print("\x1b[H")
}

func moveCursor(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
}

func getScreenSize() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	return w, h
}

type Screen struct {
	w, h int
}

func NewScreen() Screen {
	w, h := getScreenSize()
	return Screen{
		w: w,
		h: h,
	}
}

func (scr *Screen) Size() (int, int) {
	return scr.w, scr.h
}

const Esc rune = 0x1b

func runeSize(b byte) int {
	switch {
	case b & 0x80 == 0:
		return 1
	case b & 0xe0 == 0xc0:
		return 2
	case b & 0xf0 == 0xe0:
		return 3
	case b & 0xf8 == 0xf0:
		return 4
	default:
		return -1 // invalid
	}
}

func readRune() rune {
	buf := make([]byte, 1)
	_, err := io.ReadFull(os.Stdin, buf)
	if err != nil {
		panic(err)
	}
	size := runeSize(buf[0])
	if size == -1 {
		panic("Invalid UTF-8")
	}
	full := make([]byte, size)
	full[0] = buf[0]
	if size > 1 {
		_, err := io.ReadFull(os.Stdin, full[1:])
		if err != nil {
			panic(err)
		}
	}
	r, size := utf8.DecodeRune(full)
	if r == utf8.RuneError && size == 1 {
		panic("Invalid UTF-8")
	}
	return r
}

type Keyboard struct {}

func NewKeyboard() Keyboard {
	return Keyboard{}
}

func (kb *Keyboard) ReadRune() rune {
	return readRune()
}

type Editor struct {
	scr *Screen
	kb *Keyboard
	x, y int
	line *strings.Builder
}

func NewEditor(scr *Screen, kb *Keyboard) Editor {
	_, h := scr.Size()
	return Editor{
		scr: scr,
		x: 0,
		y: h / 2,
		line: new(strings.Builder),
	}
}

func (ed *Editor) Screen() *Screen {
	return ed.scr
}

func (ed *Editor) AddRune(r rune) {
	ed.line.WriteRune(r)
}

func draw(ed *Editor) {
	clearScreen()
	goHome()

	fmt.Print("Hit Esc to Exit")

	moveCursor(ed.x, ed.y)
	fmt.Print(ed.line.String())
}

func main() {
	// init
	oldState, err := enableRawMode()
	if err != nil {
		panic(err)
	}
	defer func() {
		disableRawMode(oldState)
		showCursor()
	}()

	// main
	scr := NewScreen()
	kb := NewKeyboard()
	ed := NewEditor(&scr, &kb)
	for {
		hideCursor()
		draw(&ed)
		showCursor()

		r := kb.ReadRune()
		if r == Esc {
			break
		}
		ed.AddRune(r)
	}

	// cleanup
	clearScreen()
	goHome()
}
