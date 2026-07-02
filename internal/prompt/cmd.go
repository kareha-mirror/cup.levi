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
	MoveBackwardByLine
	MoveToLine

	SaveAndQuit
	Save
	ForceSave
	Quit
	ForceQuit
	Load
	ForceLoad
	Read
	Next
	Prev

	Shell

	SaveAll
	QuitAll
	ForceQuitAll

	TabStop
	AutoIndent
	NoAutoIndent

	Open
	Newline
	Colors

	Mem   // XXX debug
	Hello // XXX debug

	Ring
	Error
)

var IsBufMove = map[Kind]struct{}{
	Next: {},
	Prev: {},
}
