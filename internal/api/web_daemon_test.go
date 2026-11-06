package api

import (
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
	src := webDaemonSource(&Ctx{Root: t.TempDir()}, read, 30*time.Second, now)
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

func waitFor(t *testing.T, ok func() bool) {
	t.Helper()
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); time.Sleep(10 * time.Millisecond) {
		if ok() {
			return
		}
	}
	t.Fatal("condition not met")
}
