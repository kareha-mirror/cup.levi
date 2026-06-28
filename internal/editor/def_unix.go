//go:build unix

package editor

func DefaultConfig() *Config {
	return &Config{
		AutoIndent: true,
		TabStop:    4,

		Colors: "standard",
		Silent: false,
		CRLF:   false, // unix
		Depth:  4,
		Shared: "xyz",

		EscapeTimeout: 100,
	}
}
