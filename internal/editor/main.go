package editor

import (
	"regexp"
	"strconv"

	"tea.kareha.org/lab/termi"
)

func (ed *Editor) ExitInsert() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	ed.lines[ed.row] = ed.ins.Line()
	ed.ins.Reset()
	ed.mode = ModeCommand
	ed.MoveLeft(1)
}

func (ed *Editor) InsertNewline() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	before := make([]string, 0, len(ed.lines)+1)
	before = append(before, ed.lines[:ed.row]...)
	var after []string
	if ed.row+1 < len(ed.lines) {
		after = ed.lines[ed.row+1:]
	} else {
		after = []string{}
	}
	lines := ed.ins.Newline()
	ed.lines = append(append(before, lines...), after...)
	ed.row++
	ed.col = 0
	// row and col are already confined
}

func (ed *Editor) Backspace() {
	if ed.mode != ModeInsert {
		panic("invalid state")
	}
	if !ed.ins.Backspace() {
		ed.Ring()
		return
	}
	ed.col--
	// col is already confined
}

var cmdRe = regexp.MustCompile("^(\\d*)([:mziaIARoOdyYxXDsScCpPrJ><\\.uUZ]*)(\\d*)([hjkl0\\$\\^|wbeWBE\\n\\+\\-G\\)\\(\\}\\{\\]\\[HML'`/?nNfFtT;,g]*)(.*?)$")
var letterRe = regexp.MustCompile("([m'`fFtT;,r])(.)$")
var letterSubRe = regexp.MustCompile("[fFtT;,]")

func (ed *Editor) Run(noNum bool, num int, op string, noSubnub bool, subnum int, mv string, letter string, replay bool) bool {
	switch op {
	case "x":
		ed.OpDelete(num)
		return true
	}

	switch op {
	case "i":
		ed.InsertBefore()
		return true
	case "a":
		ed.InsertAfter()
		return true
	case "ZZ":
		ed.MiscSaveAndQuit()
		return true
	}

	switch mv {
	case "h":
		ed.MoveLeft(num)
		return true
	case "j":
		ed.MoveDown(num)
		return true
	case "k":
		ed.MoveUp(num)
		return true
	case "l":
		ed.MoveRight(num)
		return true
	case "0":
		ed.MoveToStart()
		return true
	case "$":
		ed.MoveToEnd()
		return true
	case "^":
		ed.MoveToNonBlank()
		return true
	case "|":
		ed.MoveToColumn(num)
		return true
	}

	return false
}

func (ed *Editor) Main() {
	for !ed.quit {
		ed.Draw()

		key := termi.ReadKey()
		switch ed.mode {
		case ModeCommand:
			switch key.Kind {
			case termi.KeyRune:
				if key.Rune == termi.RuneEscape {
					ed.combuf.ClearAll()
					ed.Ring()
					continue
				}
				ed.combuf.InsertRune(key.Rune)

				comb := ed.combuf.String()
				m := cmdRe.FindStringSubmatch(comb)
				/*
					if len(m) < 1 {
						// TODO error "not (yet) a vi command [" + comb + "]"
						ed.combuf.Clear()
						continue
					}
				*/
				var numStr, op, subnumStr, mv string
				if len(m) > 0 {
					numStr, op, subnumStr, mv = m[1], m[2], m[3], m[4]
				}
				m = letterRe.FindStringSubmatch(comb)
				var letterCommand, letter string
				if len(m) > 0 {
					letterCommand, letter = m[1], m[2]
				}
				if letterCommand != "" {
					if letterCommand == "m" || letterCommand == "r" {
						op = letterCommand
						mv = ""
					} else if letterCommand == "'" || letterCommand == "`" {
						mv = letterCommand
					} else if letterSubRe.MatchString(letterCommand) {
						mv = letterCommand
					}
				}

				noNum := false
				num := 1
				if numStr == "" {
					noNum = true
				} else if numStr == "0" {
					mv = "0"
				} else {
					n, err := strconv.Atoi(numStr)
					if err != nil {
						panic(err)
					}
					num = n
				}

				noSubnum := false
				subnum := 1
				if subnumStr == "" {
					noSubnum = true
				} else if subnumStr == "0" {
					mv = "0"
				} else {
					n, err := strconv.Atoi(subnumStr)
					if err != nil {
						panic(err)
					}
					subnum = n
				}

				if ed.Run(noNum, num, op, noSubnum, subnum, mv, letter, false) {
					ed.combuf.Clear()
				}
			case termi.KeyUp:
				ed.MoveUp(1)
			case termi.KeyDown:
				ed.MoveDown(1)
			case termi.KeyRight:
				ed.MoveRight(1)
			case termi.KeyLeft:
				ed.MoveLeft(1)
			default:
				ed.Ring()
			}
		case ModeInsert:
			switch key.Kind {
			case termi.KeyRune:
				switch key.Rune {
				case termi.RuneEscape:
					ed.ExitInsert()
				case termi.RuneEnter:
					ed.InsertNewline()
				case termi.RuneBackspace:
					ed.Backspace()
				case termi.RuneDelete:
					ed.Backspace()
				default:
					ed.InsertRune(key.Rune)
				}
			case termi.KeyUp:
				ed.ExitInsert()
				ed.MoveUp(1)
			case termi.KeyDown:
				ed.ExitInsert()
				ed.MoveDown(1)
			case termi.KeyRight:
				ed.ExitInsert()
				ed.MoveRight(1)
			case termi.KeyLeft:
				ed.ExitInsert()
				ed.MoveLeft(1)
			default:
				ed.Ring()
			}
		}
	}
}
