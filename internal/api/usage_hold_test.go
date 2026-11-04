package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/usage"
)

// holdFake reports token records from a package variable.
type holdFake struct{ limitAdapter }

const holdFakeAgent = "hold-fake"

var holdFakeRecords []adapter.UsageRecord

func (holdFake) Name() string { return holdFakeAgent }
func (holdFake) UsageRecords(_ func(string) string, from, to time.Time) ([]adapter.UsageRecord, error) {
	var out []adapter.UsageRecord
	for _, r := range holdFakeRecords {
		if !r.At.Before(from) && r.At.Before(to) {
			out = append(out, r)
		}
	}
	return out, nil
}

var holdFakeOnce = func() func() {
	done := false
	return func() {
		if !done {
			adapter.Register(holdFake{})
			done = true
		}
	}
}()

// finishedRun is a finished run that burned tokens counted tokens.
func finishedRun(n int, tokens int64, at time.Time) run.Journal {
	j := pickJournal(fmt.Sprintf("cs%d", n),
		run.Entry{Type: run.TypeTokens, At: at, Session: "s", Tokens: &run.TokenUse{Session: run.TokenCounts{Output: tokens, CacheRead: 1 << 30}}},
		run.Entry{Type: run.TypeFinished, At: at, End: run.Succeeded})
	j.ID = fmt.Sprintf("done%d", n)
	return j
}

func TestUsageHoldProjection(t *testing.T) {
	holdFakeOnce()
	f := newDaemonFoldFix(t)
	f.ctx.Getenv = func(string) string { return "" }
	now := f.clock
	f.ctx.Now = func() time.Time { return now }
	common, err := gitCommonDir(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	holdFakeRecords = []adapter.UsageRecord{{At: now.Add(-time.Hour), Tokens: 600}}
	set := func(k, v string) { f.config[k] = v }
	conf := func() []Setting {
		c, err := daemonConfOf(f.ctx, DaemonRequest{})
		if err != nil {
			t.Fatal(err)
		}
		return c.rows
	}
	var runs []run.Journal
	hold := func(learned ...daemon.LearnedLimit) *usageHold {
		windows, _ := usageWindowsNow(f.ctx, common, nil)
		return usageHoldOf(conf(), learned, runs, "h", windows)
	}

	// A window with no limit never holds.
	if h := hold(); h != nil {
		t.Fatalf("no limit held: %+v", h)
	}
	// Under the line (limit 1000, reserve 20% -> line 800; use 600) it starts.
	set("rhiza.daemon.usage.session", "1000")
	if h := hold(); h != nil {
		t.Fatalf("under the line held: %+v", h)
	}
	// With no finished runs only the use so far counts.
	holdFakeRecords[0].Tokens = 801
	if h := hold(); h == nil || h.Kind != "session" || h.Projected != 801 {
		t.Fatalf("use over the line: %+v", h)
	}
	holdFakeRecords[0].Tokens = 600

	// Finished runs set the mean (cache reads do not count): the next run
	// projects 600+150=750, under; one live run too makes 900, over.
	for i := range 12 {
		tokens := int64(150)
		if i < 2 { // the two oldest are outside the last 10
			tokens = 9999
		}
		runs = append(runs, finishedRun(i, tokens, now.Add(time.Duration(i-20)*time.Minute)))
	}
	if h := hold(); h != nil {
		t.Fatalf("one projected run held: %+v", h)
	}
	// A live run (on another host, so it counts as running) adds a second
	// projected run: 600+2*150 = 900, over the 800 line, though 750 was not.
	live := pickJournal("live", run.Entry{Type: run.TypeProcess, At: now, Process: &run.Process{Host: "elsewhere", Pid: 1}})
	live.ID = "live"
	runs = append(runs, live)
	if h := hold(); h == nil || h.Projected != 900 {
		t.Fatalf("a live run plus the next should project 900: %+v", h)
	}
	runs = runs[:12]
	// Limit 900 puts the line at 720: the use so far (600) is under it, and
	// the one projected run (600+150=750) is what goes over.
	set("rhiza.daemon.usage.session", "900")
	if h := hold(); h == nil || h.Projected != 750 || h.Used != 600 {
		t.Fatalf("one projected run over the line: %+v", h)
	}

	// A set limit wins over a learned one; a learned one stands in when unset.
	set("rhiza.daemon.usage.session", "5000")
	if h := hold(daemon.LearnedLimit{Kind: "session", Tokens: 700}); h != nil {
		t.Fatalf("the setting should win over the learned limit: %+v", h)
	}
	set("rhiza.daemon.usage.session", "")
	if h := hold(daemon.LearnedLimit{Kind: "session", Tokens: 700}); h == nil {
		t.Fatal("a learned limit with no setting should hold")
	}
	// The weekly window, with its own limit, holds on its own.
	set("rhiza.daemon.usage.weekly", "700")
	if h := hold(); h == nil || h.Kind != "weekly" {
		t.Fatalf("weekly hold: %+v", h)
	}

	// After the window resets there is no current window: no hold.
	set("rhiza.daemon.usage.weekly", "")
	now = now.Add(6 * time.Hour)
	if h := hold(daemon.LearnedLimit{Kind: "session", Tokens: 700}); h != nil {
		t.Fatalf("held after the reset: %+v", h)
	}
}

// Only a hit of a window's kind moves that window: a session hit leaves the
// weekly window alone, and a weekly hit leaves the session window alone.
func TestUsageWindowsHitKinds(t *testing.T) {
	holdFakeOnce()
	f := newDaemonFoldFix(t)
	f.ctx.Getenv = func(string) string { return "" }
	t0 := f.clock
	common, err := gitCommonDir(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	holdFakeRecords = []adapter.UsageRecord{{At: t0, Tokens: 1000}, {At: t0.Add(time.Hour), Tokens: 500}}
	defer func() { holdFakeRecords = nil }()
	for _, h := range []daemon.LimitHit{
		{HitAt: t0.Add(time.Hour), ResetAt: t0.Add(3 * time.Hour)},      // session
		{HitAt: t0.Add(time.Hour), ResetAt: t0.Add(6 * 24 * time.Hour)}, // weekly
	} {
		if err := daemon.RecordLimitHit(common, h); err != nil {
			t.Fatal(err)
		}
	}
	at := func(now time.Time) map[usage.Kind]usage.Window {
		f.ctx.Now = func() time.Time { return now }
		w, _ := usageWindowsNow(f.ctx, common, nil)
		return w
	}
	// Before the session reset: the session window ends at its hit's reset,
	// and the weekly window at its own, each holding both records.
	w := at(t0.Add(2 * time.Hour))
	if s := w[usage.Session]; s.Used != 1500 || !s.Reset.Equal(t0.Add(3*time.Hour)) {
		t.Fatalf("session window %+v: want 1500 used, reset at the session hit's", s)
	}
	if wk := w[usage.Weekly]; wk.Used != 1500 || !wk.Reset.Equal(t0.Add(6*24*time.Hour)) {
		t.Fatalf("weekly window %+v: want 1500 used, reset at the weekly hit's", wk)
	}
	// After the session reset the weekly window still counts both records.
	w = at(t0.Add(4 * time.Hour))
	if _, ok := w[usage.Session]; ok {
		t.Fatalf("session window after its reset: %+v", w[usage.Session])
	}
	if wk := w[usage.Weekly]; wk.Used != 1500 {
		t.Fatalf("weekly window after a session hit %+v: want 1500 used", wk)
	}
}

// A saved window keeps its start from poll to poll: the chain is built from
// it, not from wherever the lookback falls.
func TestUsageWindowsAnchored(t *testing.T) {
	holdFakeOnce()
	f := newDaemonFoldFix(t)
	f.ctx.Getenv = func(string) string { return "" }
	b := f.clock
	common, err := gitCommonDir(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	holdFakeRecords = []adapter.UsageRecord{
		{At: b, Tokens: 100}, {At: b.Add(270 * time.Minute), Tokens: 200}, {At: b.Add(8 * time.Hour), Tokens: 400},
	}
	defer func() { holdFakeRecords = nil }()
	poll := func(now time.Time, saved []daemon.UsageWindow) (map[usage.Kind]usage.Window, []daemon.UsageWindow) {
		f.ctx.Now = func() time.Time { return now }
		return usageWindowsNow(f.ctx, common, saved)
	}
	// At b+9h the session chain is [b, b+5h) then [b+8h, b+13h).
	w, saved := poll(b.Add(9*time.Hour), nil)
	if s := w[usage.Session]; s.Used != 400 || !s.Start.Equal(b.Add(8*time.Hour)) {
		t.Fatalf("first poll session window %+v", s)
	}
	for _, s := range saved {
		if s.Kind == "session" && s.Used != 400 {
			t.Fatalf("saved session window %+v: want its use, 400, kept for health", s)
		}
	}
	// At b+10h30m a fresh 10h lookback starts after b, so its chain would run
	// [b+4h30m, b+9h30m) and find no window now; the saved start keeps it.
	if w, _ := poll(b.Add(630*time.Minute), nil); w[usage.Session].Used == 400 {
		t.Fatalf("the unanchored lookback should drift: %+v", w[usage.Session])
	}
	w, saved = poll(b.Add(630*time.Minute), saved)
	if s := w[usage.Session]; s.Used != 400 || !s.Start.Equal(b.Add(8*time.Hour)) || !s.Reset.Equal(b.Add(13*time.Hour)) {
		t.Fatalf("anchored session window %+v: want [b+8h, b+13h) with 400", s)
	}
	// Past the reset with no new record there is no window, and the saved
	// one stays so the next window starts at or after its reset.
	w, saved = poll(b.Add(14*time.Hour), saved)
	if _, ok := w[usage.Session]; ok {
		t.Fatalf("session window after the reset: %+v", w[usage.Session])
	}
	var kept bool
	for _, s := range saved {
		kept = kept || (s.Kind == "session" && s.Reset.Equal(b.Add(13*time.Hour)))
	}
	if !kept {
		t.Fatalf("saved %+v: want the reset session window kept", saved)
	}
}

// A hold stops picks, resumes, queued starts and plans, and the next poll
// after it clears starts the held change set.
func TestUsageHoldStopsLaunches(t *testing.T) {
	f := newDaemonFoldFix(t)
	hold := &usageHold{Kind: "session", Used: 900, Projected: 900, Limit: 1000, Reserve: 20, Reset: f.clock.Add(time.Hour)}
	var skips []daemon.Skip
	started, resumed, planned, err := daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 4, 0, false, "", false, hold, &skips)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.picked)+len(started)+len(resumed)+len(planned) != 0 {
		t.Fatalf("held poll launched: picked %v started %v resumed %v planned %v", f.picked, started, resumed, planned)
	}
	var found bool
	for _, sk := range skips {
		if sk.Change == f.target && sk.Reason == daemon.SkipUsageReserve && len(sk.Detail) == 1 {
			found = true
		}
	}
	if !found {
		t.Fatalf("skips %+v: want %s skipped for the usage reserve", skips, f.target)
	}

	skips = nil
	started, _, _, err = daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 4, 0, false, "", false, nil, &skips)
	if err != nil {
		t.Fatal(err)
	}
	if len(started) != 1 || started[0] != f.target {
		t.Fatalf("after the hold cleared started %v, want %s", started, f.target)
	}
}

// A held poll launches nothing the daemon would launch: no resume of an
// interrupted run, no queued start or continue, no plan request, and it
// leaves the queues as they were and a live run running. The poll after
// the hold clears launches each of them.
func TestUsageHoldStopsEveryLaunch(t *testing.T) {
	f := newDaemonFoldFix(t)
	var detached [][]string
	f.ctx.Detach = func(p Program, _ string) error {
		detached = append(detached, p.Args)
		return nil
	}
	common, err := gitCommonDir(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Unix(1700000000, 0)
	create := func(id, cs string, more ...run.Entry) {
		t.Helper()
		if err := f.rs.Create(id, run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Change: cs}}); err != nil {
			t.Fatal(err)
		}
		for _, e := range more {
			if err := f.rs.Append(id, e); err != nil {
				t.Fatal(err)
			}
		}
	}
	// An interrupted run of the accepted change set: the daemon resumes it.
	create("interrupted", f.target)
	// An interrupted run the owner queued a continue for.
	create("continued", "01M6HHMTEDFK2AWQKX50EB7001")
	// A live run, on another host so it counts as running.
	create("live", "01M6HHMTEDFK2AWQKX50EB7002", run.Entry{Type: run.TypeProcess, At: at, Process: &run.Process{Host: "elsewhere", Pid: 1}})
	queuedStart := "01M6HHMTEDFK2AWQKX50EB7003"
	if err := daemon.UpdateStarts(common, func([]daemon.Start) ([]daemon.Start, error) {
		return []daemon.Start{{Change: queuedStart, At: at}}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := daemon.UpdateContinues(common, func([]daemon.Continue) ([]daemon.Continue, error) {
		return []daemon.Continue{{Run: "continued", At: at}}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := daemon.WriteRequests(common, []daemon.Request{{ID: "req1", Kind: daemon.KindFeature, Text: "Dark mode", At: at}}); err != nil {
		t.Fatal(err)
	}
	entries := func(id string) []run.Entry {
		t.Helper()
		j, err := f.rs.Read(id)
		if err != nil {
			t.Fatal(err)
		}
		return j.Entries
	}
	liveBefore := len(entries("live"))

	hold := &usageHold{Kind: "session", Used: 900, Projected: 900, Limit: 1000, Reserve: 20, Reset: f.clock.Add(time.Hour)}
	var skips []daemon.Skip
	started, resumed, planned, err := daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 8, 0, false, "", false, hold, &skips)
	if err != nil {
		t.Fatal(err)
	}
	if len(detached)+len(started)+len(resumed)+len(planned) != 0 {
		t.Fatalf("held poll launched %q: started %v resumed %v planned %v", detached, started, resumed, planned)
	}
	for _, id := range []string{"interrupted", "continued", "live"} {
		for _, e := range entries(id) {
			if e.Type == run.TypeDaemonResumed {
				t.Fatalf("held poll journaled a resume of %s", id)
			}
		}
	}
	if got := entries("live"); len(got) != liveBefore || f.rs.Live("live") {
		t.Fatalf("held poll touched the live run: %d entries, want %d", len(got), liveBefore)
	}
	if j, _ := f.rs.Read("live"); j.State().Ended {
		t.Fatal("held poll ended the live run")
	}
	if s, err := daemon.ReadStarts(common); err != nil || len(s) != 1 || s[0].Change != queuedStart {
		t.Fatalf("held poll changed the queued starts: %v, %v", s, err)
	}
	if c, err := daemon.ReadContinues(common); err != nil || len(c) != 1 || c[0].Run != "continued" {
		t.Fatalf("held poll changed the queued continues: %v, %v", c, err)
	}
	if q, err := daemon.ReadRequests(common); err != nil || len(q) != 1 || q[0].Run != "" {
		t.Fatalf("held poll started the plan request: %v, %v", q, err)
	}

	// The hold clears: the next poll launches each of them.
	skips = nil
	started, resumed, planned, err = daemonPoll1(f.ctx, f.rs, f.pending, f.unj, 8, 0, false, "", false, nil, &skips)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(resumed) != "[interrupted continued]" || fmt.Sprint(started) != "["+queuedStart+"]" || fmt.Sprint(planned) != "[req1]" {
		t.Fatalf("after the hold: started %v resumed %v planned %v", started, resumed, planned)
	}
	if len(detached) != 4 {
		t.Fatalf("after the hold detached %q, want 4", detached)
	}
	if got := entries("live"); len(got) != liveBefore {
		t.Fatalf("the live run was touched: %d entries, want %d", len(got), liveBefore)
	}
}
