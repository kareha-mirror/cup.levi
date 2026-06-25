package editor

type MoveMeta struct {
	Linewise  bool
	FreeCol   bool
	Inclusive bool
	Locate    bool
}

var MoveMetas = map[CmdKind]MoveMeta{
	CmdMoveLeft:  {},
	CmdMoveDown:  {Linewise: true, FreeCol: true},
	CmdMoveUp:    {Linewise: true, FreeCol: true},
	CmdMoveRight: {},

	CmdMoveToStart:    {},
	CmdMoveToEnd:      {},
	CmdMoveToNonBlank: {},
	CmdMoveToColumn:   {},

	CmdMoveByWord:              {},
	CmdMoveByWordEx:            {}, // XXX debug
	CmdMoveBackwardByWord:      {},
	CmdMoveToEndOfWord:         {},
	CmdMoveByLooseWord:         {},
	CmdMoveBackwardByLooseWord: {},
	CmdMoveToEndOfLooseWord:    {},

	CmdMoveByLine:         {Linewise: true},
	CmdMoveBackwardByLine: {Linewise: true},
	CmdMoveToLastLine:     {Linewise: true},
	CmdMoveToLine:         {Linewise: true},

	CmdMoveBySentence:          {Linewise: true},
	CmdMoveBackwardBySentence:  {Linewise: true},
	CmdMoveByParagraph:         {Linewise: true},
	CmdMoveBackwardByParagraph: {Linewise: true},
	CmdMoveBySection:           {Linewise: true},
	CmdMoveBackwardBySection:   {Linewise: true},

	CmdMoveToTopOfView:         {Linewise: true},
	CmdMoveToMiddleOfView:      {Linewise: true},
	CmdMoveToBottomOfView:      {Linewise: true},
	CmdMoveToBelowTopOfView:    {Linewise: true},
	CmdMoveToAboveBottomOfView: {Linewise: true},

	CmdMoveToMark:     {},
	CmdMoveToMarkLine: {Linewise: true},

	CmdMoveBackToMark:     {},
	CmdMoveBackToMarkLine: {Linewise: true},

	CmdMoveSearchForward:        {Locate: true},
	CmdMoveSearchBackward:       {Locate: true},
	CmdMoveSearchNextMatch:      {Locate: true},
	CmdMoveSearchPrevMatch:      {Locate: true},
	CmdMoveSearchRepeatForward:  {Locate: true},
	CmdMoveSearchRepeatBackward: {Locate: true},

	CmdMoveFindForward:        {Inclusive: true},
	CmdMoveFindBackward:       {Inclusive: true},
	CmdMoveFindBeforeForward:  {Inclusive: true},
	CmdMoveFindBeforeBackward: {Inclusive: true},
	CmdMoveFindNextMatch:      {Inclusive: true},
	CmdMoveFindPrevMatch:      {Inclusive: true},
}
