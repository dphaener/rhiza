package run

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

func (f *fakeRunner) ChangedFiles(from, to string) ([]string, error) {
	f.changedRanges = append(f.changedRanges, from+".."+to)
	return f.changed, f.changedErr
}

const tierWF = leavesWF + "models = [\"opus\"]\nlow_risk_models = [\"sonnet\"]\nrisky_paths = [\".github/**\"]\n"

func verifierModels(r *fakeRunner) []string {
	var out []string
	for _, s := range r.leafSpecs {
		if s.Role == workflow.RoleVerifier {
			out = append(out, s.Model)
		}
	}
	return out
}

// journalledVerifierModels are the models of the journal's verifier
// leaf-session entries, in order, as "handle attempt model".
func journalledVerifierModels(t *testing.T, e *Engine) []string {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, en := range j.Entries {
		if en.Type == TypeLeafSession && en.Role == workflow.RoleVerifier {
			out = append(out, en.Handle+" "+attemptStr(en.Attempt)+" "+en.Model)
		}
	}
	return out
}

func runTier(t *testing.T, src string, prep func(r *fakeRunner)) (*fakeRunner, *Engine) {
	t.Helper()
	e, r, _, _ := setup(t, src)
	r.leaves = twoLeaves[:1]
	passAll(r, r.leaves)
	prep(r)
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	return r, e
}

// A low-risk first attempt starts its verifier on the low-risk ladder, a
// commit touching a risky path on models; the journal holds the choice.
func TestVerifierTierChoice(t *testing.T) {
	cases := []struct {
		name  string
		src   string
		files []string
		err   error
		want  string
	}{
		{"low risk", tierWF, []string{"internal/x/x.go"}, nil, "sonnet"},
		{"risky path", tierWF, []string{"internal/x/x.go", ".github/workflows/ci.yml"}, nil, "opus"},
		{"diff unreadable", tierWF, nil, errors.New("bad object"), "opus"},
		{"no low-risk ladder", strings.Replace(tierWF, "low_risk_models = [\"sonnet\"]\n", "", 1), []string{"internal/x/x.go"}, nil, "opus"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, e := runTier(t, c.src, func(r *fakeRunner) { r.changed, r.changedErr = c.files, c.err })
			if got := verifierModels(r); len(got) != 1 || got[0] != c.want {
				t.Fatalf("verifier models = %q, want [%s]", got, c.want)
			}
			j, err := e.Store.Read(runID)
			if err != nil {
				t.Fatal(err)
			}
			for _, en := range j.Entries {
				if en.Type == TypeLeafSession && en.Role == workflow.RoleVerifier && en.Model != c.want {
					t.Errorf("journalled model %q, want %s", en.Model, c.want)
				}
			}
			if len(r.changedRanges) > 0 && r.changedRanges[0] != fakeHead+"..C-leaf-a" {
				t.Errorf("diffed %q", r.changedRanges)
			}
		})
	}
}

// A second attempt verifies on models even when the diff is low-risk, and
// the journal holds each attempt's choice.
func TestVerifierTierRetry(t *testing.T) {
	r, e := runTier(t, tierWF, func(r *fakeRunner) {
		r.changed = []string{"internal/x/x.go"}
		r.outcome[leafSID("leaf-a", "verifier", 1)] = workflow.OutcomeFail
		r.outcome[leafSID("leaf-a", "builder", 2)] = workflow.OutcomePass
		r.outcome[leafSID("leaf-a", "verifier", 2)] = workflow.OutcomePass
	})
	if got := verifierModels(r); len(got) != 2 || got[0] != "sonnet" || got[1] != "opus" {
		t.Errorf("verifier models = %q, want [sonnet opus]", got)
	}
	if got, want := journalledVerifierModels(t, e), []string{"leaf-a 1 sonnet", "leaf-a 2 opus"}; !reflect.DeepEqual(got, want) {
		t.Errorf("journalled verifier models = %q, want %q", got, want)
	}
}

// A verifier in a fresh round (a retry after the leaf failed every
// attempt of round 1) starts on models even for a low-risk first attempt
// of the round.
func TestVerifierTierFreshRound(t *testing.T) {
	e, r, _, _ := setup(t, leafFailsWF+"models = [\"opus\"]\nlow_risk_models = [\"sonnet\"]\nrisky_paths = [\".github/**\"]\n")
	r.changed = []string{"internal/x/x.go"}
	if st := leafFailsRun(t, e, r, 0); st.Result != Failed {
		t.Fatalf("first run %q", st.Result)
	}
	for _, role := range []string{workflow.RoleBuilder, workflow.RoleVerifier} {
		r.outcome[leafSID("leaf-a", role, 4)] = workflow.OutcomePass
	}
	r.leafAnswers["QL1"] = [2]string{"use-a-map", "Use a map, not a list."}
	before := len(r.leafSpecs)
	st, err := e.Retry(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded {
		t.Fatalf("retry %q: %s", st.Result, st.Detail)
	}
	var v LeafSpec
	for _, s := range r.leafSpecs[before:] {
		if s.Leaf.Handle == "leaf-a" && s.Role == workflow.RoleVerifier {
			v = s
		}
	}
	if v.Round != 2 || v.RoundAttempt != 1 || v.Model != "opus" {
		t.Errorf("fresh round verifier spec round %d attempt %d model %q, want round 2 attempt 1 opus", v.Round, v.RoundAttempt, v.Model)
	}
	// leaf-b and leaf-c are low-risk first attempts of round 1: sonnet.
	want := []string{"leaf-b 1 sonnet", "leaf-a 4 opus", "leaf-c 1 sonnet"}
	if got := journalledVerifierModels(t, e); !reflect.DeepEqual(got, want) {
		t.Errorf("journalled verifier models = %q, want %q", got, want)
	}
}

// A checks step's verifier in review cycle 2 starts on models whatever the
// diff, and the journal holds that choice.
func TestVerifierTierReviewCycle(t *testing.T) {
	e, r, _, _ := setup(t, checksWF+"models = [\"opus\"]\nlow_risk_models = [\"sonnet\"]\nrisky_paths = [\".github/**\"]\n")
	r.leaves = fixLeaves
	a := fixLeaves[0]
	r.changed = []string{"internal/x/x.go"}
	passChecks(r, 1, threeChecks...)
	passChecks(r, 2, "alpha")
	r.findings[checkSID("alpha", "check", 1)] = []Finding{onLeaf(sev("f1", SeverityMajor), a)}
	r.outcome[checkSID("alpha", "critic", 1)] = workflow.OutcomePass
	passFixes(r, 2, a)
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if got := outcomes(st); !reflect.DeepEqual(got, []string{"review:pass"}) {
		t.Fatalf("outcomes = %v", got)
	}
	var cycles []int
	for _, s := range r.leafSpecs {
		if s.Role == workflow.RoleVerifier {
			cycles = append(cycles, s.Cycle)
		}
	}
	if !reflect.DeepEqual(cycles, []int{2}) || !reflect.DeepEqual(verifierModels(r), []string{"opus"}) {
		t.Errorf("verifier cycles %v models %q, want [2] [opus]", cycles, verifierModels(r))
	}
	if got, want := journalledVerifierModels(t, e), []string{"leaf-a 2 opus"}; !reflect.DeepEqual(got, want) {
		t.Errorf("journalled verifier models = %q, want %q", got, want)
	}
}

// A fresh round and a review cycle past 1 are risky whatever the diff.
func TestVerifierRiskyRetryStates(t *testing.T) {
	role := workflow.Role{Models: []string{"opus"}, LowRiskModels: []string{"sonnet"}, RiskyPaths: []string{".github/**"}}
	r := newFakeRunner()
	r.changed = []string{"internal/x/x.go"}
	l := LeafState{Sessions: []LeafSession{{Role: workflow.RoleBuilder, Attempt: 1, Commit: "H0"}}}
	a := leafAction{attempt: 1, commit: "C1", cycle: 0}
	spec := LeafSpec{Round: 1, RoundAttempt: 1}
	if verifierRisky(r, role, l, spec, a) {
		t.Error("first attempt, low-risk diff is risky")
	}
	if !verifierRisky(r, role, l, LeafSpec{Round: 2, RoundAttempt: 1}, a) {
		t.Error("fresh round is not risky")
	}
	a.cycle = 2
	if !verifierRisky(r, role, l, spec, a) {
		t.Error("review cycle 2 is not risky")
	}
}
