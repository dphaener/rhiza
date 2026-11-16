package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
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
	t.Parallel()
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
		if err := f.rs.Create(id, run.Entry{Type: run.TypeStarted, At: at.Add(time.Duration(i) * time.Minute), Started: &run.Started{Workflow: "ship-change-classic", Change: cs}}); err != nil {
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

	// Without an index path the fold reads a throwaway index, never the
	// log, and answers as the index does.
	byIndex, byLog := newInboxFold(f.ctx), func() *inboxFold {
		c := *f.ctx
		// No index path on purpose: the throwaway index answers as the index does.
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
	if f.logReads != 0 {
		t.Errorf("the fold over a throwaway index read the log %d times, want 0", f.logReads)
	}
}

// An overlap between change sets with no run is notified by the change
// sets' titles from the fold, never by either ID.
func TestDaemonNotifyOverlapByTitle(t *testing.T) {
	t.Parallel()
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

// TestInboxIndexedReadsNoLog: an inbox with overlapping in-flight change
// sets and a plan run stopped at its gate on a question reads the log zero
// times with an index (the workspace's projection answers the states and
// overlaps, the worktree's projection the question) and lists the same
// items as over the log. A gate-stopped run whose worktree is gone keeps
// its item, without its question's statement, and no error.
//
// Serial: it calls t.Setenv (XDG_CACHE_HOME, HOME) so the worktree's index
// goes to a temp folder, process-wide state.
func TestInboxIndexedReadsNoLog(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", cache)
	f := newDaemonFoldFix(t)
	// Two in-flight change sets that both write done-0: an overlap. The
	// second asks the question the plan run stops on.
	var sets []string
	for i := range 2 {
		cs := f.call("change.new", fmt.Sprintf(`{"intent":"Rename %d"}`, i)).(ChangeNewResponse).ChangeSet
		sets = append(sets, cs.String())
		f.call("set", fmt.Sprintf(`{"ref":"done-0","fields":["statement=Done %d"]}`, i))
		if i == 1 {
			f.call("add", `{"kind":"question","handle":"gate-ask","fields":["statement=Which name wins?"]}`)
		}
		f.call("propose", `{}`)
	}
	at := time.Date(2026, 11, 10, 1, 0, 0, 0, time.UTC)
	gated := func(rid, worktree string) {
		t.Helper()
		es := []run.Entry{
			{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan", Change: sets[1], Worktree: worktree, Plan: true}},
			{Type: run.TypeEntered, At: at.Add(time.Second), Step: "gate", Visit: 1},
			{Type: run.TypeGated, At: at.Add(2 * time.Second), Changes: []string{sets[1]}},
			{Type: run.TypeStopped, At: at.Add(3 * time.Second), Step: "gate", Detail: run.OpenStop + "gate-ask; answer each with a Decision, then rhiza run resume " + rid},
		}
		if err := f.rs.Create(rid, es[0]); err != nil {
			t.Fatal(err)
		}
		for _, e := range es[1:] {
			if err := f.rs.Append(rid, e); err != nil {
				t.Fatal(err)
			}
		}
	}
	gated("RGATE", f.ctx.Root)
	gated("RGONE", filepath.Join(t.TempDir(), "gone"))
	inbox := func(ctx *Ctx) InboxResponse {
		t.Helper()
		resp, err := callOp(t, ctx, "inbox", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		return resp.(InboxResponse)
	}
	inbox(f.ctx) // catch both indexes up with the log

	f.logReads = 0
	resp := inbox(f.ctx)
	if f.logReads != 0 {
		t.Errorf("inbox with an index read the log %d times, want 0", f.logReads)
	}
	// The other side: no index path on purpose, to compare a throwaway
	// index against the indexed inbox; its overlaps, titles and gate
	// questions are read through the fold directly.
	folded := *f.ctx
	folded.IndexPath = ""
	byLog := newInboxFold(&folded)
	defer byLog.close()
	f.logReads = 0
	ol, err := byLog.overlaps()
	if err != nil {
		t.Fatal(err)
	}
	if f.logReads != 0 {
		t.Errorf("the fold over a throwaway index read the log %d times, want 0", f.logReads)
	}
	var oi []InboxItem
	for _, it := range resp.Items {
		if it.Kind == InboxOverlap {
			// The inbox titles the overlap's change sets after the fold;
			// the titles are compared below.
			if it.OtherTitle != byLog.titleOf(it.Other) {
				t.Errorf("overlap %s: other title %q over the index, %q over the log", it.ID, it.OtherTitle, byLog.titleOf(it.Other))
			}
			it.ChangeTitle, it.OtherTitle = "", ""
			oi = append(oi, it)
		}
	}
	if !reflect.DeepEqual(oi, ol) || len(oi) != 1 {
		t.Errorf("overlaps: index %+v, log %+v", oi, ol)
	}
	for _, it := range resp.Items {
		if it.Change == "" {
			continue
		}
		if s, _ := byLog.stateOf(it.Change); it.Kind == InboxOpen && s != "proposed" {
			t.Errorf("item %s: change set %s is %q over the log, want proposed", it.ID, it.Change, s)
		}
		if it.ChangeTitle != "" && it.ChangeTitle != byLog.titleOf(it.Change) {
			t.Errorf("item %s: title %q over the index, %q over the log", it.ID, it.ChangeTitle, byLog.titleOf(it.Change))
		}
	}
	for _, rid := range []string{"RGATE", "RGONE"} {
		j, err := f.rs.Read(rid)
		if err != nil {
			t.Fatal(err)
		}
		want := inboxGateQuestions(&folded, j.State())
		var got []InboxGateQuestion
		for _, it := range resp.Items {
			if it.Kind == InboxOpen && it.Run == rid {
				got = it.Questions
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s questions: inbox over the index %+v, over the log %+v", rid, got, want)
		}
	}
	kinds := map[string]int{}
	questions := map[string][]InboxGateQuestion{}
	for _, it := range resp.Items {
		kinds[it.Kind]++
		if it.Kind == InboxOpen {
			questions[it.Run] = it.Questions
		}
	}
	if kinds[InboxOverlap] != 1 {
		t.Errorf("got kinds %v, want one overlap", kinds)
	}
	if q := questions["RGATE"]; len(q) != 1 || q[0].Handle != "gate-ask" || q[0].Statement != "Which name wins?" {
		t.Errorf("RGATE questions = %+v, want gate-ask with its statement", q)
	}
	if q, ok := questions["RGONE"]; !ok || len(q) != 1 || q[0].Handle != "gate-ask" || q[0].Statement != "" {
		t.Errorf("RGONE questions = %+v (listed %v), want gate-ask without a statement", q, ok)
	}
}
