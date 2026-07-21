# levi Commands

Categories

* Motion
* Insert
* Edit
* Mark
* Copy
* View
* Miscellaneous
* Select Current Buffer
* For Compatibility

## Motion

Motion commands move cursor.
They themself don't change text content of buffer.

### Move by Character / Move by Line

* `h` : Move cursor left by character. (MoveLeft)
* `j` : Move cursor down by line. (MoveDown)
* internal use : Move cursor here. (MoveHere)
* `k` : Move cursor up by line. (MoveUp)
* `l` : Move cursor right by character. (MoveRight)

### Move in Line

* `0` : Move cursor to start of current line. (MoveToStart)
* `$` : Move cursor to end of current line. (MoveToEnd)
* `^` : Move cursor to first non-blank character of current line. (MoveToAfterIndent)
* `<num>` `|` : Move cursor to column <num> of current line. (MoveToColumn) (Proper vi's column number is visual-based, but levi's is rune-based.)

### Move by Word / Move by Loose Word

* `w` : Move cursor forward by word. (MoveByWord)
* internal use : Move cursor forward by word used by cw. (MoveByChangeWord)
* internal use : Move cursor forward by word used by dw. (MoveByDeleteWord)
* `b` : Move cursor backward by word. (MoveBackwardByWord)
* `e` : Move cursor to end of word. (MoveToEndOfWord)
* `W` : Move cursor forward by loose word. (MoveByLooseWord)
* internal use : Move cursor forward by loose word used by cW. (MoveByChangeLooseWord)
* internal use : Move cursor forward by word used by dW. (MoveByDeleteLooseWord)
* `B` : Move cursor backward by loose word. (MoveBackwardByLooseWord)
* `E` : Move cursor to end of loose word. (MoveToEndOfLooseWord)

### Move by Line

* `Enter`, `+` : Move cursor to first non-blank character of next line. (MoveByLine)
* `-` : Move cursor to first non-blank character of previous line. (MoveBackwardByLine)
* `G` : Move cursor to first non-blank character of last line. (MoveToLastLine)
* `<num>` `G` : Move cursor to first non-blank character of line specified by <num>. (MoveToLine)

### Move by Block

* `)` : Move cursor forward by sentence. (MoveBySentence)
* `(` : Move cursor backward by sentence. (MoveBackwardBySentence)
* `}` : Move cursor forward by paragraph. (MoveByParagraph)
* `{` : Move cursor backward by paragraph. (MoveBackwardByParagraph)
* `]]` : Move cursor forward by section. (MoveBySection)
* `[[` : Move cursor backward by section. (MoveBackwardBySection)

### Move in View

* `H` : Move cursor to top of view. (MoveToTopOfView)
* `M` : Move cursor to middle of view. (MoveToMiddleOfView)
* `L` : Move cursor to bottom of view. (MoveToBottomOfView)
* `<num>` `H` : Move cursor below <num> lines from top of view. (MoveToBelowTopOfView)
* `<num>` `L` : Move cursor above <num> lines from bottom of view. (MoveToAboveBottomOfView)

### Move to Mark

* `` ` `` `<char>` : Move cursor to marked position labelled by <char>. (MoveToMark)
* `'` `<char>` : Move cursor to marked line labelled by <char>. (MoveToMarkLine)

### Move by Context

* ``` `` ``` : Move cursor to previous position in context. (BackToMark)
* `''` : Move cursor to previous line in context. (BackToMarkLine)

### Search

* `/` `<pattern>` `Enter` : Search <pattern> and move to it. (Search)
* `?` `<pattern>` `Enter` : Search <pattern> backward and move to it. (SearchBackward)
* `n` : Repeat last search operation to search next match. (SearchNext)
* `N` : Repeat last search operation to search previous match. (SearchPrev)
* `/` `Enter` : Repeat last search. (RepeatSearch)
* `?` `Enter` : Repeat last backward search. (RepeatBackwardSearch)

### Find Character

* `f` `<char>` : Find character <char> in current line and move to it. (Find)
* `F` `<char>` : Find character <char> backward in current line and move to it. (FindBackward)
* `t` `<char>` : Find character <char> in current line and move to before it. (FindBefore)
* `T` `<char>` : Find character <char> backward in current line and move before it. (FindBeforeBackward)
* `;` : Repeat find operation to find next match. (FindNext)
* `,` : Repeat find operation to find previous match. (FindPrev)

## Insert

Insert commands are commands which transit to insert mode.
They are identified by IsInsertCmd.
Insert commands which have multiplication number are identified by IsMultiInsertCmd.

### Insert

* `i` : Switches to insert mode. (Insert)
* `a` : Switches to insert mode after cursor. (InsertAfter)
* `I` : Switches to insert mode after indent of current line. (InsertAfterIndent)
* `A` : Switches to insert mode after end of current line. (InsertAfterEnd)

### Insert Line

* `o` : Inserts blank line below cursor and switches to insert mode. (InsertLine)
* `O` : Inserts blank line above cursor and switches to insert mode. (InsertLineAbove)

### Change / Substitute

* `c` `<mv>` : Changes region from cursor to destination of motion. (ChangeRegion)
* `s` `<char>` : Substitutes character under cursor. (Subst)

### Unsupported

* `R` : Switches to insert mode and overwrites current line. (Overwrite)

## Edit

Edit commands are commands which change text content of buffer.
They are identified by IsEditCmd set.

### Paste (Put)

* `"` `<reg>` `p` : Paste after cursor from register <reg>. (Paste)
* `"` `<reg>` `P` : Paste before cursor from register <reg>. (PasteBefore)

### Delete

* `x` : Delete character under cursor. (Delete)
* `X` : Delete character before cursor. (DeleteBefore)
* `d` `<mv>` : Delete region from current cursor to destination of motion <mv>. (DeleteRegion)

### Edit

* `r` `<char>` : Replace single character under cursor. (Replace)
* `J` : Join current line with next line. (Join)
* `>` `<mv>` : Indent region from current cursor to destination of motion <mv>. (IndentRegion)
* `<` `<mv>` : Outdent region from current cursor to destination of motion <mv>. (OutdentRegion)

### Restore

* `U` : Restore current line to previous state last visited. (Restore)

## Mark Commands

Most other mark commands are categorized to motion commands.

* `m` `<char>` : Mark current cursor position labelled by <char>. (Mark)

## Copy Commands

These commands copy lines or runes into registers.
They don't change text content of buffer.

### Copy (Yank)

* `y` `<mv>` : Copy region from current cursor to destination of motion <mv>. (CopyRegion)

## View Commands

View commands scroll screen.
They possibly move cursor, but are not recognized as motion commands.

### Scroll by View Height / Scroll by Line

* `Ctrl-F` : Scroll down by view height. (ViewDown)
* `Ctrl-B` : Scroll up by view height. (ViewUp)
* `Ctrl-D` : Scroll down by half view height. (ViewDownHalf)
* `Ctrl-U` : Scroll up by half view height. (ViewUpHalf)
* `Ctrl-Y` : Scroll down by line. (ViewDownLine)
* `Ctrl-E` : Scroll up by line. (ViewUpLine)

### Reposition

* `z` `Enter` : Reposition cursor line to top of view. (ViewToTop)
* `z.` : Reposition cursor line middle of view. (ViewToMiddle)
* `z-` : Reposition cursor line bottom of view. (ViewToBottom)

## Miscellaneous

* `Ctrl-G` : Show info about states of current buffer. (ShowInfo)
* `Ctrl-L` : Redraw view. (Redraw)
* `.` : Repeat last command which is repeatable. (Repeat)
* `u`: Undo last modification or redo by undoing itself. (Undo)
* `ZZ` : Save and close. (SaveAndClose)
* `Ctrl-Z` : Suspend editor process. (Suspend)

## Select Current Buffer

* `Ctrl-^`, `Ctrl-_` : Go to last visited buffer. (LastBuf)
* `<num>` `Ctrl-^`, `<num>` `Ctrl-_` : Go to buffer specified by <num>. (GoToBuf)  (Not available in nvi.  Available in Vim.)
* `:next` `Enter`, `:n` `Enter`, `zj` (levi enhancement) : Go to next buffer in list. (NextBuf)
* `:prev` `Enter`, `zk` (levi enhancement) : Go to previous buffer in list. (PrevBuf)

## For Compatibility

* internal use : Show highlighted message. (Ring)
