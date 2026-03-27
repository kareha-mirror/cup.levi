package editor

// h : Move cursor left by character.
func (ed *Editor) MoveLeft(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col -= n
	ed.Confine()
}

// j : Move cursor down by line.
func (ed *Editor) MoveDown(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.row += n
	ed.Confine()
}

// k : Move cursor up by line.
func (ed *Editor) MoveUp(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.row -= n
	ed.Confine()
}

// l : Move cursor right by character.
func (ed *Editor) MoveRight(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col += n
	ed.Confine()
}

// 0 : Move cursor to start of current line.
func (ed *Editor) MoveToStart() {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col = 0
	// col is already confined
}

// $ : Move cursor to end of current line.
func (ed *Editor) MoveToEnd() {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col = ed.RuneCount() - 1
	ed.Confine()
}

// ^ : Move cursor to first non-blank character of current line.
func (ed *Editor) MoveToNonBlank() {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	line := ed.CurrentLine()
	i := 0
	for _, r := range line {
		if r != ' ' && r != '\t' {
			break
		}
		i++
	}
	ed.col = i
	ed.Confine()
}

// <num>| : Move cursor to column <num> of current line.
// (Note: Proper vi's column number is visual-based, but levi' is rune-based.)
func (ed *Editor) MoveToColumn(n int) {
	if ed.mode != ModeCommand {
		panic("invalid state")
	}
	ed.col = n - 1
	ed.Confine()
}
