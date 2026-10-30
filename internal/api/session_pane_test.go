package api

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/session"
)

// paneRex is a fake rex for openSessionPane: it answers `block call size`
// with size (or sizeErr), `split` with splitOut (or splitErr) and `window
// new` with a window, and logs each call.
type paneRex struct {
	size     string
	sizeErr  error
	splitOut string
	splitErr error
	blocks   string
	// splitErrOnce fails only the first split with splitErr.
	splitErrOnce bool
	calls        []string
}

func (p *paneRex) exec(name string, args []string, _ string, _ []string) ([]byte, error) {
	p.calls = append(p.calls, strings.Join(args, " "))
	switch args[0] {
	case "block":
		if args[1] == "ls" {
			return []byte(p.blocks), nil
		}
		return []byte(p.size), p.sizeErr
	case "split":
		err := p.splitErr
		if p.splitErrOnce {
			p.splitErr = nil
		}
		return []byte(p.splitOut), err
	}
	return []byte(`{"window_id":"window:w2","block_ids":["block:new"]}`), nil
}

// openSessionPane splits the anchor along its longer side in pixels,
// unless the step's split fixes the direction, in which case Rex is not
// asked the size; an unreadable size splits right; an anchor Rex no
// longer has splits another pane of its window.
func TestOpenSessionPane(t *testing.T) {
	anchor := session.Tab{Window: "window:w1", Block: "block:b1"}
	landscape := `{"columns":200,"rows":50,"cell_width_px":8,"cell_height_px":16}`
	portrait := `{"columns":100,"rows":120,"cell_width_px":8,"cell_height_px":16}`
	split := func(dir string) string {
		return "split -b block:b1 --split " + dir + " --ratio 50 -c /wt --json --label verifier a #1 -- rhiza session run S1"
	}
	size := "block call -b block:b1 size"
	ls := "block ls -w window:w1 --json"
	split2 := "split -b block:b9 --split right --ratio 50 -c /wt --json --label verifier a #1 -- rhiza session run S1"
	size2 := "block call -b block:b9 size"
	for _, c := range []struct {
		name      string
		rex       paneRex
		override  string
		calls     []string
		block     string
		window    string
		wantError bool
	}{
		{"landscape splits right", paneRex{size: landscape, splitOut: `{"block_id":"block:b2","window_id":"window:w1"}`},
			"auto", []string{size, split("right")}, "block:b2", "window:w1", false},
		{"portrait splits below", paneRex{size: portrait, splitOut: `{"block_id":"block:b2"}`},
			"auto", []string{size, split("below")}, "block:b2", "window:w1", false},
		{"override below on landscape", paneRex{size: landscape, splitOut: `{"block_id":"block:b2"}`},
			"below", []string{split("below")}, "block:b2", "window:w1", false},
		{"override right on portrait", paneRex{size: portrait, splitOut: `{"block_ids":["block:b3"]}`},
			"right", []string{split("right")}, "block:b3", "window:w1", false},
		{"size unreadable splits right", paneRex{sizeErr: errors.New("rex: no such method"), splitOut: `{"block_id":"block:b2"}`},
			"auto", []string{size, split("right")}, "block:b2", "window:w1", false},
		{"anchor gone splits another pane", paneRex{size: landscape, splitErr: errors.New("rex split: block block:b1 not found"), splitErrOnce: true,
			blocks: `{"blocks":[{"block_id":"block:b1"},{"block_id":"block:b9"}]}`, splitOut: `{"block_id":"block:b2"}`},
			"auto", []string{size, split("right"), ls, size2, split2}, "block:b2", "window:w1", false},
		{"anchor gone, no other pane", paneRex{size: landscape, splitErr: errors.New("rex split: block block:b1 not found"),
			blocks: `{"blocks":[{"block_id":"block:b1"}]}`},
			"auto", []string{size, split("right"), ls}, "", "", true},
		{"split fails", paneRex{size: landscape, splitErr: errors.New("rex split: server down")},
			"auto", []string{size, split("right")}, "", "", true},
		{"split prints no block", paneRex{size: landscape, splitOut: `{}`},
			"auto", []string{size, split("right")}, "", "", true},
	} {
		rex := c.rex
		ctx := &Ctx{Exec: rex.exec}
		tab, _, err := openSessionPane(ctx, session.Tab{Title: "T"}, anchor, c.override, "verifier a #1", "/wt", "rhiza", "S1")
		if (err != nil) != c.wantError {
			t.Errorf("%s: err = %v", c.name, err)
		}
		if !reflect.DeepEqual(rex.calls, c.calls) {
			t.Errorf("%s: calls =\n%q\nwant\n%q", c.name, rex.calls, c.calls)
		}
		if !c.wantError && (tab.Block != c.block || tab.Window != c.window) {
			t.Errorf("%s: tab = %+v, want block %s window %s", c.name, tab, c.block, c.window)
		}
	}
}

// Without an anchor a session opens a new tab, as every session not of a
// leaves step does.
func TestOpenSessionPaneNewTab(t *testing.T) {
	rex := paneRex{}
	tab, args, err := openSessionPane(&Ctx{Exec: rex.exec}, session.Tab{Title: "T"}, session.Tab{}, "auto", "", "/wt", "rhiza", "S1")
	if err != nil {
		t.Fatal(err)
	}
	if want := "window new T -c /wt --json -- rhiza session run S1"; strings.Join(args, " ") != want || len(rex.calls) != 1 {
		t.Errorf("args = %q, calls = %q", args, rex.calls)
	}
	if tab.Block != "block:new" || tab.Window != "window:w2" {
		t.Errorf("tab = %+v", tab)
	}
}

// A pane closes with `rex block close`; one Rex no longer has is already
// closed, and a session with no block has none to close.
func TestCloseSessionPane(t *testing.T) {
	var calls []string
	errNext := error(nil)
	ctx := &Ctx{Exec: func(_ string, args []string, _ string, _ []string) ([]byte, error) {
		calls = append(calls, strings.Join(args, " "))
		return nil, errNext
	}}
	rec := session.Record{Tab: session.Tab{Session: "s1", Block: "block:b2"}}
	if closed, err := closeSessionPane(ctx, rec, "/"); err != nil || !closed {
		t.Errorf("close = %v, %v", closed, err)
	}
	if want := []string{"block close -s s1 block:b2"}; !reflect.DeepEqual(calls, want) {
		t.Errorf("calls = %q, want %q", calls, want)
	}
	errNext = errors.New("rex block close: block block:b2 not found")
	if closed, err := closeSessionPane(ctx, rec, "/"); err != nil || closed {
		t.Errorf("gone: close = %v, %v", closed, err)
	}
	errNext = errors.New("rex block close: server down")
	if _, err := closeSessionPane(ctx, rec, "/"); err == nil {
		t.Errorf("failure: no error")
	}
	if closed, err := closeSessionPane(ctx, session.Record{}, "/"); err != nil || closed || len(calls) != 3 {
		t.Errorf("no block: close = %v, %v, calls %q", closed, err, calls)
	}
}
