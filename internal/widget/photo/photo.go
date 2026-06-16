// Package photo provides the photo-frame widget: a slideshow of images from one
// or more directories.
//
// Pipeline (the M1.Ph4 discipline applied to heavy work): a background goroutine
// picks the next file on a timer AND decodes/resizes it to the window size
// (NewPixbufFromFileAtScale, off the main thread — big photos must not block the
// UI); the ready pixbuf is handed to the main thread via widget.OnMain, which
// wraps it in a texture and shows it.
package photo

import (
	"fmt"
	"io/fs"
	"log/slog"
	"math/rand/v2"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/olegmif/dashboard/internal/config"
	"github.com/olegmif/dashboard/internal/widget"
)

// options are the photo widget's [widget.options].
type options struct {
	Dirs      []string `toml:"dirs"`      // directories to scan for images
	Recursive bool     `toml:"recursive"` // descend into subdirectories
	Interval  int      `toml:"interval"`  // seconds between images
	Order     string   `toml:"order"`     // "sequential" (default) | "random"
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".webp": true, ".gif": true, ".bmp": true,
}

// New builds the photo-frame widget.
func New(cfg config.Widget) (gtk.Widgetter, error) {
	opts := options{Interval: 30, Order: "sequential"}
	if err := cfg.DecodeOptions(&opts); err != nil {
		slog.Warn("photo: invalid options, using defaults", "err", err)
		opts = options{Interval: 30, Order: "sequential"}
	}
	if opts.Interval <= 0 {
		opts.Interval = 30
	}

	files := scanImages(opts.Dirs, opts.Recursive)
	slog.Info("photo: scanned source",
		"dirs", opts.Dirs, "recursive", opts.Recursive, "found", len(files))

	if len(files) == 0 {
		return status(statusText(opts.Dirs, files)), nil
	}

	picture := gtk.NewPicture()
	picture.SetHExpand(true)
	picture.SetVExpand(true)
	picture.SetContentFit(gtk.ContentFitCover) // fill the window, cropping excess
	// small inset so the glass mat (frosted background) shows around the photo
	const mat = 12
	picture.SetMarginTop(mat)
	picture.SetMarginBottom(mat)
	picture.SetMarginStart(mat)
	picture.SetMarginEnd(mat)

	// root carries the frosted-glass background; the photo sits on it inset by
	// the margins above, so the glass shows as a mat (passe-partout).
	root := gtk.NewBox(gtk.OrientationVertical, 0)
	root.AddCSSClass("widget-photo")
	root.AddCSSClass("glass")
	root.Append(picture)

	paths := make(chan string)
	go rotate(files, opts, paths)
	go func() {
		for path := range paths {
			// Decode + resize to the window size off the main thread. Big photos
			// (4000×3000) must not be decoded on the UI thread.
			pb, err := gdkpixbuf.NewPixbufFromFileAtScale(path, cfg.Width, cfg.Height, true)
			if err != nil {
				slog.Warn("photo: cannot load image, skipping", "path", path, "err", err)
				continue
			}
			widget.OnMain(func() {
				picture.SetPaintable(gdk.NewTextureForPixbuf(pb))
			})
		}
	}()

	return root, nil
}

// rotate emits the next image path every opts.Interval seconds, forever.
func rotate(files []string, opts options, out chan<- string) {
	interval := time.Duration(opts.Interval) * time.Second
	for i := 0; ; i++ {
		var path string
		if opts.Order == "random" {
			path = files[rand.IntN(len(files))]
		} else {
			path = files[i%len(files)]
		}
		out <- path
		time.Sleep(interval)
	}
}

// status builds a simple centered label, shown when there are no images.
func status(text string) gtk.Widgetter {
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.AddCSSClass("widget")
	box.AddCSSClass("widget-photo")
	box.SetHAlign(gtk.AlignCenter)
	box.SetVAlign(gtk.AlignCenter)
	box.Append(gtk.NewLabel(text))
	return box
}

func statusText(dirs, files []string) string {
	switch {
	case len(dirs) == 0:
		return "фото: каталоги не заданы"
	case len(files) == 0:
		return "фото: изображения не найдены"
	default:
		return fmt.Sprintf("фото: найдено %d", len(files))
	}
}

// scanImages collects image file paths from dirs, sorted. Unreadable dirs/files
// are skipped with a warning rather than failing the widget.
func scanImages(dirs []string, recursive bool) []string {
	var out []string
	for _, dir := range dirs {
		if recursive {
			_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					slog.Warn("photo: scan error, skipping", "path", path, "err", err)
					return nil
				}
				if !d.IsDir() && isImage(path) {
					out = append(out, path)
				}
				return nil
			})
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			slog.Warn("photo: cannot read dir, skipping", "dir", dir, "err", err)
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && isImage(e.Name()) {
				out = append(out, filepath.Join(dir, e.Name()))
			}
		}
	}
	sort.Strings(out)
	return out
}

func isImage(name string) bool {
	return imageExts[strings.ToLower(filepath.Ext(name))]
}
