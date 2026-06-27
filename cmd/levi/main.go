package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"tea.kareha.org/cup/levi/internal/editor"
)

const appName = "levi"

const failure = 1

func realMain() (totalErr error) {
	// parse options
	cfgDir := flag.String("d", "", "config directory")
	flag.Parse()

	if *cfgDir == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		// default config directory
		*cfgDir = filepath.Join(dir, appName)
	}

	paths := flag.Args()

	// init editor
	ed, err := editor.Init(*cfgDir, paths)
	if err != nil {
		return err
	}
	defer func() {
		if err := ed.Finish(); err != nil {
			totalErr = errors.Join(totalErr, err)
		}
	}()

	// enter main loop
	if err := ed.Main(); err != nil {
		return err
	}

	return totalErr
}

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(failure)
	}
}
