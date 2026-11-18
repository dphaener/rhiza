package api

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
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
	// planReads and planFolds count this fixture's Ctx's plan reads and
	// the in-flight change sets they folded (Ctx.planRead, planFolded), so
	// a parallel test's reads do not show in them.
	planReads, planFolds atomic.Int64
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

func newDaemonFoldFix(t testing.TB) *daemonFoldFix { return newDaemonFoldFixN(t, 98, 0) }

// daemonFoldSeed is the snapshot of newDaemonFoldFixN's repository for one
// count of done change sets and of the change sets the log grows to, built
// once per binary: a landed milestone, done change sets landed in one step
// each, then six in flight, the first accepted (target) and the next three
// proposed (proposed), with the build's IDs kept so every copy knows them;
// then, when grown is set, the landed ones copied until the log holds grown
// change sets (growLog), so the template's index is caught up with them too.
type daemonFoldSeed struct {
	snap     *snapshot
	target   string
	proposed []string
}

// daemonFoldSeeds holds a daemonFoldSeed per count of done change sets.
var daemonFoldSeeds sync.Map // [2]int{done, grown} -> *daemonFoldSeed

// daemonFoldSeedOf returns the seed for done change sets and a log grown to
// grown change sets (none when 0), registering it on first use. The seed that wins the store is the one whose build runs, so
// its fields are the ones its snapshot's build wrote.
func daemonFoldSeedOf(done, grown int) *daemonFoldSeed {
	key := [2]int{done, grown}
	if v, ok := daemonFoldSeeds.Load(key); ok {
		return v.(*daemonFoldSeed)
	}
	seed := &daemonFoldSeed{}
	seed.snap = &snapshot{name: fmt.Sprintf("daemon-fold-%d-grown-%d", done, grown), build: func(dir string) (*Ctx, error) {
		ctx, err := buildInit(dir)
		if err != nil {
			return nil, err
		}
		call := func(name, req string) (any, error) { return runOp(ctx, name, req) }
		calls := func(reqs ...[2]string) error {
			for _, r := range reqs {
				if _, err := call(r[0], r[1]); err != nil {
					return err
				}
			}
			return nil
		}
		propose, accept, land := [2]string{"propose", `{}`}, [2]string{"accept", `{"reason":"ok"}`}, [2]string{"land", `{}`}
		// A milestone, landed.
		if err := calls(
			[2]string{"change.new", `{"intent":"Milestone"}`},
			[2]string{"add", `{"kind":"intent","handle":"org"}`},
			[2]string{"add", `{"kind":"intent","handle":"m","parent":"org","json_fields":["milestone=null"]}`},
			[2]string{"set", `{"ref":"m","fields":["milestone=1"]}`},
			propose, accept, land,
		); err != nil {
			return nil, err
		}
		// Many change sets landed in one step each, then a few in flight.
		rec := *ctx
		rec.Record = true
		for i := range done {
			if _, err := runOp(&rec, "add", fmt.Sprintf(`{"kind":"capability","handle":"done-%d"}`, i)); err != nil {
				return nil, err
			}
		}
		for i := range 6 {
			resp, err := call("change.new", fmt.Sprintf(`{"intent":"Work %d"}`, i))
			if err != nil {
				return nil, err
			}
			cs := resp.(ChangeNewResponse).ChangeSet.String()
			if _, err := call("add", fmt.Sprintf(`{"kind":"capability","handle":"cap-%d"}`, i)); err != nil {
				return nil, err
			}
			switch {
			case i == 0:
				if err := calls([2]string{"link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"m"}`, cs)}, propose, accept); err != nil {
					return nil, err
				}
				seed.target = cs
			case i < 4:
				if err := calls(propose); err != nil {
					return nil, err
				}
				seed.proposed = append(seed.proposed, cs)
			}
		}
		if grown > 0 {
			if err := growLog(filepath.Join(ctx.Root, store.Dir, store.LogDir), grown); err != nil {
				return nil, err
			}
		}
		return ctx, nil
	}}
	v, _ := daemonFoldSeeds.LoadOrStore(key, seed)
	return v.(*daemonFoldSeed)
}

// newDaemonFoldFixN is newDaemonFoldFix with done change sets landed in
// one step each and, when grown is set, the log grown to grown change sets
// by copying them directly (growLog). The repository is a copy of
// daemonFoldSeedOf(done, grown)'s snapshot, its index with it.
func newDaemonFoldFixN(t testing.TB, done, grown int) *daemonFoldFix {
	t.Helper()
	seed := daemonFoldSeedOf(done, grown)
	_, ctx := seed.snap.copy(t)
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
	f := &daemonFoldFix{call: call, config: map[string]string{}, ctx: ctx, pending: map[string]daemonPending{}, unj: map[string]daemonPending{}, clock: ctx.Now(), target: seed.target}
	proposed := seed.proposed
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
	ctx.logRead = func() { f.logReads++ }
	ctx.planRead = func() { f.planReads.Add(1) }
	ctx.planFolded = func() { f.planFolds.Add(1) }
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
	if f.folds() == 0 || f.planReads.Load() == 0 {
		t.Fatalf("first poll: %d folds, %d plan reads; want some", f.folds(), f.planReads.Load())
	}
	if f.logReads != 0 {
		t.Errorf("first poll over the index read the log %d times", f.logReads)
	}
	if len(f.picked) != 0 {
		t.Fatalf("held change set picked: %v", f.picked)
	}

	f.logReads = 0
	cached := f.folds()
	reads0, folds0 := f.planReads.Load(), f.planFolds.Load()
	f.poll(t)
	if f.logReads != 0 || f.folds() != cached {
		t.Errorf("idle poll read the log %d times and folded %d times", f.logReads, f.folds()-cached)
	}
	if n := f.planReads.Load() - reads0; n != 0 {
		t.Errorf("idle poll ran planRead %d times", n)
	}
	if n := f.planFolds.Load() - folds0; n != 0 {
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
	reads0 = f.planReads.Load()
	f.poll(t)
	if f.folds() == cached || f.planReads.Load() == reads0 {
		t.Errorf("poll after a log entry: %d folds, %d plan reads; want both", f.folds()-cached, f.planReads.Load()-reads0)
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
	f := newDaemonFoldFixN(t, 2, 0)
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
		// The log side: the change sets folded from the log, and the models
		// of a fresh projection built from it in a throwaway index.
		ws, err := openCtxWorkspace(f.ctx)
		if err != nil {
			t.Fatal(err)
		}
		list, _ := changeset.FoldLog(ws.log())
		fresh := *f.ctx
		fresh.IndexPath = ""
		fp, err := openProjection(&fresh, false)
		if err != nil {
			t.Fatal(err)
		}
		logged := daemonDerive(f.ctx, indexed.res, list, func(cs id.ID) Model { return fp.model(nil, cs) })
		if err := fp.failure(); err != nil {
			t.Fatalf("%s: fresh projection: %v", step, err)
		}
		fp.close()
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
