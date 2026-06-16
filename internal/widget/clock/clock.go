// Package clock provides the clock widget. M1.Ph3: placeholder only; the real
// ticking clock arrives in M2.
package clock

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// New builds the clock widget.
func New(_ config.Widget) (gtk.Widgetter, error) {
	return widget.Placeholder("clock"), nil
}
