//go:build windows

package config

func Default() *Config {
	return &Config{
		AutoIndent: true,
		TabStop:    4,

		Colors: "standard",
		Silent: false,
		CRLF:   true, // windows
		Depth:  4,
		Shared: "xyz",

		EscapeTimeout: 100,
	}
}
