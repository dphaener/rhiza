package api

import (
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/check"
)

// codeCheck names rhiza check's failures (not its warnings), then each
// leaf build-check found not OK; it passes with neither.
func TestCodeCheckCauses(t *testing.T) {
	cr := CheckResponse{OK: false, Findings: []CheckFinding{
		{Severity: check.SeverityWarn, Rule: "evidence-actor", Message: "just a warning"},
		{Severity: check.SeverityFail, Rule: "render", Path: ".rhiza/rendered/index.md", Message: "differs from what render writes"},
	}}
	br := RunBuildCheckResponse{Run: "R", Change: "C", Leaves: []RunBuildLeaf{
		{Handle: "alpha", OK: true}, {Handle: "beta", Why: "no Evidence verifies it"},
	}}
	got := codeCheck(cr, br)
	want := []string{
		"rhiza check: fail render .rhiza/rendered/index.md: differs from what render writes",
		"leaf beta: no Evidence verifies it",
	}
	if got.Passed || !slices.Equal(got.Causes, want) || got.ExitCode() != 1 {
		t.Errorf("codeCheck = %+v, want causes %q", got, want)
	}
	ok := codeCheck(CheckResponse{OK: true, Findings: cr.Findings[:1]}, RunBuildCheckResponse{Leaves: br.Leaves[:1]})
	if !ok.Passed || len(ok.Causes) != 0 || ok.ExitCode() != 0 {
		t.Errorf("codeCheck = %+v, want a pass", ok)
	}
}
