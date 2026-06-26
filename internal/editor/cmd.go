package editor

import (
	"tea.kareha.org/cup/levi/internal/buf"
)

type CmdKind int

type Cmd struct {
	Kind      CmdKind
	Num       int
	Letter    rune
	Pat       string
	Reg       string
	Start     buf.Loc
	End       buf.Loc
	Inclusive bool
}

const (
	CmdInvalid CmdKind = iota

	CmdMoveLeft
	CmdMoveDown
	CmdMoveUp
	CmdMoveRight

	CmdMoveToStart
	CmdMoveToEnd
	CmdMoveToNonBlank
	CmdMoveToColumn

	CmdMoveByWord
	CmdMoveByWordEx // XXX debug
	CmdMoveBackwardByWord
	CmdMoveToEndOfWord
	CmdMoveByLooseWord
	CmdMoveBackwardByLooseWord
	CmdMoveToEndOfLooseWord

	CmdMoveByLine
	CmdMoveBackwardByLine
	CmdMoveToLastLine
	CmdMoveToLine

	CmdMoveBySentence
	CmdMoveBackwardBySentence
	CmdMoveByParagraph
	CmdMoveBackwardByParagraph
	CmdMoveBySection
	CmdMoveBackwardBySection

	CmdMoveToTopOfView
	CmdMoveToMiddleOfView
	CmdMoveToBottomOfView
	CmdMoveToBelowTopOfView
	CmdMoveToAboveBottomOfView

	CmdMarkSet
	CmdMoveToMark
	CmdMoveToMarkLine

	CmdMoveBackToMark
	CmdMoveBackToMarkLine

	CmdViewDown
	CmdViewUp
	CmdViewDownHalf
	CmdViewUpHalf
	CmdViewDownLine
	CmdViewUpLine

	CmdViewToTop
	CmdViewToMiddle
	CmdViewToBottom

	CmdViewRedraw

	CmdMoveSearchForward
	CmdMoveSearchBackward
	CmdMoveSearchNextMatch
	CmdMoveSearchPrevMatch
	CmdMoveSearchRepeatForward
	CmdMoveSearchRepeatBackward

	CmdMoveFindForward
	CmdMoveFindBackward
	CmdMoveFindBeforeForward
	CmdMoveFindBeforeBackward
	CmdMoveFindNextMatch
	CmdMoveFindPrevMatch

	CmdInsertBefore
	CmdInsertAfter
	CmdInsertBeforeNonBlank
	CmdInsertAfterEnd
	CmdInsertOverwrite

	CmdInsertOpenBelow
	CmdInsertOpenAbove

	CmdOpCopyLine
	CmdOpCopyRegion
	CmdOpCopyLineRegion
	CmdOpCopyWord
	CmdOpCopyToEnd
	CmdOpCopyLineIntoReg

	CmdOpPaste
	CmdOpPasteBefore
	CmdOpPasteFromReg
	CmdOpPasteBeforeFromReg

	CmdOpDelete
	CmdOpDeleteBefore
	CmdOpDeleteLine
	CmdOpDeleteRegion
	CmdOpDeleteLineRegion
	CmdOpDeleteWord
	CmdOpDeleteToEnd

	CmdOpChangeLine
	CmdOpChangeRegion
	CmdOpChangeLineRegion
	CmdOpChangeWord
	CmdOpChangeToEnd
	CmdOpSubst
	CmdOpSubstLine

	CmdEditReplace
	CmdEditJoin
	CmdEditIndent
	CmdEditOutdent
	CmdEditIndentRegion
	CmdEditOutdentRegion

	CmdMiscShowInfo
	CmdMiscRepeat
	CmdMiscUndo
	CmdMiscRestore
	CmdMiscSaveAndQuit
	CmdMiscSuspend
)

var InsertCmds = map[CmdKind]struct{}{
	CmdInsertBefore:         {},
	CmdInsertAfter:          {},
	CmdInsertBeforeNonBlank: {},
	CmdInsertAfterEnd:       {},
	CmdInsertOverwrite:      {},

	CmdInsertOpenBelow: {},
	CmdInsertOpenAbove: {},

	CmdOpChangeLine:       {},
	CmdOpChangeRegion:     {},
	CmdOpChangeLineRegion: {},
	CmdOpChangeWord:       {},
	CmdOpChangeToEnd:      {},
	CmdOpSubst:            {},
	CmdOpSubstLine:        {},

	//CmdEditReplace: {},
}

var EditCmds = map[CmdKind]struct{}{
	CmdOpPaste:              {},
	CmdOpPasteBefore:        {},
	CmdOpPasteFromReg:       {},
	CmdOpPasteBeforeFromReg: {},

	CmdOpDelete:           {},
	CmdOpDeleteBefore:     {},
	CmdOpDeleteLine:       {},
	CmdOpDeleteRegion:     {},
	CmdOpDeleteLineRegion: {},
	CmdOpDeleteWord:       {},
	CmdOpDeleteToEnd:      {},

	CmdEditJoin:          {},
	CmdEditIndent:        {},
	CmdEditOutdent:       {},
	CmdEditIndentRegion:  {},
	CmdEditOutdentRegion: {},
}

var MultiInsertCmds = map[CmdKind]struct{}{
	CmdInsertBefore:         {},
	CmdInsertAfter:          {},
	CmdInsertBeforeNonBlank: {},
	CmdInsertAfterEnd:       {},
	CmdInsertOverwrite:      {},

	CmdInsertOpenBelow: {},
	CmdInsertOpenAbove: {},
}
