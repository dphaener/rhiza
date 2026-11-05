package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// One inbox call reads the change sets once, however many runs ask for a
// state or title: none of the log with an index, the log once without one.
// Both give the same items, overlaps among them.
func TestInboxReadsChangeSetsOnce(t *testing.T) {
	f := newDaemonFoldFix(t)
	// Two in-flight change sets that both write done-0: an overlap.
	for i := range 2 {
		f.call("change.new", fmt.Sprintf(`{"intent":"Rename %d"}`, i))
		f.call("set", fmt.Sprintf(`{"ref":"done-0","fields":["statement=Done %d"]}`, i))
		f.call("propose", `{}`)
	}
	// Failed runs, one per change set in flight, each a state and a title
	// lookup.
	at := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	var sets []string
	for _, c := range f.call("change.ls", `{}`).(ChangeLsResponse).ChangeSets {
		sets = append(sets, c.ID.String())
	}
	for i, cs := range sets {
		id := fmt.Sprintf("run-%02d", i)
		if err := f.rs.Create(id, run.Entry{Type: run.TypeStarted, At: at.Add(time.Duration(i) * time.Minute), Started: &run.Started{Workflow: "ship-change", Change: cs}}); err != nil {
			t.Fatal(err)
		}
		if err := f.rs.Append(id, run.Entry{Type: run.TypeFinished, At: at.Add(time.Duration(i)*time.Minute + time.Second), End: run.Failed, Detail: "boom"}); err != nil {
			t.Fatal(err)
		}
	}
	inbox := func() []byte {
		t.Helper()
		out, err := json.Marshal(f.call("inbox", `{}`))
		if err != nil {
			t.Fatal(err)
		}
		return out
	}

	f.logReads = 0
	indexed := inbox()
	if f.logReads != 0 {
		t.Errorf("inbox with an index read the log %d times", f.logReads)
	}
	var resp InboxResponse
	if err := json.Unmarshal(indexed, &resp); err != nil {
		t.Fatal(err)
	}
	kinds := map[string]int{}
	for _, it := range resp.Items {
		kinds[it.Kind]++
		if it.Kind == InboxFailed && it.ChangeTitle == "" {
			t.Errorf("failed item %s has no change title", it.ID)
		}
	}
	if kinds[InboxOverlap] != 1 || kinds[InboxFailed] < 5 {
		t.Errorf("got kinds %v, want one overlap and the failed runs", kinds)
	}

	// Without an index the fold reads the log once for every lookup, and
	// answers as the index does.
	byIndex, byLog := newInboxFold(f.ctx), func() *inboxFold {
		c := *f.ctx
		c.IndexPath = ""
		return newInboxFold(&c)
	}()
	defer byIndex.close()
	f.logReads = 0
	for _, cs := range append(sets, "01M6KV89S5MKWG2DHA762F8WTZ", "nope") {
		si, oki := byIndex.stateOf(cs)
		sl, okl := byLog.stateOf(cs)
		if si != sl || oki != okl {
			t.Errorf("state of %s: index %q %v, log %q %v", cs, si, oki, sl, okl)
		}
		if ti, tl := byIndex.titleOf(cs), byLog.titleOf(cs); ti != tl {
			t.Errorf("title of %s: index %q, log %q", cs, ti, tl)
		}
	}
	oi, err := byIndex.overlaps()
	if err != nil {
		t.Fatal(err)
	}
	ol, err := byLog.overlaps()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(oi, ol) || len(oi) != 1 {
		t.Errorf("overlaps: index %+v, log %+v", oi, ol)
	}
	if f.logReads != 1 {
		t.Errorf("the log fold read the log %d times, want 1", f.logReads)
	}
}

// An overlap between change sets with no run is notified by the change
// sets' titles from the fold, never by either ID.
func TestDaemonNotifyOverlapByTitle(t *testing.T) {
	f := newDaemonFoldFix(t)
	var rex bytes.Buffer
	f.ctx.Rex = &rex
	var sets []string
	for i := range 2 {
		cs := f.call("change.new", fmt.Sprintf(`{"intent":"Rename %d"}`, i)).(ChangeNewResponse).ChangeSet
		sets = append(sets, cs.String())
		f.call("set", fmt.Sprintf(`{"ref":"done-0","fields":["statement=Done %d"]}`, i))
		f.call("propose", `{}`)
	}
	if err := daemonNotify(f.ctx, f.rs, &f.state, ""); err != nil {
		t.Fatal(err)
	}
	var notices []string
	for _, r := range rexReports(t, rex.String()) {
		if r["state"] == "blocked" && strings.Contains(r["msg"], "an overlap with") {
			notices = append(notices, r["msg"])
		}
	}
	if len(notices) != 1 || !strings.HasPrefix(notices[0], "Rename ") || !strings.Contains(notices[0], "an overlap with Rename ") {
		t.Fatalf("overlap notices %q, want one naming both change sets by title", notices)
	}
	for _, cs := range sets {
		if strings.Contains(notices[0], cs) {
			t.Errorf("notice %q names change set %s by ID", notices[0], cs)
		}
	}
}
