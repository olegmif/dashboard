// Package deskwindow pins a GTK toplevel window to the desktop as a widget,
// hiding all windowing-system specifics behind a Backend.
//
// The "pin a window to the background at a fixed spot" operation has no portable
// form — every windowing system does it differently (override-redirect on X11,
// the wlr-layer-shell protocol on Wayland, per-WM rules, ...). Backend is the
// single seam where that lives, so the rest of the app stays WM-agnostic. GTK
// itself is constant across backends, so coupling to it here is fine.
package deskwindow

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// Geometry is an absolute placement in the root coordinate space, in pixels.
// On a multi-monitor X11 setup the left monitor usually starts at X=0.
type Geometry struct {
	X, Y, Width, Height int
}

// Backend places a toplevel window as a desktop widget and keeps it there.
type Backend interface {
	// Name reports the backend kind, for logging.
	Name() string
	// Pin wires the window's lifecycle so it ends up pinned at geo. Call it
	// before the window is shown.
	Pin(win *gtk.ApplicationWindow, geo Geometry)
}
