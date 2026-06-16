// Package config loads the dashboard's TOML configuration.
package config

import "github.com/BurntSushi/toml"

// Config is the whole dashboard configuration: a list of widgets, each in its
// own desktop window.
type Config struct {
	Widgets []Widget `toml:"widget"`
}

// Widget describes one widget-window: which widget to show (Type), an optional
// per-window identifier (Instance, used as WM_CLASS instance / window title),
// and its placement on screen in root-space pixels.
type Widget struct {
	Type     string `toml:"type"`
	Instance string `toml:"instance"`
	X        int    `toml:"x"`
	Y        int    `toml:"y"`
	Width    int    `toml:"width"`
	Height   int    `toml:"height"`
}

// Default is used when no config file is present or it lists no widgets.
func Default() Config {
	return Config{Widgets: []Widget{
		{Type: "clock", Instance: "clock", X: 40, Y: 40, Width: 300, Height: 120},
	}}
}

// Load reads the TOML file at path. A missing file surfaces as an error
// satisfying os.IsNotExist; the caller decides how to fall back.
func Load(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	return cfg, err
}
