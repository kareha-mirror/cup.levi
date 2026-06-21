package editor

///////////////////
// View Commands //
///////////////////

//
// Scroll by View Height / Scroll by Line
//

// Ctrl-F : Scroll down by view height.
func (ed *Editor) ViewDown(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewDown")
}

// Ctrl-B : Scroll up by view height.
func (ed *Editor) ViewUp(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewUp")
}

// Ctrl-D : Scroll down by half view height.
func (ed *Editor) ViewDownHalf(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewDownHalf")
}

// Ctrl-U : Scroll up by half view height.
func (ed *Editor) ViewUpHalf(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewUpHalf")
}

// Ctrl-Y : Scroll down by line.
func (ed *Editor) ViewDownLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewDownLine")
}

// Ctrl-E : Scroll up by line.
func (ed *Editor) ViewUpLine(n int) {
	ed.EnsureCommand()
	ed.Unimplemented("ViewUpLine")
}

//
// Reposition
//

// z Enter : Reposition cursor line to top of view.
func (ed *Editor) ViewToTop() {
	ed.EnsureCommand()
	ed.Unimplemented("ViewToTop")
}

// z. : Reposition cursor line middle of view.
func (ed *Editor) ViewToMiddle() {
	ed.EnsureCommand()
	ed.Unimplemented("ViewToMiddle")
}

// z- : Reposition cursor line bottom of view.
func (ed *Editor) ViewToBottom() {
	ed.EnsureCommand()
	ed.Unimplemented("ViewToBottom")
}

//
// Redraw
//

// Ctrl-L : Redraw view.
func (ed *Editor) ViewRedraw() {
	ed.EnsureCommand()
	ed.redraw = true
}
