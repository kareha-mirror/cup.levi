package buf

import (
	"fmt"
	"unicode/utf8"
)

func (b *Buf) Info() string {
	path := b.Path
	if path == "" {
		path = "(memory)"
	}
	modified := "unmodified"
	if b.Modified {
		modified = "modified"
	}
	info := "empty file"
	numLines := b.NumLines()
	if numLines > 0 {
		numBytes := numLines
		numRunes := numLines
		if b.CRLF {
			numBytes *= 2
			numRunes *= 2
		}
		for _, line := range b.Lines {
			numBytes += len(line)
			numRunes += utf8.RuneCountInString(line)
		}
		info = fmt.Sprintf(
			"line %d of %d [%d%%] %d bytes, %d runes.",
			b.Loc.Row+1, numLines, 100*(b.Loc.Row+1)/numLines,
			numBytes, numRunes,
		)
	}
	return fmt.Sprintf("%s: %s: %s", path, modified, info)
}
