package runekind

func IsBlank(r rune) bool {
	return r == ' ' || r == '\t'
}

func IsWord(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		r == '_' ||
		(r >= '0' && r <= '9') ||
		(r >= 'A' && r <= 'Z')
}

func IsSymbol(r rune) bool {
	return !IsBlank(r) && !IsWord(r) && r < 0x80
}

func IsOther(r rune) bool {
	return r >= 0x80
}

type RuneKind int

const (
	Blank RuneKind = iota
	Word
	Symbol
	Other
)

func Kind(r rune) RuneKind {
	switch {
	case IsBlank(r):
		return Blank
	case IsWord(r):
		return Word
	case IsSymbol(r):
		return Symbol
	default:
		return Other
	}
}
