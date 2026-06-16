// Package clock provides the clock widget.
//
// M1.Ph4 hosts the reference pattern for every live widget here, since a clock
// is the natural ticking example. M2 will flesh it out (formatting, font, date);
// for now it just ticks once a second.
package clock

import (
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// New builds the clock widget.
//
// Reference pattern (copy this shape for telegram/weather/sysicons):
//   - a background goroutine PRODUCES values and sends them over a channel; it
//     must never touch GTK directly;
//   - a consumer goroutine drains the channel and marshals each value onto the
//     GTK main thread via widget.OnMain to update the widget.
func New(_ config.Widget) (gtk.Widgetter, error) {
	frame := gtk.NewFrame("clock")
	frame.AddCSSClass("widget")
	frame.AddCSSClass("widget-clock")

	label := gtk.NewLabel("")
	label.SetHExpand(true)
	label.SetVExpand(true)
	frame.SetChild(label)

	ticks := make(chan string)
	go produce(ticks)
	go func() {
		// now is per-iteration (Go 1.22+), so each OnMain closure captures its
		// own value even though it runs asynchronously on the main thread.
		for now := range ticks {
			widget.OnMain(func() { label.SetText(now) })
		}
	}()

	return frame, nil
}

// produce emits the current time once a second. It runs for the app's lifetime
// (a desktop widget never stops), so no cancellation is wired up.
func produce(out chan<- string) {
	for {
		out <- time.Now().Format("15:04:05")
		time.Sleep(time.Second)
	}
}
