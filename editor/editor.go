package editor

import (
	"tea.kareha.org/cup/levi/hook"

	"tea.kareha.org/cup/levi/internal/editor"
)

type Editor struct {
	ed *editor.Editor
}

func WithHooks(
	cfgDir string, paths []string, hooks hook.Hooks,
) (*Editor, error) {
	ed, err := editor.WithHooks(cfgDir, paths, hooks)
	if err != nil {
		return nil, err
	}
	return &Editor{ed: ed}, nil
}

func Init(cfgDir string, paths []string) (*Editor, error) {
	return WithHooks(cfgDir, paths, hook.DefaultHooks())
}

func (ed *Editor) Finish() error {
	return ed.ed.Finish()
}

func (ed *Editor) Main() error {
	return ed.ed.Main()
}

func ClearSharedRegisters(cfgDir string) error {
	return editor.ClearSharedRegisters(cfgDir)
}
