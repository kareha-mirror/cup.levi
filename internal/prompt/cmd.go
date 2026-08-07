package prompt

type Kind int

type Cmd struct {
	Kind Kind
	Num  int
	Name string
}

const (
	Invalid Kind = iota

	MoveByLine
	MoveBackByLine
	MoveToLine

	Save
	SaveAndClose
	WriteAndClose
	Write
	ForceWrite
	Close
	ForceClose
	Load
	ForceLoad
	Read
	Next
	Prev

	Shell

	SaveAll
	ForceSaveAll
	CloseAll
	ForceCloseAll

	TabStop
	AutoIndent
	NoAutoIndent

	Open
	Newline
	Colors
	DetectIndent
	NoDetectIndent

	Mem   // XXX debug
	Hello // XXX debug

	Ring
	Error
)

var IsBufMove = map[Kind]struct{}{
	Next: {},
	Prev: {},
}
