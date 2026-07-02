package editor

import (
	"tea.kareha.org/cup/levi/internal/cmd"
)

func (ed *Editor) Run(c cmd.Pair, replay bool) (bool, bool) {
	ed.Commit()

	// * motion commands are delegated to RunMove
	// * other commands are dispatched in this method
	// * compound commands also refer to RunMove

	//
	// Motion Commands
	//

	if c.Op.Kind == cmd.Invalid {
		if attr, ok := cmd.AttrOf[c.Mv.Kind]; ok {
			if loc, ok := ed.RunMove(c.Mv, 1); ok {
				b := ed.Buf()
				if attr.Linewise {
					if attr.FreeCol {
						loc.Col = b.ConfineFreeColInclusive(loc.Row)
					}
				} else {
					loc = b.ConfineInclusive(loc)
					b.VirtCol = loc.Col
				}
				b.Loc = loc
				if b.Loc.Col < b.ViewLoc.Col {
					b.ViewLoc.Col = 0
				}
				if attr.Locate {
					ed.Locate()
				}
			}
			return false, true
		}

		ed.Notice("Not a levi command [%s]", ed.cmdInp.String())
		return false, true
	}

	//
	// Other Commands
	//

	return ed.RunOp(c, replay)
}
