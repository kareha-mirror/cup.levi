package buf

import (
	"sort"
	"strings"
)

func countStartingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r != ' ' {
			return n
		}
		n++
	}
	return n
}

func (b *Buf) DetectIndent() {
	list := []int(nil)
	i := 0
	prev := 0
	for _, line := range b.Lines {
		if i >= 300 {
			break
		}
		if strings.HasPrefix(line, "\t") {
			b.Indent = "\t"
			b.IndentDetected = true
			return
		}
		n := countStartingSpaces(line)
		delta := n - prev
		if delta >= 2 {
			list = append(list, delta)
		}
		prev = n
		i++
	}
	if len(list) < 1 {
		return
	}
	sort.Ints(list)
	n := list[len(list)/2]
	b.Indent = strings.Repeat(" ", n)
	b.IndentDetected = true
}
