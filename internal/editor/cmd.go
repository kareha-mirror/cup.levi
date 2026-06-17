package editor

type CmdKind int

type Loc struct {
	Row int
	Col int
}

type Cmd struct {
	Kind     CmdKind
	Num      int
	Letter   rune
	Pat      string
	Reg      rune
	Start    Loc
	End      Loc
	StartRow int
	EndRow   int
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
	CmdMarkMoveTo
	CmdMarkMoveToLine

	CmdMarkBack
	CmdMarkBackToLine

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

	CmdSearchForward
	CmdSearchBackward
	CmdSearchNextMatch
	CmdSearchPrevMatch
	CmdSearchRepeatForward
	CmdSearchRepeatBackward

	CmdFindForward
	CmdFindBackward
	CmdFindBeforeForward
	CmdFindBeforeBackward
	CmdFindNextMatch
	CmdFindPrevMatch

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

var RepeatableCmds = map[CmdKind]bool{
	CmdInsertBefore:         true,
	CmdInsertAfter:          true,
	CmdInsertBeforeNonBlank: true,
	CmdInsertAfterEnd:       true,
	CmdInsertOverwrite:      true,

	CmdInsertOpenBelow: true,
	CmdInsertOpenAbove: true,

	CmdOpCopyLine:        true,
	CmdOpCopyRegion:      true,
	CmdOpCopyLineRegion:  true,
	CmdOpCopyWord:        true,
	CmdOpCopyToEnd:       true,
	CmdOpCopyLineIntoReg: true,

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

	CmdOpChangeLine:       true,
	CmdOpChangeRegion:     true,
	CmdOpChangeLineRegion: true,
	CmdOpChangeWord:       true,
	CmdOpChangeToEnd:      true,
	CmdOpSubst:            true,
	CmdOpSubstLine:        true,

	CmdEditReplace:       true,
	CmdEditJoin:          true,
	CmdEditIndent:        true,
	CmdEditOutdent:       true,
	CmdEditIndentRegion:  true,
	CmdEditOutdentRegion: true,

	CmdMiscUndo: true,
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
