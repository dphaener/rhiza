package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/dphaener/rhiza/internal/web"
)

// healthRefreshEvery is how old the web UI's held health report may get
// before a request starts a new refresh in the background.
const healthRefreshEvery = 30 * time.Second

// healthCache holds the last health report the web UI has, refreshed in the
// background: never on a request's path, and no more than one at a time.
type healthCache struct {
	read  func() []web.HealthPill
	every time.Duration
	now   func() time.Time

	mu        sync.Mutex
	pills     []web.HealthPill
	have      bool
	at        time.Time
	refreshes bool
}

// get returns the held report (have is false until the first has landed) and
// starts a refresh when there is none or it is older than the interval.
func (h *healthCache) get() (pills []web.HealthPill, have bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.refreshes && (!h.have || h.now().Sub(h.at) >= h.every) {
		h.refreshes = true
		go h.refresh()
	}
	return h.pills, h.have
}

// fingerprint is the held report as the /events stream compares it: it
// moves when a report lands that differs from the one held, so an open page
// re-fetches its daemon section. Like get, it starts a refresh when the
// report is missing or old, so an open stream keeps the report fresh.
func (h *healthCache) fingerprint() string {
	pills, have := h.get()
	return fmt.Sprintf("%v|%v", have, pills)
}

func (h *healthCache) refresh() {
	pills := h.read()
	h.mu.Lock()
	h.pills, h.have, h.at, h.refreshes = pills, true, h.now(), false
	h.mu.Unlock()
}

// WebHealth is the web UI's health report, held in memory and refreshed in
// the background. The Overview's daemon section reads it, and the /events
// stream's daemon event moves when a new report lands.
type WebHealth struct{ cache *healthCache }

// NewWebHealth holds the health report of the repository base works in.
func NewWebHealth(base *Ctx) *WebHealth {
	return newWebHealth(func() []web.HealthPill { return readHealthPills(base) }, healthRefreshEvery, time.Now)
}

func newWebHealth(read func() []web.HealthPill, every time.Duration, now func() time.Time) *WebHealth {
	return &WebHealth{cache: &healthCache{read: read, every: every, now: now}}
}

// WebDaemonSource reads the daemon and its health for the web UI's Overview,
// from the repository base works in. Running, paused and slots are read as the
// files are; health is the last report health holds.
func WebDaemonSource(base *Ctx, health *WebHealth) web.DaemonSource {
	cache := health.cache
	return func() web.DaemonView {
		c := quietCtx(base)
		var v web.DaemonView
		if out, err := runDaemonStatus(c, DaemonStatusRequest{}); err == nil {
			s := out.(DaemonStatusResponse)
			v.Running, v.Limit = s.Running, s.MaxRuns
			for _, run := range s.Runs {
				if run.State == "running" {
					v.Slots++
				}
			}
			for _, set := range s.Settings {
				if set.Key == "rhiza.daemon.paused" {
					v.Paused = set.Value == "true"
				}
			}
		}
		pills, have := cache.get()
		v.Health, v.HealthPending = pills, !have
		return v
	}
}

func quietCtx(base *Ctx) *Ctx {
	c := *base
	c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
	c.Stdin, c.Rex, c.Start, c.Detach, c.Interrupts = nil, nil, nil, nil, nil
	return &c
}

func readHealthPills(base *Ctx) []web.HealthPill {
	var pills []web.HealthPill
	if out, err := runHealth(quietCtx(base), nil); err == nil {
		for _, p := range out.(HealthResponse).Parts {
			pills = append(pills, web.HealthPill{Part: p.Part, State: p.State, Summary: p.Summary})
		}
	}
	return pills
}
