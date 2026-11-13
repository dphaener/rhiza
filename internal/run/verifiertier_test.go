package run

import (
	"errors"
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

// A second attempt verifies on models even when the diff is low-risk.
func TestVerifierTierRetry(t *testing.T) {
	r, _ := runTier(t, tierWF, func(r *fakeRunner) {
		r.changed = []string{"internal/x/x.go"}
		r.outcome[leafSID("leaf-a", "verifier", 1)] = workflow.OutcomeFail
		r.outcome[leafSID("leaf-a", "builder", 2)] = workflow.OutcomePass
		r.outcome[leafSID("leaf-a", "verifier", 2)] = workflow.OutcomePass
	})
	if got := verifierModels(r); len(got) != 2 || got[0] != "sonnet" || got[1] != "opus" {
		t.Errorf("verifier models = %q, want [sonnet opus]", got)
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
