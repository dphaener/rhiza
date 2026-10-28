package api

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// statusLeafRunner is a run.Runner and run.LeafRunner whose leaf sessions
// are fakes: each answers its scripted polls in turn, then exits with its
// outcome (pass unless fail names the leaf). Command steps exit with code.
type statusLeafRunner struct {
	leaves []run.Leaf
	polls  map[string][]run.AgentStatus // by session
	fail   map[string]bool              // by leaf handle: every session exits 1
	code   int
	// store is the run's journal, where each check session records an
	// empty findings round as it starts, so a check that passes has one.
	store *run.Store
}

func statusSID(spec run.LeafSpec) string {
	return spec.Leaf.Handle + "-" + spec.Role + "-" + strconv.Itoa(spec.Attempt)
}

func (f *statusLeafRunner) Command(workflow.Step, []string) (int, string, error) {
	return f.code, "", nil
}
func (f *statusLeafRunner) Holds(workflow.Step, []string) (bool, error) { return true, nil }
func (f *statusLeafRunner) StartAgent(workflow.Step, string, []string) (string, error) {
	return "", fmt.Errorf("no agents here")
}
func (f *statusLeafRunner) AgentStatus(sid string) (run.AgentStatus, error) {
	if p := f.polls[sid]; len(p) > 0 {
		f.polls[sid] = p[1:]
		return p[0], nil
	}
	for h := range f.fail {
		if strings.HasPrefix(sid, h+"-") {
			return run.AgentStatus{Done: true, Outcome: "fail", Detail: "exit 1"}, nil
		}
	}
	return run.AgentStatus{Done: true, Outcome: "pass", Detail: "exit 0"}, nil
}
func (f *statusLeafRunner) StartArbiter(workflow.Step, string, []string) (string, error) {
	return "", fmt.Errorf("no arbiter here")
}
func (f *statusLeafRunner) Answer(string, string) (run.Answer, error) { return run.Answer{}, nil }
func (f *statusLeafRunner) AddQuestion(string, string) (run.Question, error) {
	return run.Question{}, fmt.Errorf("no questions here")
}
func (f *statusLeafRunner) Deliver(string, string, string, string) error { return nil }
func (f *statusLeafRunner) AgentEnded(workflow.Step, string) error       { return nil }
func (f *statusLeafRunner) Leaves(string) ([]run.Leaf, error)            { return f.leaves, nil }
func (f *statusLeafRunner) StartLeaf(spec run.LeafSpec) (string, error)  { return statusSID(spec), nil }
func (f *statusLeafRunner) Head() (string, error)                        { return "H0", nil }
func (f *statusLeafRunner) LeafCommit(leaf run.Leaf, _ string) (string, string, error) {
	return "C-" + leaf.Handle, "", nil
}
func (f *statusLeafRunner) Evidence(run.Leaf, string) (bool, string, error) { return true, "", nil }
func (f *statusLeafRunner) CommitEvidence(leaf run.Leaf, _ string) (string, error) {
	return "E-" + leaf.Handle, nil
}
func (f *statusLeafRunner) OpenCheckout(leaf run.Leaf, attempt int, _ string) (string, error) {
	return "/checkouts/" + leaf.Handle, nil
}
func (f *statusLeafRunner) RemoveCheckout(string) error { return nil }
func (f *statusLeafRunner) AddLeafQuestion(leaf run.Leaf, _ string) (run.Question, error) {
	return run.Question{ID: "Q-" + leaf.Handle, Handle: "leaf-failed-" + leaf.Handle}, nil
}
func (f *statusLeafRunner) LeafAnswer(string) (string, string, bool, error) {
	return "", "", false, nil
}
func (f *statusLeafRunner) ClosePane(string) error { return nil }
func (f *statusLeafRunner) StartCheck(spec run.CheckSpec) (string, error) {
	if f.store != nil {
		e := run.Entry{Type: run.TypeFindings, Step: spec.Step.Name, Check: spec.Check.Name, Cycle: spec.Cycle}
		if err := f.store.Append("R1", e); err != nil {
			return "", err
		}
	}
	return spec.Check.Name + "-" + spec.Role + "-" + strconv.Itoa(spec.Cycle), nil
}
func (f *statusLeafRunner) OpenCheckCheckout(check, role string, cycle int, _ string) (string, error) {
	return "/checkouts/" + check + "-" + role + "-" + strconv.Itoa(cycle), nil
}
func (f *statusLeafRunner) StartFix(spec run.FixSpec) (string, error) {
	return "fixer-" + strconv.Itoa(spec.Cycle), nil
}
func (f *statusLeafRunner) ReviewCommit(string) (string, string, error) { return "C-review", "", nil }
func (f *statusLeafRunner) AddReviewQuestion(string) (run.Question, error) {
	return run.Question{ID: "Q-review", Handle: "review-failed-1"}, nil
}

// statusLeafRunner must stay a run.CheckRunner, or checks steps over it
// end with an error outcome.
var _ run.CheckRunner = (*statusLeafRunner)(nil)

// statusRun runs workflow src over f with Rex reporting as run.start
// wires it, and returns the reports written.
func statusRun(t *testing.T, src string, f *statusLeafRunner) []map[string]string {
	t.Helper()
	def, err := workflow.Parse("wf", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	store := run.NewStore(filepath.Join(t.TempDir(), "runs"))
	now := noSleep{}.Now()
	if err := run.Start(store, "R1", run.Started{Workflow: "wf", Definition: src, Change: "CS1", Title: "Status runs"}, now, def); err != nil {
		t.Fatal(err)
	}
	f.store = store
	rex := &bytes.Buffer{}
	e := &run.Engine{Store: store, Runner: f, Clock: noSleep{}, Now: func() time.Time { return now }}
	reportToRex(e, rex)
	if _, err := e.Run("R1"); err != nil {
		t.Fatal(err)
	}
	reports := rexReports(t, rex.String())
	idRe := regexp.MustCompile(`^[A-Za-z0-9_.+-]{1,32}(/[A-Za-z0-9_.+-]{1,32})*$`)
	for _, r := range reports {
		line := fmt.Sprintf("state=%s", r["state"])
		for _, k := range []string{"id", "kind", "progress", "title", "msg"} {
			if v, ok := r[k]; ok {
				line += fmt.Sprintf(" %s=%q", k, v)
			}
		}
		t.Log(line)
		if id, ok := r["id"]; ok && !idRe.MatchString(id) {
			t.Errorf("report %v: id %q does not match the id grammar", r, id)
		}
		if r["app"] != "" && r["app"] != "rhiza" {
			t.Errorf("report %v: app is not rhiza", r)
		}
	}
	return reports
}

const statusLeavesWF = `[step.build]
kind = "leaves"
[step.build.builder]
prompt = "b"
[step.build.verifier]
prompt = "v"
`

// first is the index of the first report from from on that ok holds for,
// or -1.
func first(reports []map[string]string, from int, ok func(map[string]string) bool) int {
	for i := max(from, 0); i < len(reports); i++ {
		if ok(reports[i]) {
			return i
		}
	}
	return -1
}

// A run over a build of three leaves keeps Rex's program status current:
// root working with the step first; per leaf a child record, builder then
// verifier, cleared once the leaf passes; root progress rising to 100; a
// blocked child while a session waits on a question; and root done last.
func TestProgramStatusBuild(t *testing.T) {
	f := &statusLeafRunner{
		leaves: []run.Leaf{
			{ID: "L1", Handle: "leaf-a"},
			{ID: "L2", Handle: "a-leaf-handle-longer-than-thirty-two-bytes", DependsOn: []string{"L1"}},
			{ID: "L3", Handle: "leaf-c"},
		},
		polls: map[string][]run.AgentStatus{
			"a-leaf-handle-longer-than-thirty-two-bytes-verifier-1": {
				{Blocked: &run.Block{Kind: "question", Question: "Which?"}}, {Blocked: &run.Block{Kind: "question"}}, {},
			},
		},
	}
	reports := statusRun(t, statusLeavesWF, f)
	if len(reports) == 0 {
		t.Fatal("no reports")
	}
	root := func(r map[string]string) bool { _, ok := r["id"]; return !ok }
	if r := reports[0]; !root(r) || r["state"] != "working" || !strings.HasPrefix(r["msg"], "step build") || r["title"] != "Status runs" {
		t.Errorf("first report %v, want root working at step build", r)
	}
	for n, lf := range f.leaves {
		id := "leaf/" + strconv.Itoa(n+1)
		b := first(reports, 0, func(r map[string]string) bool {
			return r["id"] == id && r["state"] == "working" && r["title"] == lf.Handle && r["msg"] == lf.Handle+": builder working"
		})
		v := first(reports, b, func(r map[string]string) bool {
			return r["id"] == id && r["title"] == lf.Handle && strings.HasPrefix(r["msg"], lf.Handle+": verifier ")
		})
		c := first(reports, v, func(r map[string]string) bool { return r["id"] == id && r["state"] == "clear" })
		if b < 0 || v < 0 || c < 0 {
			t.Errorf("leaf %s: builder at %d, verifier at %d, clear at %d; want all, in that order", id, b, v, c)
		}
	}
	blocked := first(reports, 0, func(r map[string]string) bool {
		return r["id"] == "leaf/2" && r["state"] == "blocked" && r["kind"] == "question" &&
			strings.HasSuffix(r["msg"], ": verifier blocked on a question")
	})
	if blocked < 0 {
		t.Error("no blocked child record for leaf/2's verifier")
	}
	if i := first(reports, blocked, func(r map[string]string) bool { return root(r) && r["state"] == "blocked" }); i < 0 {
		t.Error("the owner's blocked notice is no longer the root record while leaf/2 waits")
	}
	var progress []int
	for _, r := range reports {
		if p, ok := r["progress"]; ok {
			if !root(r) {
				t.Errorf("progress on a child record %v", r)
			}
			n, err := strconv.Atoi(p)
			if err != nil {
				t.Fatal(err)
			}
			progress = append(progress, n)
		}
	}
	for i := 1; i < len(progress); i++ {
		if progress[i] < progress[i-1] {
			t.Errorf("progress falls: %v", progress)
		}
	}
	if len(progress) == 0 || progress[0] != 0 || progress[len(progress)-1] != 100 {
		t.Errorf("progress %v, want from 0 rising to 100", progress)
	}
	if i := first(reports, 0, func(r map[string]string) bool {
		return root(r) && strings.Contains(r["msg"], "2 of 3 leaves passed")
	}); i < 0 {
		t.Error("no root msg counts the leaves passed")
	}
	last := reports[len(reports)-1]
	if !root(last) || last["state"] != "done" || last["msg"] != "run succeeded" || last["title"] != "Status runs" {
		t.Errorf("last report %v, want root done", last)
	}
}

// A run whose leaf fails every attempt, or whose command step fails, ends
// with root error naming the failure on one line.
func TestProgramStatusError(t *testing.T) {
	t.Run("leaf", func(t *testing.T) {
		f := &statusLeafRunner{leaves: []run.Leaf{{ID: "L1", Handle: "leaf-a"}}, fail: map[string]bool{"leaf-a": true}}
		reports := statusRun(t, statusLeavesWF, f)
		last := reports[len(reports)-1]
		if _, ok := last["id"]; ok || last["state"] != "error" || !strings.HasPrefix(last["msg"], "run failed: ") {
			t.Errorf("last report %v, want root error", last)
		}
		if i := first(reports, 0, func(r map[string]string) bool {
			return r["id"] == "leaf/1" && r["msg"] == "leaf-a: builder working (attempt 3)"
		}); i < 0 {
			t.Error("no child record for the third attempt's builder")
		}
	})
	t.Run("command", func(t *testing.T) {
		f := &statusLeafRunner{code: 1}
		reports := statusRun(t, "[step.check]\nkind = \"command\"\nrun = \"false\"\n", f)
		if len(reports) != 2 {
			t.Fatalf("reports %v, want working then error", reports)
		}
		if r := reports[0]; r["state"] != "working" || r["msg"] != "step check" {
			t.Errorf("first report %v", r)
		}
		if r := reports[1]; r["state"] != "error" || r["msg"] != `run failed: step "check" ended with outcome "fail"` {
			t.Errorf("last report %v", r)
		}
	})
}

// rexRecord keeps progress to working and blocked, writes no app on a
// clear, and puts id on a child record.
func TestRexRecord(t *testing.T) {
	p := 40
	for _, c := range []struct {
		id, state string
		progress  *int
		want      string
	}{
		{"", "working", &p, "\x1b]7501;state=working:progress=40:app=rhiza\x1b\\"},
		{"", "done", &p, "\x1b]7501;state=done:app=rhiza\x1b\\"},
		{"leaf/1", "clear", nil, "\x1b]7501;state=clear:id=leaf/1\x1b\\"},
		{"leaf/2", "working", nil, "\x1b]7501;state=working:id=leaf/2:app=rhiza\x1b\\"},
	} {
		if got := string(rexRecord(c.id, c.state, "", "", "", c.progress)); got != c.want {
			t.Errorf("rexRecord(%q, %s) = %q, want %q", c.id, c.state, got, c.want)
		}
	}
}

const statusChecksWF = `[step.review]
kind = "checks"
[[step.review.check]]
name = "alpha"
prompt = "a"
[[step.review.check]]
name = "beta"
prompt = "b"
[step.review.critic]
prompt = "c"
[step.review.fixer]
prompt = "f"
[step.review.verifier]
prompt = "v"
`

// A run over a checks step keeps Rex's program status current: a child
// record per live check, its cycle named past cycle 1, and the root
// counting the checks passed up to 100.
func TestProgramStatusChecks(t *testing.T) {
	f := &statusLeafRunner{polls: map[string][]run.AgentStatus{
		"beta-check-1": {{Done: true, Outcome: "fail", Detail: "exit 1"}},
	}}
	reports := statusRun(t, statusChecksWF, f)
	root := func(r map[string]string) bool { _, ok := r["id"]; return !ok }
	if i := first(reports, 0, func(r map[string]string) bool {
		return r["id"] == "check/1" && r["title"] == "alpha" && r["msg"] == "alpha: check working"
	}); i < 0 {
		t.Error("no child record for check alpha")
	}
	if i := first(reports, 0, func(r map[string]string) bool {
		return r["id"] == "check/2" && r["msg"] == "beta: check working (cycle 2)"
	}); i < 0 {
		t.Error("no child record for check beta's cycle 2")
	}
	if i := first(reports, 0, func(r map[string]string) bool {
		return root(r) && strings.Contains(r["msg"], "1 of 2 checks passed") && r["progress"] == "50"
	}); i < 0 {
		t.Error("no root msg counts the checks passed")
	}
	last := reports[len(reports)-1]
	if !root(last) || last["state"] != "done" {
		t.Errorf("last report %v, want root done", last)
	}
}
