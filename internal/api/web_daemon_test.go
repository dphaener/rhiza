package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/web"
)

// The Overview and its daemon fragment answer at once while health blocks:
// the pills show pending, concurrent requests start one refresh, the report
// shows once it lands, and an old report is served while one new refresh
// starts.
func TestWebDaemonHealthCachedInBackground(t *testing.T) {
	var started atomic.Int32
	release := make(chan struct{})
	var mu sync.Mutex
	state := "ok"
	read := func() []web.HealthPill {
		started.Add(1)
		<-release
		mu.Lock()
		defer mu.Unlock()
		return []web.HealthPill{{Part: "model", State: state}}
	}
	clock := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { mu.Lock(); defer mu.Unlock(); return clock }
	src := WebDaemonSource(&Ctx{Root: t.TempDir()}, newWebHealth(read, 30*time.Second, now))
	h := web.HandlerWith(src, nil)
	get := func(path string) string {
		start := time.Now()
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if d := time.Since(start); d > 400*time.Millisecond {
			t.Errorf("%s took %v", path, d)
		}
		return w.Body.String()
	}

	for _, p := range []string{"/", "/fragments/daemon"} {
		body := get(p)
		if !strings.Contains(body, "rz-pill--pending") || strings.Contains(body, "rz-pill--failed\" title=\"model") {
			t.Errorf("%s: want pending pills:\n%s", p, body)
		}
	}
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() { defer wg.Done(); get("/fragments/daemon") }()
	}
	wg.Wait()
	// The refresh runs on its own goroutine, so wait for it to start, then
	// check after a pause that no second one followed.
	waitFor(t, func() bool { return started.Load() >= 1 })
	time.Sleep(50 * time.Millisecond)
	if n := started.Load(); n != 1 {
		t.Fatalf("%d refreshes started, want 1", n)
	}

	release <- struct{}{}
	waitFor(t, func() bool { return strings.Contains(get("/fragments/daemon"), `title="model: ok`) })
	if strings.Contains(get("/"), "rz-pill--pending") {
		t.Error("pending after the report landed")
	}

	mu.Lock()
	clock = clock.Add(time.Minute)
	state = "warn"
	mu.Unlock()
	for range 5 {
		if body := get("/fragments/daemon"); !strings.Contains(body, `title="model: ok`) {
			t.Errorf("old report not served:\n%s", body)
		}
	}
	waitFor(t, func() bool { return started.Load() == 2 })
	time.Sleep(50 * time.Millisecond)
	if n := started.Load(); n != 2 {
		t.Fatalf("%d refreshes started, want 2", n)
	}
	release <- struct{}{}
	waitFor(t, func() bool { return strings.Contains(get("/fragments/daemon"), `title="model: warn`) })
}

// A health report landing in the background sends the open page a daemon
// event, so its pending panel is replaced live; a refresh that finds the same
// report sends none.
func TestWebDaemonHealthLandingSendsDaemonEvent(t *testing.T) {
	release := make(chan struct{})
	var started atomic.Int32
	read := func() []web.HealthPill {
		started.Add(1)
		<-release
		return []web.HealthPill{{Part: "model", State: "ok"}}
	}
	var mu sync.Mutex
	clock := time.Date(2026, 11, 6, 0, 0, 0, 0, time.UTC)
	now := func() time.Time { mu.Lock(); defer mu.Unlock(); return clock }
	ctx := &Ctx{Root: t.TempDir()}
	health := newWebHealth(read, 30*time.Second, now)
	mux := http.NewServeMux()
	mux.Handle("/events", web.Events(WebEventSources(ctx, health), 10*time.Millisecond))
	mux.Handle("/", web.HandlerWith(WebDaemonSource(ctx, health), nil))

	page := httptest.NewRecorder()
	mux.ServeHTTP(page, httptest.NewRequest("GET", "/fragments/daemon", nil))
	if !strings.Contains(page.Body.String(), "rz-pill--pending") {
		t.Fatalf("want the pending panel:\n%s", page.Body.String())
	}

	rc, stop := context.WithCancel(context.Background())
	defer stop()
	l := &streamRecorder{rec: httptest.NewRecorder()}
	go mux.ServeHTTP(l, httptest.NewRequest("GET", "/events", nil).WithContext(rc))
	waitFor(t, func() bool { return strings.Contains(l.body(), ": ready") })
	const landed = "event: daemon\ndata: {\"id\":\"health\"}\n\n"
	time.Sleep(50 * time.Millisecond)
	if strings.Contains(l.body(), landed) {
		t.Fatalf("daemon event before the report landed:\n%s", l.body())
	}

	release <- struct{}{}
	waitFor(t, func() bool { return strings.Contains(l.body(), landed) })
	page = httptest.NewRecorder()
	mux.ServeHTTP(page, httptest.NewRequest("GET", "/fragments/daemon", nil))
	if b := page.Body.String(); strings.Contains(b, "rz-pill--pending") || !strings.Contains(b, `title="model: ok`) {
		t.Errorf("the re-fetched panel is not the report:\n%s", b)
	}

	// An old report: the stream starts a refresh; the same report again
	// sends no event.
	mu.Lock()
	clock = clock.Add(time.Minute)
	mu.Unlock()
	waitFor(t, func() bool { return started.Load() == 2 })
	release <- struct{}{}
	time.Sleep(100 * time.Millisecond)
	if n := strings.Count(l.body(), landed); n != 1 {
		t.Errorf("%d daemon events for health, want 1:\n%s", n, l.body())
	}
}

// streamRecorder is a ResponseRecorder an event stream can write to while
// the test reads it.
type streamRecorder struct {
	mu  sync.Mutex
	rec *httptest.ResponseRecorder
}

func (l *streamRecorder) Header() http.Header { return l.rec.Header() }
func (l *streamRecorder) WriteHeader(c int)   { l.rec.WriteHeader(c) }
func (l *streamRecorder) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rec.Write(b)
}
func (l *streamRecorder) Flush() {}
func (l *streamRecorder) body() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rec.Body.String()
}

func waitFor(t *testing.T, ok func() bool) {
	t.Helper()
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); time.Sleep(10 * time.Millisecond) {
		if ok() {
			return
		}
	}
	t.Fatal("condition not met")
}
