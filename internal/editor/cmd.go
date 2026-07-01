package editor

type CmdKind int

type Cmd struct {
	Kind CmdKind
	Num  int
	Rune rune
	Pat  string
}

type CmdPair struct {
	Reg rune
	Op  Cmd
	Mv  Cmd
}

const (
	InvalidCmd CmdKind = iota

	//
	// Motion Commands
	//

	// Motion commands move cursor.
	// They themself don't change text content of buffer.

	MoveLeft
	MoveDown
	MoveHere // debug
	MoveUp
	MoveRight

	MoveToStart
	MoveToEnd
	MoveToAfterIndent
	MoveToColumn

	MoveByWord
	MoveByChangeWord
	MoveByDeleteWord
	MoveBackwardByWord
	MoveToEndOfWord
	MoveByLooseWord
	MoveBackwardByLooseWord
	MoveToEndOfLooseWord

	MoveByLine
	MoveBackwardByLine
	MoveToLastLine
	MoveToLine

	MoveBySentence
	MoveBackwardBySentence
	MoveByParagraph
	MoveBackwardByParagraph
	MoveBySection
	MoveBackwardBySection

	MoveToTopOfView
	MoveToMiddleOfView
	MoveToBottomOfView
	MoveToBelowTopOfView
	MoveToAboveBottomOfView

	MoveToMark
	MoveToMarkLine

	BackToMark
	BackToMarkLine

	Search
	SearchBackward
	SearchNext
	SearchPrev
	RepeatSearch
	RepeatBackwardSearch

	Find
	FindBackward
	FindBefore
	FindBeforeBackward
	FindNext
	FindPrev

	//
	// Insert Commands
	//

	// Insert commands are commands which transit to insert mode.
	// They are identified by IsInsertCmd.
	// Insert commands which can have multiplexer number
	// are identified by IsMultiInsertCmd.

	Insert
	InsertAfter
	InsertAfterIndent
	InsertAfterEnd

	InsertLine
	InsertLineAbove

	ChangeRegion
	Subst

	Overwrite // unsupported

	//
	// Edit Commands
	//

	// Edit commands are commands which change text content of buffer.
	// They are identified by IsEditCmd set.

	Paste
	PasteBefore

	Delete
	DeleteBefore
	DeleteRegion

	Replace
	Join
	IndentRegion
	OutdentRegion

	Restore

	//
	// Mark Commands
	//

	// Most other mark commands are categorized to motion commands.

	Mark

	//
	// Copy Commands
	//

	// These commands copy lines or runes into registers.
	// They don't change text content of buffer.

	CopyRegion

	//
	// View Commands
	//

	// View commands scroll screen.
	// They possibly move cursor, but are not used as motion commands.

	ViewDown
	ViewUp
	ViewDownHalf
	ViewUpHalf
	ViewDownLine
	ViewUpLine

	ViewToTop
	ViewToMiddle
	ViewToBottom

	Redraw

	//
	// Miscellaneous Commands
	//

	ShowInfo
	Repeat
	Undo
	SaveAndClose
	Suspend

	LastBuf
	GoToBuf
	NextBuf
	PrevBuf
)

var IsInsertCmd = map[CmdKind]struct{}{
	Insert:            {},
	InsertAfter:       {},
	InsertAfterIndent: {},
	InsertAfterEnd:    {},

	InsertLine:      {},
	InsertLineAbove: {},

	ChangeRegion: {},
	Subst:        {},
}

var IsMultiInsertCmd = map[CmdKind]struct{}{
	Insert:            {},
	InsertAfter:       {},
	InsertAfterIndent: {},
	InsertAfterEnd:    {},

	InsertLine:      {},
	InsertLineAbove: {},
}

var IsEditCmd = map[CmdKind]struct{}{
	Paste:       {},
	PasteBefore: {},

	Delete:       {},
	DeleteBefore: {},
	DeleteRegion: {},

	Join:          {},
	IndentRegion:  {},
	OutdentRegion: {},

	Replace: {},

	Restore: {},
}

var IsModifyingCmd = map[CmdKind]struct{}{}

func init() {
	for c := range IsInsertCmd {
		IsModifyingCmd[c] = struct{}{}
	}
	for c := range IsEditCmd {
		IsModifyingCmd[c] = struct{}{}
	}
}

var IsBufMoveCmd = map[CmdKind]struct{}{
	LastBuf: {},
	GoToBuf: {},
	NextBuf: {},
	PrevBuf: {},
}
