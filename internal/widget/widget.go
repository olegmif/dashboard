// Package widget defines the contract every dashboard widget satisfies and a
// registry that maps a widget's type name (from config) to its builder. Adding
// a new widget type means writing a builder and registering it — the main loop
// stays unaware of concrete widgets.
package widget

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
)

// Builder constructs a widget's root GTK widget from its config entry. This is
// all the main loop knows about any widget.
type Builder func(cfg config.Widget) (gtk.Widgetter, error)

// Registry maps widget type names to builders.
type Registry struct {
	builders map[string]Builder
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{builders: make(map[string]Builder)}
}

// Register associates a type name with a builder. A later Register for the same
// type wins.
func (r *Registry) Register(typ string, b Builder) {
	r.builders[typ] = b
}

// Build constructs the widget for cfg.Type, or errors if the type is unknown.
func (r *Registry) Build(cfg config.Widget) (gtk.Widgetter, error) {
	b, ok := r.builders[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("unknown widget type %q", cfg.Type)
	}
	return b(cfg)
}

// Placeholder builds a labeled stub used by a widget until it is implemented
// (M2–M6). It is a framed box titled with the widget type.
func Placeholder(typ string) gtk.Widgetter {
	frame := gtk.NewFrame(typ)
	frame.AddCSSClass("widget")
	frame.AddCSSClass("widget-" + typ)

	label := gtk.NewLabel("заглушка")
	label.SetHExpand(true)
	label.SetVExpand(true)
	frame.SetChild(label)
	return frame
}
