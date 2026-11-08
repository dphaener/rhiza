package api

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// openPRGH is a fake gh with one pull request that exists once pr create
// has run. It keeps the title and body each create and edit were given.
type openPRGH struct {
	t      *testing.T
	exec   func(string, []string, string, []string) ([]byte, error)
	origin string
	branch string
	exists bool
	state  string
	title  string
	body   string
	calls  []string
}

func (g *openPRGH) run(name string, args []string, dir string, env []string) ([]byte, error) {
	if name != "gh" {
		return g.exec(name, args, dir, env)
	}
	k := strings.Join(args, " ")
	g.calls = append(g.calls, k)
	file := func() string {
		for i, a := range args {
			if a == "--body-file" {
				b, err := os.ReadFile(args[i+1])
				if err != nil {
					g.t.Fatal(err)
				}
				return string(b)
			}
		}
		g.t.Fatalf("gh %s: no --body-file", k)
		return ""
	}
	switch {
	case strings.HasPrefix(k, "pr view "):
		if !g.exists {
			return nil, errors.New("gh " + k + ": no pull requests found for branch \"" + g.branch + "\"")
		}
		// A branch origin lacks has no head yet: "" (a foreign PR seen
		// before anything is pushed).
		head, _ := g.exec("git", []string{"--git-dir", g.origin, "rev-parse", "--verify", "--quiet", "refs/heads/" + g.branch}, "", nil)
		return []byte(`{"number":7,"url":"https://github.com/o/r/pull/7","state":"` + g.state + `","body":` + strconv.Quote(g.body) + `,"headRefOid":"` + strings.TrimSpace(string(head)) + `","statusCheckRollup":[]}`), nil
	case strings.HasPrefix(k, "pr create "):
		if g.exists {
			g.t.Fatalf("gh %s: a second pull request", k)
		}
		for i, a := range args {
			if a == "--title" {
				g.title = args[i+1]
			}
		}
		g.exists, g.state, g.body = true, "OPEN", file()
		return []byte("https://github.com/o/r/pull/7\n"), nil
	case strings.HasPrefix(k, "pr edit 7 "):
		g.body = file()
		return nil, nil
	}
	g.t.Fatalf("unexpected gh call %q", k)
	return nil, nil
}

// openPRFixture is a real git repository on cs/ship-it with a bare origin,
// a run R1 of the change set whose journal holds the given entries after
// its start, and a fake gh.
func openPRFixture(t *testing.T, review bool, ref ...string) (*Ctx, *openPRGH, string, func(...string) string) {
	t.Helper()
	tmp := t.TempDir()
	root := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	ctx := testCtx(root)
	exec := hermeticExec(t, tmp)
	git := func(args ...string) string {
		t.Helper()
		out, err := exec("git", args, root, nil)
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}
	ctx.Git = func(args ...string) (string, error) {
		out, err := exec("git", args, root, nil)
		return string(out), err
	}
	origin := filepath.Join(tmp, "origin.git")
	if _, err := exec("git", []string{"init", "-q", "--bare", origin}, tmp, nil); err != nil {
		t.Fatal(err)
	}
	git("init", "-q")
	git("remote", "add", "origin", origin)
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	req := map[string]string{"intent": "Ship it", "title": "Ship it"}
	if len(ref) > 0 {
		req["ref"] = ref[0]
	}
	reqJSON, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	cs := newChange(t, ctx, string(reqJSON)).ChangeSet
	git("add", "-A")
	git("commit", "-q", "-m", "model")
	git("checkout", "-q", "-b", "cs/ship-it")

	g := &openPRGH{t: t, exec: exec, origin: origin, branch: "cs/ship-it"}
	ctx.Exec = g.run

	common := filepath.Join(root, ".git")
	store := run.NewStore(run.Dir(common))
	entries := []run.Entry{
		{Type: run.TypeStarted, Started: &run.Started{Workflow: "review", Change: cs.String(), Title: "Ship it"}},
	}
	if review {
		entries = append(entries,
			run.Entry{Type: run.TypeCheckCycle, Step: "review", Cycle: 1, Checks: []string{"alpha"}},
			run.Entry{Type: run.TypeFindings, Step: "review", Check: "alpha", Cycle: 1, Findings: []run.Finding{}},
			run.Entry{Type: run.TypeCheck, Step: "review", Check: "alpha", Cycle: 1, Outcome: run.CheckPassed},
		)
	}
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	for i, e := range entries {
		e.At = now
		if i == 0 {
			e.Step = "review"
			if err := store.Create("R1", e); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := store.Append("R1", e); err != nil {
			t.Fatal(err)
		}
	}
	ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R1"
		}
		return ""
	}
	return ctx, g, origin, git
}

// The acceptance of leaf ship-open-pr, against a real git remote and a fake
// gh: a first call pushes the branch and creates one PR titled by the
// change set whose body is the review summary; a second call after a new
// commit pushes it and edits the body without creating another; and the
// PR is recorded in the ship state file for the steps after it.
func TestShipOpenPR(t *testing.T) {
	ctx, g, origin, git := openPRFixture(t, true)
	remoteHead := func() string {
		out, err := g.exec("git", []string{"--git-dir", origin, "rev-parse", "cs/ship-it"}, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(out))
	}
	resp, err := callOp(t, ctx, "ship.open-pr", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(ShipOpenPRResponse)
	first := git("rev-parse", "HEAD")
	if !r.Created || !r.Pushed || r.PR != 7 || r.Head != first || remoteHead() != first {
		t.Fatalf("first call = %+v, origin at %s", r, remoteHead())
	}
	if g.title != "Ship it" || !strings.HasPrefix(g.body, "# Review of run R1\n") || !strings.Contains(g.body, "### alpha: passed") {
		t.Errorf("title %q, body:\n%s", g.title, g.body)
	}
	if got := human(t, r); got != "Opened PR #7 (pushed "+first[:7]+"): https://github.com/o/r/pull/7.\n" {
		t.Errorf("human = %q", got)
	}
	state := filepath.Join(ctx.Root, ".git", "rhiza", "ship", r.Change+".json")
	if st, _ := readShipState(state); st.PR != 7 || st.URL != "https://github.com/o/r/pull/7" || st.Branch != "cs/ship-it" || st.Pushed != first || st.Run != "R1" {
		t.Errorf("state = %+v", st)
	}

	// A new commit and a new line in the summary: pushed, body replaced.
	if err := os.WriteFile(filepath.Join(ctx.Root, "more.txt"), []byte("more\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git("add", "more.txt")
	git("commit", "-q", "-m", "more")
	store := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	if err := store.Append("R1", run.Entry{Type: run.TypeCodeCheck, Step: "code", Outcome: run.CodeCheckPass, At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	before := g.body
	resp, err = callOp(t, ctx, "ship.open-pr", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r = resp.(ShipOpenPRResponse)
	second := git("rev-parse", "HEAD")
	if r.Created || !r.Pushed || r.PR != 7 || remoteHead() != second || second == first {
		t.Fatalf("second call = %+v, origin at %s", r, remoteHead())
	}
	if g.body == before || !strings.Contains(g.body, "step code: pass") {
		t.Errorf("body not replaced:\n%s", g.body)
	}
	creates := 0
	for _, c := range g.calls {
		if strings.HasPrefix(c, "pr create ") {
			creates++
		}
	}
	if creates != 1 {
		t.Errorf("%d pr create calls in %q", creates, g.calls)
	}
	if st, _ := readShipState(state); st.Pushed != second || st.Head != "" {
		t.Errorf("state after second call = %+v", st)
	}
	if got := human(t, r); !strings.HasPrefix(got, "Updated PR #7 (pushed ") {
		t.Errorf("human = %q", got)
	}

	// ship check then reads the recorded PR.
	cr, err := callOp(t, ctx, "ship.check", `{"pr_open":true}`)
	if err != nil {
		t.Fatal(err)
	}
	if chk := cr.(ShipCheckResponse); !chk.Passed || chk.PR != 7 {
		t.Errorf("ship check = %+v", cr)
	}
}

// Outside a run, or in one with no review journaled, ship open-pr exits
// non-zero naming why, and pushes and opens nothing.
func TestShipOpenPRRefusals(t *testing.T) {
	ctx, g, origin, _ := openPRFixture(t, false)
	_, err := callOp(t, ctx, "ship.open-pr", `{}`)
	wantExit(t, err, 2)
	if !strings.Contains(err.Error(), "no review journaled") {
		t.Errorf("no review: %v", err)
	}

	ctx.Getenv = func(string) string { return "" }
	_, err = callOp(t, ctx, "ship.open-pr", `{}`)
	wantExit(t, err, 2)
	if !strings.Contains(err.Error(), "works inside a run") || !strings.Contains(err.Error(), envRun) {
		t.Errorf("outside a run: %v", err)
	}
	if len(g.calls) != 0 {
		t.Errorf("gh called: %q", g.calls)
	}
	if out, _ := g.exec("git", []string{"--git-dir", origin, "branch", "--list"}, "", nil); strings.TrimSpace(string(out)) != "" {
		t.Errorf("origin has branches: %s", out)
	}
}

// A plan run with no review round is refused like a review run with none:
// ship open-pr opens nothing for a plan nobody has reviewed.
func TestShipOpenPRRefusesPlanWithoutRound(t *testing.T) {
	ctx, g, _, _ := openPRFixture(t, false)
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cs := w.changeSets()[0].ID.String()
	store := run.NewStore(run.Dir(filepath.Join(ctx.Root, ".git")))
	err = store.Create("R2", run.Entry{
		Type: run.TypeStarted, Step: "interview", At: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC),
		Started: &run.Started{Workflow: "plan", Change: cs, Title: "Ship it"},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx.Getenv = func(k string) string {
		if k == envRun {
			return "R2"
		}
		return ""
	}
	_, err = callOp(t, ctx, "ship.open-pr", `{}`)
	wantExit(t, err, 2)
	if !strings.Contains(err.Error(), "no review journaled") {
		t.Errorf("plan run with no round: %v", err)
	}
	if len(g.calls) != 0 {
		t.Errorf("gh called: %q", g.calls)
	}
}

// TestShipOpenPRRefusesForeignPR: a pull request on the branch whose body
// names another change set, or none, is refused whatever its state, with a
// conflict naming it; nothing is edited, created or written to ship state.
func TestShipOpenPRRefusesForeignPR(t *testing.T) {
	for _, body := range []string{"# Plan\n\nChange set 01M6T1ZR0000000000000000AB: Other.\n", "Fixes a bug."} {
		for _, state := range []string{"OPEN", "MERGED", "CLOSED"} {
			// onOrigin: the foreign PR's commit is ahead on origin's branch,
			// which open-pr must not pull in; otherwise origin lacks the
			// branch and this run's commit must not be pushed onto the PR.
			for _, onOrigin := range []bool{true, false} {
				ctx, g, _, git := openPRFixture(t, true)
				g.exists, g.state, g.body = true, state, body
				local := git("rev-parse", "HEAD")
				if onOrigin {
					git("checkout", "-q", "-b", "foreign")
					git("commit", "-q", "--allow-empty", "-m", "foreign")
					git("push", "-q", "origin", "foreign:cs/ship-it")
					git("checkout", "-q", "cs/ship-it")
				}
				originHead := func() string {
					out, _ := g.exec("git", []string{"--git-dir", g.origin, "rev-parse", "--verify", "--quiet", "refs/heads/cs/ship-it"}, "", nil)
					return strings.TrimSpace(string(out))
				}
				before := originHead()
				_, err := callOp(t, ctx, "ship.open-pr", `{}`)
				if err == nil || !strings.Contains(err.Error(), "#7") {
					t.Errorf("%q, %s, %v: err = %v, want a conflict naming #7", body, state, onOrigin, err)
					continue
				}
				wantExit(t, err, 3)
				for _, c := range g.calls {
					if strings.HasPrefix(c, "pr create") || strings.HasPrefix(c, "pr edit") {
						t.Errorf("%q, %s, %v: gh %s", body, state, onOrigin, c)
					}
				}
				if g.body != body {
					t.Errorf("%q, %s, %v: body changed to %q", body, state, onOrigin, g.body)
				}
				if h := git("rev-parse", "HEAD"); h != local {
					t.Errorf("%q, %s, %v: HEAD moved to %s, want %s: the foreign PR's commits were pulled in", body, state, onOrigin, h, local)
				}
				if h := originHead(); h != before {
					t.Errorf("%q, %s, %v: origin's cs/ship-it moved from %q to %q: pushed onto the foreign PR", body, state, onOrigin, before, h)
				}
				ents, _ := filepath.Glob(filepath.Join(ctx.Root, ".git", "rhiza", "ship", "*.json"))
				if len(ents) != 0 {
					t.Errorf("%q, %s, %v: ship state written: %v", body, state, onOrigin, ents)
				}
			}
		}
	}
}

// The acceptance of leaf ship-pr-body-ref: a change set's ref is listed in
// the PR body in one inline code span with its backticks and line breaks
// escaped, and a change set without a ref gets no ref line.
func TestShipOpenPRBodyRef(t *testing.T) {
	ref := "https://x.test/1 `a` @mention\n# heading"
	ctx, g, _, _ := openPRFixture(t, true, ref)
	if _, err := callOp(t, ctx, "ship.open-pr", `{}`); err != nil {
		t.Fatal(err)
	}
	want := "Ref: `https://x.test/1 \\x60a\\x60 @mention\\n# heading`\n"
	if !strings.Contains(g.body, "\n"+want) || strings.Count(g.body, "Ref:") != 1 {
		t.Errorf("body lacks the escaped ref line %q:\n%s", want, g.body)
	}
	for _, l := range strings.Split(g.body, "\n") {
		if strings.Contains(l, "heading") && !strings.HasPrefix(l, "Ref: `") {
			t.Errorf("ref text escaped its line: %q", l)
		}
	}

	ctx, g, _, _ = openPRFixture(t, true)
	if _, err := callOp(t, ctx, "ship.open-pr", `{}`); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(g.body, "Ref:") {
		t.Errorf("body for a change set without a ref has a ref line:\n%s", g.body)
	}
}
