package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWidgetDecodeOptions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `
[[widget]]
type = "clock"
[widget.options]
format = "15:04"

[[widget]]
type = "photo"
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Widgets) != 2 {
		t.Fatalf("got %d widgets, want 2", len(cfg.Widgets))
	}

	type clockOptions struct {
		Format string `toml:"format"`
	}

	// Widget WITH options: the table is decoded onto the defaults.
	withOpts := clockOptions{Format: "default"}
	if err := cfg.Widgets[0].DecodeOptions(&withOpts); err != nil {
		t.Fatalf("DecodeOptions (with): %v", err)
	}
	if withOpts.Format != "15:04" {
		t.Errorf("format = %q, want %q", withOpts.Format, "15:04")
	}

	// Widget WITHOUT options: defaults are preserved (no-op).
	noOpts := clockOptions{Format: "default"}
	if err := cfg.Widgets[1].DecodeOptions(&noOpts); err != nil {
		t.Fatalf("DecodeOptions (without): %v", err)
	}
	if noOpts.Format != "default" {
		t.Errorf("format = %q, want %q (defaults must survive)", noOpts.Format, "default")
	}
}
