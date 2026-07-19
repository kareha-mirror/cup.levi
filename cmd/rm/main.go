package main

import "os"

// go run ./cmd/rm simulates Unix command rm -rf to be used in Taskfile.
func main() {
	for _, path := range os.Args[1:] {
		_ = os.RemoveAll(path)
	}
}
