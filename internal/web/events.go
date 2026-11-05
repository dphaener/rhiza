package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"
)

// The kinds of event the stream names. Each says which part of the page's
// data moved; the page re-fetches it. The event's data is {"id": ...}: the
// one thing that moved (a log entry, a run, an inbox item), or "daemon".
const (
	EventChange = "change" // a change set moved or landed: the log grew
	EventRun    = "run"    // a run's journal moved
	EventLog    = "log"    // a run's driver log grew
	EventInbox  = "inbox"  // an inbox item arrived or left
	EventDaemon = "daemon" // the daemon's status changed: paused, slots or health
)

// Source reports one kind of thing as it is now: each thing's ID with a
// fingerprint that differs whenever the thing changed.
type Source func() map[string]string

// Sources are the Sources of an events stream, by event kind.
type Sources map[string]Source

// PollEvery is how often a stream looks for changes.
const PollEvery = time.Second

// Events serves GET /events: a server-sent event stream. A connection
// starts from how things are then and sends an event, named by kind, for
// each thing that is new, changed or gone after it, so an open page knows to
// re-fetch. It looks every `every`, and ends with the request.
func Events(src Sources, every time.Duration) http.Handler {
	kinds := make([]string, 0, len(src))
	for k := range src {
		kinds = append(kinds, k)
	}
	slices.Sort(kinds)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "only GET is allowed", http.StatusMethodNotAllowed)
			return
		}
		fl, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming is not supported", http.StatusInternalServerError)
			return
		}
		seen := map[string]map[string]string{}
		for _, k := range kinds {
			seen[k] = src[k]()
		}
		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("X-Accel-Buffering", "no")
		fmt.Fprint(w, ": ready\n\n")
		fl.Flush()
		t := time.NewTicker(every)
		defer t.Stop()
		for {
			select {
			case <-r.Context().Done():
				return
			case <-t.C:
			}
			sent := false
			for _, k := range kinds {
				now := src[k]()
				for _, id := range changed(seen[k], now) {
					b, _ := json.Marshal(map[string]string{"id": id})
					fmt.Fprintf(w, "event: %s\ndata: %s\n\n", k, b)
					sent = true
				}
				seen[k] = now
			}
			if !sent {
				fmt.Fprint(w, ": keep-alive\n\n")
			}
			fl.Flush()
		}
	})
}

// changed is the IDs in either map whose fingerprint differs, sorted.
func changed(before, now map[string]string) []string {
	var ids []string
	for id, f := range now {
		if old, ok := before[id]; !ok || old != f {
			ids = append(ids, id)
		}
	}
	for id := range before {
		if _, ok := now[id]; !ok {
			ids = append(ids, id)
		}
	}
	slices.Sort(ids)
	return ids
}
