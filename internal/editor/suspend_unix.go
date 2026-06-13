//go: build unix

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
	termi.Cooked()
	fmt.Print(termi.ShowCursor)
	fmt.Print(termi.ResetAlternate)
	ed.redraw = true

	err = p.Signal(syscall.SIGTSTP)
	if err != nil {
		ed.Error("Cannot send signal")
		fmt.Print(termi.SetAlternate)
		termi.Raw()
		return
	}

	fmt.Print(termi.SetAlternate)
	termi.Raw()
}
