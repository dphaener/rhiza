package api

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/session"
)

func TestParsePRView(t *testing.T) {
	in := `{"number":12,"url":"https://github.com/o/r/pull/12","state":"OPEN","headRefOid":"abc1234def",
	"statusCheckRollup":[
	  {"__typename":"CheckRun","name":"test","status":"COMPLETED","conclusion":"SUCCESS"},
	  {"__typename":"CheckRun","name":"lint","status":"IN_PROGRESS","conclusion":""},
	  {"__typename":"StatusContext","context":"ci/legacy","state":"FAILURE"},
	  {"__typename":"StatusContext","context":"ci/pending","state":"PENDING"},
	  {"__typename":"StatusContext","context":"ci/expected","state":"EXPECTED"}
	]}`
	v, err := parsePRView([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if v.Number != 12 || v.URL != "https://github.com/o/r/pull/12" || v.State != "OPEN" || v.Head != "abc1234def" {
		t.Errorf("parsePRView = %+v", v)
	}
	want := []ShipCheck{
		{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"},
		{Name: "lint", Status: "IN_PROGRESS"},
		{Name: "ci/legacy", Status: "COMPLETED", Conclusion: "FAILURE"},
		{Name: "ci/pending", Status: "PENDING"},
		{Name: "ci/expected", Status: "EXPECTED"},
	}
	if !slices.Equal(v.Checks, want) {
		t.Errorf("checks = %+v, want %+v", v.Checks, want)
	}

	v, err = parsePRView([]byte(`{"number":3,"state":"OPEN","headRefOid":"x","statusCheckRollup":[]}`))
	if err != nil || v.Checks == nil || len(v.Checks) != 0 {
		t.Errorf("empty rollup: %+v, %v", v, err)
	}
	if _, err := parsePRView([]byte(`not json`)); err == nil {
		t.Error("parsePRView accepted bad JSON")
	}
}

func TestParseRunList(t *testing.T) {
	got, err := parseRunList([]byte(`[{"name":"CI","status":"completed","conclusion":"success"},{"name":"Rhiza","status":"in_progress","conclusion":""}]`))
	if err != nil {
		t.Fatal(err)
	}
	want := []ShipCheck{{Name: "CI", Status: "completed", Conclusion: "success"}, {Name: "Rhiza", Status: "in_progress"}}
	if !slices.Equal(got, want) {
		t.Errorf("parseRunList = %+v, want %+v", got, want)
	}
	if got, err := parseRunList([]byte(`[]`)); err != nil || got == nil || len(got) != 0 {
		t.Errorf("empty list: %v, %v", got, err)
	}
	if _, err := parseRunList([]byte(`{}`)); err == nil {
		t.Error("parseRunList accepted an object")
	}
}

func TestChecksVerdict(t *testing.T) {
	ok := ShipCheck{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"}
	cases := []struct {
		name           string
		checks         []ShipCheck
		settled, green bool
		failed         []string
	}{
		{"Empty", nil, false, false, nil},
		{"Green", []ShipCheck{ok, {Name: "run", Status: "completed", Conclusion: "success"}}, true, true, nil},
		{"NeutralSkipped", []ShipCheck{ok, {Name: "n", Status: "COMPLETED", Conclusion: "NEUTRAL"}, {Name: "s", Status: "completed", Conclusion: "skipped"}}, true, true, nil},
		{"InProgress", []ShipCheck{ok, {Name: "lint", Status: "IN_PROGRESS"}}, false, false, nil},
		{"Queued", []ShipCheck{{Name: "lint", Status: "queued"}}, false, false, nil},
		{"StatusPending", []ShipCheck{ok, {Name: "ctx", Status: "PENDING"}}, false, false, nil},
		{"Failed", []ShipCheck{ok, {Name: "lint", Status: "COMPLETED", Conclusion: "FAILURE"}, {Name: "e2e", Status: "completed", Conclusion: "cancelled"}}, true, false, []string{"lint", "e2e"}},
		{"StatusError", []ShipCheck{{Name: "ctx", Status: "COMPLETED", Conclusion: "ERROR"}}, true, false, []string{"ctx"}},
		{"FailedAndPending", []ShipCheck{{Name: "lint", Status: "COMPLETED", Conclusion: "FAILURE"}, {Name: "e2e", Status: "IN_PROGRESS"}}, false, false, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			settled, green, failed := checksVerdict(c.checks)
			if settled != c.settled || green != c.green || !slices.Equal(failed, c.failed) {
				t.Errorf("checksVerdict = %v, %v, %v; want %v, %v, %v", settled, green, failed, c.settled, c.green, c.failed)
			}
		})
	}
}

// fakeGH answers gh calls by their joined argv; an answer starting with
// "!" fails with the rest as gh's message.
type fakeGH struct {
	answers map[string]string
	calls   []string
	dirs    []string
}

func (f *fakeGH) run(name string, args []string, dir string, _ []string) ([]byte, error) {
	k := strings.Join(args, " ")
	f.calls = append(f.calls, name+" "+k)
	f.dirs = append(f.dirs, dir)
	if name != "gh" {
		return nil, errors.New("not gh")
	}
	a, ok := f.answers[k]
	if !ok {
		return nil, errors.New("gh " + k + ": unexpected call")
	}
	if msg, failed := strings.CutPrefix(a, "!"); failed {
		return nil, errors.New("gh " + k + ": " + msg)
	}
	return []byte(a), nil
}

const (
	shipPRView   = "pr view cs/ship-it --json number,url,state,headRefOid,statusCheckRollup"
	shipMergeArg = "pr merge 7 --squash --match-head-commit abc1234def"
	shipMergedPR = "pr view 7 --json mergeCommit,state"
	shipRunList  = "run list --commit fff0000aaa --json name,status,conclusion"
)

// shipFixture is a session fixture with the change set "Ship it" current
// and a fake gh.
func shipFixture(t *testing.T) (*sessionFixture, *fakeGH, string) {
	t.Helper()
	f := newSessionFixture(t, "ship-it")
	cs := newChange(t, f.ctx, `{"intent":"Ship it","title":"Ship it"}`).ChangeSet
	g := &fakeGH{answers: map[string]string{}}
	f.ctx.Exec = g.run
	return f, g, filepath.Join(f.common, "rhiza", "ship", cs.String()+".json")
}

func openPR(rollup string) string {
	return `{"number":7,"url":"https://github.com/o/r/pull/7","state":"OPEN","headRefOid":"abc1234def","statusCheckRollup":[` + rollup + `]}`
}

const greenRollup = `{"name":"test","status":"COMPLETED","conclusion":"SUCCESS"},{"context":"ci/x","state":"SUCCESS"}`

func shipCheck(t *testing.T, f *sessionFixture, req string) ShipCheckResponse {
	t.Helper()
	resp, err := callOp(t, f.ctx, "ship.check", req)
	if err != nil {
		t.Fatalf("ship.check %s: %v", req, err)
	}
	return resp.(ShipCheckResponse)
}

func TestShipCheckNoPR(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = "!no pull requests found for branch \"cs/ship-it\""
	for _, req := range []string{`{}`, `{"settled":true}`} {
		r := shipCheck(t, f, req)
		if r.ExitCode() != 1 || r.Passed || r.Written || !strings.Contains(r.Detail, "no open pull request") {
			t.Errorf("%s: %+v", req, r)
		}
		if got := human(t, r); !strings.Contains(got, "no open pull request for branch cs/ship-it") {
			t.Errorf("human = %q", got)
		}
	}
	if g.dirs[0] != f.worktree {
		t.Errorf("gh ran in %s, want the worktree %s", g.dirs[0], f.worktree)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("state file written: %v", err)
	}
}

func TestShipCheckGhFails(t *testing.T) {
	f, g, _ := shipFixture(t)
	g.answers[shipPRView] = "!HTTP 401: Bad credentials"
	_, err := callOp(t, f.ctx, "ship.check", `{}`)
	wantExit(t, err, 1)
	if !strings.Contains(err.Error(), "Bad credentials") {
		t.Errorf("err = %v", err)
	}
}

func TestShipCheckClosedPR(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = strings.Replace(openPR(greenRollup), `"OPEN"`, `"CLOSED"`, 1)
	r := shipCheck(t, f, `{"settled":true}`)
	if r.Passed || r.Settled || !strings.Contains(r.Detail, "is CLOSED, not open") {
		t.Errorf("closed PR: %+v", r)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("state file written: %v", err)
	}
}

func TestShipCheckPending(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(`{"name":"test","status":"IN_PROGRESS","conclusion":""}`)
	for _, req := range []string{`{}`, `{"settled":true}`} {
		r := shipCheck(t, f, req)
		if r.Passed || r.Settled || r.ExitCode() != 1 {
			t.Errorf("%s: %+v", req, r)
		}
		if got := human(t, r); got != "PR #7 at abc1234: not settled: 0 of 1 checks finished\n" {
			t.Errorf("human = %q", got)
		}
	}
	g.answers[shipPRView] = openPR(``)
	if r := shipCheck(t, f, `{"settled":true}`); r.Passed || !strings.Contains(r.Detail, "no checks reported yet") {
		t.Errorf("no checks: %+v", r)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("state file written: %v", err)
	}
}

func TestShipCheckFailed(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(`{"name":"test","status":"COMPLETED","conclusion":"FAILURE"},{"name":"lint","status":"COMPLETED","conclusion":"SUCCESS"},{"context":"ci/x","state":"ERROR"}`)
	if r := shipCheck(t, f, `{"settled":true}`); !r.Passed || r.ExitCode() != 0 || r.Written {
		t.Errorf("settled: %+v", r)
	}
	r := shipCheck(t, f, `{}`)
	if r.Passed || r.ExitCode() != 1 || !slices.Equal(r.Failed, []string{"test", "ci/x"}) || r.Written {
		t.Errorf("failed: %+v", r)
	}
	if got := human(t, r); got != "PR #7 at abc1234: failed: test, ci/x\n" {
		t.Errorf("human = %q", got)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("state file written: %v", err)
	}
}

func TestShipCheckGreenThenMerge(t *testing.T) {
	f, g, state := shipFixture(t)

	_, err := callOp(t, f.ctx, "ship.merge", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza ship check")
	_, err = callOp(t, f.ctx, "ship.check", `{"main":true}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza ship merge")

	g.answers[shipPRView] = openPR(greenRollup)
	if r := shipCheck(t, f, `{"settled":true}`); !r.Passed || r.Written {
		t.Errorf("settled: %+v", r)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("--settled wrote the state file: %v", err)
	}
	r := shipCheck(t, f, `{}`)
	if !r.Passed || !r.Green || !r.Written || r.ExitCode() != 0 {
		t.Fatalf("green: %+v", r)
	}
	if got := human(t, r); got != "PR #7 at abc1234: 2 checks passed\n" {
		t.Errorf("human = %q", got)
	}
	st, err := readShipState(state)
	if err != nil {
		t.Fatal(err)
	}
	if st.PR != 7 || st.URL != "https://github.com/o/r/pull/7" || st.Branch != "cs/ship-it" || st.Head != "abc1234def" || len(st.Checks) != 2 {
		t.Errorf("state = %+v", st)
	}

	g.answers[shipMergeArg] = ""
	g.answers[shipMergedPR] = `{"mergeCommit":{"oid":"fff0000aaa"},"state":"MERGED"}`
	resp, err := callOp(t, f.ctx, "ship.merge", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	m := resp.(ShipMergeResponse)
	if m.PR != 7 || m.Head != "abc1234def" || m.MergeCommit != "fff0000aaa" || m.AlreadyMerged {
		t.Errorf("merge = %+v", m)
	}
	if got := human(t, m); got != "Merged PR #7 at abc1234 as fff0000.\n" {
		t.Errorf("human = %q", got)
	}
	for _, c := range g.calls {
		if strings.Contains(c, "--admin") || strings.Contains(c, "--auto") || strings.Contains(c, "--delete-branch") {
			t.Errorf("gh call %q", c)
		}
	}
	if !slices.Contains(g.calls, "gh "+shipMergeArg) {
		t.Errorf("no pinned merge in %q", g.calls)
	}
	n := len(g.calls)
	resp, err = callOp(t, f.ctx, "ship.merge", `{}`)
	if err != nil || !resp.(ShipMergeResponse).AlreadyMerged || len(g.calls) != n {
		t.Errorf("second merge: %+v, %v, calls %q", resp, err, g.calls[n:])
	}

	g.answers[shipRunList] = `[{"name":"CI","status":"in_progress","conclusion":""}]`
	if r := shipCheck(t, f, `{"main":true,"settled":true}`); r.Passed || r.Settled {
		t.Errorf("main pending: %+v", r)
	}
	g.answers[shipRunList] = `[{"name":"CI","status":"completed","conclusion":"success"}]`
	r = shipCheck(t, f, `{"main":true}`)
	if !r.Passed || !r.Written || r.Commit != "fff0000aaa" {
		t.Errorf("main green: %+v", r)
	}
	if got := human(t, r); got != "main at fff0000: 1 check passed\n" {
		t.Errorf("human = %q", got)
	}
	if st, _ := readShipState(state); len(st.Main) != 1 || st.MergeCommit != "fff0000aaa" || st.Head != "abc1234def" {
		t.Errorf("state = %+v", st)
	}
}

func TestShipMergeNotMerged(t *testing.T) {
	f, g, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Branch: "cs/ship-it", Head: "abc1234def"}); err != nil {
		t.Fatal(err)
	}
	g.answers[shipMergeArg] = "!Head branch was modified. Review and try the merge again."
	_, err := callOp(t, f.ctx, "ship.merge", `{}`)
	wantExit(t, err, 1)
	if !strings.Contains(err.Error(), "Head branch was modified") {
		t.Errorf("err = %v", err)
	}
	g.answers[shipMergeArg] = ""
	g.answers[shipMergedPR] = `{"mergeCommit":null,"state":"OPEN"}`
	_, err = callOp(t, f.ctx, "ship.merge", `{}`)
	wantExit(t, err, 1)
	if st, _ := readShipState(state); st.MergeCommit != "" {
		t.Errorf("state = %+v", st)
	}
}

func TestShipTargetWithoutWorktree(t *testing.T) {
	f, g, _ := shipFixture(t)
	if err := os.RemoveAll(f.worktree); err != nil {
		t.Fatal(err)
	}
	g.answers[shipPRView] = "!no pull requests found for branch \"cs/ship-it\""
	shipCheck(t, f, `{}`)
	if main := filepath.Dir(f.common); g.dirs[0] != main {
		t.Errorf("gh ran in %s, want the main checkout %s", g.dirs[0], main)
	}
}

func TestShipStateRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "cs.json")
	if st, err := readShipState(path); err != nil || st.PR != 0 {
		t.Fatalf("missing state: %+v, %v", st, err)
	}
	want := shipState{PR: 3, URL: "u", Branch: "cs/x", Head: "h", Checks: []ShipCheck{{Name: "a", Status: "COMPLETED", Conclusion: "SUCCESS"}}, MergeCommit: "m"}
	if err := writeShipState(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := readShipState(path)
	if err != nil || got.PR != 3 || got.Head != "h" || got.MergeCommit != "m" || !slices.Equal(got.Checks, want.Checks) {
		t.Errorf("round trip = %+v, %v", got, err)
	}
	entries, _ := os.ReadDir(filepath.Dir(path))
	if len(entries) != 1 {
		t.Errorf("left %d files, want 1", len(entries))
	}
}

func TestShipEvidenceHandle(t *testing.T) {
	if got := shipEvidenceHandle(7, "ship-it"); got != "pr-7-ship-it" {
		t.Errorf("handle = %q", got)
	}
	long := strings.Repeat("abcdefghi-", 10)
	got := shipEvidenceHandle(1234, long)
	if len(got) > maxHandle || !strings.HasPrefix(got, "pr-1234-abcdefghi-") || strings.HasSuffix(got, "-") || strings.Contains(got, "--") {
		t.Errorf("long handle = %q (%d bytes)", got, len(got))
	}
}

func TestShipEvidenceFields(t *testing.T) {
	st := shipState{PR: 7, URL: "https://github.com/o/r/pull/7", Head: "abc1234def",
		Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"}, {Name: "ci/x", Status: "COMPLETED", Conclusion: "NEUTRAL"}}}
	want := []string{
		"statement=Pull request #7 (https://github.com/o/r/pull/7) passed CI at abc1234",
		"result=2 checks passed: test SUCCESS, ci/x NEUTRAL",
		"pr=https://github.com/o/r/pull/7",
		"head=abc1234def",
		"status=verified",
	}
	if got := shipEvidenceFields(st); !slices.Equal(got, want) {
		t.Errorf("fields =\n%q\nwant\n%q", got, want)
	}
	st.Checks = st.Checks[:1]
	if got := shipEvidenceFields(st)[1]; got != "result=1 check passed: test SUCCESS" {
		t.Errorf("one check: %q", got)
	}
}

// shipTestActor is the handle of the actor testCtx's git user.email
// (test@example.com) records on the first write.
const shipTestActor = "test"

// shipLandFixture is a ship fixture whose change set adds a capability and
// is accepted, with the PR's checks recorded and git answering as on the
// PR branch.
func shipLandFixture(t *testing.T) (*sessionFixture, string, id.ID) {
	t.Helper()
	f, _, state := shipFixture(t)
	capID := mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"shipping"}`).Nodes[0]
	for _, v := range []string{"propose", "accept"} {
		if _, err := callOp(t, f.ctx, v, `{}`); err != nil {
			t.Fatalf("%s: %v", v, err)
		}
	}
	if err := writeShipState(state, shipState{PR: 7, URL: "https://github.com/o/r/pull/7", Branch: "cs/ship-it", Head: "abc1234def",
		Checks: []ShipCheck{{Name: "test", Status: "COMPLETED", Conclusion: "SUCCESS"}}}); err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{
		"symbolic-ref --short HEAD":           "cs/ship-it\n",
		"add -- .rhiza":                       "",
		"diff --cached --name-only -- .rhiza": ".rhiza/rendered/x.md\n",
		"commit -m " + shipLandMessage + " -m Change: " + currentOf(t, f.ctx).String() + "\nActor: " + shipTestActor + " -- .rhiza": "",
		"rev-parse HEAD": "beef000111\n",
		"rev-parse --abbrev-ref --symbolic-full-name @{upstream}": "origin/cs/ship-it\n",
		"rev-list --count @{upstream}..HEAD":                      "1\n",
		"push":                                                    "",
	} {
		f.git.out[k] = v
	}
	return f, state, capID
}

func TestShipLand(t *testing.T) {
	f, state, capID := shipLandFixture(t)
	cs := currentOf(t, f.ctx)
	resp, err := callOp(t, f.ctx, "ship.land", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipLandResponse)
	if !r.EvidenceWritten || r.Evidence != "pr-7-ship-it" || r.EvidenceChange != cs.String() || !r.Landed ||
		!slices.Equal(r.Verifies, []id.ID{capID}) || r.Commit != "beef000111" || r.Pushed != "beef000111" || len(r.Rendered) == 0 {
		t.Fatalf("land = %+v", r)
	}
	w := ctxWorkspace(t, f.ctx)
	if c, _ := w.changeSet(cs); c.State != changeset.StateLanded {
		t.Errorf("change set is %s", c.State)
	}
	m := w.model(id.ID{})
	ev, err := m.Resolve("pr-7-ship-it")
	if err != nil || ev.Kind != op.KindEvidence || fieldText(ev.Fields["pr"]) != "https://github.com/o/r/pull/7" ||
		!strings.Contains(fieldText(ev.Fields["result"]), "test SUCCESS") {
		t.Fatalf("evidence = %+v, %v", ev, err)
	}
	if _, err := m.LinkBetween(ev.ID, capID, op.LinkVerifies); err != nil {
		t.Errorf("verifies link: %v", err)
	}
	if st, _ := readShipState(state); st.Pushed != "beef000111" || st.Head != "abc1234def" {
		t.Errorf("state = %+v", st)
	}
	for _, c := range f.git.calls {
		if strings.Contains(c, "--force") || strings.Contains(c, "--no-verify") {
			t.Errorf("git call %q", c)
		}
	}
	if got := human(t, r); !strings.Contains(got, "Recorded PR #7 as evidence pr-7-ship-it, verifying 1 capability.") ||
		!strings.Contains(got, "Committed beef000.") || !strings.Contains(got, "Pushed beef000.") {
		t.Errorf("human = %q", got)
	}

	// Again, as a resumed run would (landing cleared the current change
	// set; a step names it by RHIZA_CHANGE): nothing is left to do.
	f.ctx.ChangeRef = cs.String()
	f.git.out["diff --cached --name-only -- .rhiza"] = ""
	f.git.out["rev-list --count @{upstream}..HEAD"] = "0\n"
	resp, err = callOp(t, f.ctx, "ship.land", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(ShipLandResponse)
	if r.EvidenceWritten || r.Landed || r.Commit != "" || r.Pushed != "" {
		t.Errorf("second land = %+v", r)
	}
	if got := human(t, r); !strings.Contains(got, "already recorded") || !strings.Contains(got, "Nothing to commit.") {
		t.Errorf("human = %q", got)
	}
}

func TestShipLandAfterAgentLanded(t *testing.T) {
	f, _, capID := shipLandFixture(t)
	cs := currentOf(t, f.ctx)
	if _, err := callOp(t, f.ctx, "land", `{}`); err != nil {
		t.Fatal(err)
	}
	f.ctx.ChangeRef = cs.String()
	delete(f.git.out, "rev-parse --abbrev-ref --symbolic-full-name @{upstream}")
	f.git.out["push -u origin cs/ship-it"] = ""
	resp, err := callOp(t, f.ctx, "ship.land", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipLandResponse)
	if !r.EvidenceWritten || r.Landed || r.EvidenceChange == cs.String() || r.EvidenceChange == "" || r.Pushed == "" {
		t.Fatalf("land = %+v", r)
	}
	w := ctxWorkspace(t, f.ctx)
	rec, _ := id.Parse(r.EvidenceChange)
	if c, _ := w.changeSet(rec); c.State != changeset.StateLanded || c.Open.Kind != changeset.KindRecord {
		t.Errorf("record change set: %s %s", c.State, c.Open.Kind)
	}
	m := w.model(id.ID{})
	ev, err := m.Resolve("pr-7-ship-it")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.LinkBetween(ev.ID, capID, op.LinkVerifies); err != nil {
		t.Errorf("verifies link: %v", err)
	}
	if !slices.Contains(f.git.calls, "push -u origin cs/ship-it") {
		t.Errorf("no push -u in %q", f.git.calls)
	}
}

func TestShipLandRefusals(t *testing.T) {
	f, state, _ := shipLandFixture(t)
	f.git.out["symbolic-ref --short HEAD"] = "main\n"
	_, err := callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 3)
	if !strings.Contains(err.Error(), "on main") {
		t.Errorf("err = %v", err)
	}
	if err := os.Remove(state); err != nil {
		t.Fatal(err)
	}
	_, err = callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza ship check")
	f.ctx.DryRun = true
	_, err = callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 2)
	if c, _ := ctxWorkspace(t, f.ctx).changeSet(currentOf(t, f.ctx)); c.State != changeset.StateAccepted {
		t.Errorf("change set is %s", c.State)
	}
}

// shipExec answers gh and git by "<name> <args>"; "!" fails with the rest.
type shipExec struct {
	answers map[string]string
	calls   []string
	dirs    []string
}

func (s *shipExec) run(name string, args []string, dir string, _ []string) ([]byte, error) {
	k := name + " " + strings.Join(args, " ")
	s.calls = append(s.calls, k)
	s.dirs = append(s.dirs, dir)
	a, ok := s.answers[k]
	if !ok {
		return nil, errors.New(k + ": unexpected call")
	}
	if msg, failed := strings.CutPrefix(a, "!"); failed {
		return nil, errors.New(k + ": " + msg)
	}
	return []byte(a), nil
}

func TestShipCloseOut(t *testing.T) {
	f, _, state := shipFixture(t)
	x := &shipExec{answers: map[string]string{}}
	f.ctx.Exec = x.run
	_, err := callOp(t, f.ctx, "ship.close-out", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza ship merge")
	if len(x.calls) != 0 {
		t.Errorf("calls without a merge: %q", x.calls)
	}
	if err := writeShipState(state, shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}); err != nil {
		t.Fatal(err)
	}
	main := filepath.Dir(f.common)
	x.answers = map[string]string{
		"git worktree remove " + f.worktree:                    "",
		"git rev-parse --verify --quiet refs/heads/cs/ship-it": "2222\n",
		"git branch -D cs/ship-it":                             "Deleted branch cs/ship-it\n",
		"git push origin --delete cs/ship-it":                  "!error: unable to delete 'cs/ship-it': remote ref does not exist",
		"git symbolic-ref --short HEAD":                        "main\n",
		"git status --porcelain --untracked-files=no":          "",
		"git pull --ff-only origin main":                       "",
	}
	resp, err := callOp(t, f.ctx, "ship.close-out", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipCloseOutResponse)
	want := []ShipStep{{Step: "remove-worktree", Result: "done"}, {Step: "delete-branch", Result: "done"},
		{Step: "delete-remote-branch", Result: "already"}, {Step: "update-main", Result: "done"}}
	if !r.MainUpdated || r.ExitCode() != 0 || !slices.Equal(r.Steps, want) {
		t.Errorf("close-out = %+v", r)
	}
	for i, d := range x.dirs {
		if d != main {
			t.Errorf("%q ran in %s, want the main checkout %s", x.calls[i], d, main)
		}
	}
	for _, c := range x.calls {
		if strings.Contains(c, "--force") || strings.Contains(c, " -f ") {
			t.Errorf("forced: %q", c)
		}
	}
	got := human(t, r)
	for _, s := range []string{"Removed worktree " + f.worktree + ".", "Deleted branch cs/ship-it.", "origin/cs/ship-it was already gone.", "Updated main in " + main + "."} {
		if !strings.Contains(got, s) {
			t.Errorf("human = %q, want %q", got, s)
		}
	}

	// The worktree and branch gone, main on another branch.
	if err := os.RemoveAll(f.worktree); err != nil {
		t.Fatal(err)
	}
	x.answers["git rev-parse --verify --quiet refs/heads/cs/ship-it"] = "!"
	x.answers["git symbolic-ref --short HEAD"] = "feature\n"
	resp, err = callOp(t, f.ctx, "ship.close-out", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(ShipCloseOutResponse)
	if r.MainUpdated || r.ExitCode() != 1 || !strings.Contains(r.Detail, "main was not updated") || !strings.Contains(r.Detail, "on feature, not main") {
		t.Errorf("not on main = %+v", r)
	}
	if r.Steps[0].Result != "already" || r.Steps[1].Result != "already" {
		t.Errorf("steps = %+v", r.Steps)
	}
	if got := human(t, r); !strings.Contains(got, "Worktree "+f.worktree+" was already gone.") || !strings.Contains(got, "main was not updated") {
		t.Errorf("human = %q", got)
	}
}

func TestShipCloseOutDirtyWorktree(t *testing.T) {
	f, _, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}); err != nil {
		t.Fatal(err)
	}
	x := &shipExec{answers: map[string]string{
		"git worktree remove " + f.worktree: "!fatal: '" + f.worktree + "' contains modified or untracked files, use --force to delete it",
	}}
	f.ctx.Exec = x.run
	_, err := callOp(t, f.ctx, "ship.close-out", `{}`)
	wantExit(t, err, 1)
	if !strings.Contains(err.Error(), "contains modified or untracked files") || len(x.calls) != 1 {
		t.Errorf("err = %v, calls %q", err, x.calls)
	}
}

func TestShipCommand(t *testing.T) {
	f, _, _ := shipFixture(t)
	f.ctx.Start = nil
	resp, err := callOp(t, f.ctx, "ship.command", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipCommandResponse)
	main := filepath.Dir(f.common)
	if r.Ran || r.ExitCode() != 0 || r.Dir != main {
		t.Errorf("no script = %+v", r)
	}
	if got := human(t, r); got != "no ship command (.rhiza/workflows/ship.sh)\n" {
		t.Errorf("human = %q", got)
	}

	script := filepath.Join(main, ".rhiza", "workflows", "ship.sh")
	if err := os.MkdirAll(filepath.Dir(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(script, []byte("#!/bin/sh\nmake install\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var started []Program
	f.ctx.Start = func(p Program) (Process, error) {
		started = append(started, p)
		return &fakeProcess{out: strings.NewReader("installed\n"), code: 4}, nil
	}
	var stderr strings.Builder
	f.ctx.Stderr = &stderr
	resp, err = callOp(t, f.ctx, "ship.command", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(ShipCommandResponse)
	if !r.Ran || r.ExitCode() != 4 || r.Script != script {
		t.Errorf("script = %+v", r)
	}
	if len(started) != 1 || started[0].Name != "sh" || started[0].Dir != main || started[0].Args[len(started[0].Args)-1] != script {
		t.Errorf("started %+v", started)
	}
	if stderr.String() != "installed\n" {
		t.Errorf("stderr = %q", stderr.String())
	}
	if got := human(t, r); got != "Ran .rhiza/workflows/ship.sh: exit 4.\n" {
		t.Errorf("human = %q", got)
	}
}

func TestShipCheckWaitsForPushedHead(t *testing.T) {
	f, g, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Head: "0ld0000000", Pushed: "abc1234def"}); err != nil {
		t.Fatal(err)
	}
	g.answers[shipPRView] = strings.Replace(openPR(greenRollup), "abc1234def", "0ld0000000", 1)
	if r := shipCheck(t, f, `{"settled":true}`); r.Passed || r.Settled || !strings.Contains(r.Detail, "not the pushed abc1234") {
		t.Errorf("old head: %+v", r)
	}
	g.answers[shipPRView] = openPR(greenRollup)
	if r := shipCheck(t, f, `{}`); !r.Passed || !r.Written {
		t.Errorf("pushed head: %+v", r)
	}
	if st, _ := readShipState(state); st.Head != "abc1234def" {
		t.Errorf("state = %+v", st)
	}
}

// --pr-open passes on an open PR whatever its CI, records nothing, and
// fails with no PR or a closed one.
func TestShipCheckPROpen(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = "!no pull requests found for branch \"cs/ship-it\""
	if r := shipCheck(t, f, `{"pr_open":true}`); r.Passed || r.ExitCode() != 1 {
		t.Errorf("no PR: %+v", r)
	}
	g.answers[shipPRView] = openPR(`{"name":"test","status":"IN_PROGRESS","conclusion":""}`)
	r := shipCheck(t, f, `{"pr_open":true}`)
	if !r.Passed || r.ExitCode() != 0 || r.Written || r.Settled {
		t.Errorf("open PR with pending CI: %+v", r)
	}
	if got := human(t, r); got != "PR #7 at abc1234: open from cs/ship-it\n" {
		t.Errorf("human = %q", got)
	}
	g.answers[shipPRView] = strings.Replace(openPR(greenRollup), `"OPEN"`, `"MERGED"`, 1)
	if r := shipCheck(t, f, `{"pr_open":true}`); r.Passed {
		t.Errorf("merged PR passed --pr-open: %+v", r)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("state file written: %v", err)
	}
	for _, req := range []string{`{"pr_open":true,"main":true}`, `{"pr_open":true,"settled":true}`} {
		_, err := callOp(t, f.ctx, "ship.check", req)
		wantExit(t, err, 2)
	}
}

// Close-out closes the Rex tab of each of the change set's sessions after
// the worktree is gone; a tab Rex no longer has is already closed, and
// other change sets' tabs are left alone.
func TestShipCloseOutClosesTabs(t *testing.T) {
	f, _, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}); err != nil {
		t.Fatal(err)
	}
	cs := strings.TrimSuffix(filepath.Base(state), ".json")
	store := session.NewStore(session.Dir(f.common))
	for _, rec := range []session.Record{
		{ID: "11111111-1111-4111-8111-111111111111", Change: cs, Tab: session.Tab{Session: "session:s1", Window: "window:w1"}},
		{ID: "22222222-2222-4222-8222-222222222222", Change: cs, Tab: session.Tab{Session: "session:s1", Window: "window:gone"}},
		{ID: "33333333-3333-4333-8333-333333333333", Change: "OTHER", Tab: session.Tab{Window: "window:other"}},
	} {
		if err := store.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	x := &shipExec{answers: map[string]string{
		"git worktree remove " + f.worktree:                    "",
		"rex window close -s session:s1 window:w1":             "",
		"rex window close -s session:s1 window:gone":           `!Error: window target "window:gone" not found`,
		"git rev-parse --verify --quiet refs/heads/cs/ship-it": "!",
		"git push origin --delete cs/ship-it":                  "",
		"git symbolic-ref --short HEAD":                        "main\n",
		"git status --porcelain --untracked-files=no":          "",
		"git pull --ff-only origin main":                       "",
	}}
	f.ctx.Exec = x.run
	resp, err := callOp(t, f.ctx, "ship.close-out", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipCloseOutResponse)
	want := []ShipStep{{Step: "remove-worktree", Result: "done"},
		{Step: "close-tab", Result: "done", Tab: "window:w1"}, {Step: "close-tab", Result: "already", Tab: "window:gone"},
		{Step: "delete-branch", Result: "already"}, {Step: "delete-remote-branch", Result: "done"}, {Step: "update-main", Result: "done"}}
	if !r.MainUpdated || !slices.Equal(r.Steps, want) {
		t.Errorf("close-out steps = %+v", r.Steps)
	}
	if i := slices.Index(x.calls, "rex window close -s session:s1 window:w1"); i != 1 {
		t.Errorf("tab closed at call %d, want right after the worktree is removed: %q", i, x.calls)
	}
	got := human(t, r)
	for _, s := range []string{"Closed Rex tab window:w1.", "Rex tab window:gone was already closed."} {
		if !strings.Contains(got, s) {
			t.Errorf("human = %q, want %q", got, s)
		}
	}

	// Rex failing for another reason fails the step.
	x.answers["rex window close -s session:s1 window:w1"] = "!Error: cannot connect to the Rex server"
	_, err = callOp(t, f.ctx, "ship.close-out", `{}`)
	wantExit(t, err, 1)
	if !strings.Contains(err.Error(), "cannot connect") {
		t.Errorf("err = %v", err)
	}
}
