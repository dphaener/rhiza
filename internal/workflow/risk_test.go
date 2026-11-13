package workflow

import (
	"reflect"
	"strings"
	"testing"
)

const tierVerifier = "[step.build.verifier]\nprompt = \"verify {leaf}\"\nmodels = [\"opus\"]\nlow_risk_models = [\"sonnet\"]\nrisky_paths = [\"**/testdata/**\", \"*_test.go\", \".github/**\"]\n"

func TestVerifierTierLoads(t *testing.T) {
	src := leavesBase + "[step.build.builder]\nprompt = \"build {leaf}\"\n" + tierVerifier
	d, err := Parse("w", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	v := d.Steps["build"].Verifier
	if !reflect.DeepEqual(v.LowRiskModels, []string{"sonnet"}) || len(v.RiskyPaths) != 3 {
		t.Errorf("verifier = %+v", v)
	}
	// a checks step's verifier takes them too
	src = checksBase + "[step.review.critic]\nprompt = \"c\"\n[step.review.fixer]\nprompt = \"f\"\n" +
		strings.ReplaceAll(tierVerifier, "step.build", "step.review")
	if _, err := Parse("w", []byte(src)); err != nil {
		t.Errorf("checks verifier: %v", err)
	}
	// a role setting neither has none
	d, err = Parse("w", []byte(leavesBase+leavesRoles))
	if err != nil {
		t.Fatal(err)
	}
	if v := d.Steps["build"].Verifier; v.LowRiskModels != nil || v.RiskyPaths != nil {
		t.Errorf("untiered verifier = %+v", v)
	}
}

func TestVerifierTierRefused(t *testing.T) {
	verifier := "[step.build.verifier]\nprompt = \"v\"\n"
	for _, tc := range []struct{ name, src, want string }{
		{"builder", leavesBase + "[step.build.builder]\nprompt = \"b\"\nrisky_paths = [\"a\"]\n" + verifier, "only for the verifier role"},
		{"builder-low", leavesBase + "[step.build.builder]\nprompt = \"b\"\nlow_risk_models = [\"sonnet\"]\n" + verifier, "only for the verifier role"},
		{"fixer", checksBase + "[step.review.critic]\nprompt = \"c\"\n[step.review.fixer]\nprompt = \"f\"\nlow_risk_models = [\"sonnet\"]\n[step.review.verifier]\nprompt = \"v\"\n", "only for the verifier role"},
		{"critic", checksBase + "[step.review.critic]\nprompt = \"c\"\nrisky_paths = [\"a\"]\n[step.review.fixer]\nprompt = \"f\"\n[step.review.verifier]\nprompt = \"v\"\n", "only for the verifier role"},
		{"check", strings.Replace(checksBase, "prompt = \"check {check}\"\n", "prompt = \"check {check}\"\nrisky_paths = [\"a\"]\n", 1) + checksRoles, "only for the verifier role; a check cannot have them"},
		{"check-low", strings.Replace(checksBase, "prompt = \"check {check}\"\n", "prompt = \"check {check}\"\nlow_risk_models = [\"sonnet\"]\n", 1) + checksRoles, "only for the verifier role; a check cannot have them"},
		{"agent step", "[step.a]\nkind = \"agent\"\nprompt = \"go\"\nlow_risk_models = [\"sonnet\"]\n", "belong to a leaves or checks step's verifier role; an agent step has no verifier"},
		{"agent step paths", "[step.a]\nkind = \"agent\"\nprompt = \"go\"\nrisky_paths = [\"a\"]\n", "an agent step has no verifier"},
		{"leaves step", leavesBase + "risky_paths = [\"a\"]\n[step.build.builder]\nprompt = \"b\"\n" + verifier, "set them under [step.build.verifier]"},
		{"blank model", leavesBase + "[step.build.builder]\nprompt = \"b\"\n[step.build.verifier]\nprompt = \"v\"\nlow_risk_models = [\"sonnet\", \" \"]\n", "low_risk_models entry 2 is blank"},
		{"empty glob", leavesBase + "[step.build.builder]\nprompt = \"b\"\n[step.build.verifier]\nprompt = \"v\"\nrisky_paths = [\"\"]\n", "risky_paths entry 1 is empty"},
		{"bad glob", leavesBase + "[step.build.builder]\nprompt = \"b\"\n[step.build.verifier]\nprompt = \"v\"\nrisky_paths = [\"[\"]\n", "not a valid glob"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("w", []byte(tc.src))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %v, want one containing %q", err, tc.want)
			}
		})
	}
}

func TestRisky(t *testing.T) {
	r := Role{RiskyPaths: []string{"**/testdata/**", "*_test.go", ".github/**", "internal/*/x.go"}}
	for file, want := range map[string]bool{
		"internal/a/testdata/b.txt": true,
		"testdata/b.txt":            true,
		"main_test.go":              true,
		"internal/run/run_test.go":  true,
		".github/workflows/ci.yml":  true,
		"internal/a/x.go":           true,
		"internal/a/b/x.go":         false,
		"internal/run/leaves.go":    false,
		"docs/github/ci.yml":        false,
		"internal/x/x.go":           true,
	} {
		if got := r.Risky(file); got != want {
			t.Errorf("Risky(%q) = %v, want %v", file, got, want)
		}
	}
	if (Role{}).Risky("a.go") {
		t.Error("a role without globs is never risky")
	}
}
