package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
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
		return usageHoldOf(f.ctx, common, conf(), learned, runs, "h")
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
	set("rhiza.daemon.usage.session", "700")
	if h := hold(); h == nil || h.Projected != 750 {
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
