//go:build unix

package editor

import (
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

var stdinFile *os.File
var stdoutFile *os.File
var stderrFile *os.File

func setup(cmd *exec.Cmd) {
	stdin, err := unix.Dup(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	stdout, err := unix.Dup(syscall.Stdout)
	if err != nil {
		panic(err)
	}
	stderr, err := unix.Dup(syscall.Stderr)
	if err != nil {
		panic(err)
	}
	stdinFile = os.NewFile(uintptr(stdin), "(stdin)")
	stdoutFile = os.NewFile(uintptr(stdout), "(stdout)")
	stderrFile = os.NewFile(uintptr(stderr), "(stderr)")
	cmd.Stdin = stdinFile
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile
}

func terminate() {
	stdinFile.Close()
	stdoutFile.Close()
	stderrFile.Close()
}
