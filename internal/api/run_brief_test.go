package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/workflow"
)

// envOf is the value launchEnv sets for name, and whether it sets it.
func envOf(rec session.Record, name string) (string, bool) {
	set, _ := launchEnv(rec)
	for _, kv := range set {
		if v, ok := strings.CutPrefix(kv, name+"="); ok {
			return v, true
		}
	}
	return "", false
}

// An agent step's, a builder's and a checks step's verifier's session each
// get their own directory with brief.md and slice.json, slice.json the
// bytes `rhiza slice --json` prints, RHIZA_BRIEF and RHIZA_SLICE set, and
// RHIZA_HANDBACK only once handback.md exists; a session started outside a
// run gets none.
func TestSessionBrief(t *testing.T) {
	t.Parallel()
	f, r := paneRunner(t)
	// The run's worktree holds the change set, as a real one does.
	csID, err := id.Parse(r.change)
	if err != nil {
		t.Fatal(err)
	}
	if err := copyChangeSetLog(f.ctx.Root, f.worktree, csID); err != nil {
		t.Fatal(err)
	}
	// slice is what rhiza slice --json prints for the change set now.
	slice := func() []byte {
		sc := *testCtx(t, f.worktree)
		sc.ChangeRef = r.change
		sl, err := callOp(t, &sc, "slice", `{}`)
		if err != nil {
			t.Fatal(err)
		}
		b, err := canon.Marshal(sl)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	leaf := run.Leaf{Handle: "leaf-a"}
	starts := map[string]func() (string, error){
		"agent": func() (string, error) {
			return r.StartAgent(workflow.Step{Name: "plan", Prompt: "plan"}, "", nil)
		},
		"builder": func() (string, error) {
			return r.StartLeaf(run.LeafSpec{Step: workflow.Step{Name: "build", Builder: &workflow.Role{Prompt: "b"}},
				Role: workflow.RoleBuilder, Leaf: leaf, Attempt: 1})
		},
		"verifier": func() (string, error) {
			return r.StartCheck(run.CheckSpec{Step: workflow.Step{Name: "review", Verifier: &workflow.Role{Prompt: "v"}},
				Role: workflow.RoleVerifier, Check: workflow.Check{Name: "alpha"}, Cycle: 1})
		},
		// A checks step's leaf verifier: StartLeaf with a cycle, as
		// run's checks step starts it.
		"checks verifier": func() (string, error) {
			return r.StartLeaf(run.LeafSpec{Step: workflow.Step{Name: "checks", Verifier: &workflow.Role{Prompt: "v"}},
				Role: workflow.RoleVerifier, Leaf: leaf, Attempt: 1, Cycle: 1})
		},
	}
	// wantBase is the directory name a start must get, where it is fixed.
	wantBase := map[string]string{
		"checks verifier": run.LeafFixSessionName(workflow.RoleVerifier, leaf.Handle, 1, 1),
	}
	dirs := map[string]bool{}
	for name, start := range starts {
		want := slice()
		sid, err := start()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		rec, err := r.sessions.Read(sid)
		if err != nil {
			t.Fatal(err)
		}
		brief, ok := envOf(rec, envBrief)
		if !ok || filepath.Base(brief) != "brief.md" {
			t.Fatalf("%s: RHIZA_BRIEF = %q, %v", name, brief, ok)
		}
		dir := filepath.Dir(brief)
		if dirs[dir] || !strings.HasPrefix(dir, filepath.Join(f.common, "rhiza", "runs", "R1")) {
			t.Errorf("%s: directory %s is shared or outside the run's", name, dir)
		}
		dirs[dir] = true
		if w, ok := wantBase[name]; ok && filepath.Base(dir) != w {
			t.Errorf("%s: directory %s, want one named %s", name, dir, w)
		}
		if sp, _ := envOf(rec, envSlice); sp != filepath.Join(dir, "slice.json") {
			t.Errorf("%s: RHIZA_SLICE = %q", name, sp)
		}
		got, err := os.ReadFile(filepath.Join(dir, "slice.json"))
		if err != nil || string(got) != string(want) {
			t.Errorf("%s: slice.json differs from rhiza slice --json (%v)\n got %s\nwant %s", name, err, got, want)
		}
		md, err := os.ReadFile(brief)
		if err != nil || !strings.Contains(string(md), "Readable tabs") || !strings.HasSuffix(string(md), run.ContractLastLine+"\n") {
			t.Errorf("%s: brief.md %q (%v)", name, md, err)
		}
		if _, ok := envOf(rec, envHandback); ok {
			t.Errorf("%s: RHIZA_HANDBACK set with no handback.md", name)
		}
		if err := os.WriteFile(filepath.Join(dir, "handback.md"), []byte("x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if hb, ok := envOf(rec, envHandback); !ok || hb != filepath.Join(dir, "handback.md") {
			t.Errorf("%s: RHIZA_HANDBACK = %q, %v with handback.md", name, hb, ok)
		}
	}
	if rec := (session.Record{ID: "S-out", Change: r.change}); true {
		for _, n := range []string{envBrief, envSlice, envHandback} {
			if _, ok := envOf(rec, n); ok {
				t.Errorf("a session outside a run got %s", n)
			}
		}
	}
}

// An unwritable runs directory fails the start with its path and opens no
// tab.
func TestSessionBriefUnwritable(t *testing.T) {
	t.Parallel()
	f, r := paneRunner(t)
	if err := os.MkdirAll(filepath.Join(f.common, "rhiza"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.common, "rhiza", "runs"), []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.rex.calls = nil
	_, err := r.StartAgent(workflow.Step{Name: "plan", Prompt: "plan"}, "", nil)
	if err == nil || !strings.Contains(err.Error(), filepath.Join(f.common, "rhiza", "runs")) {
		t.Fatalf("error = %v, want one naming the runs directory", err)
	}
	if len(f.rex.calls) != 0 {
		t.Errorf("rex called: %q", f.rex.calls)
	}
}

// brief.md lists every decision that applies with its statement, even when
// the default slice budget cuts the Decisions section to a count, and the
// change set's own Decisions, which the slice holds in focus.
func TestSessionBriefDecisions(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "brief-decisions")
	newChange(t, f.ctx, `{"intent":"Base"}`)
	var writes []string
	writes = append(writes, `{"op":"add","kind":"capability","handle":"base","fields":["statement=The base"]}`)
	long := strings.Repeat("because it holds ", 40)
	for i := range 40 {
		h := fmt.Sprintf("landed-%02d", i)
		writes = append(writes,
			fmt.Sprintf(`{"op":"add","kind":"decision","handle":%q,"fields":["statement=Landed decision %02d %s"]}`, h, i, long),
			fmt.Sprintf(`{"op":"link","from":"base","type":"decided-by","to":%q}`, h))
	}
	if _, err := applyFile(t, f.ctx, `{"writes":[`+strings.Join(writes, ",")+`]}`); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"propose", "accept", "land"} {
		move(t, f.ctx, s, `{}`)
	}
	cs := newChange(t, f.ctx, `{"intent":"Readable tabs"}`).ChangeSet
	if _, err := applyFile(t, f.ctx, `{"writes":[
	  {"op":"add","kind":"capability","handle":"child","parent":"base","fields":["statement=The child"]},
	  {"op":"add","kind":"decision","handle":"own","fields":["statement=Own decision of the change set"]},
	  {"op":"link","from":"child","type":"decided-by","to":"own"}
	]}`); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: f.ctx, change: cs.String(), run: "R1", worktree: filepath.Join(t.TempDir(), "none")}
	dir := filepath.Join(t.TempDir(), "session")
	if err := r.fillSession(dir, nil); err != nil {
		t.Fatal(err)
	}
	js, err := os.ReadFile(filepath.Join(dir, "slice.json"))
	if err != nil || !strings.Contains(string(js), `"section": "decisions"`) {
		t.Fatalf("slice.json does not cut decisions, so this test proves nothing (%v): %s", err, js)
	}
	md, err := os.ReadFile(filepath.Join(dir, "brief.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Own decision of the change set"}
	for i := range 40 {
		want = append(want, fmt.Sprintf("Landed decision %02d", i))
	}
	for _, w := range want {
		if !strings.Contains(string(md), w) {
			t.Errorf("brief.md lacks %q:\n%s", w, md)
		}
	}
}

// brief.md lists the change set's own decisions and constraints, which the
// slice holds in focus, each by its statement.
func TestSessionBriefOwnDecisionsConstraints(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "brief-own")
	cs := newChange(t, f.ctx, `{"intent":"Readable tabs"}`).ChangeSet
	if _, err := applyFile(t, f.ctx, `{"writes":[
	  {"op":"add","kind":"capability","handle":"alpha","fields":["statement=Alpha"]},
	  {"op":"add","kind":"capability","handle":"beta","fields":["statement=Beta"]},
	  {"op":"add","kind":"decision","handle":"dec1","fields":["statement=Use maps."]},
	  {"op":"add","kind":"constraint","handle":"con1","fields":["statement=No network."]}
	]}`); err != nil {
		t.Fatal(err)
	}
	r := &stepRunner{ctx: f.ctx, change: cs.String(), run: "R1", worktree: filepath.Join(t.TempDir(), "none")}
	dir := filepath.Join(t.TempDir(), "session")
	if err := r.fillSession(dir, nil); err != nil {
		t.Fatal(err)
	}
	md, err := os.ReadFile(filepath.Join(dir, "brief.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range []string{"## Decisions\n\n- Use maps.", "## Constraints\n\n- No network."} {
		if !strings.Contains(string(md), w) {
			t.Errorf("brief.md lacks %q:\n%s", w, md)
		}
	}
}

// A capability takes files and must_not, which read back and reach brief.md
// under the leaf's acceptance; (failure case) a field "file" is refused,
// naming the contract fields among the known ones.
func TestSessionBriefLeafContract(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "brief-contract")
	cs := newChange(t, f.ctx, `{"intent":"Readable tabs"}`).ChangeSet
	mustWrite(t, f.ctx, "add", `{"kind":"capability","handle":"alpha","fields":["statement=Alpha","acceptance=go test a",`+
		`"files=internal/a.go","must_not=touch b"]}`)
	n := nodeAt(t, f.ctx, "alpha")
	if got, want := string(n.Fields["files"])+" "+string(n.Fields["must_not"]), `"internal/a.go" "touch b"`; got != want {
		t.Errorf("files, must_not = %s, want %s", got, want)
	}
	_, err := callOp(t, f.ctx, "add", `{"kind":"capability","handle":"beta","fields":["statement=Beta","file=x.go"]}`)
	wantExit(t, err, 2)
	if err == nil || !strings.Contains(err.Error(), "known fields: acceptance, built_in, failure_waiver, files, must_not, statement") {
		t.Errorf("field file: %v", err)
	}
	r := &stepRunner{ctx: f.ctx, change: cs.String(), run: "R1", worktree: filepath.Join(t.TempDir(), "none")}
	dir := filepath.Join(t.TempDir(), "session")
	if err := r.fillSession(dir, nil); err != nil {
		t.Fatal(err)
	}
	md, err := os.ReadFile(filepath.Join(dir, "brief.md"))
	if err != nil {
		t.Fatal(err)
	}
	if w := "Acceptance: go test a\nFiles: internal/a.go\nMust not: touch b\n\n"; !strings.Contains(string(md), w) {
		t.Errorf("brief.md lacks %q:\n%s", w, md)
	}
	if strings.Contains(string(md), "Failure waiver:") {
		t.Errorf("brief.md lists a failure waiver the leaf lacks:\n%s", md)
	}
}

// retryBriefWorkflow's agent step is followed by a review-check-shaped
// command step that, failing, retries the run from the agent step.
const retryBriefWorkflow = `
[step.review]
kind = "agent"
prompt = "review {change}"
on.pass = "check"

[step.check]
kind = "command"
run = "rhiza run review-check"
retry_from = "review"
`

// A run fails at a review-check-shaped command step after a findings
// round; retried, it restarts the agent step, whose new session's brief.md
// holds the refusal's output and the round's findings under Previous
// attempt. (failure case) The session started before the retry, in a run
// not yet retried, gets no such section.
func TestRunRetryBriefPreviousAttempt(t *testing.T) {
	t.Parallel()
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "w", retryBriefWorkflow)
	const refusal = "review-check: refused, 1 major finding is open"
	checkCode := 1
	f.ctx.Stderr = io.Discard
	f.ctx.Start = func(p Program) (Process, error) {
		if len(p.Args) == 2 && strings.TrimPrefix(p.Args[1], "exec 2>&1\n") == "rhiza run review-check" && checkCode != 0 {
			return &fakeProcess{out: strings.NewReader(refusal + "\n"), code: checkCode}, nil
		}
		return &fakeProcess{out: strings.NewReader("")}, nil
	}
	runs := run.NewStore(run.Dir(f.common))
	sessions := session.NewStore(session.Dir(f.common))
	tabs := 0
	// Every agent session exits 0 as its tab opens; the first records a
	// findings round as a plan reviewer would.
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		f.rex.run(name, args, dir, env)
		if len(args) > 1 && args[0] == "block" && args[1] == "close" {
			return nil, nil
		}
		if tabs++; tabs == 1 {
			list, err := runs.List()
			if err != nil || len(list) != 1 {
				t.Fatalf("runs %v, %v", list, err)
			}
			if err := runs.Append(list[0].ID, run.Entry{Type: run.TypeFindings, At: f.ctx.Now(),
				Findings: []run.Finding{{ID: "f1", Kind: "acceptance", Severity: "major", Reason: "leaf alpha has no acceptance"}}}); err != nil {
				t.Fatal(err)
			}
		}
		code := 0
		if _, err := sessions.Update(args[len(args)-1], func(r *session.Record) {
			r.Status.Exit = &session.Report{State: session.StateExited, Code: &code}
		}); err != nil {
			t.Fatal(err)
		}
		return []byte(rexWindowJSON), nil
	}
	f.ctx.repoSession = true
	resp, err := callOp(t, f.ctx, "run.start", `{"foreground":true,"workflow":"w","permission_mode":"manual"}`)
	if err != nil {
		t.Fatal(err)
	}
	info := resp.(RunResponse).Run
	if info.State != run.Failed {
		t.Fatalf("first attempt = %+v", info)
	}
	brief := func(visit int) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(run.SessionDir(f.common, info.ID, "review", visit, run.AgentSessionName()), briefFile))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	if b := brief(1); strings.Contains(b, "Previous attempt") {
		t.Errorf("the session started before the retry has a previous attempt:\n%s", b)
	}
	checkCode = 0
	resp, err = callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"`+info.ID+`","retry":true}`)
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(RunResponse); r.Run.State != run.Succeeded || r.Retry != "review" {
		t.Fatalf("retried run = %+v", r)
	}
	b := brief(2)
	for _, w := range []string{
		"## Previous attempt\n", "- Failed step: check\n- Outcome: fail\n- Failing command: `rhiza run review-check`\n",
		"### Output tail\n\n```\n" + refusal + "\n```\n",
		"### Findings of the latest round\n\n- plan review, major: leaf alpha has no acceptance\n",
	} {
		if !strings.Contains(b, w) {
			t.Errorf("the restarted step's brief.md lacks %q:\n%s", w, b)
		}
	}
	if b := brief(1); strings.Contains(b, "Previous attempt") {
		t.Errorf("the retry rewrote the earlier session's brief.md:\n%s", b)
	}
}
