package api

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestFeedFollowTicks drives follow with a replaced ticker: it starts at the
// end of the feed, emits a change set landing between two ticks once, and
// stops on an interrupt.
func TestFeedFollowTicks(t *testing.T) {
	ctx := initCtx(t)
	ada := recordAda(t, ctx, "ada@example.com")
	rec := *ctx
	rec.Record = true
	rec.Actor = ada
	land := func(handle string) {
		t.Helper()
		if _, err := callOp(t, &rec, "add", fmt.Sprintf(`{"kind":"capability","handle":%q}`, handle)); err != nil {
			t.Fatal(err)
		}
	}
	land("before")

	ticks := make(chan time.Time)
	old := feedTicker
	feedTicker = func() (<-chan time.Time, func()) { return ticks, func() {} }
	t.Cleanup(func() { feedTicker = old })
	stop := make(chan struct{})
	out := &syncBuf{}
	f := *ctx
	f.Out, f.JSON = out, true
	f.Interrupts = func() (<-chan struct{}, func()) { return stop, func() {} }
	done := make(chan error, 1)
	go func() {
		resp, err := callOp(t, &f, "feed", `{"follow":true}`)
		if err == nil {
			if s, ok := resp.(Streamer); !ok || !s.Streamed() {
				err = fmt.Errorf("response %T is not streamed", resp)
			}
		}
		done <- err
	}()
	tick := func() {
		t.Helper()
		select {
		case ticks <- time.Now():
		case err := <-done:
			t.Fatalf("follow stopped early: %v", err)
		}
	}

	tick() // the first poll is over: nothing landed since the start
	if got := out.String(); got != "" {
		t.Fatalf("emitted before any land: %q", got)
	}
	land("between")
	tick()
	tick() // a second poll must not repeat it
	tick()
	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 1 || !strings.Contains(lines[0], `"title":"add capability between"`) || !strings.HasPrefix(lines[0], `{"actor":`) {
		t.Fatalf("want one NDJSON line for 'between', got %q", out.String())
	}
	close(stop)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("follow did not stop on interrupt")
	}
}

// TestFeedFollowLimit emits the backlog from a since and stops at the limit.
func TestFeedFollowLimit(t *testing.T) {
	ctx := initCtx(t)
	ada := recordAda(t, ctx, "ada@example.com")
	rec := *ctx
	rec.Record = true
	rec.Actor = ada
	for _, h := range []string{"a", "b", "c", "d"} {
		if _, err := callOp(t, &rec, "add", fmt.Sprintf(`{"kind":"capability","handle":%q}`, h)); err != nil {
			t.Fatal(err)
		}
	}
	out := &syncBuf{}
	f := *ctx
	f.Out = out
	if _, err := callOp(t, &f, "feed", `{"follow":true,"since":"2026-10-24T00:00:00Z","limit":3}`); err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(out.String(), "\n"); n != 3 || strings.Count(out.String(), "- ") != 3 {
		t.Fatalf("want 3 human lines, got %q", out.String())
	}
}

// syncBuf is a buffer a test reads while follow writes to it.
type syncBuf struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuf) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuf) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
