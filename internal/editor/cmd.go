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

	CmdMoveToNonBlankOfNextLine
	CmdMoveToNonBlankOfPrevLine
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

	CmdPromptMoveToLine

	CmdPromptSaveAndQuit
	CmdPromptSave
	CmdPromptForceSave
	CmdPromptQuit
	CmdPromptForceQuit
	CmdPromptOpen
	CmdPromptForceOpen
	CmdPromptRead
	CmdPromptNext
	CmdPromptPrev

	CmdPromptShell

	CmdPromptSaveAll
	CmdPromptQuitAll
	CmdPromptForceQuitAll
)
