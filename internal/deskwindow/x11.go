//go:build linux

package deskwindow

// X11 override-redirect backend.
//
// An override-redirect window bypasses the window manager entirely: dwm (and
// every EWMH WM) ignores such windows in manage(), so they are never tiled,
// never stack-managed, and stay on all tags — exactly the conky-style desktop
// widget behaviour. The catch under GTK4: the toolkit exposes neither the X
// window id nor the override-redirect flag, so we drop to Xlib via cgo.
//
// Timing matters: override-redirect must be set BEFORE the window is mapped, or
// the WM will already have grabbed it. We set it at "realize" (the GdkSurface —
// and its X window — exist, but mapping happens after), then move/resize/lower
// at "map", after GTK has placed the surface itself.

// #cgo pkg-config: gtk4 x11
// #include <gdk/gdk.h>
// #include <gdk/x11/gdkx.h>
// #include <X11/Xlib.h>
//
// // Returns 0 on success, 1 if the surface is not an X11 surface.
// static int dw_set_override_redirect(void *surface_ptr) {
//     GdkSurface *surface = GDK_SURFACE(surface_ptr);
//     if (!GDK_IS_X11_SURFACE(surface)) return 1;
//     Display *xd = GDK_DISPLAY_XDISPLAY(gdk_surface_get_display(surface));
//     Window   xw = gdk_x11_surface_get_xid(surface);
//     XSetWindowAttributes attrs;
//     attrs.override_redirect = True;
//     XChangeWindowAttributes(xd, xw, CWOverrideRedirect, &attrs);
//     XFlush(xd);
//     return 0;
// }
//
// // Returns 0 on success, 1 if the surface is not an X11 surface.
// static int dw_place(void *surface_ptr, int x, int y, int w, int h) {
//     GdkSurface *surface = GDK_SURFACE(surface_ptr);
//     if (!GDK_IS_X11_SURFACE(surface)) return 1;
//     Display *xd = GDK_DISPLAY_XDISPLAY(gdk_surface_get_display(surface));
//     Window   xw = gdk_x11_surface_get_xid(surface);
//     XMoveResizeWindow(xd, xw, x, y, w, h);
//     XLowerWindow(xd, xw);
//     XFlush(xd);
//     return 0;
// }
import "C"

import (
	"log/slog"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// New returns the desktop-window backend for the current session. Only the X11
// override-redirect backend exists today; backend selection (e.g. Wayland
// layer-shell) will branch here later.
func New() Backend { return &x11Override{} }

type x11Override struct{}

func (*x11Override) Name() string { return "x11/override-redirect" }

func (*x11Override) Pin(win *gtk.ApplicationWindow, geo Geometry) {
	win.ConnectRealize(func() {
		if rc := C.dw_set_override_redirect(surfacePointer(win)); rc != 0 {
			slog.Warn("deskwindow: surface is not X11; override-redirect skipped — " +
				"window will be managed by the WM")
		}
	})
	win.ConnectMap(func() {
		C.dw_place(surfacePointer(win),
			C.int(geo.X), C.int(geo.Y), C.int(geo.Width), C.int(geo.Height))
	})
}

// surfacePointer returns the native GdkSurface* of a realized window.
// gtk.BaseWidget is needed because ApplicationWindow embeds *coreglib.Object
// directly, so a bare win.Native() would resolve to Object.Native() (uintptr)
// instead of Widget.Native() (*NativeSurface).
func surfacePointer(win *gtk.ApplicationWindow) unsafe.Pointer {
	surface := gdk.BaseSurface(gtk.BaseWidget(win).Native().Surface())
	return unsafe.Pointer(surface.Native())
}
