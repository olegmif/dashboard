// Package clock provides the clock widget: current time and (optionally) date,
// updated every second.
//
// It also hosts the reference pattern for every live widget (introduced in
// M1.Ph4): a background goroutine produces values over a channel; a consumer
// marshals each onto the GTK main thread via widget.OnMain. Never touch GTK
// widgets from a goroutine.
package clock

import (
	"log/slog"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// options are the clock's [widget.options]. Format and DateFormat use Go's
// reference-time layout ("Mon Jan 2 15:04:05 2006"). An empty DateFormat hides
// the date row.
type options struct {
	Format     string `toml:"format"`
	DateFormat string `toml:"date_format"`
}

func defaults() options {
	return options{Format: "15:04:05", DateFormat: "Monday, 2 January"}
}

// New builds the clock widget.
func New(cfg config.Widget) (gtk.Widgetter, error) {
	opts := defaults()
	if err := cfg.DecodeOptions(&opts); err != nil {
		slog.Warn("clock: invalid options, using defaults", "err", err)
		opts = defaults()
	}

	// content holds the labels and is centred inside the glass panel.
	content := gtk.NewBox(gtk.OrientationVertical, 0)
	content.SetHExpand(true)
	content.SetVExpand(true)
	content.SetHAlign(gtk.AlignCenter)
	content.SetVAlign(gtk.AlignCenter)

	timeLabel := gtk.NewLabel("")
	timeLabel.AddCSSClass("clock-time")
	content.Append(timeLabel)

	var dateLabel *gtk.Label
	if opts.DateFormat != "" {
		dateLabel = gtk.NewLabel("")
		dateLabel.AddCSSClass("clock-date")
		content.Append(dateLabel)
	}

	// root fills the window and carries the frosted-glass background.
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.AddCSSClass("widget-clock")
	box.AddCSSClass("glass")
	box.Append(content)

	ticks := make(chan tick)
	go produce(ticks, opts)
	go func() {
		for t := range ticks {
			widget.OnMain(func() {
				timeLabel.SetText(t.time)
				if dateLabel != nil {
					dateLabel.SetText(t.date)
				}
			})
		}
	}()

	return box, nil
}

type tick struct{ time, date string }

// produce emits the formatted time/date once a second. It runs for the app's
// lifetime (a desktop widget never stops), so no cancellation is wired up.
func produce(out chan<- tick, opts options) {
	for {
		now := time.Now()
		out <- tick{time: now.Format(opts.Format), date: now.Format(opts.DateFormat)}
		time.Sleep(time.Second)
	}
}
