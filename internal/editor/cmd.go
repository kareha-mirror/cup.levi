package editor

type CmdKind int

type Cmd struct {
	Kind CmdKind
	Num  int
	Ltr  rune
	Pat  string
}

type CmdPair struct {
	Reg string
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
	MoveUp
	MoveRight

	MoveToStart
	MoveToEnd
	MoveToNonBlank
	MoveToColumn

	MoveByWord
	MoveByChangeWord
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

	SearchForward
	SearchBackward
	SearchNext
	SearchPrev
	RepeatSearchForward
	RepeatSearchBackward

	FindForward
	FindBackward
	FindBeforeForward
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

	InsertBefore
	InsertAfter
	InsertBeforeNonBlank
	InsertAfterEnd
	Overwrite

	OpenBelow
	OpenAbove

	ChangeLine
	ChangeRegion
	ChangeWord
	ChangeToEnd
	Subst
	SubstLine

	//
	// Edit Commands
	//

	// Edit commands are commands which change text content of buffer.
	// They are identified by IsEditCmd set.

	Paste
	PasteBefore

	Delete
	DeleteBefore
	DeleteLine
	DeleteRegion
	DeleteWord
	DeleteToEnd

	Replace
	Join
	Indent
	Outdent
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

	CopyLine
	CopyRegion
	CopyWord
	CopyToEnd

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
	// Miscellaneous commands
	//

	ShowInfo
	Repeat
	Undo
	SaveAndQuit
	Suspend
)

var IsInsertCmd = map[CmdKind]struct{}{
	InsertBefore:         {},
	InsertAfter:          {},
	InsertBeforeNonBlank: {},
	InsertAfterEnd:       {},
	Overwrite:            {},

	OpenBelow: {},
	OpenAbove: {},

	ChangeLine:   {},
	ChangeRegion: {},
	ChangeWord:   {},
	ChangeToEnd:  {},
	Subst:        {},
	SubstLine:    {},
}

var IsMultiInsertCmd = map[CmdKind]struct{}{
	InsertBefore:         {},
	InsertAfter:          {},
	InsertBeforeNonBlank: {},
	InsertAfterEnd:       {},
	Overwrite:            {},

	OpenBelow: {},
	OpenAbove: {},
}

var IsEditCmd = map[CmdKind]struct{}{
	Paste:       {},
	PasteBefore: {},

	Delete:       {},
	DeleteBefore: {},
	DeleteLine:   {},
	DeleteRegion: {},
	DeleteWord:   {},
	DeleteToEnd:  {},

	Join:          {},
	Indent:        {},
	Outdent:       {},
	IndentRegion:  {},
	OutdentRegion: {},

	Replace: {},

	Restore: {},
}
