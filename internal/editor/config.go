package editor

type Config struct {
	AutoIndent bool
	TabWidth   int
}

func DefaultConfig() *Config {
	return &Config{
		AutoIndent: true,
		TabWidth:   4,
	}
}
