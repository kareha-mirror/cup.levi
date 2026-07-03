package cmd

// Attributes of motion command.
type Attr struct {
	Linewise  bool
	FreeCol   bool
	Inclusive bool
	Locate    bool
}

// Map of attributes of motion commands.
var AttrOf = map[Kind]Attr{
	MoveLeft:          {},
	MoveDown:          {Linewise: true, FreeCol: true},
	MoveHere:          {Linewise: true, FreeCol: true},
	MoveUp:            {Linewise: true, FreeCol: true},
	MoveRight:         {},
	MoveToStart:       {},
	MoveToEnd:         {},
	MoveToAfterIndent: {},
	MoveToColumn:      {},

	MoveByWord:              {},
	MoveByChangeWord:        {},
	MoveByDeleteWord:        {},
	MoveBackwardByWord:      {},
	MoveToEndOfWord:         {},
	MoveByLooseWord:         {},
	MoveBackwardByLooseWord: {},
	MoveToEndOfLooseWord:    {},

	MoveByLine:         {Linewise: true},
	MoveBackwardByLine: {Linewise: true},
	MoveToLastLine:     {Linewise: true},
	MoveToLine:         {Linewise: true, Locate: true},

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

	MoveToMark:     {Locate: true},
	MoveToMarkLine: {Linewise: true, Locate: true},
	BackToMark:     {Locate: true},
	BackToMarkLine: {Linewise: true, Locate: true},

	Search:               {Locate: true},
	SearchBackward:       {Locate: true},
	SearchNext:           {Locate: true},
	SearchPrev:           {Locate: true},
	RepeatSearch:         {Locate: true},
	RepeatBackwardSearch: {Locate: true},

	Find:               {Inclusive: true},
	FindBackward:       {Inclusive: true},
	FindBefore:         {Inclusive: true},
	FindBeforeBackward: {Inclusive: true},
	FindNext:           {Inclusive: true},
	FindPrev:           {Inclusive: true},
}
