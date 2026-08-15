package buf

import (
	"sort"
	"strings"
)

const sampleLines = 300

func countStartingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r != ' ' {
			break
		}
		n++
	}
	return n
}

func (b *Buf) DetectIndent() {
	tabs := 0
	spaces := []int(nil)

	i := 0
	prev := 0
	for _, line := range b.Lines {
		if i >= sampleLines {
			break
		}
		if strings.HasPrefix(line, "\t") {
			tabs++
			continue
		}
		n := countStartingSpaces(line)
		delta := n - prev
		if delta >= 2 {
			spaces = append(spaces, delta)
		}
		prev = n
		i++
	}

	if len(spaces) < 1 {
		if tabs < 1 {
			return
		}
		b.Indent = "\t"
		b.IndentDetected = true
		return
	}

	sort.Ints(spaces)
	n := spaces[len(spaces)/2]
	count := 0
	for _, k := range spaces {
		if k == n {
			count++
		}
	}

	if tabs >= count {
		b.Indent = "\t"
	} else {
		b.Indent = strings.Repeat(" ", n)
	}
	b.IndentDetected = true
}
