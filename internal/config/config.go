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
// its placement in root-space pixels, and an opaque per-widget options table.
type Widget struct {
	Type     string `toml:"type"`
	Instance string `toml:"instance"`
	X        int    `toml:"x"`
	Y        int    `toml:"y"`
	Width    int    `toml:"width"`
	Height   int    `toml:"height"`

	// Options is the raw [widget.options] subtable, decoded lazily per widget
	// type via DecodeOptions. Kept opaque here so config stays unaware of any
	// specific widget's settings.
	Options toml.Primitive `toml:"options"`

	meta *toml.MetaData // set by Load; needed to decode Options later
}

// DecodeOptions decodes this widget's [widget.options] subtable into target
// (a pointer to a typed options struct). A widget with no options is a no-op
// and leaves target untouched, so initialise target with defaults first:
//
//	opts := clockOptions{Format: "15:04:05"}
//	_ = cfg.DecodeOptions(&opts)
func (w Widget) DecodeOptions(target any) error {
	if w.meta == nil {
		return nil
	}
	return w.meta.PrimitiveDecode(w.Options, target)
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
	meta, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return cfg, err
	}
	// Stash a private copy of the metadata in each widget so it can decode its
	// own Options later. A per-widget copy avoids any shared-context races, and
	// widgets are decoded sequentially on the main thread anyway.
	for i := range cfg.Widgets {
		m := meta
		cfg.Widgets[i].meta = &m
	}
	return cfg, nil
}
