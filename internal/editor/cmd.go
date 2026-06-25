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
	Reg       rune
	Start     buf.Loc
	End       buf.Loc
	StartRow  int
	EndRow    int
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

var MoveCmds = map[CmdKind]struct{}{
	CmdMoveLeft:  {},
	CmdMoveDown:  {},
	CmdMoveUp:    {},
	CmdMoveRight: {},

	CmdMoveToStart:    {},
	CmdMoveToEnd:      {},
	CmdMoveToNonBlank: {},
	CmdMoveToColumn:   {},

	CmdMoveByWord:              {},
	CmdMoveBackwardByWord:      {},
	CmdMoveToEndOfWord:         {},
	CmdMoveByLooseWord:         {},
	CmdMoveBackwardByLooseWord: {},
	CmdMoveToEndOfLooseWord:    {},

	CmdMoveByLine:         {},
	CmdMoveBackwardByLine: {},
	CmdMoveToLastLine:     {},
	CmdMoveToLine:         {},

	CmdMoveBySentence:          {},
	CmdMoveBackwardBySentence:  {},
	CmdMoveByParagraph:         {},
	CmdMoveBackwardByParagraph: {},
	CmdMoveBySection:           {},
	CmdMoveBackwardBySection:   {},

	CmdMoveToTopOfView:         {},
	CmdMoveToMiddleOfView:      {},
	CmdMoveToBottomOfView:      {},
	CmdMoveToBelowTopOfView:    {},
	CmdMoveToAboveBottomOfView: {},

	CmdMoveToMark:     {},
	CmdMoveToMarkLine: {},

	CmdMoveBackToMark:     {},
	CmdMoveBackToMarkLine: {},

	CmdMoveSearchForward:        {},
	CmdMoveSearchBackward:       {},
	CmdMoveSearchNextMatch:      {},
	CmdMoveSearchPrevMatch:      {},
	CmdMoveSearchRepeatForward:  {},
	CmdMoveSearchRepeatBackward: {},

	CmdMoveFindForward:        {},
	CmdMoveFindBackward:       {},
	CmdMoveFindBeforeForward:  {},
	CmdMoveFindBeforeBackward: {},
	CmdMoveFindNextMatch:      {},
	CmdMoveFindPrevMatch:      {},
}

var InsertCmds = map[CmdKind]bool{
	CmdInsertBefore:         true,
	CmdInsertAfter:          true,
	CmdInsertBeforeNonBlank: true,
	CmdInsertAfterEnd:       true,
	CmdInsertOverwrite:      true,

	CmdInsertOpenBelow: true,
	CmdInsertOpenAbove: true,

	CmdOpChangeLine:       true,
	CmdOpChangeRegion:     true,
	CmdOpChangeLineRegion: true,
	CmdOpChangeWord:       true,
	CmdOpChangeToEnd:      true,
	CmdOpSubst:            true,
	CmdOpSubstLine:        true,

	//CmdEditReplace:       true,
}

var EditCmds = map[CmdKind]bool{
	CmdOpPaste:        true,
	CmdOpPasteBefore:  true,
	CmdOpPasteFromReg: true,

	CmdOpDelete:           true,
	CmdOpDeleteBefore:     true,
	CmdOpDeleteLine:       true,
	CmdOpDeleteRegion:     true,
	CmdOpDeleteLineRegion: true,
	CmdOpDeleteWord:       true,
	CmdOpDeleteToEnd:      true,

	CmdEditJoin:          true,
	CmdEditIndent:        true,
	CmdEditOutdent:       true,
	CmdEditIndentRegion:  true,
	CmdEditOutdentRegion: true,
}

var MultiInsertCmds = map[CmdKind]bool{
	CmdInsertBefore:         true,
	CmdInsertAfter:          true,
	CmdInsertBeforeNonBlank: true,
	CmdInsertAfterEnd:       true,
	CmdInsertOverwrite:      true,

	CmdInsertOpenBelow: true,
	CmdInsertOpenAbove: true,
}
