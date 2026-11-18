package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// A plan request whose run never journals (run start failed in its tab
// before journaling) goes back in the queue once daemonStartGrace passed,
// and is started again, with a new run, daemonUnjournaledHold after its
// first start.
func TestDaemonPlanRequestsRequeuesUnjournaled(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	f.ctx.repoSession = true // the rhiza-<repo> session is daemon-repo-session's concern
	f.ctx.Self = "/bin/rhiza"
	var detached [][]string
	f.ctx.Detach = func(p Program, logPath string) error {
		detached = append(detached, p.Args)
		return nil
	}
	if _, err := callOp(t, f.ctx, "request", `{"text":"Dark mode"}`); err != nil {
		t.Fatal(err)
	}
	rs := run.NewStore(run.Dir(f.common))
	poll := func() []string {
		t.Helper()
		detached = nil
		got, err := daemonPlanRequests(f.ctx, rs, nil, 0, 2, "auto", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	request := func() daemon.Request {
		t.Helper()
		reqs, err := daemon.ReadRequests(f.common)
		if err != nil || len(reqs) != 1 {
			t.Fatalf("requests %v, %v", reqs, err)
		}
		return reqs[0]
	}
	if got := poll(); len(got) != 1 || len(detached) != 1 {
		t.Fatalf("first poll started %v with rex calls %q, want one start", got, detached)
	}
	first := request()
	if first.Run == "" {
		t.Fatal("started request has no run")
	}
	now = now.Add(daemonStartGrace - time.Second)
	if got := poll(); len(got) != 0 || request().Run != first.Run {
		t.Fatalf("within the grace: started %v, run %q, want none and %q", got, request().Run, first.Run)
	}
	now = now.Add(time.Second)
	if got := poll(); len(got) != 0 || len(detached) != 0 {
		t.Fatalf("after the grace: started %v, want none until the hold passed", got)
	}
	q := request()
	if q.Run != "" || !q.RetryAt.Equal(first.StartedAt.Add(daemonUnjournaledHold)) {
		t.Fatalf("after the grace: run %q, retry at %v, want back in the queue until %v", q.Run, q.RetryAt, first.StartedAt.Add(daemonUnjournaledHold))
	}
	if state, _ := daemon.RequestState(q, nil); state != daemon.RequestQueued {
		t.Errorf("after the grace: state %s, want queued", state)
	}
	now = first.StartedAt.Add(daemonUnjournaledHold)
	if got := poll(); len(got) != 1 || len(detached) != 1 {
		t.Fatalf("after the hold: started %v, want the request again", got)
	}
	if q := request(); q.Run == "" || q.Run == first.Run || !q.RetryAt.IsZero() {
		t.Errorf("after the hold: run %q (first %q), retry at %v, want a new run", q.Run, first.Run, q.RetryAt)
	}
}

func TestRequestPlanRunArgsPassRef(t *testing.T) {
	t.Parallel()
	for _, ref := range []string{"https://example.com/issues/7", "--dry-run", "--change=01M6VY7JMM5NPB3QVG1ECZH1YW", "-h"} {
		with := planRunArgs(daemon.Request{ID: "r1", Ref: ref}, "plan", "default", false, "manual")
		if !slices.Contains(with, "--ref="+ref) || slices.Contains(with, ref) {
			t.Fatalf("with ref %q: want the one word --ref=%s in %q", ref, ref, with)
		}
	}
	without := planRunArgs(daemon.Request{ID: "r2"}, "plan", "default", false, "manual")
	if slices.ContainsFunc(without, func(a string) bool { return strings.HasPrefix(a, "--ref") }) {
		t.Fatalf("without a ref: %q", without)
	}
}

// A queued milestone request is demand while no other milestone request's
// run is live or unjournaled: a started milestone request stays in the
// queue after its run ends, and blocks the next one only while its run
// runs, or while its start is within daemonStartGrace and not journaled.
func TestDaemonMilestoneWaitsAfterEndedRun(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	host, _ := os.Hostname()
	rs := run.NewStore(run.Dir(f.common))
	const first = "01M6YRUNFIRST0000000000000"
	write := func(firstRun string) {
		t.Helper()
		if err := daemon.WriteRequests(f.common, []daemon.Request{
			{ID: "r1", Kind: daemon.KindMilestone, Text: "First", At: now.Add(-time.Hour), Run: firstRun, StartedAt: now.Add(-time.Second)},
			{ID: "r2", Kind: daemon.KindMilestone, Text: "Second", At: now.Add(-time.Minute)},
		}); err != nil {
			t.Fatal(err)
		}
	}
	waits := func() bool {
		t.Helper()
		runs, err := rs.List()
		if err != nil {
			t.Fatal(err)
		}
		return daemonMilestoneWaits(f.ctx, rs, runs, host)
	}

	// The first's start has not journaled, within the grace: it blocks.
	write(first)
	if waits() {
		t.Fatal("first start unjournaled within the grace: second waits, want blocked")
	}
	// Its run journals and runs (a process on another host is running).
	if err := rs.Create(first, run.Entry{Type: run.TypeStarted, At: now, Started: &run.Started{Workflow: "plan-milestone"}}); err != nil {
		t.Fatal(err)
	}
	if err := rs.Append(first, run.Entry{Type: run.TypeProcess, At: now, Process: &run.Process{Host: host + "-elsewhere", Pid: 1}}); err != nil {
		t.Fatal(err)
	}
	if waits() {
		t.Fatal("first run live: second waits, want blocked")
	}
	// Its run ends: the request stays in the queue with its run set, and
	// the second is demand again.
	if err := rs.Append(first, run.Entry{Type: run.TypeFinished, At: now, End: run.Succeeded}); err != nil {
		t.Fatal(err)
	}
	if !waits() {
		t.Fatal("first run finished: second blocked, want it to wait for a slot")
	}
}

// A queued milestone request waits only while another milestone request's
// plan run holds the owner: live at the interview step, stopped at a
// sign-off (or gate) wait, or its start unjournaled within daemonStartGrace.
// Live at any later step it blocks nothing.
func TestDaemonMilestoneWaitsOwnerSteps(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "req")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	host, _ := os.Hostname()
	rs := run.NewStore(run.Dir(f.common))
	const first = "01M6YRUNOWNER00000000000000"
	write := func(startedAt time.Time) {
		t.Helper()
		if err := daemon.WriteRequests(f.common, []daemon.Request{
			{ID: "r1", Kind: daemon.KindMilestone, Text: "First", At: now.Add(-time.Hour), Run: first, StartedAt: startedAt},
			{ID: "r2", Kind: daemon.KindMilestone, Text: "Second", At: now.Add(-time.Minute)},
		}); err != nil {
			t.Fatal(err)
		}
	}
	waits := func() bool {
		t.Helper()
		runs, err := rs.List()
		if err != nil {
			t.Fatal(err)
		}
		return daemonMilestoneWaits(f.ctx, rs, runs, host)
	}
	append1 := func(e run.Entry) {
		t.Helper()
		e.At = now
		if err := rs.Append(first, e); err != nil {
			t.Fatal(err)
		}
	}

	// Unjournaled within the grace: the second waits.
	write(now.Add(-time.Second))
	if waits() {
		t.Fatal("unjournaled within the grace: second can start, want it to wait")
	}
	// Journaled, live (another host) at the interview: still waits.
	if err := rs.Create(first, run.Entry{Type: run.TypeStarted, At: now, Started: &run.Started{Workflow: "plan-milestone"}}); err != nil {
		t.Fatal(err)
	}
	append1(run.Entry{Type: run.TypeProcess, Process: &run.Process{Host: host + "-elsewhere", Pid: 1}})
	append1(run.Entry{Type: run.TypeEntered, Step: "interview"})
	if waits() {
		t.Fatal("live at the interview: second can start, want it to wait")
	}
	// Live at a later step: the second is demand.
	for _, step := range []string{"milestone-check", "review"} {
		append1(run.Entry{Type: run.TypeEntered, Step: step})
		if !waits() {
			t.Fatalf("live at %s: second blocked, want it to start", step)
		}
	}
	// Stopped at the gate, waiting on the owner: it holds the owner again.
	append1(run.Entry{Type: run.TypeEntered, Step: "gate"})
	append1(run.Entry{Type: run.TypeStopped, Detail: "waiting for the owner"})
	if waits() {
		t.Fatal("stopped at the gate: second can start, want it to wait")
	}
}

// rhiza request --note, repeated, queues the request with its notes in
// order, lines kept and outer blank space trimmed, and request ls shows
// them; a note that is blank, too long or not UTF-8 is a usage error that
// queues nothing.
func TestRequestNotes(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "req")
	multi := "\r\n  1. Keep the brief short\r\n2. Name the files\r\n\r\n- and the tests\n\n"
	body, err := json.Marshal(map[string]any{"text": "Dark mode", "notes": []string{"  First idea ", multi}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, f.ctx, "request", string(body)); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, f.ctx, "request", `{"text":"No notes"}`); err != nil {
		t.Fatal(err)
	}
	want := []string{"First idea", "1. Keep the brief short\n2. Name the files\n\n- and the tests"}
	reqs, err := daemon.ReadRequests(f.common)
	if err != nil || len(reqs) != 2 {
		t.Fatalf("requests %v, %v", reqs, err)
	}
	if !slices.Equal(reqs[0].Notes, want) || reqs[1].Notes != nil {
		t.Fatalf("requests.json notes %q and %q, want %q and none", reqs[0].Notes, reqs[1].Notes, want)
	}
	resp, err := callOp(t, f.ctx, "request.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	ls := resp.(RequestLsResponse)
	if !slices.Equal(ls.Requests[0].Notes, want) || ls.Requests[1].Notes != nil {
		t.Fatalf("request ls notes %q and %q, want %q and none", ls.Requests[0].Notes, ls.Requests[1].Notes, want)
	}
	out, err := json.Marshal(ls)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Requests []map[string]json.RawMessage `json:"requests"`
	}
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	if _, ok := doc.Requests[1]["notes"]; ok {
		t.Errorf("request ls --json: a notes key for the request without notes in %s", out)
	}
	if text := human(t, resp); !strings.Contains(text, "Dark mode [2 notes]") || strings.Contains(text, "No notes [") {
		t.Errorf("request ls text:\n%s\nwant [2 notes] on the first request only", text)
	}

	for name, note := range map[string]string{
		"empty":     "",
		"blank":     " \r\n\t\n ",
		"too long":  strings.Repeat("x", daemon.MaxRequestNote+1),
		"not UTF-8": "bad \xff byte",
	} {
		// Built as the request itself: JSON would replace the bad bytes.
		o, _ := Lookup("request")
		_, err := o.Run(f.ctx, RequestRequest{Text: "Refused " + name, Notes: []string{"fine", note}})
		if e := AsError(err); e == nil || e.Code != CodeUsage || !strings.Contains(e.Message, "--note 2") {
			t.Errorf("%s note: %v, want a usage error naming --note 2", name, err)
		}
		if reqs, err := daemon.ReadRequests(f.common); err != nil || len(reqs) != 2 {
			t.Errorf("%s note: %d requests, %v; want nothing queued", name, len(reqs), err)
		}
	}
}

// rhiza request note adds a note to a queued request, named by its ID or a
// unique prefix, after those --note gave it; --dry-run writes nothing. A
// request whose plan run is drawn, an unknown ID, an ambiguous or short
// prefix and an empty note are usage errors that leave requests.json as it
// was.
func TestRequestNoteAdd(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "req")
	resp, err := callOp(t, f.ctx, "request", `{"text":"Dark mode","notes":["First idea"]}`)
	if err != nil {
		t.Fatal(err)
	}
	id := resp.(RequestResponse).Request.ID
	body, err := json.Marshal(map[string]string{"request": id, "text": "  Second\r\nidea  "})
	if err != nil {
		t.Fatal(err)
	}
	resp, err = callOp(t, f.ctx, "request.note", string(body))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"First idea", "Second\nidea"}
	if got := resp.(RequestNoteResponse).Request.Notes; !slices.Equal(got, want) {
		t.Fatalf("request note printed notes %q, want %q", got, want)
	}
	if text := human(t, resp); !strings.Contains(text, "Note 2:\n  Second\n  idea\n") {
		t.Errorf("request note text:\n%s\nwant the second note's lines", text)
	}
	ls, err := callOp(t, f.ctx, "request.ls", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if got := ls.(RequestLsResponse).Requests[0].Notes; !slices.Equal(got, want) {
		t.Fatalf("request ls notes %q, want %q", got, want)
	}

	// A unique prefix resolves; a dry run writes nothing.
	f.ctx.DryRun = true
	resp, err = callOp(t, f.ctx, "request.note", `{"request":"`+id[:6]+`","text":"Third"}`)
	f.ctx.DryRun = false
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(RequestNoteResponse); !r.DryRun || len(r.Request.Notes) != 3 || r.Request.ID != id {
		t.Fatalf("dry run by prefix: %+v", r)
	}
	resp, err = callOp(t, f.ctx, "request.note", `{"request":"`+id[:6]+`","text":"Third"}`)
	if err != nil {
		t.Fatal(err)
	}
	if got := resp.(RequestNoteResponse).Request.Notes; len(got) != 3 || got[2] != "Third" {
		t.Fatalf("prefix note: notes %q", got)
	}

	// Refusals: set up a drawn request and two sharing a prefix.
	reqs, err := daemon.ReadRequests(f.common)
	if err != nil {
		t.Fatal(err)
	}
	at := reqs[0].At
	reqs = append(reqs,
		daemon.Request{ID: "01ZZAAAA0000000000000000AA", Kind: daemon.KindFeature, Text: "Twin one", At: at},
		daemon.Request{ID: "01ZZAAAA0000000000000000BB", Kind: daemon.KindFeature, Text: "Twin two", At: at},
		daemon.Request{ID: "01YYDRAWN00000000000000000", Kind: daemon.KindFeature, Text: "Drawn", At: at, Run: "01YYRUN0000000000000000000", StartedAt: at},
	)
	if err := daemon.WriteRequests(f.common, reqs); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(daemon.Dir(f.common), daemon.RequestsFile))
	if err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"drawn":     `{"request":"01YYDRAWN00000000000000000","text":"Late"}`,
		"unknown":   `{"request":"01XXXXXXXXXXXXXXXXXXXXXXXX","text":"Lost"}`,
		"ambiguous": `{"request":"01ZZAAAA","text":"Which"}`,
		"short":     `{"request":"` + id[:3] + `","text":"Short"}`,
		"empty":     `{"request":"` + id + `","text":" \n "}`,
	} {
		_, err := callOp(t, f.ctx, "request.note", body)
		if e := AsError(err); e == nil || e.Code != CodeUsage {
			t.Errorf("%s: %v, want a usage error", name, err)
		}
		after, err := os.ReadFile(filepath.Join(daemon.Dir(f.common), daemon.RequestsFile))
		if err != nil || string(after) != string(before) {
			t.Errorf("%s: requests.json changed (%v)", name, err)
		}
	}
	// Not UTF-8, built as the request itself: JSON would replace the bytes.
	o, _ := Lookup("request.note")
	if _, err := o.Run(f.ctx, RequestNoteRequest{Request: id, Text: "bad \xff"}); AsError(err) == nil || AsError(err).Code != CodeUsage {
		t.Errorf("not UTF-8: %v, want a usage error", err)
	}
}
