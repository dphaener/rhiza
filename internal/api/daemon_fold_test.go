package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	// pickReads is logReads when the last detach ran: the reads the poll
	// took before the start dropped the cache.
	pickReads int
	events    string
	cmd       string
	target    string
	clock     time.Time
	skips     []daemon.Skip
	call      func(name, req string) any
	refs      string
	config    map[string]string
}

func newDaemonFoldFix(t testing.TB) *daemonFoldFix {
	t.Helper()
	ctx := testCtx(t.TempDir())
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
	for i := range 98 {
		o, _ := Lookup("add")
		r, err := DecodeRequest(o, []byte(fmt.Sprintf(`{"kind":"capability","handle":"done-%d"}`, i)))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := o.Run(&rec, r); err != nil {
			t.Fatal(err)
		}
	}
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
		}
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
		if len(args) == 3 && args[0] == "config" && args[1] == "--get" {
			if v, ok := f.config[args[2]]; ok {
				return v, nil
			}
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
		f.pickReads = f.logReads
		return nil
	}
	ctx.Self = "rhiza"
	ctx.repoSession = true
	f.rs = run.NewStore(run.Dir(common))
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
	if _, _, _, err := daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 4, 0, false, "", false, nil, &skips); err != nil {
		t.Fatal(err)
	}
	if err := daemonNotify(f.ctx, f.rs, &f.state, ""); err != nil {
		t.Fatal(err)
	}
	if err := daemonEvents(f.ctx, &f.state, f.cmd); err != nil {
		t.Fatal(err)
	}
}

func (f *daemonFoldFix) eventLines() []string {
	data, _ := os.ReadFile(f.events)
	return strings.Fields(string(data))
}

// A poll with the log, refs and settings unchanged reuses the last fold:
// no log read, no plan read, no change set folded. The clock-driven hold
// still lapses, and a lifecycle move reads again and runs the event
// command.
func TestDaemonIdlePollSkipsFold(t *testing.T) {
	f := newDaemonFoldFix(t)
	f.unj[f.target] = daemonPending{at: f.clock}

	f.poll(t) // the first poll folds and seeds the events
	if f.logReads == 0 || planReadCount.Load() == 0 {
		t.Fatalf("first poll: %d log reads, %d plan reads; want some", f.logReads, planReadCount.Load())
	}
	if len(f.picked) != 0 {
		t.Fatalf("held change set picked: %v", f.picked)
	}

	f.logReads = 0
	reads0, folds0 := planReadCount.Load(), planFoldCount.Load()
	f.poll(t)
	if f.logReads != 0 {
		t.Errorf("idle poll read the log %d times", f.logReads)
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
	reads0 = planReadCount.Load()
	f.poll(t)
	if f.logReads == 0 || planReadCount.Load() == reads0 {
		t.Errorf("poll after a log entry: %d log reads, %d plan reads; want both", f.logReads, planReadCount.Load()-reads0)
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
// the daemon starts or resumes drops the cache: the next poll reads the log
// again and acts on the change.
func TestDaemonPollCacheInvalidates(t *testing.T) {
	f := newDaemonFoldFix(t)
	f.unj[f.target] = daemonPending{at: f.clock}
	f.config["rhiza.autonomy.ship"] = "owner"

	// readsNext polls and returns the log reads it took.
	readsNext := func() int {
		f.logReads = 0
		f.poll(t)
		return f.logReads
	}
	idle := func(what string) {
		t.Helper()
		f.backdate(t, time.Now().Add(-time.Hour))
		readsNext() // refill
		if n := readsNext(); n != 0 {
			t.Fatalf("%s: idle poll read the log %d times", what, n)
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
		t.Fatalf("after an accept: %d log reads, picked %v, skips %+v; want reads and %s skipped", f.logReads, f.picked, f.skips, late1)
	}
	idle("accept")

	// A setting: the same cached log, now picked under the new value. The
	// pick itself drops the cache, so count only the reads taken before it:
	// those come from the setting changing the stamp.
	f.config["rhiza.autonomy.ship"] = "auto"
	f.pickReads = -1
	readsNext()
	if len(f.picked) != 1 || f.picked[0] != late1 {
		t.Fatalf("after daemon set: picked %v; want %s picked", f.picked, late1)
	}
	if f.pickReads <= 0 {
		t.Fatalf("after daemon set: the poll read the log %d times before the pick; want reads", f.pickReads)
	}
	idle("pick")

	// The daemon starts a run.
	if err := daemonStart(f.ctx, late1, false, ""); err != nil {
		t.Fatal(err)
	}
	if n := readsNext(); n == 0 {
		t.Error("poll after the daemon started a run did not read the log")
	}
	idle("start")

	// A moved ref.
	f.refs = "abc123 refs/heads/main\n"
	if n := readsNext(); n == 0 {
		t.Error("poll after a ref moved did not read the log")
	}
	idle("ref")

	// The daemon resumes a run.
	if err := daemonResume(f.ctx, "run1"); err != nil {
		t.Fatal(err)
	}
	if n := readsNext(); n == 0 {
		t.Error("poll after the daemon resumed a run did not read the log")
	}
}
