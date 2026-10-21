package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"tea.kareha.org/cup/levi/editor"
	"tea.kareha.org/cup/termi/lock"
)

const (
	appName = "levi"
	failure = 1
)

func realMain() (totalErr error) {
	var cfgDir string
	var unlock bool
	var clear bool

	// parse options
	flag.StringVar(&cfgDir, "d", "", "config directory")
	flag.BoolVar(&unlock, "unlock", false, "unlock shared lock")
	flag.BoolVar(&clear, "clear", false, "clear shared registers")
	flag.Parse()

	if cfgDir == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		// default config directory
		cfgDir = filepath.Join(dir, appName)
	}

	quickExit := false

	if unlock {
		err := lock.ForceUnlock(cfgDir)
		if err != nil {
			return err
		}
		fmt.Println("Shared lock has been removed.")
		quickExit = true
	}

	if clear {
		err := editor.ClearSharedRegisters(cfgDir)
		if err != nil {
			return err
		}
		fmt.Println("Shared registers have been cleared.")
		quickExit = true
	}

	if quickExit {
		// quick exit when unlock or clear is issued
		return nil
	}

	paths := flag.Args()

	// initialize editor
	ed, err := editor.Init(cfgDir, paths)
	if err != nil {
		return err
	}
	defer func() {
		if err := ed.Finish(); err != nil {
			totalErr = errors.Join(totalErr, err)
		}
	}()

	// enter main loop of editor
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
