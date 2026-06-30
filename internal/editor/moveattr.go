package editor

type MoveAttr struct {
	Linewise  bool
	FreeCol   bool
	Inclusive bool
	Locate    bool
}

var MoveAttrs = map[CmdKind]MoveAttr{
	MoveLeft:  {},
	MoveDown:  {Linewise: true, FreeCol: true},
	MoveUp:    {Linewise: true, FreeCol: true},
	MoveRight: {},

	MoveToStart:    {},
	MoveToEnd:      {},
	MoveToNonBlank: {},
	MoveToColumn:   {},

	MoveByWord:              {},
	MoveByChangeWord:        {},
	MoveBackwardByWord:      {},
	MoveToEndOfWord:         {},
	MoveByLooseWord:         {},
	MoveBackwardByLooseWord: {},
	MoveToEndOfLooseWord:    {},

	MoveByLine:         {Linewise: true},
	MoveBackwardByLine: {Linewise: true},
	MoveToLastLine:     {Linewise: true},
	MoveToLine:         {Linewise: true},

	MoveBySentence:          {Linewise: true},
	MoveBackwardBySentence:  {Linewise: true},
	MoveByParagraph:         {Linewise: true},
	MoveBackwardByParagraph: {Linewise: true},
	MoveBySection:           {Linewise: true},
	MoveBackwardBySection:   {Linewise: true},

	MoveToTopOfView:         {Linewise: true},
	MoveToMiddleOfView:      {Linewise: true},
	MoveToBottomOfView:      {Linewise: true},
	MoveToBelowTopOfView:    {Linewise: true},
	MoveToAboveBottomOfView: {Linewise: true},

	MoveToMark:     {},
	MoveToMarkLine: {Linewise: true},

	BackToMark:     {},
	BackToMarkLine: {Linewise: true},

	SearchForward:        {Locate: true},
	SearchBackward:       {Locate: true},
	SearchNext:           {Locate: true},
	SearchPrev:           {Locate: true},
	RepeatSearchForward:  {Locate: true},
	RepeatSearchBackward: {Locate: true},

	FindForward:        {Inclusive: true},
	FindBackward:       {Inclusive: true},
	FindBeforeForward:  {Inclusive: true},
	FindBeforeBackward: {Inclusive: true},
	FindNext:           {Inclusive: true},
	FindPrev:           {Inclusive: true},
}
