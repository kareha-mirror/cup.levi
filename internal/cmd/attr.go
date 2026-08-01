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
	MoveLeft:            {},
	MoveDown:            {Linewise: true, FreeCol: true},
	MoveHere:            {Linewise: true, FreeCol: true},
	MoveUp:              {Linewise: true, FreeCol: true},
	MoveRight:           {},
	MoveToStart:         {},
	MoveToEnd:           {},
	MoveToFirstNonBlank: {},
	MoveToColumn:        {},

	MoveByWord:          {},
	MoveByChangeWord:    {},
	MoveByDeleteWord:    {},
	MoveBackByWord:      {},
	MoveToEndOfWord:     {},
	MoveByBigword:       {},
	MoveByChangeBigword: {},
	MoveByDeleteBigword: {},
	MoveBackByBigword:   {},
	MoveToEndOfBigword:  {},

	MoveByLine:     {Linewise: true},
	MoveBackByLine: {Linewise: true},
	MoveToLastLine: {Linewise: true},
	MoveToLine:     {Linewise: true, Locate: true},

	MoveBySentence:      {},
	MoveBackBySentence:  {},
	MoveByParagraph:     {Linewise: true},
	MoveBackByParagraph: {Linewise: true},
	MoveBySection:       {Linewise: true},
	MoveBackBySection:   {Linewise: true},

	MoveToTopOfView:         {Linewise: true},
	MoveToMiddleOfView:      {Linewise: true},
	MoveToBottomOfView:      {Linewise: true},
	MoveToBelowTopOfView:    {Linewise: true},
	MoveToAboveBottomOfView: {Linewise: true},

	MoveToMark:     {Locate: true},
	MoveToMarkLine: {Linewise: true, Locate: true},
	BackToMark:     {Locate: true},
	BackToMarkLine: {Linewise: true, Locate: true},

	Search:           {Locate: true},
	SearchBack:       {Locate: true},
	SearchNext:       {Locate: true},
	SearchPrev:       {Locate: true},
	RepeatSearch:     {Locate: true},
	RepeatBackSearch: {Locate: true},

	Find:           {Inclusive: true},
	FindBack:       {Inclusive: true},
	FindBefore:     {Inclusive: true},
	FindBeforeBack: {Inclusive: true},
	FindNext:       {Inclusive: true},
	FindPrev:       {Inclusive: true},
}
