// Package telegram provides the Telegram chat widget. M1.Ph3: placeholder only;
// the real chat cards arrive in M5–M6.
package telegram

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// New builds the Telegram widget.
func New(_ config.Widget) (gtk.Widgetter, error) {
	return widget.Placeholder("telegram"), nil
}
