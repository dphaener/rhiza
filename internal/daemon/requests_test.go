package daemon

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	if err != nil || len(got) != 2 || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, %v; want %v", got, err, want)
	}
	if got[0].Intent() != "Fix Login fails" || got[1].Intent() != "x" {
		t.Errorf("intents %q, %q", got[0].Intent(), got[1].Intent())
	}
}

// A request's notes round-trip in order, lines kept; a request without
// notes writes no notes key, so its shape is as before notes.
func TestRequestsNotesRoundTrip(t *testing.T) {
	common := t.TempDir()
	at := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	want := []Request{
		{ID: "a", Kind: KindFeature, Text: "Dark mode", At: at, Notes: []string{"First", "1. one\n2. two\n\n- three"}},
		{ID: "b", Kind: KindBug, Text: "No notes", At: at},
	}
	if err := WriteRequests(common, want); err != nil {
		t.Fatal(err)
	}
	got, err := ReadRequests(common)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, %v; want %#v", got, err, want)
	}
	b, err := os.ReadFile(filepath.Join(Dir(common), RequestsFile))
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Requests []map[string]json.RawMessage `json:"requests"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	if _, ok := doc.Requests[0]["notes"]; !ok {
		t.Errorf("request with notes: no notes key in %s", b)
	}
	if _, ok := doc.Requests[1]["notes"]; ok {
		t.Errorf("request without notes: a notes key in %s", b)
	}
}

func TestNormalizeRequestNote(t *testing.T) {
	for in, want := range map[string]string{
		"plain":                          "plain",
		"\n\n  1. one\r\n2. two\r\n\n\t": "1. one\n2. two",
		"a\n\n  b  ":                     "a\n\n  b",
	} {
		if got, err := NormalizeRequestNote(in); err != nil || got != want {
			t.Errorf("NormalizeRequestNote(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	if got, err := NormalizeRequestNote(strings.Repeat("x", MaxRequestNote)); err != nil || len(got) != MaxRequestNote {
		t.Errorf("a note of %d bytes: %d, %v; want it kept", MaxRequestNote, len(got), err)
	}
	for _, bad := range []string{"", " \n\t\r\n ", strings.Repeat("x", MaxRequestNote+1), "bad \xff byte"} {
		if got, err := NormalizeRequestNote(bad); err == nil {
			t.Errorf("NormalizeRequestNote(%.20q) = %.20q, want refused", bad, got)
		}
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

// Writers that update the queue at once each keep their update: none
// overwrites another's.
func TestUpdateRequestsConcurrent(t *testing.T) {
	common := t.TempDir()
	const n = 20
	errs := make(chan error, n)
	for i := range n {
		go func() {
			errs <- UpdateRequests(common, func(reqs []Request) ([]Request, error) {
				time.Sleep(time.Millisecond)
				return append(reqs, Request{ID: string(rune('a' + i))}), nil
			})
		}()
	}
	for range n {
		if err := <-errs; err != nil {
			t.Fatal(err)
		}
	}
	got, err := ReadRequests(common)
	if err != nil || len(got) != n {
		t.Fatalf("got %d requests (%v), want %d", len(got), err, n)
	}
}
