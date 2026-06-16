// Package sysicons provides the system-status icon row. M1.Ph3: placeholder
// only; the real icon row arrives in M3.
package sysicons

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// New builds the system-status icon row widget.
func New(_ config.Widget) (gtk.Widgetter, error) {
	return widget.Placeholder("sysicons"), nil
}
