# levi Prompt Commands

## Motion

* `:+` `<num>` `Enter` : Move cursor to first non-blank character of next line. (MoveByLine)
* `:-` `<num>` `Enter` : Move cursor to first non-blank character of previous line. (MoveBackwardByLine)
* `:` `<num>` `Enter` : Move cursor to first non-blank character of line specifined by <num>. (MoveToLine)

## Save / Load

* `:wq` `Enter` : Save current file and quit. (SaveAndQuit)
* `:w` `Enter` : Save current file. (Save)
* `:w!` `Enter` : Force save current file. (ForceSave)
* `:q` `Enter` : Quit editor. (Quit)
* `:q!` `Enter` : Force quit editor. (ForceQuit)
* `:e` `Enter` : Load file. (Load)
* `:e!` `Enter` : Force load file. (ForceLoad)
* `:r` `Enter` : Read file and insert to current buffer. (Read)
* `:n` `Enter` : Go to next buffer in list. (Next)
* `:prev` `Enter` : Go to previous buffer in list. (Prev)

## Shell

* `:sh` `Enter` : Execute shell. (Shell)

## Save All

* `:wa` `Enter` : Save all files. (SaveAll)
* `:wa!` `Enter` : Force save all files. (ForceSaveAll)
* `:qa` `Enter` : Close all files and quit editor. (QuitAll)
* `:qa!` `Enter` : Force close all files and quit editor. (ForceQuitAll)

## Settings

* `:set` `ts=<num>` `Enter` : Set tab stop size. (TabStop)
* `:set` `ai` `Enter` : Set auto indent enabled. (AutoIndent)
* `:set` `noai` `Enter` : Set auto indent disabled. (NoAutoIndent)

## levi Enhanced

* `:open` `Enter` : Open file in new buffer. (Open)
* `:newline` `Enter` : Set newline type. (Newline)
* `:colors` `Enter` : Set colorscheme. (Colors)

## Debug

* `:mem` `Enter` : Show memory usage. (Mem)
* `:hello` `Enter` : Used by debug. (Hello)

## For Compatibility

* internal use : Show highlighted message. (Ring)
* internal use : Show error message. (Error)
