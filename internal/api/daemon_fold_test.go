package api

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/store"
)

// daemonFoldFix is a repo with a large log, an index, a fake git and
// detach, and one accepted change set the daemon would pick.
type daemonFoldFix struct {
	ctx      *Ctx
	rs       *run.Store
	pending  map[string]daemonPending
	unj      map[string]daemonPending
	state    daemon.State
	picked   []string
	logReads int
	// pickFolds is folds() when the last detach ran: the folds the poll
	// took before the start dropped the cache.
	pickFolds int
	events    string
	cmd       string
	target    string
	clock     time.Time
	skips     []daemon.Skip
	call      func(name, req string) any
	refs      string
	config    map[string]string
	runsDir   string
}

func newDaemonFoldFix(t testing.TB) *daemonFoldFix { return newDaemonFoldFixN(t, 98, nil) }

// newDaemonFoldFixN is newDaemonFoldFix with done change sets landed in
// one step each. grow, when set, is called with the log folder once the
// change sets are written, to add more of them directly.
func newDaemonFoldFixN(t testing.TB, done int, grow func(logDir string)) *daemonFoldFix {
	t.Helper()
	ctx := testCtx(t, t.TempDir())
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	call := func(name, req string) any {
		o, _ := Lookup(name)
		r, err := DecodeRequest(o, []byte(req))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := o.Run(ctx, r)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	call("init", `{}`)
	f := &daemonFoldFix{call: call, config: map[string]string{}, ctx: ctx, pending: map[string]daemonPending{}, unj: map[string]daemonPending{}, clock: ctx.Now()}
	// A milestone, landed.
	call("change.new", `{"intent":"Milestone"}`)
	call("add", `{"kind":"intent","handle":"org"}`)
	call("add", `{"kind":"intent","handle":"m","parent":"org","json_fields":["milestone=null"]}`)
	call("set", `{"ref":"m","fields":["milestone=1"]}`)
	f.lifecycle(call, "propose", "accept", "land")
	// Many change sets landed in one step each, then a few in flight.
	rec := *ctx
	rec.Record = true
	for i := range done {
		o, _ := Lookup("add")
		r, err := DecodeRequest(o, []byte(fmt.Sprintf(`{"kind":"capability","handle":"done-%d"}`, i)))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := o.Run(&rec, r); err != nil {
			t.Fatal(err)
		}
	}
	var proposed []string
	for i := range 6 {
		cs := call("change.new", fmt.Sprintf(`{"intent":"Work %d"}`, i)).(ChangeNewResponse).ChangeSet
		call("add", fmt.Sprintf(`{"kind":"capability","handle":"cap-%d"}`, i))
		switch {
		case i == 0:
			call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"m"}`, cs.String()))
			f.lifecycle(call, "propose", "accept")
			f.target = cs.String()
		case i < 4:
			call("propose", `{}`)
			proposed = append(proposed, cs.String())
		}
	}
	if grow != nil {
		grow(filepath.Join(ctx.Root, store.Dir, store.LogDir))
	}
	// The log is old news: a folder changed in the last clock ticks is not
	// a stamp to trust.
	old := time.Now().Add(-time.Hour)
	f.backdate(t, old)
	common := t.TempDir()
	ctx.Git = func(args ...string) (string, error) {
		if len(args) > 0 && args[0] == "rev-parse" {
			return common, nil
		}
		if len(args) > 0 && args[0] == "for-each-ref" {
			return f.refs, nil
		}
		if out, err, ok := gitConfigFake(f.config, args); ok && err == nil {
			return out, nil
		}
		if len(args) == 3 && args[0] == "config" && args[1] == "--get-regexp" {
			var out string
			for k, v := range f.config {
				out += k + " " + v + "\n"
			}
			return out, nil
		}
		return "", fmt.Errorf("no git")
	}
	ctx.Detach = func(p Program, _ string) error {
		f.picked = append(f.picked, p.Args[3])
		f.pickFolds = f.folds()
		return nil
	}
	ctx.Self = "rhiza"
	ctx.repoSession = true
	f.runsDir = run.Dir(common)
	f.rs = run.NewStore(f.runsDir)
	// The proposed change sets are covered (a failed plan-review run each),
	// so a poll starts no plan review beside what the tests watch.
	for i, cs := range proposed {
		rid := fmt.Sprintf("review%d", i)
		at := ctx.Now().Add(-time.Hour)
		if err := f.rs.Create(rid, run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan-review", Change: cs}}); err != nil {
			t.Fatal(err)
		}
		if err := f.rs.Append(rid, run.Entry{Type: run.TypeFinished, At: at, End: run.Failed, Detail: "fixture"}); err != nil {
			t.Fatal(err)
		}
	}
	f.events = filepath.Join(t.TempDir(), "events")
	f.cmd = `printf '%s\n' "$RHIZA_EVENT" >> ` + f.events
	w, err := openWorkspace(ctx.Root, ctx.Now, ctx.IDs, false)
	if err != nil {
		t.Fatal(err)
	}
	st := w.st
	ctx.readLog = func() store.Log { f.logReads++; return st.ReadLog() }
	return f
}

func (f *daemonFoldFix) lifecycle(call func(string, string) any, moves ...string) {
	for _, m := range moves {
		req := `{}`
		if m == "accept" {
			req = `{"reason":"ok"}`
		}
		call(m, req)
	}
}

func (f *daemonFoldFix) backdate(t testing.TB, at time.Time) {
	t.Helper()
	root := filepath.Join(f.ctx.Root, store.Dir, store.LogDir)
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err == nil && d.IsDir() {
			err = os.Chtimes(p, at, at)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}

// poll runs the daemon's whole poll closure once.
func (f *daemonFoldFix) poll(t testing.TB) {
	t.Helper()
	var skips []daemon.Skip
	defer func() { f.skips = skips }()
	if _, _, _, err := daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 0, 4, 0, false, "", false, nil, nil, &skips); err != nil {
		t.Fatal(err)
	}
	if err := daemonNotify(f.ctx, f.rs, &f.state, ""); err != nil {
		t.Fatal(err)
	}
	if err := daemonEvents(f.ctx, &f.state, f.cmd); err != nil {
		t.Fatal(err)
	}
}

// folds is the folds the daemon's fold cache has taken.
func (f *daemonFoldFix) folds() int {
	if f.ctx.daemonFold == nil {
		return 0
	}
	return f.ctx.daemonFold.folds
}

func (f *daemonFoldFix) eventLines() []string {
	data, _ := os.ReadFile(f.events)
	return strings.Fields(string(data))
}

// A poll with the log, refs and settings unchanged reuses the last fold:
// no fold, no plan read, no change set folded; over the index no poll reads
// the log. The clock-driven hold
// still lapses, and a lifecycle move reads again and runs the event
// command.
func TestDaemonIdlePollSkipsFold(t *testing.T) {
	t.Parallel()
	f := newDaemonFoldFix(t)
	f.unj[f.target] = daemonPending{at: f.clock}

	f.poll(t) // the first poll folds and seeds the events
	if f.folds() == 0 || planReadCount.Load() == 0 {
		t.Fatalf("first poll: %d folds, %d plan reads; want some", f.folds(), planReadCount.Load())
	}
	if f.logReads != 0 {
		t.Errorf("first poll over the index read the log %d times", f.logReads)
	}
	if len(f.picked) != 0 {
		t.Fatalf("held change set picked: %v", f.picked)
	}

	f.logReads = 0
	cached := f.folds()
	reads0, folds0 := planReadCount.Load(), planFoldCount.Load()
	f.poll(t)
	if f.logReads != 0 || f.folds() != cached {
		t.Errorf("idle poll read the log %d times and folded %d times", f.logReads, f.folds()-cached)
	}
	if n := planReadCount.Load() - reads0; n != 0 {
		t.Errorf("idle poll ran planRead %d times", n)
	}
	if n := planFoldCount.Load() - folds0; n != 0 {
		t.Errorf("idle poll folded %d in-flight change sets", n)
	}
	if len(f.picked) != 0 {
		t.Fatalf("held change set picked early: %v", f.picked)
	}

	// The 15-minute hold lapses by the clock alone.
	now := f.clock.Add(daemonUnjournaledHold + time.Minute)
	f.ctx.Now = func() time.Time { return now }
	f.logReads = 0
	f.poll(t)
	if len(f.picked) != 1 {
		t.Fatalf("after the hold picked = %v, want the change set once; skips %+v", f.picked, f.skips)
	}

	// A lifecycle move is a new log entry: the next poll reads and plans again.
	before := len(f.eventLines())
	o, _ := Lookup("change.new")
	r, _ := DecodeRequest(o, []byte(`{"intent":"Fresh"}`))
	if _, err := o.Run(f.ctx, r); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"propose"} {
		o, _ := Lookup(s)
		r, _ := DecodeRequest(o, []byte(`{}`))
		if _, err := o.Run(f.ctx, r); err != nil {
			t.Fatal(err)
		}
	}
	f.logReads = 0
	cached = f.folds()
	reads0 = planReadCount.Load()
	f.poll(t)
	if f.folds() == cached || planReadCount.Load() == reads0 {
		t.Errorf("poll after a log entry: %d folds, %d plan reads; want both", f.folds()-cached, planReadCount.Load()-reads0)
	}
	if f.logReads != 0 {
		t.Errorf("poll after a log entry read the log %d times", f.logReads)
	}
	if got := f.eventLines(); len(got) != before+1 || got[len(got)-1] != "proposed" {
		t.Errorf("events = %v, want one more, proposed", got)
	}
}

// BenchmarkDaemonIdlePoll reports the cost of a poll with nothing changed.
func BenchmarkDaemonIdlePoll(b *testing.B) {
	f := newDaemonFoldFix(b)
	f.poll(b)
	b.ResetTimer()
	for range b.N {
		f.poll(b)
	}
}

// Each of a new log entry, a moved ref, a changed daemon setting and a run
// the daemon starts or resumes drops the cache: the next poll folds again
// and acts on the change.
func TestDaemonPollCacheInvalidates(t *testing.T) {
	t.Parallel()
	f := newDaemonFoldFix(t)
	f.unj[f.target] = daemonPending{at: f.clock}
	f.config["rhiza.autonomy.ship"] = "owner"

	// readsNext polls and returns the folds it took.
	readsNext := func() int {
		n := f.folds()
		f.poll(t)
		return f.folds() - n
	}
	idle := func(what string) {
		t.Helper()
		f.backdate(t, time.Now().Add(-time.Hour))
		readsNext() // refill
		if n := readsNext(); n != 0 {
			t.Fatalf("%s: idle poll folded %d times", what, n)
		}
	}
	accept := func(i int) string {
		cs := f.call("change.new", fmt.Sprintf(`{"intent":"Late %d"}`, i)).(ChangeNewResponse).ChangeSet
		f.call("add", fmt.Sprintf(`{"kind":"capability","handle":"late-%d"}`, i))
		f.call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"m"}`, cs.String()))
		f.lifecycle(f.call, "propose", "accept")
		return cs.String()
	}
	skipped := func(cs string) bool {
		for _, s := range f.skips {
			if s.Change == cs && s.Reason == daemon.SkipShipOwner {
				return true
			}
		}
		return false
	}
	idle("start")

	// A log entry: a change set accepted. The owner gates shipping, so it
	// is seen and skipped, not picked.
	late1 := accept(1)
	if n := readsNext(); n == 0 || !skipped(late1) || len(f.picked) != 0 {
		t.Fatalf("after an accept: %d folds, picked %v, skips %+v; want a fold and %s skipped", n, f.picked, f.skips, late1)
	}
	idle("accept")

	// A setting: the same cached log, now picked under the new value. The
	// pick itself drops the cache, so count only the folds taken before it:
	// those come from the setting changing the stamp.
	f.config["rhiza.autonomy.ship"] = "auto"
	f.pickFolds = -1
	before := f.folds()
	readsNext()
	if len(f.picked) != 1 || f.picked[0] != late1 {
		t.Fatalf("after daemon set: picked %v; want %s picked", f.picked, late1)
	}
	if f.pickFolds <= before {
		t.Fatalf("after daemon set: the poll folded %d times before the pick; want a fold", f.pickFolds-before)
	}
	idle("pick")

	// The daemon starts a run.
	if err := daemonStart(f.ctx, nil, late1, false, ""); err != nil {
		t.Fatal(err)
	}
	if n := readsNext(); n == 0 {
		t.Error("poll after the daemon started a run did not fold again")
	}
	idle("start")

	// A moved ref.
	f.refs = "abc123 refs/heads/main\n"
	if n := readsNext(); n == 0 {
		t.Error("poll after a ref moved did not fold again")
	}
	idle("ref")

	// The daemon resumes a run.
	if err := daemonResume(f.ctx, "run1"); err != nil {
		t.Fatal(err)
	}
	if n := readsNext(); n == 0 {
		t.Error("poll after the daemon resumed a run did not fold again")
	}
}

// TestDaemonFoldIndexedMatchesLog: over the index the poll's fold reads the
// log zero times and gives the candidates, proposed list, plan skips,
// overlaps and lifecycle moves the fold over the log gives, for proposed
// change sets with and without leaves, an accepted one and two in-flight
// change sets that set the same node. An entry file that vanishes (the
// index rebuilds) changes the stamp, and the next poll still matches the
// log.
func TestDaemonFoldIndexedMatchesLog(t *testing.T) {
	t.Parallel()
	f := newDaemonFoldFixN(t, 2, nil)
	bare := f.call("change.new", `{"intent":"No leaves"}`).(ChangeNewResponse).ChangeSet
	f.call("add", `{"kind":"intent","handle":"bare-intent"}`)
	f.call("propose", `{}`)
	for _, v := range []string{"A", "B"} {
		f.call("change.new", fmt.Sprintf(`{"intent":"Rename %s"}`, v))
		f.call("set", fmt.Sprintf(`{"ref":"done-0","fields":["statement=%s"]}`, v))
		f.call("propose", `{}`)
	}
	f.backdate(t, time.Now().Add(-time.Hour))

	same := func(step string) *daemonFolded {
		t.Helper()
		f.logReads = 0
		before := f.folds()
		f.poll(t)
		if f.folds() == before {
			t.Fatalf("%s: the poll did not fold", step)
		}
		if f.logReads != 0 {
			t.Errorf("%s: the poll over the index read the log %d times, want 0", step, f.logReads)
		}
		indexed := f.ctx.daemonFold.folded
		if indexed == nil {
			t.Fatalf("%s: no cached fold", step)
		}
		// The log side: the change sets folded from the log and the models
		// built over them, as the fold without an index does.
		ws, err := openCtxWorkspace(f.ctx)
		if err != nil {
			t.Fatal(err)
		}
		list, _ := changeset.FoldLog(ws.log())
		logged := daemonDerive(f.ctx, indexed.res, list, newModelBase(list).model)
		if f.logReads == 0 {
			t.Errorf("%s: the fold over the log did not read it", step)
		}
		for _, c := range []struct {
			what   string
			ix, lg any
		}{
			{"candidates", indexed.cands, logged.cands},
			{"proposed", indexed.proposed, logged.proposed},
			{"plan skips", indexed.planSkips, logged.planSkips},
			{"overlaps", indexed.overlaps, logged.overlaps},
			{"moves", indexed.moves, logged.moves},
		} {
			if !reflect.DeepEqual(c.ix, c.lg) {
				t.Errorf("%s: %s over the index %+v, over the log %+v", step, c.what, c.ix, c.lg)
			}
		}
		if len(indexed.sets) != len(logged.sets) {
			t.Errorf("%s: %d change sets over the index, %d over the log", step, len(indexed.sets), len(logged.sets))
		}
		for cs, c := range logged.sets {
			if got, ok := indexed.changeSet(cs); !ok || got.State != c.State {
				t.Errorf("%s: change set %s over the index %q (%v), over the log %q", step, cs, got.State, ok, c.State)
			}
		}
		return indexed
	}
	got := same("all")
	if !slices.Contains(got.cands, f.target) {
		t.Errorf("candidates %v lack the accepted change set %s", got.cands, f.target)
	}
	if len(got.proposed) != 3 || slices.Contains(got.proposed, bare.String()) {
		t.Errorf("proposed = %v, want the three with leaves, not %s", got.proposed, bare)
	}
	if len(got.overlaps) == 0 {
		t.Error("no overlap for the two change sets that set done-0")
	}
	if len(got.moves) == 0 {
		t.Error("no lifecycle moves")
	}

	// A proposed change set's propose entry vanishes (a rebase): the index
	// rebuilds, and it is no longer proposed.
	gone := got.proposed[0]
	dir := filepath.Join(f.ctx.Root, store.Dir, store.LogDir, gone)
	es, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	removed := false
	for _, e := range es {
		p := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), `"type": "propose"`) {
			if err := os.Remove(p); err != nil {
				t.Fatal(err)
			}
			removed = true
		}
	}
	if !removed {
		t.Fatalf("no propose entry in %s", dir)
	}
	f.backdate(t, time.Now().Add(-time.Hour))
	got = same("vanished")
	if slices.Contains(got.proposed, gone) {
		t.Errorf("proposed %v still lists %s after its propose entry vanished", got.proposed, gone)
	}
	if !slices.Contains(got.cands, f.target) {
		t.Errorf("after the rebuild candidates %v lack %s", got.cands, f.target)
	}
}
