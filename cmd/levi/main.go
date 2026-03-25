package main

import (
	"tea.kareha.org/lab/levi/internal/editor"
)

func main() {
	ed := editor.Init()
	defer ed.Finish()
	ed.Main()
}
