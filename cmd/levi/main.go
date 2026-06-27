package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"tea.kareha.org/cup/levi/internal/editor"
)

const appName = "levi"

const failure = 1

func fatal(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(failure)
	// never returns
}

func defaultConfigDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		fatal(err)
	}
	return filepath.Join(dir, appName)
}

func main() {
	// parse options
	cfgDir := flag.String("d", "", "config directory")
	flag.Parse()
	if *cfgDir == "" {
		*cfgDir = defaultConfigDir()
	}
	paths := flag.Args()

	// init editor
	ed, err := editor.Init(*cfgDir, paths)
	if err != nil {
		fatal(err)
	}
	defer func() {
		if err := ed.Finish(); err != nil {
			fatal(err)
		}
	}()

	// enter main loop
	ed.Main()
}
