package cmd

type Kind int

type Cmd struct {
	Kind Kind
	Num  int
	Rune rune
	Pat  string
}

type Pair struct {
	Reg rune
	Op  Cmd
	Mv  Cmd
}

const (
	Invalid Kind = iota

	//
	// Motion Commands
	//

	// Motion commands move cursor.
	// They themself don't change text content of buffer.

	MoveLeft
	MoveDown
	MoveHere
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
	// Insert commands which have multiplication number
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
	// They possibly move cursor, but are not recognized as motion commands.

	ViewDown
	ViewUp
	ViewDownHalf
	ViewUpHalf
	ViewDownLine
	ViewUpLine

	ViewToTop
	ViewToMiddle
	ViewToBottom

	//
	// Miscellaneous Commands
	//

	ShowInfo
	Redraw
	Repeat
	Undo
	SaveAndClose
	Suspend

	// Select Current Buffer

	LastBuf
	GoToBuf
	NextBuf
	PrevBuf

	// For Compatibility

	Ring
)

var IsInsert = map[Kind]struct{}{
	Insert:            {},
	InsertAfter:       {},
	InsertAfterIndent: {},
	InsertAfterEnd:    {},

	InsertLine:      {},
	InsertLineAbove: {},

	ChangeRegion: {},
	Subst:        {},
}

var IsMultiInsert = map[Kind]struct{}{
	Insert:            {},
	InsertAfter:       {},
	InsertAfterIndent: {},
	InsertAfterEnd:    {},

	InsertLine:      {},
	InsertLineAbove: {},
}

var IsEdit = map[Kind]struct{}{
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

var IsModifying = map[Kind]struct{}{}

func init() {
	for c := range IsInsert {
		IsModifying[c] = struct{}{}
	}
	for c := range IsEdit {
		IsModifying[c] = struct{}{}
	}
}

var IsBufMove = map[Kind]struct{}{
	LastBuf: {},
	GoToBuf: {},
	NextBuf: {},
	PrevBuf: {},
}
