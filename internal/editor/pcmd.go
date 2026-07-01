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
	PcmdLoad
	PcmdForceLoad
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

	PcmdOpen
	PcmdNewline
	PcmdColors

	PcmdMem   // XXX debug
	PcmdHello // XXX debug
)

var IsBufMovePcmd = map[PcmdKind]struct{}{
	PcmdNext: {},
	PcmdPrev: {},
}
