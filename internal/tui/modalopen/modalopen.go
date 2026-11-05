// Package modalopen opens a Rhiza app as a rexuikit Modal: each app's Rex
// action runs `rhiza tui open <app>`, which opens the app's layer through
// rexuikit's surface over Rex's control connection, centred over the
// window of the block the action ran on and sized to the cells of the app's
// content, then follows the content's size as it changes (an inbox list
// growing or shrinking) with Handle.Resize until the modal closes.
package modalopen

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/dphaener/rexuikit/client"
	"github.com/dphaener/rexuikit/surface"

	"github.com/dphaener/rhiza/internal/tui/modalshell"
)

// Request is one modal to open.
type Request struct {
	SessionID string
	// BlockID is the block the action ran on; the modal opens over its
	// window. Empty, or a block no window holds, means the active window.
	BlockID string
	// Label is the block's label, which Rex shows as the layer's title.
	Label string
	// Argv is the app's program; the modal closes when it exits.
	Argv []string
	Cwd  string
	// Cols and Rows are the content's natural size the modal opens at.
	Cols, Rows int
	// Sizes delivers the content's size each time it changes; the modal is
	// resized to each. Nil means the modal keeps its first size.
	Sizes <-chan client.Size
}

// Open opens r's modal on rex, a connection attached to r.SessionID, and
// resizes it to each size from r.Sizes. It returns once the modal is gone
// (its program exited or someone closed it) or ctx ends.
func Open(ctx context.Context, rex surface.Rex, r Request) error {
	if r.Cols <= 0 || r.Rows <= 0 {
		return errors.New("modalopen: no content size")
	}
	s := surface.New(rex, r.SessionID)
	window, err := windowOf(ctx, rex, r.SessionID, r.BlockID)
	if err != nil {
		return err
	}
	h, err := s.Open(ctx, surface.Spec{
		Kind:     surface.Modal,
		Argv:     r.Argv,
		Cwd:      r.Cwd,
		Label:    r.Label,
		WindowID: window,
		Cells:    client.Size{Columns: r.Cols, Rows: r.Rows},
		Anchor:   surface.Centered,
		Focus:    true,
	})
	if err != nil {
		return err
	}
	sizes := r.Sizes
	for {
		select {
		case <-h.Done():
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case sz, ok := <-sizes:
			if !ok {
				sizes = nil
				continue
			}
			if err := resize(ctx, s, h, sz); err != nil {
				select {
				case <-h.Done():
					return nil
				default:
					return err
				}
			}
		}
	}
}

// resize fits h's layer to sz cells over its window as it is now, centred
// as the modal opened.
func resize(ctx context.Context, s *surface.Surface, h *surface.Handle, sz client.Size) error {
	if sz.Columns <= 0 || sz.Rows <= 0 {
		return nil
	}
	wc, _, err := s.Window(ctx, h.WindowID())
	if err != nil {
		return err
	}
	b, _ := surface.BoundsForCells(wc, sz.Columns, sz.Rows, surface.Centered)
	if b == h.Bounds() {
		return nil
	}
	return h.Resize(ctx, b)
}

// windowOf returns the window holding block, in any of its layers, or ""
// (the active window) when no window does.
func windowOf(ctx context.Context, rex surface.Rex, session, block string) (string, error) {
	if block == "" {
		return "", nil
	}
	v, err := rex.View(ctx, session)
	if err != nil {
		return "", err
	}
	for _, w := range v.Windows {
		for _, l := range w.Layers {
			if holds(l.Layout, block) {
				return w.WindowID, nil
			}
			for _, b := range l.Blocks {
				if b.BlockID == block {
					return w.WindowID, nil
				}
			}
		}
	}
	return "", nil
}

func holds(n *client.LayoutNode, block string) bool {
	if n == nil {
		return false
	}
	if n.BlockID == block {
		return true
	}
	return n.Split != nil && (holds(&n.Split.Before, block) || holds(&n.Split.After, block))
}

// WatchSize polls path, a size file modalshell writes, every interval and
// sends each size that differs from the last one sent. The channel closes
// when ctx ends.
func WatchSize(ctx context.Context, path string, every time.Duration) <-chan client.Size {
	out := make(chan client.Size)
	go func() {
		defer close(out)
		t := time.NewTicker(every)
		defer t.Stop()
		var last client.Size
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
			cols, rows, err := modalshell.ReadSize(path)
			if err != nil {
				continue
			}
			sz := client.Size{Columns: cols, Rows: rows}
			if sz == last {
				continue
			}
			select {
			case out <- sz:
				last = sz
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// ServerAddr is the address of the Rex server: REX_SERVER when set, else
// the per-user socket Rex listens on under $XDG_RUNTIME_DIR (Linux) when it
// exists, else rexuikit's default (macOS).
func ServerAddr(getenv func(string) string) (string, error) {
	if a := getenv("REX_SERVER"); a != "" {
		return a, nil
	}
	if d := getenv("XDG_RUNTIME_DIR"); d != "" {
		p := filepath.Join(d, "rex", "server.sock")
		if _, err := os.Stat(p); err == nil {
			return (&url.URL{Scheme: "unix", Path: p}).String(), nil
		}
	}
	return client.ServerAddr(getenv)
}
