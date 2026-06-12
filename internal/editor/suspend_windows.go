//go:build windows

package editor

func (ed *Editor) MiscSuspend() {
	ed.EnsureCommand()
	ed.Unimplemented("MiscSuspend")
}
