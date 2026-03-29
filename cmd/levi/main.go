package main

import (
	"os"

	"tea.kareha.org/cup/levi/internal/editor"
)

func main() {
	ed := editor.Init(os.Args)
	defer ed.Finish()
	ed.Main()
}
