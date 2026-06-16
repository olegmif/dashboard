// Package config loads the dashboard's TOML configuration.
package config

import "github.com/BurntSushi/toml"

// Config is the whole dashboard configuration. It grows per milestone; for M1
// it only carries window placement.
type Config struct {
	Window Window `toml:"window"`
}

// Window is the desktop-widget placement on the left monitor, in root-space
// pixels (left monitor usually starts at x=0).
type Window struct {
	X      int `toml:"x"`
	Y      int `toml:"y"`
	Width  int `toml:"width"`
	Height int `toml:"height"`
}

// Default is used when no config file is present.
func Default() Config {
	return Config{Window: Window{X: 0, Y: 0, Width: 1920, Height: 1080}}
}

// Load reads the TOML file at path on top of Default. Missing keys keep their
// default. Returns the error if the file exists but cannot be parsed; a missing
// file is the caller's concern (use os.IsNotExist on the error).
func Load(path string) (Config, error) {
	cfg := Default()
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}
