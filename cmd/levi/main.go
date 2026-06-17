package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"tea.kareha.org/cup/levi/internal/editor"
)

const appName = "levi"

func fatal(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func getConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		fatal(err)
	}
	return filepath.Join(dir, appName)
}

func main() {
	configDir := flag.String("d", "", "config directory")
	flag.Parse()
	if *configDir == "" {
		*configDir = getConfigDir()
	}
	args := flag.Args()

	ed, err := editor.Init(*configDir, args)
	if err != nil {
		fatal(err)
	}
	defer func() {
		err := ed.Finish()
		if err != nil {
			fatal(err)
		}
	}()

	ed.Main()
}
