package editor

type KillMode int

const (
	KillNone = iota
	KillLines
	KillRunes
)

type Reg struct {
	Mode   KillMode
	Killed []string
	Shared bool
}

type Regs struct {
	Map map[string]*Reg
}

func (regs *Regs) SetReg(name string, reg *Reg) {
	if regs.Map == nil {
		regs.Map = make(map[string]*Reg)
	}
	regs.Map[name] = reg
}

func (regs *Regs) Mode(name string) KillMode {
	reg, ok := regs.Map[name]
	if !ok {
		return KillNone
	}
	return reg.Mode
}

func (regs *Regs) Killed(name string) []string {
	reg, ok := regs.Map[name]
	if !ok {
		return nil
	}
	return reg.Killed
}

func (regs *Regs) Shared(name string) bool {
	reg, ok := regs.Map[name]
	if !ok {
		return false
	}
	return reg.Shared
}

func (regs *Regs) SetShared(name string, shared bool) {
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillNone,
			Killed: nil,
			Shared: shared,
		}
		regs.SetReg(name, reg)
		return
	}
	reg.Shared = shared
}

func (regs *Regs) SetLines(name string, killed []string) {
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillLines,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	reg.Mode = KillLines
	reg.Killed = append([]string{}, killed...)
}

func (regs *Regs) SetRunes(name string, killed []string) {
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillRunes,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	reg.Mode = KillRunes
	reg.Killed = append([]string{}, killed...)
}

func (regs *Regs) AddLines(name string, killed []string) {
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillLines,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	if reg.Mode == KillLines {
		reg.Killed = append(reg.Killed, killed...)
		return
	}
	reg.Mode = KillLines
	reg.Killed = append([]string{}, killed...)
}

func (regs *Regs) AddRunes(name string, killed []string) {
	reg, ok := regs.Map[name]
	if !ok {
		reg = &Reg{
			Mode:   KillRunes,
			Killed: append([]string{}, killed...),
			Shared: false,
		}
		regs.SetReg(name, reg)
		return
	}
	if reg.Mode == KillRunes {
		lines := append([]string{}, reg.Killed[:len(reg.Killed)-1]...)
		line := reg.Killed[len(reg.Killed)-1] + killed[0]
		lines = append(lines, line)
		if 1 < len(killed) {
			lines = append(lines, killed[1:]...)
		}
		reg.Killed = lines
		return
	}
	reg.Mode = KillRunes
	reg.Killed = append([]string{}, killed...)
}

func (regs *Regs) LoadConfig(cfg *Config) {
	for _, reg := range regs.Map {
		reg.Shared = false
	}
	for _, r := range cfg.Shared {
		regs.SetShared(string([]rune{r}), true)
	}
}
