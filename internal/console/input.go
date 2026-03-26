package console

import (
	"io"
	"os"
	"unicode/utf8"
)

type Key int

const (
	KeyNormal = iota
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
)

const RuneEscape rune = 0x1b
const RuneEnter rune = '\r'
const RuneBackspace rune = '\b'
const RuneDelete rune = 0x7f

var buf []rune = make([]rune, 0)

func runeSize(b byte) int {
	switch {
	case b&0x80 == 0:
		return 1
	case b&0xe0 == 0xc0:
		return 2
	case b&0xf0 == 0xe0:
		return 3
	case b&0xf8 == 0xf0:
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

func ReadKey() (Key, rune) {
	if len(buf) > 0 {
		r := buf[0]
		buf = buf[1:]
		return KeyNormal, r
	}
	r := readRune()
	if r != RuneEscape {
		return KeyNormal, r
	}
	r2 := readRune()
	if r2 != '[' {
		buf = append(buf, r2)
		return KeyNormal, r
	}
	r3 := readRune()
	switch r3 {
	case 'A':
		return KeyUp, 0
	case 'B':
		return KeyDown, 0
	case 'C':
		return KeyRight, 0
	case 'D':
		return KeyLeft, 0
	}
	buf = append(buf, r2)
	buf = append(buf, r3)
	return KeyNormal, r
}
