// Command dashboard renders desktop widgets on the left monitor under dwm.
//
// M1.Ph1: a minimal, undecorated GTK4 window. Nothing more — this exists to
// prove the gotk4 + GTK4 + cgo toolchain works before any widget logic.
package main

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const appID = "com.olegmif.dashboard"

func main() {
	app := gtk.NewApplication(appID, gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("dashboard")
	win.SetDefaultSize(400, 300)
	win.SetDecorated(false) // no titlebar/borders — desktop-widget style
	win.Present()
}
