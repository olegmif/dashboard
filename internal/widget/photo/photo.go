// Package photo provides the photo-frame widget. M1.Ph3: placeholder only; the
// real slideshow arrives in M2.
package photo

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// New builds the photo-frame widget.
func New(_ config.Widget) (gtk.Widgetter, error) {
	return widget.Placeholder("photo"), nil
}
