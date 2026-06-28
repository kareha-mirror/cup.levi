package editor

type CmdKind int

type Cmd struct {
	Kind   CmdKind
	Num    int
	Letter rune
	Pat    string
}

type CmdPair struct {
	Reg  string
	Main Cmd
	Sub  Cmd
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
	CmdOpCopyWord
	CmdOpCopyToEnd

	CmdOpPaste
	CmdOpPasteBefore

	CmdOpDelete
	CmdOpDeleteBefore
	CmdOpDeleteLine
	CmdOpDeleteRegion
	CmdOpDeleteWord
	CmdOpDeleteToEnd

	CmdOpChangeLine
	CmdOpChangeRegion
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

	CmdOpChangeLine:   {},
	CmdOpChangeRegion: {},
	CmdOpChangeWord:   {},
	CmdOpChangeToEnd:  {},
	CmdOpSubst:        {},
	CmdOpSubstLine:    {},

	//CmdEditReplace: {},
}

var EditCmds = map[CmdKind]struct{}{
	CmdOpPaste:       {},
	CmdOpPasteBefore: {},

	CmdOpDelete:       {},
	CmdOpDeleteBefore: {},
	CmdOpDeleteLine:   {},
	CmdOpDeleteRegion: {},
	CmdOpDeleteWord:   {},
	CmdOpDeleteToEnd:  {},

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
