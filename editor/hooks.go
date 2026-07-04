package editor

import (
	"os"
	"time"
)

type FileInfo struct {
	size    int64
	modTime time.Time
}

func NewFileInfo(size int64, modTime time.Time) FileInfo {
	return FileInfo{
		size:    size,
		modTime: modTime,
	}
}

func (fi *FileInfo) Size() int64 {
	return fi.size
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.modTime
}

func DefaultStat(name string) (FileInfo, error) {
	info, err := os.Stat(name)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		modTime: info.ModTime(),
		size:    info.Size(),
	}, nil
}

type Hooks struct {
	ReadFile  func(string) ([]byte, error)
	Stat      func(string) (FileInfo, error)
	WriteFile func(string, []byte, os.FileMode) error

	ReadKillFile  func(string) ([]byte, error)
	WriteKillFile func(string, []byte, os.FileMode) error
}

func DefaultHooks() Hooks {
	return Hooks{
		ReadFile:  os.ReadFile,
		Stat:      DefaultStat,
		WriteFile: os.WriteFile,

		ReadKillFile:  os.ReadFile,
		WriteKillFile: os.WriteFile,
	}
}
