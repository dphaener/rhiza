package api

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/session"
)

// limitPane is pane text that ends at Claude Code's usage-limit message.
const limitPane = "working on it\nClaude usage limit reached|1700000000\n> "

// session.pane reads the session's block with `rex capture` in its
// worktree, and asks the session's own adapter whether the text shows a
// usage-limit prompt.
func TestSessionPaneCapturesAndAsksTheAdapter(t *testing.T) {
	f := newMessageFixture(t)
	f.rex.out = limitPane
	resp, err := callOp(t, f.ctx, "session.pane", `{"session":"`+f.sid[:8]+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionPaneResponse)
	if r.Session != f.sid || r.Text != limitPane || !r.UsageLimit {
		t.Errorf("response %+v", r)
	}
	if want := [][]string{{"rex", "capture", "-b", "block:b1"}}; !reflect.DeepEqual(f.rex.calls, want) {
		t.Errorf("calls = %q, want %q", f.rex.calls, want)
	}
	if want := []string{f.worktree}; !reflect.DeepEqual(f.rex.dirs, want) {
		t.Errorf("dirs = %q, want %q", f.rex.dirs, want)
	}

	f.rex.out = "all fine\n> "
	resp, err = callOp(t, f.ctx, "session.pane", `{"session":"`+f.sid+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionPaneResponse); r.UsageLimit || r.Text != "all fine\n> " {
		t.Errorf("no limit: response %+v", r)
	}
}

// An agent whose adapter cannot read a pane never reports a usage limit,
// whatever the text says.
func TestSessionPaneAdapterWithoutPaneLimiter(t *testing.T) {
	f := newMessageFixture(t)
	if _, err := f.st.Update(f.sid, func(r *session.Record) { r.Agent = askAgent }); err != nil {
		t.Fatal(err)
	}
	f.rex.out = limitPane
	resp, err := callOp(t, f.ctx, "session.pane", `{"session":"`+f.sid+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionPaneResponse); r.UsageLimit || r.Text != limitPane {
		t.Errorf("response %+v", r)
	}
}

// A capture rex fails is an error naming the session.
func TestSessionPaneRexFails(t *testing.T) {
	f := newMessageFixture(t)
	f.rex.err = errors.New("rex capture: block block:b1 not found")
	_, err := callOp(t, f.ctx, "session.pane", `{"session":"`+f.sid+`"}`)
	if err == nil || !strings.Contains(AsError(err).Message, "rex could not read session "+f.sid) {
		t.Fatalf("err = %v", err)
	}
}

// session.key presses escape or enter in the session's block, in any case.
func TestSessionKeyPresses(t *testing.T) {
	f := newMessageFixture(t)
	for _, key := range []string{"escape", " Enter "} {
		f.rex.calls, f.rex.dirs = nil, nil
		resp, err := callOp(t, f.ctx, "session.key", `{"session":"`+f.sid[:8]+`","key":"`+key+`"}`)
		if err != nil {
			t.Fatalf("%s: %v", key, err)
		}
		want := strings.ToLower(strings.TrimSpace(key))
		if r := resp.(SessionKeyResponse); r.Session != f.sid || r.Key != want || r.DryRun {
			t.Errorf("%s: response %+v", key, r)
		}
		if calls := [][]string{{"rex", "send-key", "-b", "block:b1", want}}; !reflect.DeepEqual(f.rex.calls, calls) {
			t.Errorf("%s: calls = %q, want %q", key, f.rex.calls, calls)
		}
		if !reflect.DeepEqual(f.rex.dirs, []string{f.worktree}) {
			t.Errorf("%s: dirs = %q", key, f.rex.dirs)
		}
	}
}

// Only escape and enter are keys: anything else is a usage error and
// nothing is sent.
func TestSessionKeyRefusesOtherKeys(t *testing.T) {
	f := newMessageFixture(t)
	for _, key := range []string{"", "tab", "y", "ctrl-c", "escape enter"} {
		_, err := callOp(t, f.ctx, "session.key", `{"session":"`+f.sid+`","key":"`+key+`"}`)
		wantExit(t, err, 2)
	}
	f.noRex(t, "other keys")
}

// A dry run names the key and sends nothing.
func TestSessionKeyDryRun(t *testing.T) {
	f := newMessageFixture(t)
	f.ctx.DryRun = true
	resp, err := callOp(t, f.ctx, "session.key", `{"session":"`+f.sid+`","key":"escape"}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(SessionKeyResponse); r.Session != f.sid || r.Key != "escape" || !r.DryRun {
		t.Errorf("response %+v", r)
	}
	f.noRex(t, "dry run")
}

// A session that has exited, or has no tab, is a conflict for both
// operations, and rex is never run; an unknown session is not found.
func TestSessionPaneAndKeyRefuseSessions(t *testing.T) {
	for _, c := range []struct {
		name   string
		update func(*session.Record)
	}{
		{"exited", func(r *session.Record) {
			code := 0
			r.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: time.Now().UTC()}
		}},
		{"tabless", func(r *session.Record) { r.Tab.Block = "" }},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := newMessageFixture(t)
			if _, err := f.st.Update(f.sid, c.update); err != nil {
				t.Fatal(err)
			}
			_, err := callOp(t, f.ctx, "session.pane", `{"session":"`+f.sid+`"}`)
			wantExit(t, err, 3)
			_, err = callOp(t, f.ctx, "session.key", `{"session":"`+f.sid+`","key":"escape"}`)
			wantExit(t, err, 3)
			f.noRex(t, c.name)
		})
	}
	f := newMessageFixture(t)
	_, err := callOp(t, f.ctx, "session.pane", `{"session":"deadbeef"}`)
	wantExit(t, err, 4)
	_, err = callOp(t, f.ctx, "session.key", `{"session":"deadbeef","key":"enter"}`)
	wantExit(t, err, 4)
	f.noRex(t, "unknown session")
}
