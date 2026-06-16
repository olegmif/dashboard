// Command dashboard renders desktop widgets on the left monitor under dwm.
//
// M1.Ph3: multi-window framework. config.toml lists widgets; the main loop
// builds one pinned, undecorated window per widget via the deskwindow backend.
// Widget content is still placeholder stubs (real content lands in M2–M6).
package main

import (
	_ "embed"
	"errors"
	"io/fs"
	"log/slog"
	"os"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/deskwindow"
	"github.com/olegmif/dashboard/internal/widget"
	"github.com/olegmif/dashboard/internal/widget/clock"
	"github.com/olegmif/dashboard/internal/widget/photo"
	"github.com/olegmif/dashboard/internal/widget/sysicons"
	"github.com/olegmif/dashboard/internal/widget/telegram"
	"github.com/olegmif/dashboard/internal/widget/weather"
)

const (
	appID      = "com.olegmif.dashboard"
	configPath = "config.toml"
)

//go:embed style.css
var styleCSS string

func main() {
	app := gtk.NewApplication(appID, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	cfg := loadConfig()
	applyCSS()

	registry := buildRegistry()
	backend := deskwindow.New()

	for _, w := range cfg.Widgets {
		if w.Width <= 0 || w.Height <= 0 {
			slog.Warn("widget has non-positive size; skipped",
				"type", w.Type, "instance", w.Instance)
			continue
		}
		ui, err := registry.Build(w)
		if err != nil {
			slog.Warn("skipping widget", "instance", w.Instance, "err", err)
			continue
		}

		instance := w.Instance
		if instance == "" {
			instance = w.Type
		}

		win := gtk.NewApplicationWindow(app)
		win.SetTitle("dashboard:" + instance)
		win.SetDecorated(false)
		win.SetDefaultSize(w.Width, w.Height)
		win.SetChild(ui)

		geo := deskwindow.Geometry{X: w.X, Y: w.Y, Width: w.Width, Height: w.Height}
		slog.Info("pinning widget",
			"type", w.Type, "instance", instance, "geometry", geo, "backend", backend.Name())
		backend.Pin(win, geo, instance)
		win.Present()
	}
}

// loadConfig reads config.toml, falling back to defaults if it is missing,
// unparseable, or lists no widgets.
func loadConfig() config.Config {
	cfg, err := config.Load(configPath)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		slog.Warn("no config.toml; using defaults", "path", configPath)
		cfg = config.Default()
	case err != nil:
		slog.Error("cannot parse config.toml; using defaults", "err", err)
		cfg = config.Default()
	}
	if len(cfg.Widgets) == 0 {
		slog.Warn("config lists no widgets; using defaults")
		cfg = config.Default()
	}
	return cfg
}

// applyCSS installs the embedded stylesheet for every window on the display.
func applyCSS() {
	provider := gtk.NewCSSProvider()
	provider.LoadFromString(styleCSS)
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(), provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

// buildRegistry wires widget type names to their builders. Adding a widget type
// is one line here plus its package.
func buildRegistry() *widget.Registry {
	r := widget.NewRegistry()
	r.Register("clock", clock.New)
	r.Register("photo", photo.New)
	r.Register("telegram", telegram.New)
	r.Register("weather", weather.New)
	r.Register("sysicons", sysicons.New)
	return r
}
