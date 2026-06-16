// Package weather provides the weather widget. M1.Ph3: placeholder only; the
// real Open-Meteo widget arrives in M4.
package weather

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// New builds the weather widget.
func New(_ config.Widget) (gtk.Widgetter, error) {
	return widget.Placeholder("weather"), nil
}
