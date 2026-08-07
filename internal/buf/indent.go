package buf

import (
	"sort"
	"strings"
)

func countSpaces(s string) int {
	i := 0
	for _, r := range s {
		if r != ' ' {
			return i
		}
		i++
	}
	return i
}

func (b *Buf) DetectIndent() {
	spaces := []int(nil)
	n := 0
	prev := 0
	for _, line := range b.Lines {
		if n >= 300 {
			break
		}
		if strings.HasPrefix(line, "\t") {
			b.Indent = "\t"
			b.IndentDetected = true
			return
		}
		i := countSpaces(line)
		delta := i - prev
		if delta >= 2 {
			spaces = append(spaces, delta)
		}
		prev = i
		n++
	}
	if len(spaces) < 1 {
		return
	}
	sort.Ints(spaces)
	i := spaces[len(spaces)/2]
	b.Indent = strings.Repeat(" ", i)
	b.IndentDetected = true
}
