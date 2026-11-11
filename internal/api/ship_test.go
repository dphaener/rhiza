package api

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
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
	// body is the pull request body pr view answers carry when they hold
	// none of their own: the change set's 'Change set' line.
	body string
}

func (f *fakeGH) run(name string, args []string, dir string, _ []string) ([]byte, error) {
	k := strings.Join(args, " ")
	f.calls = append(f.calls, name+" "+k)
	f.dirs = append(f.dirs, dir)
	if name == "git" {
		k = "git " + k
	} else if name != "gh" {
		return nil, errors.New("not gh")
	}
	a, ok := f.answers[k]
	if !ok {
		return nil, errors.New("gh " + k + ": unexpected call")
	}
	if msg, failed := strings.CutPrefix(a, "!"); failed {
		return nil, errors.New("gh " + k + ": " + msg)
	}
	if f.body != "" && strings.HasPrefix(k, "pr view cs/") && strings.HasPrefix(a, "{") && !strings.Contains(a, `"body"`) {
		a = `{"body":` + strconv.Quote(f.body) + `,` + a[1:]
	}
	return []byte(a), nil
}

const (
	shipPRView   = "pr view cs/ship-it --json number,url,state,isDraft,headRefOid,body,mergeable,statusCheckRollup"
	shipMergeArg = "pr merge 7 --squash --match-head-commit abc1234def"
	shipDraftPR  = "pr view 7 --json isDraft,state"
	shipMergedPR = "pr view 7 --json mergeCommit,state"

	// What ship merge asks git before it merges, and ship check --main
	// asks to compare trees.
	gitFetch      = "git fetch -q origin"
	gitOriginMain = "git rev-parse --verify --quiet refs/remotes/origin/main"
	gitHoldsMain  = "git merge-base --is-ancestor refs/remotes/origin/main abc1234def"
	gitMergeTree  = "git rev-parse --verify --quiet fff0000aaa^{tree}"
	gitHeadTree   = "git rev-parse --verify --quiet abc1234def^{tree}"
)

// mainHeld makes git answer that the head holds origin/main.
func mainHeld(g *fakeGH) {
	g.answers[gitFetch] = ""
	g.answers[gitOriginMain] = "0123456789abcdef"
	g.answers[gitHoldsMain] = ""
	g.answers[shipDraftPR] = `{"isDraft":false,"state":"OPEN"}`
}

// shipFixture is a session fixture with the change set "Ship it" current
// and a fake gh.
func shipFixture(t *testing.T) (*sessionFixture, *fakeGH, string) {
	t.Helper()
	f := newSessionFixture(t, "ship-it")
	cs := newChange(t, f.ctx, `{"intent":"Ship it","title":"Ship it"}`).ChangeSet
	g := &fakeGH{answers: map[string]string{}}
	g.body = "# Plan\n\nChange set " + cs.String() + ": Ship it. Workflow plan.\n"
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
	// The red head is recorded, with no green head or checks.
	st, err := readShipState(state)
	if err != nil {
		t.Fatal(err)
	}
	if st.RedHead != "abc1234def" || len(st.RedChecks) != 3 || st.Head != "" || len(st.Checks) != 0 {
		t.Errorf("state: %+v", st)
	}
}

const infraRollup = `{"name":"test","status":"COMPLETED","conclusion":"FAILURE","detailsUrl":"https://github.com/o/r/actions/runs/42/job/9"},{"context":"lint","state":"SUCCESS"}`

// --settled --changed does not settle on the red head the failing check
// recorded, names the infra failure, and settles on a new head or on a
// newer attempt of a check run once it finishes; with no baseline it is
// --settled.
func TestShipCheckSettledChanged(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(infraRollup)
	g.answers["run rerun 42 --failed"] = ""
	// No baseline: as --settled (the head's one rerun is spent first).
	if r := shipCheck(t, f, `{"settled":true,"changed":true}`); r.Passed || len(r.Reran) != 1 {
		t.Fatalf("no baseline, first: %+v", r)
	}
	if r := shipCheck(t, f, `{"settled":true,"changed":true}`); !r.Passed || !r.Settled {
		t.Fatalf("no baseline: %+v", r)
	}
	if r := shipCheck(t, f, `{}`); r.Passed || r.Written {
		t.Fatalf("failing check: %+v", r)
	}
	st, err := readShipState(state)
	if err != nil {
		t.Fatal(err)
	}
	want := []ShipRedCheck{{Name: "test", ID: "9"}, {Name: "lint", ID: "lint"}}
	if st.RedHead != "abc1234def" || !slices.Equal(st.RedChecks, want) || st.Head != "" || st.Checks != nil {
		t.Fatalf("state: %+v", st)
	}
	r := shipCheck(t, f, `{"settled":true,"changed":true}`)
	if r.Passed || r.Settled || r.ExitCode() != 1 || !strings.Contains(r.Detail, "abc1234") || !strings.Contains(r.Detail, "test") || !strings.Contains(r.Detail, "no failing test") {
		t.Errorf("unchanged red head: %+v", r)
	}
	// A rerun of the check run (a new job ID) that is still going.
	g.answers[shipPRView] = openPR(`{"name":"test","status":"IN_PROGRESS","detailsUrl":"https://github.com/o/r/actions/runs/42/job/10"},{"context":"lint","state":"SUCCESS"}`)
	if r := shipCheck(t, f, `{"settled":true,"changed":true}`); r.Passed {
		t.Errorf("rerun running: %+v", r)
	}
	// Finished: settled, green or not.
	g.answers[shipPRView] = openPR(`{"name":"test","status":"COMPLETED","conclusion":"SUCCESS","detailsUrl":"https://github.com/o/r/actions/runs/42/job/10"},{"context":"lint","state":"SUCCESS"}`)
	if r := shipCheck(t, f, `{"settled":true,"changed":true}`); !r.Passed || !r.Settled {
		t.Errorf("rerun finished: %+v", r)
	}
	// A newer attempt, same job ID.
	g.answers[shipPRView] = openPR(`{"name":"test","status":"COMPLETED","conclusion":"FAILURE","detailsUrl":"https://github.com/o/r/actions/runs/42/attempts/2/job/9"},{"context":"lint","state":"SUCCESS"}`)
	if r := shipCheck(t, f, `{"settled":true,"changed":true}`); !r.Passed {
		t.Errorf("new attempt: %+v", r)
	}
	// A new head.
	g.answers[shipPRView] = strings.Replace(openPR(greenRollup), "abc1234def", "def5678abc", 1)
	if r := shipCheck(t, f, `{"settled":true,"changed":true}`); !r.Passed {
		t.Errorf("new head: %+v", r)
	}
	// --changed takes --settled.
	if _, err := callOp(t, f.ctx, "ship.check", `{"changed":true}`); err == nil {
		t.Error("--changed without --settled: no error")
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

	// Main moved after the checks passed: nothing is merged.
	g.answers[gitFetch] = ""
	g.answers[gitOriginMain] = "0123456789abcdef"
	g.answers[gitHoldsMain] = "!exit status 1"
	g.answers[shipDraftPR] = `{"isDraft":false,"state":"OPEN"}`
	resp, err := callOp(t, f.ctx, "ship.merge", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if m := resp.(ShipMergeResponse); !m.Moved || m.ExitCode() != 1 || m.MergeCommit != "" {
		t.Errorf("moved merge = %+v", m)
	}
	if slices.Contains(g.calls, "gh "+shipMergeArg) {
		t.Errorf("merged although main moved: %q", g.calls)
	}
	if st, _ := readShipState(state); st.MergeCommit != "" {
		t.Errorf("state = %+v", st)
	}

	mainHeld(g)
	g.answers[shipMergeArg] = ""
	g.answers[shipMergedPR] = `{"mergeCommit":{"oid":"fff0000aaa"},"state":"MERGED"}`
	resp, err = callOp(t, f.ctx, "ship.merge", `{}`)
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

	// Main is confirmed by the trees alone, never by asking gh for runs,
	// unless they differ: then CI is dispatched on the merge commit, and
	// with no workflow to dispatch the step stops.
	g.answers[gitMergeTree] = "tree1111"
	g.answers[gitHeadTree] = "tree2222"
	g.answers["git ls-tree --name-only refs/remotes/origin/main .github/workflows/"] = ""
	r = shipCheck(t, f, `{"main":true,"settled":true}`)
	if r.Passed || r.ExitCode() != 75 || !strings.Contains(r.Stop, "fff0000") || !strings.Contains(r.Stop, "no dispatchable workflow was found") {
		t.Errorf("trees differ: %+v", r)
	}
	if got := human(t, r); !strings.HasPrefix(got, "stop: merge commit fff0000: ") {
		t.Errorf("human = %q", got)
	}
	g.answers[gitHeadTree] = "tree1111"
	r = shipCheck(t, f, `{"main":true}`)
	if !r.Passed || r.Commit != "fff0000aaa" || r.ExitCode() != 0 {
		t.Errorf("main green: %+v", r)
	}
	if got := human(t, r); got != "main at fff0000: the merge commit's tree is the tested head's\n" {
		t.Errorf("human = %q", got)
	}
	for _, c := range g.calls {
		if strings.HasPrefix(c, "gh run ") {
			t.Errorf("gh call %q for main", c)
		}
	}
}

func TestShipMergeNotMerged(t *testing.T) {
	f, g, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Branch: "cs/ship-it", Head: "abc1234def"}); err != nil {
		t.Fatal(err)
	}
	mainHeld(g)
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

func TestShipCloseOutForget(t *testing.T) {
	f, _, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}); err != nil {
		t.Fatal(err)
	}
	x := &shipExec{answers: map[string]string{
		"git worktree remove " + f.worktree:                    "",
		"git rev-parse --verify --quiet refs/heads/cs/ship-it": "!",
		"git push origin --delete cs/ship-it":                  "",
		"git symbolic-ref --short HEAD":                        "main\n",
		"git status --porcelain --untracked-files=no":          "",
		"git pull --ff-only origin main":                       "",
	}}
	f.ctx.Exec = x.run
	if _, err := callOp(t, f.ctx, "ship.close-out", `{"forget":true}`); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("state file kept: %v", err)
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
	g.answers["git fetch -q origin cs/ship-it"] = ""
	g.answers["git merge-base --is-ancestor 0ld0000000 abc1234def"] = ""
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

// A head someone else pushed on top of the run's is followed; one that
// drops the run's commits stays unsettled and is named.
func TestShipCheckFollowsMovedHead(t *testing.T) {
	f, g, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Head: "0ld0000000", Pushed: "0ld0000000"}); err != nil {
		t.Fatal(err)
	}
	g.answers["git fetch -q origin cs/ship-it"] = ""
	g.answers[shipPRView] = openPR(greenRollup)

	if r := shipCheck(t, f, `{"settled":true}`); r.Passed || r.Settled || !strings.Contains(r.Detail, "drops the pushed 0ld0000") || !strings.Contains(r.Detail, "abc1234") {
		t.Errorf("dropped head: %+v", r)
	}
	g.answers["git merge-base --is-ancestor 0ld0000000 abc1234def"] = ""
	if r := shipCheck(t, f, `{"settled":true}`); !r.Passed || !r.Settled || r.Written {
		t.Errorf("settled on a followed head: %+v", r)
	}
	if r := shipCheck(t, f, `{}`); !r.Passed || !r.Written {
		t.Errorf("followed head: %+v", r)
	}
	st, _ := readShipState(state)
	if st.Head != "abc1234def" || st.Pushed != "abc1234def" {
		t.Errorf("state = %+v", st)
	}
	mainHeld(g)
	g.answers[shipMergeArg] = ""
	g.answers[shipMergedPR] = `{"state":"MERGED","mergeCommit":{"oid":"fff0000aaa"}}`
	resp, err := callOp(t, f.ctx, "ship.merge", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if m := resp.(ShipMergeResponse); m.Head != "abc1234def" {
		t.Errorf("merge = %+v", m)
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

func TestGoModule(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"module github.com/dphaener/rhiza\n\ngo 1.25\n", rhizaModule},
		{"// a comment\nmodule github.com/dphaener/rhiza // Rhiza\n", rhizaModule},
		{"module \"github.com/dphaener/rhiza\"\n", rhizaModule},
		{"\tmodule   github.com/dphaener/rhiza  \r\n", rhizaModule},
		{"module example.com/other\n", "example.com/other"},
		{"module github.com/dphaener/rhiza/v2\n", "github.com/dphaener/rhiza/v2"},
		{"modules x\ngo 1.25\n", ""},
		{"// module github.com/dphaener/rhiza\n", ""},
		{"", ""},
	} {
		if got := goModule([]byte(c.in)); got != c.want {
			t.Errorf("goModule(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestIsRhizaRepo(t *testing.T) {
	for _, c := range []struct {
		mod  string // "" writes no go.mod
		want bool
	}{
		{"", false},
		{"module github.com/dphaener/rhiza\n", true},
		{"module example.com/other\n", false},
		{"module github.com/dphaener/rhiza/v2\n", false},
	} {
		dir := t.TempDir()
		if c.mod != "" {
			if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(c.mod), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		got, err := isRhizaRepo(dir)
		if err != nil || got != c.want {
			t.Errorf("isRhizaRepo(%q) = %v, %v; want %v", c.mod, got, err, c.want)
		}
	}
}

// rhizaRepoFixture is a ship land fixture whose root is Rhiza's own repo,
// with programs answered by a shipExec, and the head build's path.
func rhizaRepoFixture(t *testing.T) (*sessionFixture, *shipExec, string) {
	t.Helper()
	f, _, _ := shipLandFixture(t)
	if err := os.WriteFile(filepath.Join(f.ctx.Root, "go.mod"), []byte("module "+rhizaModule+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	x := &shipExec{answers: map[string]string{}}
	f.ctx.Exec = x.run
	bin := filepath.Join(f.common, "rhiza", "ship", currentOf(t, f.ctx).String(), "rhiza")
	return f, x, bin
}

func TestShipLandWithHeadBuild(t *testing.T) {
	f, x, bin := rhizaRepoFixture(t)
	cs := currentOf(t, f.ctx)
	f.ctx.ActorRef = shipTestActor
	x.answers["go build -o "+bin+" ./cmd/rhiza"] = ""
	x.answers[bin+" land --change "+cs.String()+" --json --actor "+shipTestActor] = `{}`
	x.answers[bin+" render --json --actor "+shipTestActor] = `{"written":[".rhiza/rendered/a.md"],"removed":[".rhiza/rendered/stale.md"],"unchanged":3}`
	resp, err := callOp(t, f.ctx, "ship.land", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipLandResponse)
	if !r.Landed || !r.EvidenceWritten || r.Commit != "beef000111" || r.Pushed != "beef000111" ||
		!slices.Equal(r.Rendered, []string{".rhiza/rendered/a.md", ".rhiza/rendered/stale.md"}) {
		t.Fatalf("land = %+v", r)
	}
	if len(x.calls) != 3 || !strings.HasPrefix(x.calls[0], "go build") {
		t.Fatalf("calls = %q", x.calls)
	}
	for _, d := range x.dirs {
		if d != f.ctx.Root {
			t.Errorf("ran in %s, not %s", d, f.ctx.Root)
		}
	}
	if fi, err := os.Stat(filepath.Dir(bin)); err != nil || !fi.IsDir() {
		t.Errorf("head build folder: %v", err)
	}
	// The head build landed it; in-process nothing did.
	if c, _ := ctxWorkspace(t, f.ctx).changeSet(cs); c.State != changeset.StateAccepted {
		t.Errorf("change set is %s; only the head build should land it", c.State)
	}
}

func TestShipLandHeadBuildFails(t *testing.T) {
	f, x, bin := rhizaRepoFixture(t)
	cs := currentOf(t, f.ctx)
	x.answers["go build -o "+bin+" ./cmd/rhiza"] = "!cmd/rhiza/main.go:1:1: syntax error"
	_, err := callOp(t, f.ctx, "ship.land", `{}`)
	wantExit(t, err, 1)
	if !strings.Contains(err.Error(), "could not build rhiza from the PR head") || !strings.Contains(err.Error(), "syntax error") {
		t.Errorf("err = %v", err)
	}
	if len(x.calls) != 1 {
		t.Errorf("calls = %q", x.calls)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "commit") || strings.HasPrefix(c, "push") || strings.HasPrefix(c, "add") {
			t.Errorf("git %q after a failed build", c)
		}
	}
	w := ctxWorkspace(t, f.ctx)
	if c, _ := w.changeSet(cs); c.State != changeset.StateAccepted {
		t.Errorf("change set is %s", c.State)
	}
	if _, err := w.model(id.ID{}).Resolve("pr-7-ship-it"); err == nil {
		t.Error("evidence written after a failed build")
	}
}

func TestShipRenderInProcess(t *testing.T) {
	f, _, _ := shipLandFixture(t)
	x := &shipExec{answers: map[string]string{}}
	f.ctx.Exec = x.run
	t2, err := findShipTarget(f.ctx)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := shipRender(f.ctx, t2)
	if err != nil {
		t.Fatal(err)
	}
	if len(rendered) == 0 || len(x.calls) != 0 {
		t.Errorf("rendered %q with calls %q", rendered, x.calls)
	}
}

func TestShipActorFlags(t *testing.T) {
	a, err := id.Parse("01M5PWGEF85Y185JZ2ZGCNH40N")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		ctx  Ctx
		want []string
	}{
		{Ctx{}, nil},
		{Ctx{ActorRef: "me", ForRef: "owner"}, []string{"--actor", "me", "--for", "owner"}},
		{Ctx{Actor: a}, []string{"--actor", a.String()}},
		{Ctx{For: a}, []string{"--for", a.String()}},
	} {
		if got := shipActorFlags(&c.ctx); !slices.Equal(got, c.want) {
			t.Errorf("shipActorFlags(%+v) = %q, want %q", c.ctx, got, c.want)
		}
	}
}

// The stale-render test hook applies only with RHIZA_TEST_NOW set too.
func TestRenderStaleHookInert(t *testing.T) {
	env := func(m map[string]string) func(string) string { return func(k string) string { return m[k] } }
	for _, c := range []struct {
		getenv func(string) string
		want   string
	}{
		{nil, ""},
		{env(nil), ""},
		{env(map[string]string{"RHIZA_TEST_RENDER_STALE": "x"}), ""},
		{env(map[string]string{"RHIZA_TEST_NOW": "2026-10-24T12:00:00Z"}), ""},
		{env(map[string]string{"RHIZA_TEST_NOW": "2026-10-24T12:00:00Z", "RHIZA_TEST_RENDER_STALE": "x"}), "x"},
	} {
		if got := testRenderStale(&Ctx{Getenv: c.getenv}); got != c.want {
			t.Errorf("testRenderStale = %q, want %q", got, c.want)
		}
	}
}

func TestSplitConflicts(t *testing.T) {
	cases := []struct {
		name            string
		rendered        string
		paths           []string
		inRendered, oth []string
	}{
		{"None", ".rhiza/rendered/", nil, nil, nil},
		{"RenderedOnly", ".rhiza/rendered/", []string{".rhiza/rendered/index.md", ".rhiza/rendered/capability/x.md"},
			[]string{".rhiza/rendered/index.md", ".rhiza/rendered/capability/x.md"}, nil},
		{"CodeOnly", ".rhiza/rendered/", []string{"main.go", "go.mod"}, nil, []string{"main.go", "go.mod"}},
		{"Mixed", ".rhiza/rendered/", []string{".rhiza/rendered/index.md", "feature.txt"},
			[]string{".rhiza/rendered/index.md"}, []string{"feature.txt"}},
		// The model and the folder beside rendered are not rendered.
		{"LogAndLookalike", ".rhiza/rendered/", []string{".rhiza/log/a/b.json", ".rhiza/rendered.md", ".rhiza/rendered-old/x.md", ".rhiza/repo.json"},
			nil, []string{".rhiza/log/a/b.json", ".rhiza/rendered.md", ".rhiza/rendered-old/x.md", ".rhiza/repo.json"}},
		// A project in a subfolder: its rendered files sit under the prefix.
		{"Prefix", "sub/.rhiza/rendered/", []string{"sub/.rhiza/rendered/index.md", ".rhiza/rendered/index.md"},
			[]string{"sub/.rhiza/rendered/index.md"}, []string{".rhiza/rendered/index.md"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in, other := splitConflicts(c.rendered, c.paths)
			if !slices.Equal(in, c.inRendered) || !slices.Equal(other, c.oth) {
				t.Errorf("splitConflicts = %v, %v; want %v, %v", in, other, c.inRendered, c.oth)
			}
		})
	}
}

func TestGitLines(t *testing.T) {
	if got := gitLines("a.go\n\n b.md \n"); !slices.Equal(got, []string{"a.go", "b.md"}) {
		t.Errorf("gitLines = %q", got)
	}
	if got := gitLines(""); len(got) != 0 {
		t.Errorf("gitLines(\"\") = %q", got)
	}
}

// Git answers for a ship rebase on cs/ship-it whose origin main is at
// main000 and whose origin branch is at abc1234def.
const (
	rebaseOnto     = "main000"
	rebaseOrigin   = "abc1234def"
	rebaseConflict = "diff --name-only --diff-filter=U"
	rebasePush     = "push origin HEAD:refs/heads/cs/ship-it"
	renderedSpec   = ":(top,literal).rhiza/rendered/"
)

// shipRebaseFixture is shipLandFixture with git answering a branch that
// origin's main has moved past; the merge itself stops (no answer).
func shipRebaseFixture(t *testing.T) (*sessionFixture, string) {
	t.Helper()
	f, state, _ := shipLandFixture(t)
	for k, v := range map[string]string{
		"fetch -q origin": "",
		"rev-parse --verify --quiet refs/remotes/origin/main":       rebaseOnto + "\n",
		"rev-parse --verify --quiet refs/remotes/origin/cs/ship-it": rebaseOrigin + "\n",
		"merge-base --is-ancestor " + rebaseOrigin + " HEAD":        "",
		"rev-parse HEAD":                               "beef000111\n",
		"merge --abort":                                "",
		"add -A -- " + renderedSpec:                    "",
		"diff --cached --name-only -- " + renderedSpec: "",
		"rm -q -f --ignore-unmatch -- :(top,literal).rhiza/rendered/index.md": "",
	} {
		f.git.out[k] = v
	}
	return f, state
}

const rebaseMerge = "merge --no-edit -m Merge origin/main into cs/ship-it " + rebaseOnto

func shipRebase(t *testing.T, f *sessionFixture) ShipRebaseResponse {
	t.Helper()
	resp, err := callOp(t, f.ctx, "ship.rebase", `{}`)
	if err != nil {
		t.Fatalf("ship.rebase: %v", err)
	}
	return resp.(ShipRebaseResponse)
}

func TestShipRebaseNotBehind(t *testing.T) {
	f, state := shipRebaseFixture(t)
	f.git.out["merge-base --is-ancestor "+rebaseOnto+" HEAD"] = ""
	r := shipRebase(t, f)
	if r.Behind || r.Rebased || r.Pushed != "" || r.Onto != rebaseOnto || r.ExitCode() != 0 {
		t.Fatalf("rebase = %+v", r)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "merge ") || strings.HasPrefix(c, "push") {
			t.Errorf("git call %q when not behind", c)
		}
	}
	if st, _ := readShipState(state); st.Head != "abc1234def" || len(st.Checks) != 1 || st.Pushed != "" {
		t.Errorf("state changed: %+v", st)
	}
	if got := human(t, r); got != "cs/ship-it already holds origin/main at main000; nothing to merge.\n" {
		t.Errorf("human = %q", got)
	}
}

func TestShipRebaseClean(t *testing.T) {
	f, state := shipRebaseFixture(t)
	f.git.out[rebaseMerge] = ""
	f.git.out[rebasePush] = ""
	r := shipRebase(t, f)
	if !r.Behind || !r.Rebased || len(r.Settled) != 0 || r.Pushed != "beef000111" || r.Commit != "" || r.ExitCode() != 0 {
		t.Fatalf("rebase = %+v", r)
	}
	// The next ship check waits for CI on the pushed head.
	if st, _ := readShipState(state); st.Pushed != "beef000111" || st.Head != "" || st.Checks != nil || st.PR != 7 {
		t.Errorf("state = %+v", st)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "push") && c != rebasePush {
			t.Errorf("push %q; want only %q", c, rebasePush)
		}
		if strings.Contains(c, "force") || strings.Contains(c, "lease") || strings.HasPrefix(c, "rebase") {
			t.Errorf("git call %q forces or rebases", c)
		}
		if strings.Contains(c, "main:") || strings.Contains(c, ":refs/heads/main") {
			t.Errorf("git call %q reaches main", c)
		}
	}
}

func TestShipRebaseRenderedConflict(t *testing.T) {
	f, state := shipRebaseFixture(t)
	f.git.out[rebaseConflict] = ".rhiza/rendered/index.md\n"
	f.git.out["diff --cached --name-only"] = ".rhiza/rendered/index.md\n"
	f.git.out["commit -q --no-edit"] = ""
	f.git.out[rebasePush] = ""
	r := shipRebase(t, f)
	if !r.Rebased || !slices.Equal(r.Settled, []string{".rhiza/rendered/index.md"}) || len(r.Conflicts) != 0 || r.Pushed != "beef000111" {
		t.Fatalf("rebase = %+v", r)
	}
	// The conflict was settled in order: drop the file, render, commit the
	// merge, push.
	var order []string
	for _, c := range f.git.calls {
		switch {
		case strings.HasPrefix(c, "rm "), c == "commit -q --no-edit", c == rebasePush:
			order = append(order, strings.Fields(c)[0])
		}
	}
	if !slices.Equal(order, []string{"rm", "commit", "push"}) {
		t.Errorf("git order = %v", order)
	}
	if slices.Contains(f.git.calls, "merge --abort") {
		t.Error("a rendered-only conflict aborted the merge")
	}
	if st, _ := readShipState(state); st.Pushed != "beef000111" || st.Head != "" {
		t.Errorf("state = %+v", st)
	}
	if got := human(t, r); !strings.Contains(got, "Settled 1 rendered file by rendering again.") || !strings.Contains(got, "Pushed beef000.") {
		t.Errorf("human = %q", got)
	}
}

func TestShipRebaseCodeConflict(t *testing.T) {
	f, state := shipRebaseFixture(t)
	cs := currentOf(t, f.ctx)
	f.git.out[rebaseConflict] = ".rhiza/rendered/index.md\nfeature.txt\n"
	f.git.out["merge-base HEAD "+rebaseOnto] = "base000\n"
	f.git.out["log --format=%H%x09%s%x09%(trailers:key=Change,valueonly,separator=%x2C) base000.."+rebaseOnto] =
		"aaa1111\tAdd other\tCS_OTHER\nbbb2222\tTouch up\t\n"
	r := shipRebase(t, f)
	if r.Rebased || !slices.Equal(r.Conflicts, []string{".rhiza/rendered/index.md", "feature.txt"}) || r.ExitCode() != 1 ||
		!strings.Contains(r.Detail, "feature.txt") || r.Pushed != "" {
		t.Fatalf("rebase = %+v", r)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "push") || strings.HasPrefix(c, "rm ") || strings.HasPrefix(c, "merge --abort") || strings.HasPrefix(c, "rebase") {
			t.Errorf("git call %q after a code conflict", c)
		}
	}
	// No question was written.
	if _, err := ctxWorkspace(t, f.ctx).model(cs).Resolve("rebase-conflict-ship-it"); err == nil {
		t.Error("a question was raised")
	}
	// The context is in the ship state.
	st, _ := readShipState(state)
	mc := st.MergeConflict
	if mc == nil || mc.Change != cs.String() || mc.Branch != "cs/ship-it" || mc.Onto != rebaseOnto ||
		!slices.Equal(mc.Files, r.Conflicts) || len(mc.MainCommits) != 2 ||
		mc.MainCommits[0].Commit != "aaa1111" || !slices.Equal(mc.MainCommits[0].Changes, []string{"CS_OTHER"}) ||
		len(mc.MainCommits[1].Changes) != 0 || st.Head != "abc1234def" || st.Pushed != "" {
		t.Errorf("state = %+v", st)
	}
	if got := human(t, r); !strings.Contains(got, "- feature.txt\n") || !strings.Contains(got, "left in progress") {
		t.Errorf("human = %q", got)
	}

	// Again, with the merge in progress: the conflicts are reported from
	// the state, nothing is merged.
	f.git.out["rev-parse --verify --quiet MERGE_HEAD"] = "cafe000\n"
	n := len(f.git.calls)
	r = shipRebase(t, f)
	if len(r.Conflicts) != 2 || r.MergeConflict == nil || r.ExitCode() != 1 {
		t.Errorf("again = %+v", r)
	}
	for _, c := range f.git.calls[n:] {
		if strings.HasPrefix(c, "merge ") || strings.HasPrefix(c, "push") {
			t.Errorf("git call %q with a merge in progress", c)
		}
	}
}

// ship check --conflicts passes only for a recorded unresolved merge with
// fewer than 3 conflict fixes (merge commits with the run's Review:
// trailer) on the branch.
func TestShipCheckConflicts(t *testing.T) {
	f, _ := shipRebaseFixture(t)
	fixes := ""
	inner := f.ctx.Exec
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name == "git" && strings.Join(args, " ") == "log --merges --format=%H -E --grep=^Review: RUN1$ HEAD" {
			return []byte(fixes), nil
		}
		return inner(name, args, dir, env)
	}
	f.ctx.Getenv = func(k string) string {
		if k == "RHIZA_RUN" {
			return "RUN1"
		}
		return ""
	}

	// No recorded conflict: fails.
	r := shipCheck(t, f, `{"conflicts":true}`)
	if r.Passed || r.ExitCode() != 1 || !strings.Contains(r.Detail, "no merge conflict is recorded") {
		t.Errorf("no conflict: %+v", r)
	}

	// Record an unresolved conflict.
	f.git.out[rebaseConflict] = "feature.txt\n"
	f.git.out["merge-base HEAD "+rebaseOnto] = "base000\n"
	if r := shipRebase(t, f); len(r.Conflicts) != 1 {
		t.Fatalf("rebase = %+v", r)
	}

	// No fixes yet, and two: passes.
	for _, out := range []string{"", "a1\nb2\n"} {
		fixes = out
		r = shipCheck(t, f, `{"conflicts":true}`)
		if !r.Passed || r.ExitCode() != 0 || !slices.Equal(r.Conflicts, []string{"feature.txt"}) {
			t.Errorf("fixes %q: %+v", out, r)
		}
	}

	// Three fixes: fails, naming why.
	fixes = "a1\nb2\nc3\n"
	r = shipCheck(t, f, `{"conflicts":true}`)
	if r.Passed || r.ExitCode() != 1 || !strings.Contains(r.Detail, "3 conflict fixes") {
		t.Errorf("three fixes: %+v", r)
	}

	// Once the merge is resolved it fails too.
	fixes = ""
	inner2 := f.ctx.Exec
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name == "git" && strings.Join(args, " ") == "merge-base --is-ancestor "+rebaseOnto+" HEAD" {
			return nil, nil
		}
		return inner2(name, args, dir, env)
	}
	r = shipCheck(t, f, `{"conflicts":true}`)
	if r.Passed || !strings.Contains(r.Detail, "is resolved") {
		t.Errorf("resolved: %+v", r)
	}
}

// A fixer's merge, committed on top of origin's branch and holding main,
// is rendered, pushed plainly and recorded.
func TestShipRebaseFixerMerge(t *testing.T) {
	f, state := shipRebaseFixture(t)
	f.git.out["merge-base --is-ancestor "+rebaseOnto+" HEAD"] = ""
	f.git.out["rev-list --merges --count "+rebaseOrigin+"..HEAD"] = "1\n"
	f.git.out[rebasePush] = ""
	r := shipRebase(t, f)
	if !r.Behind || !r.Rebased || r.Pushed != "beef000111" || r.ExitCode() != 0 {
		t.Fatalf("rebase = %+v", r)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "merge ") {
			t.Errorf("git call %q merging again", c)
		}
	}
	if st, _ := readShipState(state); st.Pushed != "beef000111" || st.Head != "" || st.MergeConflict != nil {
		t.Errorf("state = %+v", st)
	}
}

func TestShipRebasePushFails(t *testing.T) {
	f, state := shipRebaseFixture(t)
	f.git.out[rebaseMerge] = ""
	f.git.out["reset -q --keep beef000111"] = ""
	if _, err := callOp(t, f.ctx, "ship.rebase", `{}`); err == nil {
		t.Fatal("ship.rebase passed with a failed push")
	}
	if !slices.Contains(f.git.calls, "reset -q --keep beef000111") {
		t.Errorf("the branch was not put back; calls %q", f.git.calls)
	}
	if st, _ := readShipState(state); st.Pushed != "" || st.Head != "abc1234def" {
		t.Errorf("state changed: %+v", st)
	}
}

func TestShipRebaseRefusals(t *testing.T) {
	f, _ := shipRebaseFixture(t)
	f.git.out["symbolic-ref --short HEAD"] = "main\n"
	_, err := callOp(t, f.ctx, "ship.rebase", `{}`)
	wantExit(t, err, 3)
	if !strings.Contains(err.Error(), "on main") {
		t.Errorf("err = %v", err)
	}
	// origin's branch has a commit the local branch lacks: no merge, no
	// push.
	f, _ = shipRebaseFixture(t)
	delete(f.git.out, "merge-base --is-ancestor "+rebaseOrigin+" HEAD")
	_, err = callOp(t, f.ctx, "ship.rebase", `{}`)
	wantExit(t, err, 3)
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "merge ") || strings.HasPrefix(c, "push") {
			t.Errorf("git call %q with origin's branch ahead", c)
		}
	}
	f.ctx.DryRun = true
	_, err = callOp(t, f.ctx, "ship.rebase", `{}`)
	wantExit(t, err, 2)
}

// shipFollowFixture is a rebase fixture whose origin branch holds more
// than HEAD: origin's branch is not an ancestor of HEAD, and HEAD is an
// ancestor of it.
func shipFollowFixture(t *testing.T) *sessionFixture {
	t.Helper()
	f, _ := shipRebaseFixture(t)
	delete(f.git.out, "merge-base --is-ancestor "+rebaseOrigin+" HEAD")
	f.git.out["merge-base --is-ancestor HEAD "+rebaseOrigin] = ""
	return f
}

func TestShipFollowFastForwards(t *testing.T) {
	f := shipFollowFixture(t)
	f.git.out["status --porcelain --untracked-files=no"] = ""
	f.git.out["merge --ff-only "+rebaseOrigin] = ""
	// The fake cannot show origin's branch now held by HEAD, so the
	// rebase that follows stops; the fast-forward is what matters.
	callOp(t, f.ctx, "ship.rebase", `{}`)
	if !slices.Contains(f.git.calls, "merge --ff-only "+rebaseOrigin) {
		t.Errorf("no fast-forward; calls %q", f.git.calls)
	}
}

func TestShipFollowRefusesDirtyWorktree(t *testing.T) {
	f := shipFollowFixture(t)
	f.git.out["status --porcelain --untracked-files=no"] = " M feature.txt\n"
	_, err := callOp(t, f.ctx, "ship.rebase", `{}`)
	if err == nil || !strings.Contains(err.Error(), "feature.txt") || !strings.Contains(err.Error(), "uncommitted") {
		t.Fatalf("err = %v", err)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "merge ") || strings.HasPrefix(c, "reset") || strings.HasPrefix(c, "push") {
			t.Errorf("git call %q moved the worktree", c)
		}
	}
}

func TestShipFollowRefusesLocalCommits(t *testing.T) {
	f := shipFollowFixture(t)
	delete(f.git.out, "merge-base --is-ancestor HEAD "+rebaseOrigin)
	f.git.out["rev-list --oneline "+rebaseOrigin+"..HEAD"] = "abc1234 Local only commit\n"
	_, err := callOp(t, f.ctx, "ship.rebase", `{}`)
	if err == nil || !strings.Contains(err.Error(), "Local only commit") {
		t.Fatalf("err = %v", err)
	}
	for _, c := range f.git.calls {
		if strings.HasPrefix(c, "merge ") || strings.HasPrefix(c, "reset") || strings.HasPrefix(c, "push") {
			t.Errorf("git call %q moved the worktree", c)
		}
	}
}

// A plan run's merge leaves a merge commit in the change set's ship state.
// A later ship run's open-pr starts over from it: ship check records the
// new PR's head and ship merge merges it instead of reporting the plan's
// merge as already done.
func TestShipMergeAfterPlanMerge(t *testing.T) {
	ctx, g, _, git := openPRFixture(t, true)
	git("push", "-q", "origin", "main")
	var merges []string
	base := g.run
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		k := strings.Join(args, " ")
		switch {
		case name == "gh" && strings.HasPrefix(k, "pr merge "):
			merges = append(merges, k)
			return nil, nil
		case name == "gh" && k == "pr view 7 --json mergeCommit,state":
			return []byte(`{"mergeCommit":{"oid":"bbb2222ccc"},"state":"MERGED"}`), nil
		}
		out, err := base(name, args, dir, env)
		if name == "gh" && strings.HasPrefix(k, "pr view cs/ship-it") {
			out = []byte(strings.Replace(string(out), `"statusCheckRollup":[]`, `"statusCheckRollup":[{"name":"test","status":"COMPLETED","conclusion":"SUCCESS"}]`, 1))
		}
		return out, err
	}
	resp, err := callOp(t, ctx, "ship.open-pr", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	cs := resp.(ShipOpenPRResponse).Change
	state := filepath.Join(ctx.Root, ".git", "rhiza", "ship", cs+".json")
	if err := writeShipState(state, shipState{PR: 3, MergeCommit: "aaa1111bbb", Head: "old", Main: []ShipCheck{{Name: "x"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, ctx, "ship.open-pr", `{}`); err != nil {
		t.Fatal(err)
	}
	if st, _ := readShipState(state); st.MergeCommit != "" || st.Main != nil {
		t.Errorf("open-pr kept the plan's merge: %+v", st)
	}
	cr, err := callOp(t, ctx, "ship.check", `{}`)
	if err != nil || !cr.(ShipCheckResponse).Passed {
		t.Fatalf("check: %+v, %v", cr, err)
	}
	resp, err = callOp(t, ctx, "ship.merge", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if m := resp.(ShipMergeResponse); m.AlreadyMerged || m.MergeCommit != "bbb2222ccc" || len(merges) != 1 {
		t.Errorf("merge = %+v, merges %q", m, merges)
	}
}

func TestShipCheckNoRunIDNotRerun(t *testing.T) {
	f, g, _ := shipFixture(t)
	g.answers[shipPRView] = openPR(`{"name":"ext","status":"COMPLETED","conclusion":"FAILURE","detailsUrl":"https://ci.example/job/1"}`)
	r := shipCheck(t, f, `{"settled":true}`)
	if !r.Settled || r.Green || len(r.Reran) != 0 {
		t.Errorf("no run id: %+v", r)
	}
}

const rerunRollup = `{"name":"test","status":"COMPLETED","conclusion":"FAILURE","detailsUrl":"https://github.com/o/r/actions/runs/42/job/1"}`

// A budget of 2 reruns a head's failed run twice, recording the count
// before each gh call, then fails on the third red; a new head starts over.
func TestShipCheckRerunBudget(t *testing.T) {
	f, g, state := shipFixture(t)
	g.answers[shipPRView] = openPR(rerunRollup)
	const rerun = "gh run rerun 42 --failed"
	g.answers["run rerun 42 --failed"] = ""
	count := func() int {
		n := 0
		for _, c := range g.calls {
			if c == rerun {
				n++
			}
		}
		return n
	}
	for i := 1; i <= 2; i++ {
		r := shipCheck(t, f, `{"settled":true,"reruns":2}`)
		if r.Settled || len(r.Reran) != 1 || count() != i {
			t.Fatalf("rerun %d: %+v, gh reruns %d", i, r, count())
		}
		if st, _ := readShipState(state); st.RerunCount != i || st.RerunHead != "abc1234def" {
			t.Fatalf("state after rerun %d: %+v", i, st)
		}
	}
	r := shipCheck(t, f, `{"settled":true,"reruns":2}`)
	if !r.Settled || r.Green || len(r.Reran) != 0 || count() != 2 {
		t.Errorf("third red: %+v, gh reruns %d", r, count())
	}
	// A new head resets the count.
	g.answers[shipPRView] = strings.Replace(openPR(rerunRollup), "abc1234def", "def5678abc", 1)
	r = shipCheck(t, f, `{"settled":true,"reruns":2}`)
	if r.Settled || count() != 3 {
		t.Errorf("new head: %+v, gh reruns %d", r, count())
	}
	if st, _ := readShipState(state); st.RerunCount != 1 || st.RerunHead != "def5678abc" {
		t.Errorf("state after new head: %+v", st)
	}
}

// Without the flag a head is rerun exactly once.
func TestShipCheckRerunDefaultOnce(t *testing.T) {
	f, g, _ := shipFixture(t)
	g.answers[shipPRView] = openPR(rerunRollup)
	g.answers["run rerun 42 --failed"] = ""
	if r := shipCheck(t, f, `{"settled":true}`); r.Settled || len(r.Reran) != 1 {
		t.Fatalf("first: %+v", r)
	}
	if r := shipCheck(t, f, `{"settled":true}`); !r.Settled || r.Green || len(r.Reran) != 0 {
		t.Errorf("second: %+v", r)
	}
}

func TestShipCheckRerunsUsage(t *testing.T) {
	f, _, _ := shipFixture(t)
	for _, req := range []string{`{"reruns":0}`, `{"reruns":-1}`} {
		_, err := callOp(t, f.ctx, "ship.check", req)
		var e *Error
		if !errors.As(err, &e) || e.Code != CodeUsage {
			t.Errorf("%s: err = %v, want a usage error", req, err)
		}
	}
}

func TestShipCloseOutRecordsStaleMain(t *testing.T) {
	f, _, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Head: "abc1234def", MergeCommit: "fff0000aaa"}); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(f.worktree); err != nil {
		t.Fatal(err)
	}
	x := &shipExec{answers: map[string]string{
		"git rev-parse --verify --quiet refs/heads/cs/ship-it": "!",
		"git push origin --delete cs/ship-it":                  "!error: unable to delete 'cs/ship-it': remote ref does not exist",
		"git symbolic-ref --short HEAD":                        "main\n",
		"git status --porcelain --untracked-files=no":          " M README.md\n",
		"git pull --ff-only origin main":                       "",
	}}
	f.ctx.Exec = x.run
	items := func() []InboxItem { return inboxMainStale(f.ctx) }
	// health reads the main checkout with real git (ctx.Git); only the
	// close-out record is asserted, since the fixture's branch and status
	// are not this test's concern.
	const stale = "the last close-out could not pull"
	mainDetail := func(p HealthPart) string { return strings.Join(p.Detail, "\n") }

	if got := items(); len(got) != 0 {
		t.Fatalf("items before = %+v", got)
	}
	if h := healthMain(f.ctx); strings.Contains(h.Summary, stale) {
		t.Fatalf("health main before = %+v", h)
	}
	resp, err := callOp(t, f.ctx, "ship.close-out", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if resp.(ShipCloseOutResponse).MainUpdated {
		t.Fatal("dirty main was updated")
	}
	ms, ok := daemon.ReadMainStale(f.common)
	if !ok || !strings.Contains(ms.Detail, "has uncommitted changes") {
		t.Fatalf("record = %+v, %v", ms, ok)
	}
	got := items()
	if len(got) != 1 || got[0].Kind != InboxMainStale || !strings.Contains(got[0].Detail, "has uncommitted changes") {
		t.Fatalf("items = %+v", got)
	}
	h := healthMain(f.ctx)
	if h.State != HealthWarn || !strings.Contains(h.Summary, stale) ||
		!strings.Contains(mainDetail(h), "change set "+ms.Change) || !strings.Contains(mainDetail(h), "has uncommitted changes") {
		t.Fatalf("health main after a failed pull = %+v", h)
	}

	// A second failure keeps the item's identity.
	if _, err := callOp(t, f.ctx, "ship.close-out", `{}`); err != nil {
		t.Fatal(err)
	}
	if again := items(); len(again) != 1 || again[0].ID != got[0].ID {
		t.Errorf("items after a repeat = %+v, want %s", again, got[0].ID)
	}

	x.answers["git status --porcelain --untracked-files=no"] = ""
	resp, err = callOp(t, f.ctx, "ship.close-out", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.(ShipCloseOutResponse).MainUpdated {
		t.Fatalf("clean main not updated: %+v", resp)
	}
	if _, ok := daemon.ReadMainStale(f.common); ok {
		t.Error("record kept after a successful pull")
	}
	if left := items(); len(left) != 0 {
		t.Errorf("items after = %+v", left)
	}
	if h := healthMain(f.ctx); strings.Contains(h.Summary, stale) || strings.Contains(mainDetail(h), "has uncommitted changes") {
		t.Errorf("health main after a successful pull = %+v", h)
	}
}

func TestParseFailedTests(t *testing.T) {
	nested := `test	Run tests	2026-11-04T10:00:00Z --- FAIL: TestScripts (12.00s)
test	Run tests	2026-11-04T10:00:00Z     --- FAIL: TestScripts/review-check-worktree (1.20s)
test	Run tests	2026-11-04T10:00:00Z FAIL	github.com/dphaener/rhiza/cmd/rhiza	13.1s
`
	if got, want := parseFailedTests(nested), []string{"TestScripts/review-check-worktree"}; !slices.Equal(got, want) {
		t.Errorf("nested = %v, want %v", got, want)
	}
	two := "--- FAIL: TestA (0.00s)\n--- FAIL: TestB (0.01s)\n--- FAIL: TestA (0.00s)\n"
	if got, want := parseFailedTests(two), []string{"TestA", "TestB"}; !slices.Equal(got, want) {
		t.Errorf("two = %v, want %v", got, want)
	}
	if got := parseFailedTests("Error: setup-go failed\nprocess completed with exit code 1\n"); len(got) != 0 {
		t.Errorf("infra log = %v, want none", got)
	}
	p := "test\tRun specs\t2026-11-04T10:00:00Z "
	rspec := p + "Failed examples:\n" + p + "\n" +
		p + "rspec ./spec/views/tasks/list_spec.rb:290 # Tasks list shows rows\n" +
		p + "rspec ./spec/views/tasks/list_spec.rb:290 # Tasks list shows rows\n" +
		p + "\x1b[31mrspec ./spec/models/x_spec.rb:5\x1b[0m # X works\n" +
		p + "rspec ./spec/system/garden_location_spec.rb[1:1:11] # Garden location saves\n" +
		p + "rspec ./spec/models/y_spec.rb[1:2,1:3]\n"
	if got, want := parseFailedTests(rspec), []string{"rspec ./spec/views/tasks/list_spec.rb:290", "rspec ./spec/models/x_spec.rb:5", "rspec ./spec/system/garden_location_spec.rb[1:1:11]", "rspec ./spec/models/y_spec.rb[1:2,1:3]"}; !slices.Equal(got, want) {
		t.Errorf("rspec = %v, want %v", got, want)
	}
	mini := p + "bin/rails test test/models/x_test.rb:12\n" + p + "rails test test/models/y_test.rb:3\n"
	if got, want := parseFailedTests(mini), []string{"bin/rails test test/models/x_test.rb:12", "rails test test/models/y_test.rb:3"}; !slices.Equal(got, want) {
		t.Errorf("minitest = %v, want %v", got, want)
	}
	py := p + "FAILED tests/test_x.py::test_y - AssertionError: nope\n"
	if got, want := parseFailedTests(py), []string{"pytest tests/test_x.py::test_y"}; !slices.Equal(got, want) {
		t.Errorf("pytest = %v, want %v", got, want)
	}
	for _, none := range []string{p + "1 failure\n", p + "FAIL src/x.test.js\n", p + "Error: runner lost contact\n", p + "see rspec ./spec/x_spec.rb:1 above\n"} {
		if got := parseFailedTests(none); len(got) != 0 {
			t.Errorf("%q = %v, want none", none, got)
		}
	}
}

// --test-failures passes only with a failing test recorded for the head,
// reads the state file alone, and runs no gh.
func TestShipCheckTestFailures(t *testing.T) {
	f, g, state := shipFixture(t)
	r := shipCheck(t, f, `{"test_failures":true}`)
	if r.Passed || r.ExitCode() != 1 || r.Detail == "" {
		t.Errorf("no state: %+v", r)
	}
	if err := writeShipState(state, shipState{RerunHead: "abc1234def"}); err != nil {
		t.Fatal(err)
	}
	if r := shipCheck(t, f, `{"test_failures":true}`); r.Passed {
		t.Errorf("infra-only failure passed: %+v", r)
	}
	st := shipState{RerunHead: "abc1234def", CIFindings: []ShipFinding{{Test: "TestScripts/x", Reason: "failed"}}}
	if err := writeShipState(state, st); err != nil {
		t.Fatal(err)
	}
	r = shipCheck(t, f, `{"test_failures":true}`)
	if !r.Passed || r.ExitCode() != 0 || !slices.Equal(r.Tests, []string{"TestScripts/x"}) || len(r.Findings) != 1 {
		t.Errorf("recorded failure: %+v", r)
	}
	if len(g.calls) != 0 {
		t.Errorf("ran %v; want no gh", g.calls)
	}
	if _, err := callOp(t, f.ctx, "ship.check", `{"test_failures":true,"settled":true}`); err == nil {
		t.Error("--test-failures with --settled was accepted")
	}
}

// TestParsePRViewMultiJob feeds the rollup of one Actions run split into
// lint, script shards, rest, budgets and the aggregate check job.
func TestParsePRViewMultiJob(t *testing.T) {
	job := func(name, status, conclusion string) string {
		return `{"__typename":"CheckRun","name":"` + name + `","status":"` + status + `","conclusion":"` + conclusion +
			`","detailsUrl":"https://github.com/o/r/actions/runs/4242/job/` + name + `"}`
	}
	rollup := func(jobs ...string) []byte {
		return []byte(`{"number":9,"state":"OPEN","headRefOid":"abc","statusCheckRollup":[` + strings.Join(jobs, ",") + `]}`)
	}
	verdict := func(b []byte) (v prView, settled, green bool, failed []string) {
		t.Helper()
		v, err := parsePRView(b)
		if err != nil {
			t.Fatal(err)
		}
		settled, green, failed = checksVerdict(v.Checks)
		return
	}

	v, settled, green, failed := verdict(rollup(
		job("lint", "COMPLETED", "SUCCESS"),
		job("scripts (1/4)", "COMPLETED", "FAILURE"),
		job("scripts (2/4)", "COMPLETED", "SUCCESS"),
		job("rest", "COMPLETED", "SUCCESS"),
		job("budgets", "COMPLETED", "SUCCESS"),
		job("check", "COMPLETED", "CANCELLED"),
	))
	if !settled || green || !slices.Equal(failed, []string{"scripts (1/4)", "check"}) {
		t.Errorf("failed shard: settled=%v green=%v failed=%v", settled, green, failed)
	}
	for _, c := range v.Checks {
		if c.RunID != "4242" {
			t.Errorf("check %q RunID = %q, want 4242", c.Name, c.RunID)
		}
	}

	_, settled, green, failed = verdict(rollup(
		job("lint", "COMPLETED", "SUCCESS"),
		job("scripts (1/4)", "COMPLETED", "SUCCESS"),
		job("rest", "COMPLETED", "SUCCESS"),
		job("budgets", "COMPLETED", "SUCCESS"),
		job("check", "COMPLETED", "SUCCESS"),
	))
	if !settled || !green || len(failed) != 0 {
		t.Errorf("all passed: settled=%v green=%v failed=%v", settled, green, failed)
	}

	_, settled, green, _ = verdict(rollup(
		job("lint", "COMPLETED", "SUCCESS"),
		job("scripts (1/4)", "IN_PROGRESS", ""),
		job("check", "QUEUED", ""),
	))
	if settled || green {
		t.Errorf("in progress: settled=%v green=%v", settled, green)
	}
}

// TestShipCheckRefusesForeignPR: ship check, with and without --settled,
// refuses a pull request on the branch, open, merged or closed, whose body
// names another change set or none, naming it, and writes no ship state.
func TestShipCheckRefusesForeignPR(t *testing.T) {
	bodies := map[string]string{
		"another change set": "Change set 01M6T1ZR0000000000000000AB: Other.",
		"no change set line": "Fixes a bug.",
		"empty body":         "",
	}
	for name, body := range bodies {
		for _, state := range []string{"OPEN", "MERGED", "CLOSED"} {
			for _, req := range []string{`{}`, `{"settled":true}`} {
				f, g, statePath := shipFixture(t)
				g.answers[shipPRView] = `{"body":` + strconv.Quote(body) + `,"number":345,"url":"https://github.com/o/r/pull/345","state":"` + state + `","headRefOid":"abc1234def","statusCheckRollup":[` + greenRollup + `]}`
				_, err := callOp(t, f.ctx, "ship.check", req)
				if err == nil || !strings.Contains(err.Error(), "#345") {
					t.Errorf("%s, %s, %s: err = %v, want a conflict naming #345", name, state, req, err)
					continue
				}
				wantExit(t, err, 3)
				if _, serr := os.Stat(statePath); !errors.Is(serr, os.ErrNotExist) {
					t.Errorf("%s, %s, %s: state file written: %v", name, state, req, serr)
				}
			}
		}
	}
}

const (
	gitMergeBaseAbc = "git merge-base refs/remotes/origin/main abc1234def"
	gitDiffAbc      = "git diff --name-only base111 abc1234def"
)

// A red head after its reruns names the count, the failing tests and, when
// the diff is .rhiza-only, that fact: never a claim about the cause.
func TestShipCheckFailureDetailNotes(t *testing.T) {
	const logCall = "run view 42 --log-failed"
	const testLog = "test\tRun\t--- FAIL: TestGarden (0.1s)\n"
	cases := []struct{ name, diff, check, log, want string }{
		{"rhiza only", ".rhiza/log/x.json\n.rhiza/rendered/a.md\n", "test", testLog,
			"failed: test after 2 reruns; failing tests: TestGarden; the PR changes only .rhiza/"},
		{"other path", ".rhiza/log/x.json\nmain.go\n", "test", testLog,
			"failed: test after 2 reruns; failing tests: TestGarden"},
		{"reads rhiza", ".rhiza/rendered/a.md\n", "lint", "lint\tRun\trhiza check: .rhiza/rendered/a.md is stale\n",
			"failed: lint after 2 reruns; the PR changes only .rhiza/"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, g, _ := shipFixture(t)
			roll := strings.Replace(rerunRollup, `"name":"test"`, `"name":"`+c.check+`"`, 1)
			g.answers[shipPRView] = openPR(roll)
			g.answers["run rerun 42 --failed"] = ""
			g.answers[logCall] = c.log
			g.answers[gitMergeBaseAbc] = "base111\n"
			g.answers[gitDiffAbc] = c.diff
			for i := 0; i < 2; i++ {
				if r := shipCheck(t, f, `{"settled":true,"reruns":2}`); r.Settled {
					t.Fatalf("rerun %d: %+v", i, r)
				}
			}
			r := shipCheck(t, f, `{"settled":true,"reruns":2}`)
			if r.Detail != c.want {
				t.Errorf("detail = %q, want %q", r.Detail, c.want)
			}
			for _, bad := range []string{"unrelated", "flaky", "on main"} {
				if strings.Contains(r.Detail, bad) {
					t.Errorf("detail %q contains %q", r.Detail, bad)
				}
			}
		})
	}
}

func withRunEnv(f *sessionFixture, run *string) {
	prev := f.ctx.Getenv
	f.ctx.Getenv = func(k string) string {
		if k == envRun {
			return *run
		}
		if prev != nil {
			return prev(k)
		}
		return ""
	}
}

func TestShipCheckNamesRun(t *testing.T) {
	f, g, state := shipFixture(t)
	run := "01RUNA"
	withRunEnv(f, &run)
	g.answers[shipPRView] = openPR(greenRollup)
	shipCheck(t, f, `{}`)
	if st, _ := readShipState(state); st.Run != "01RUNA" {
		t.Errorf("run = %q, want 01RUNA", st.Run)
	}
	// A second run's check on the same change set takes the record over.
	run = "01RUNB"
	shipCheck(t, f, `{}`)
	if st, _ := readShipState(state); st.Run != "01RUNB" {
		t.Errorf("run = %q, want 01RUNB", st.Run)
	}
}

func TestShipLandNamesRun(t *testing.T) {
	f, state, _ := shipLandFixture(t)
	run := "01RUNA"
	withRunEnv(f, &run)
	if _, err := callOp(t, f.ctx, "ship.land", `{}`); err != nil {
		t.Fatal(err)
	}
	if st, _ := readShipState(state); st.Run != "01RUNA" {
		t.Errorf("run = %q, want 01RUNA", st.Run)
	}
}

func TestShipRebaseNamesRun(t *testing.T) {
	f, state := shipRebaseFixture(t)
	f.git.out[rebaseMerge] = ""
	f.git.out[rebasePush] = ""
	run := "01RUNA"
	withRunEnv(f, &run)
	shipRebase(t, f)
	if st, _ := readShipState(state); st.Run != "01RUNA" {
		t.Errorf("run = %q, want 01RUNA", st.Run)
	}
}

func TestShipLandHeadBuildErrorDetail(t *testing.T) {
	cases := []struct {
		name, out string
		err       error
		want      []string
	}{
		{"json body", `{"error":{"code":"conflict","message":"cannot land change set X: node N changed"}}`,
			errors.New("exit status 3"), []string{"conflict", "cannot land change set X: node N changed", "exit status 3"}},
		{"stderr only", "", errors.New("rhiza land: boom on stderr"), []string{"boom on stderr"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, x, bin := rhizaRepoFixture(t)
			x.answers["go build -o "+bin+" ./cmd/rhiza"] = ""
			f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
				if name == bin {
					return []byte(c.out), c.err
				}
				return x.run(name, args, dir, env)
			}
			_, err := callOp(t, f.ctx, "ship.land", `{}`)
			if err == nil {
				t.Fatal("no error")
			}
			for _, w := range c.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("err = %v, want %q", err, w)
				}
			}
		})
	}
}

func TestShipCheckDraftStops(t *testing.T) {
	for _, req := range []string{`{}`, `{"settled":true}`} {
		f, g, state := shipFixture(t)
		g.answers[shipPRView] = strings.Replace(openPR(greenRollup), `"state":"OPEN"`, `"state":"OPEN","isDraft":true`, 1)
		resp, err := callOp(t, f.ctx, "ship.check", req)
		if err != nil {
			t.Fatalf("%s: %v", req, err)
		}
		c := resp.(ShipCheckResponse)
		if c.ExitCode() != 75 || !strings.Contains(c.Stop, "#") || !strings.Contains(c.Stop, "draft") || c.Passed {
			t.Errorf("%s: check = %+v", req, c)
		}
		var b strings.Builder
		c.Human(&b)
		if !strings.HasPrefix(b.String(), "stop: ") {
			t.Errorf("%s: human = %q", req, b.String())
		}
		if st, _ := readShipState(state); st.Head != "" {
			t.Errorf("%s: recorded head: %+v", req, st)
		}
	}
}

func TestShipMergeDraftStops(t *testing.T) {
	f, g, state := shipFixture(t)
	if err := writeShipState(state, shipState{PR: 7, Branch: "cs/ship-it", Head: "abc1234def"}); err != nil {
		t.Fatal(err)
	}
	mainHeld(g)
	g.answers[shipDraftPR] = `{"isDraft":true,"state":"OPEN"}`
	resp, err := callOp(t, f.ctx, "ship.merge", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	m := resp.(ShipMergeResponse)
	if m.ExitCode() != 75 || !strings.Contains(m.Stop, "draft") || m.MergeCommit != "" {
		t.Errorf("merge = %+v", m)
	}
	for _, c := range g.calls {
		if strings.Contains(c, "pr merge") || strings.Contains(c, "pr ready") {
			t.Errorf("gh call %q on a draft", c)
		}
	}
	if st, _ := readShipState(state); st.MergeCommit != "" || st.Head != "abc1234def" {
		t.Errorf("state = %+v", st)
	}
}

func attemptRollup(job string, conclusion string) string {
	status := "COMPLETED"
	if conclusion == "" {
		status = "IN_PROGRESS"
	}
	return `{"name":"test","status":"` + status + `","conclusion":"` + conclusion + `","detailsUrl":"https://github.com/o/r/actions/runs/42/job/` + job + `"},{"context":"lint","state":"SUCCESS"}`
}

// --infra passes only on a recorded red head with no failing test.
func TestShipCheckInfra(t *testing.T) {
	f, g, _ := shipFixture(t)
	if r := shipCheck(t, f, `{"infra":true}`); r.Passed || r.ExitCode() != 1 {
		t.Errorf("no red recorded: %+v", r)
	}
	g.answers[shipPRView] = openPR(infraRollup)
	g.answers["run rerun 42 --failed"] = ""
	g.answers["run view 42 --log-failed"] = "setup\tapt\tTemporary failure resolving archive.ubuntu.com\n"
	shipCheck(t, f, `{"settled":true}`) // spends the head's one rerun
	if r := shipCheck(t, f, `{}`); r.Passed {
		t.Fatalf("failing check: %+v", r)
	}
	if r := shipCheck(t, f, `{"infra":true}`); !r.Passed || r.ExitCode() != 0 {
		t.Errorf("infra red: %+v", r)
	}
	g.answers["run view 42 --log-failed"] = "test\tRun\t--- FAIL: TestGarden (0.1s)\n"
	if r := shipCheck(t, f, `{}`); r.Passed {
		t.Fatalf("failing check: %+v", r)
	}
	if r := shipCheck(t, f, `{"infra":true}`); r.Passed || r.ExitCode() != 1 {
		t.Errorf("red with a failing test: %+v", r)
	}
	g.answers[shipPRView] = openPR(greenRollup)
	if r := shipCheck(t, f, `{}`); !r.Passed {
		t.Fatalf("green: %+v", r)
	}
	if r := shipCheck(t, f, `{"infra":true}`); r.Passed {
		t.Errorf("green head: %+v", r)
	}
	if _, err := callOp(t, f.ctx, "ship.check", `{"infra":true,"settled":true}`); err == nil {
		t.Error("--infra with --settled: no error")
	}
}

// --settled --changed --rerun-infra reruns with backoff from the state
// file, and settles on a green attempt, a failing test or a new head.
func TestShipCheckRerunInfraBackoff(t *testing.T) {
	f, g, _ := shipFixture(t)
	now := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	g.answers["run rerun 42 --failed"] = ""
	g.answers["run view 42 --log-failed"] = "setup\tapt\tTemporary failure resolving archive.ubuntu.com\n"
	g.answers[shipPRView] = openPR(attemptRollup("9", "FAILURE"))
	shipCheck(t, f, `{}`) // spends the head's one rerun
	if r := shipCheck(t, f, `{}`); r.Passed {
		t.Fatalf("record red: %+v", r)
	}
	g.calls = nil
	const req = `{"settled":true,"changed":true,"rerun_infra":true}`
	reruns := func() int {
		n := 0
		for _, c := range g.calls {
			if c == "gh run rerun 42 --failed" {
				n++
			}
		}
		return n
	}
	job := 9
	for i, wait := range []time.Duration{5, 10, 20, 30, 30} {
		d := wait * time.Minute
		now = now.Add(d - time.Second)
		before := reruns()
		r := shipCheck(t, f, req)
		if r.Passed || r.Settled || reruns() != before || !strings.Contains(r.Detail, "infra failure") || !strings.Contains(r.Detail, "next rerun at") {
			t.Fatalf("step %d early: %+v", i, r)
		}
		now = now.Add(time.Second)
		r = shipCheck(t, f, req)
		if r.Passed || reruns() != before+1 || len(r.Reran) != 1 {
			t.Fatalf("step %d due: %+v", i, r)
		}
		// The rerun attempt fails on infra again (a new job ID).
		job++
		g.answers[shipPRView] = openPR(attemptRollup(strconv.Itoa(job), "FAILURE"))
		now = now.Add(time.Second)
		if r := shipCheck(t, f, req); r.Passed || r.Settled {
			t.Fatalf("step %d attempt: %+v", i, r)
		}
		now = now.Add(-time.Second)
	}
	if n := reruns(); n != 5 {
		t.Fatalf("reruns = %d", n)
	}
	// The backoff survives a fresh process: it is in the state file.
	f.ctx.Now = func() time.Time { return now.Add(time.Minute) }
	if r := shipCheck(t, f, req); r.Passed || reruns() != 5 {
		t.Errorf("fresh process: %+v", r)
	}
	// A running attempt is not settled.
	g.answers[shipPRView] = openPR(attemptRollup("99", ""))
	if r := shipCheck(t, f, req); r.Passed {
		t.Errorf("running: %+v", r)
	}
	// Green attempt settles.
	g.answers[shipPRView] = openPR(attemptRollup("99", "SUCCESS"))
	if r := shipCheck(t, f, req); !r.Passed || !r.Settled {
		t.Errorf("green: %+v", r)
	}
	// Attempt red with a failing test settles.
	g.answers[shipPRView] = openPR(attemptRollup("100", "FAILURE"))
	g.answers["run view 42 --log-failed"] = "test\tRun\t--- FAIL: TestGarden (0.1s)\n"
	if r := shipCheck(t, f, req); !r.Passed || !r.Settled {
		t.Errorf("failing test: %+v", r)
	}
	// A new head settles.
	g.answers[shipPRView] = strings.Replace(openPR(greenRollup), "abc1234def", "def5678abc", 1)
	if r := shipCheck(t, f, req); !r.Passed {
		t.Errorf("new head: %+v", r)
	}
}
