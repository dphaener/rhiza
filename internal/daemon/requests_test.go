package daemon

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

func TestRequestsRoundTrip(t *testing.T) {
	common := t.TempDir()
	got, err := ReadRequests(common)
	if err != nil || len(got) != 0 {
		t.Fatalf("no file: %v, %v", got, err)
	}
	at := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	want := []Request{{ID: "a", Kind: KindBug, Text: "Login fails", At: at}, {ID: "b", Kind: KindFeature, Text: "x", At: at, Run: "r", StartedAt: at}}
	if err := WriteRequests(common, want); err != nil {
		t.Fatal(err)
	}
	got, err = ReadRequests(common)
	if err != nil || len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v, %v; want %v", got, err, want)
	}
	if got[0].Intent() != "Fix: Login fails" || got[1].Intent() != "x" {
		t.Errorf("intents %q, %q", got[0].Intent(), got[1].Intent())
	}
}

func TestRequestState(t *testing.T) {
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	journal := func(id string, end string) run.Journal {
		es := []run.Entry{{Type: run.TypeStarted, At: now, Step: "a", Started: &run.Started{Workflow: "plan", Change: "cs-" + id}}}
		if end != "" {
			es = append(es, run.Entry{Type: run.TypeFinished, At: now, End: end})
		}
		return run.Journal{ID: id, Entries: es}
	}
	runs := []run.Journal{journal("live", ""), journal("ok", run.Succeeded), journal("bad", run.Failed)}
	for _, tc := range []struct {
		run, state, change string
	}{
		{"", RequestQueued, ""},
		{"none", RequestPlanning, ""},
		{"live", RequestPlanning, "cs-live"},
		{"ok", RequestPlanned, "cs-ok"},
		{"bad", RequestFailed, "cs-bad"},
	} {
		state, change := RequestState(Request{ID: "q", Run: tc.run}, runs)
		if state != tc.state || change != tc.change {
			t.Errorf("run %q: %s %q, want %s %q", tc.run, state, change, tc.state, tc.change)
		}
	}
}
