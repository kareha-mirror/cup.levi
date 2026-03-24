package console

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

var state *term.State

func Raw() {
	if state != nil {
		term.Restore(int(os.Stdin.Fd()), state)
		state = nil
	}
	s, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	state = s
}

func Cooked() {
	if state == nil {
		panic("state is nil")
	}
	term.Restore(int(os.Stdin.Fd()), state)
}

func Clear() {
	fmt.Print("\x1b[2J")
}

func HomeCursor() {
	fmt.Print("\x1b[H")
}

func MoveCursor(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

func HideCursor() {
	fmt.Print("\x1b[?25l")
}

func ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func Size() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	return w, h
}

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

func ReadRune() rune {
	buf := make([]byte, 1)
	_, err := io.ReadFull(os.Stdin, buf)
	if err != nil {
		panic(err)
	}
	expected := runeSize(buf[0])
	if expected == -1 {
		panic("Invalid UTF-8 head")
	}
	full := make([]byte, expected)
	full[0] = buf[0]
	if expected > 1 {
		_, err := io.ReadFull(os.Stdin, full[1:])
		if err != nil {
			panic(err)
		}
	}
	r, size := utf8.DecodeRune(full)
	if r == utf8.RuneError && size == 1 {
		panic("Invalid UTF-8 body")
	}
	return r
}

func Print(s string) {
	fmt.Print(s)
}

func Printf(format string, a ...any) (n int, err error) {
	return fmt.Printf(format, a...)
}
