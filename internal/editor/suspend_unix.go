//go:build unix

package editor

import (
	"fmt"
	"os"
	"syscall"

	"tea.kareha.org/cup/termi"
)

func (ed *Editor) MiscSuspend() {
	ed.EnsureCommand()

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		ed.Error("Cannot find process")
		return
	}

	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	termi.StopInput()
	fmt.Print(termi.ResetAlternate)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)
	ed.redraw = true

	ed.suspended.Store(true)
	err = p.Signal(syscall.SIGTSTP)
	if err != nil {
		ed.Error("Cannot send signal")
		ed.suspended.Store(false)
		termi.Raw()
		fmt.Print(termi.SetAlternate)
		termi.StartInput()
		return
	}

	_ = <-ed.resume
	ed.suspended.Store(false)
}
