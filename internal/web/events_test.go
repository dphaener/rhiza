package web_test

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// lockedRecorder is a recorder that can be read while a handler writes it.
type lockedRecorder struct {
	mu  sync.Mutex
	rec *httptest.ResponseRecorder
}

func (l *lockedRecorder) Header() http.Header { return l.rec.Header() }
func (l *lockedRecorder) WriteHeader(c int)   { l.rec.WriteHeader(c) }
func (l *lockedRecorder) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rec.Write(b)
}
func (l *lockedRecorder) Flush() {}
func (l *lockedRecorder) body() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rec.Body.String()
}

// waitFor waits for the stream to hold s.
func waitFor(t *testing.T, l *lockedRecorder, s string) {
	t.Helper()
	for end := time.Now().Add(10 * time.Second); time.Now().Before(end); time.Sleep(5 * time.Millisecond) {
		if strings.Contains(l.body(), s) {
			return
		}
	}
	t.Fatalf("no %q in the stream:\n%s", s, l.body())
}

func TestEventsNameWhatChanged(t *testing.T) {
	ctx := fixtureCtx(t)
	out, err := ctx.Git("rev-parse", "--git-common-dir")
	if err != nil {
		t.Fatal(err)
	}
	common := filepath.Join(ctx.Root, strings.TrimSpace(out))
	h := web.Handler(nil, web.Events(api.WebEventSources(ctx), 10*time.Millisecond))

	rc, stop := context.WithCancel(context.Background())
	defer stop()
	l := &lockedRecorder{rec: httptest.NewRecorder()}
	done := make(chan struct{})
	go func() {
		defer close(done)
		h.ServeHTTP(l, httptest.NewRequest("GET", "/events", nil).WithContext(rc))
	}()
	waitFor(t, l, ": ready")
	if ct := l.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("content type %q", ct)
	}

	write := func(path, s string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		f.WriteString(s)
	}
	expect := func(event, id string) {
		t.Helper()
		waitFor(t, l, "event: "+event+"\ndata: {\"id\":\""+id+"\"}\n\n")
	}

	// A log entry.
	write(filepath.Join(ctx.Root, ".rhiza", "log", "01ZZZZZZZZZZZZZZZZZZZZZZZZ", "a.json"), "{}\n")
	expect("change", "01ZZZZZZZZZZZZZZZZZZZZZZZZ")
	// Another entry in the same change set: it moved or landed.
	moved := "event: change\ndata: {\"id\":\"01ZZZZZZZZZZZZZZZZZZZZZZZZ\"}\n\n"
	write(filepath.Join(ctx.Root, ".rhiza", "log", "01ZZZZZZZZZZZZZZZZZZZZZZZZ", "b.json"), "{}\n")
	for end := time.Now().Add(10 * time.Second); strings.Count(l.body(), moved) < 2; time.Sleep(5 * time.Millisecond) {
		if time.Now().After(end) {
			t.Fatalf("no change event for an entry appended to an existing change set:\n%s", l.body())
		}
	}
	// A run journal moves.
	write(filepath.Join(common, "rhiza", "runs", "run1.jsonl"), "{}\n")
	expect("run", "run1")
	write(filepath.Join(common, "rhiza", "runs", "run1.jsonl"), "{}\n")
	// Driver log lines.
	write(filepath.Join(common, "rhiza", "logs", "run1.log"), "hello\n")
	expect("log", "run1")
	// The daemon's status changes.
	st := daemon.State{Runs: []daemon.LiveRun{}, MaxRuns: 2}
	if err := daemon.WriteState(common, st); err != nil {
		t.Fatal(err)
	}
	expect("daemon", "daemon")
	// An inbox item arrives, then leaves.
	st.DiskLow = &daemon.DiskLow{Since: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC), Free: 1, Min: 2}
	if err := daemon.WriteState(common, st); err != nil {
		t.Fatal(err)
	}
	items := api.DiskLowInboxIDs(st.DiskLow)
	expect("inbox", items[0])
	n := strings.Count(l.body(), "event: inbox")
	st.DiskLow = nil
	if err := daemon.WriteState(common, st); err != nil {
		t.Fatal(err)
	}
	for end := time.Now().Add(10 * time.Second); strings.Count(l.body(), "event: inbox") == n; time.Sleep(5 * time.Millisecond) {
		if time.Now().After(end) {
			t.Fatal("no inbox event when the item left")
		}
	}

	stop()
	<-done
	// The stream carries nothing before the changes: only named events.
	sc := bufio.NewScanner(strings.NewReader(l.body()))
	for sc.Scan() {
		if line := sc.Text(); line != "" && !strings.HasPrefix(line, ":") && !strings.HasPrefix(line, "event: ") && !strings.HasPrefix(line, "data: ") {
			t.Errorf("stray line %q", line)
		}
	}
}

func TestEventsRefuseWrites(t *testing.T) {
	h := web.Handler(nil, web.Events(web.Sources{}, time.Second))
	for _, m := range []string{"POST", "PUT", "DELETE"} {
		if rec := serve(h, m, "/events"); rec.Code != 405 {
			t.Errorf("%s /events: %d", m, rec.Code)
		}
	}
}

// The run page reads run.status through the daemon's read-only API, so a
// real run journal shows on /runs/<id> and its fragment.
func TestRunPageThroughReadOnlyAPI(t *testing.T) {
	ctx := fixtureCtx(t)
	out, err := ctx.Git("rev-parse", "--git-common-dir")
	if err != nil {
		t.Fatal(err)
	}
	common := filepath.Join(ctx.Root, strings.TrimSpace(out))
	if err := run.NewStore(run.Dir(common)).Create("run1", run.Entry{Type: run.TypeStarted, At: time.Now(), Started: &run.Started{Workflow: "build", Title: "Build the thing"}}); err != nil {
		t.Fatal(err)
	}
	h := web.Handler(api.ReadOnlyHandler(ctx))
	if rec := serve(h, "GET", "/runs/run1"); rec.Code != 200 || !strings.Contains(rec.Body.String(), "Build the thing") {
		t.Errorf("page: %d\n%s", rec.Code, rec.Body)
	}
	if body := serve(h, "GET", "/fragments/run-detail?run=run1").Body.String(); strings.Contains(body, "No such run") {
		t.Errorf("fragment:\n%s", body)
	}
}

// A run's log tail follows the driver log: lines appended to it move the
// stream's log event and show in the tail fragment; a run without a log
// reads as such.
func TestRunLogTailFollowsDriverLog(t *testing.T) {
	ctx := fixtureCtx(t)
	out, err := ctx.Git("rev-parse", "--git-common-dir")
	if err != nil {
		t.Fatal(err)
	}
	common := filepath.Join(ctx.Root, strings.TrimSpace(out))
	if err := run.NewStore(run.Dir(common)).Create("run1", run.Entry{Type: run.TypeStarted, At: time.Now(), Started: &run.Started{Workflow: "build"}}); err != nil {
		t.Fatal(err)
	}
	h := web.Handler(api.ReadOnlyHandler(ctx), web.Events(api.WebEventSources(ctx), 10*time.Millisecond))
	frag := func() string { return serve(h, "GET", "/fragments/run-log?run=run1").Body.String() }

	if body := frag(); !strings.Contains(body, "No log") || strings.Contains(body, "rz-term") {
		t.Errorf("run without a log:\n%s", body)
	}

	rc, stop := context.WithCancel(context.Background())
	defer stop()
	l := &lockedRecorder{rec: httptest.NewRecorder()}
	done := make(chan struct{})
	go func() {
		defer close(done)
		h.ServeHTTP(l, httptest.NewRequest("GET", "/events", nil).WithContext(rc))
	}()
	waitFor(t, l, ": ready")

	logPath := filepath.Join(common, "rhiza", "logs", "run1.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		t.Fatal(err)
	}
	appendLines := func(s string) {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		f.WriteString(s)
	}
	appendLines("first line\nsecond <line>\n")
	waitFor(t, l, "event: log\ndata: {\"id\":\"run1\"}\n\n")
	body := frag()
	for _, want := range []string{"rz-term", "first line", "second &lt;line&gt;", "rz-cursor"} {
		if !strings.Contains(body, want) {
			t.Errorf("tail lacks %q:\n%s", want, body)
		}
	}
	n := strings.Count(l.body(), "event: log")
	appendLines("third line\n")
	for end := time.Now().Add(10 * time.Second); strings.Count(l.body(), "event: log") == n; time.Sleep(5 * time.Millisecond) {
		if time.Now().After(end) {
			t.Fatal("no log event for appended lines")
		}
	}
	if body := frag(); !strings.Contains(body, "third line") {
		t.Errorf("tail lacks the appended line:\n%s", body)
	}
	stop()
	<-done
}
