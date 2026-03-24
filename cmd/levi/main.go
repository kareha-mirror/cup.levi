package main

import (
	"tea.kareha.org/lab/levi/internal/console"
	"tea.kareha.org/lab/levi/internal/editor"
)

func main() {
	// init
	console.Raw()
	defer func() {
		console.Cooked()
		console.ShowCursor()
	}()

	// main
	editor.New().Main()

	// cleanup
	console.Clear()
	console.HomeCursor()
}
