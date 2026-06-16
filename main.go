// Command dashboard renders desktop widgets on the left monitor under dwm.
//
// M1.Ph2: an undecorated window pinned to a fixed rectangle on the left monitor
// via the deskwindow backend (X11 override-redirect). Geometry comes from
// config.toml. The window content is still just a placeholder label.
package main

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/deskwindow"
)

const (
	appID      = "com.olegmif.dashboard"
	configPath = "config.toml"
)

func main() {
	app := gtk.NewApplication(appID, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	cfg, err := config.Load(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		slog.Warn("no config.toml; using defaults", "path", configPath)
	} else if err != nil {
		slog.Error("cannot parse config.toml; using defaults", "err", err)
		cfg = config.Default()
	}

	win := gtk.NewApplicationWindow(app)
	win.SetTitle("dashboard")
	win.SetDecorated(false)
	win.SetDefaultSize(cfg.Window.Width, cfg.Window.Height)
	win.SetChild(gtk.NewLabel("dashboard — M1.Ph2")) // visible placeholder

	geo := deskwindow.Geometry{
		X:      cfg.Window.X,
		Y:      cfg.Window.Y,
		Width:  cfg.Window.Width,
		Height: cfg.Window.Height,
	}
	backend := deskwindow.New()
	slog.Info("pinning window", "backend", backend.Name(), "geometry", geo)
	backend.Pin(win, geo)

	win.Present()
}
