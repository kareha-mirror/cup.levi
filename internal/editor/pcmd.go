package editor

type PcmdKind int

type Pcmd struct {
	Kind PcmdKind
	Num  int
	Name string
}

const (
	PcmdInvalid PcmdKind = iota

	PcmdMoveByLine
	PcmdMoveBackwardByLine
	PcmdMoveToLine

	PcmdSaveAndQuit
	PcmdSave
	PcmdForceSave
	PcmdQuit
	PcmdForceQuit
	PcmdOpen
	PcmdForceOpen
	PcmdRead
	PcmdNext
	PcmdPrev

	PcmdShell

	PcmdSaveAll
	PcmdQuitAll
	PcmdForceQuitAll

	PcmdTabStop
	PcmdAutoIndent
	PcmdNoAutoIndent

	PcmdColors

	PcmdHello // XXX debug
)
