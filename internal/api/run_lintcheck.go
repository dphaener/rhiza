package api

// This file holds run.lint-check: it reads the model lint hits the run's
// latest code-check entry journalled and fails when there are any, so a
// review can route them to a fixer before its checks run. It runs no lint
// itself and writes nothing.

import (
	"fmt"
	"io"
	"reflect"
)

func init() {
	register(Operation{
		Name:     "run.lint-check",
		Summary:  "Pass when the latest code-check entry of the run RHIZA_RUN names has no model lint hits, else fail listing each hit",
		Request:  reflect.TypeFor[RunLintCheckRequest](),
		Response: reflect.TypeFor[RunLintCheckResponse](),
		Run:      runRunLintCheck,
	})
}

const runLintCheckNext = "rhiza describe run.lint-check"

// RunLintCheckRequest takes nothing: the run is RHIZA_RUN's.
type RunLintCheckRequest struct{}

// RunLintCheckResponse lists the hits of the latest code check. It is
// printed whether or not there are any; hits exit 1.
type RunLintCheckResponse struct {
	Run    string       `json:"run" desc:"The run's full ID"`
	Change string       `json:"change" desc:"The run's change set"`
	Passed bool         `json:"passed" desc:"True if the latest code-check entry has no lint hits, or the run has none"`
	Lints  []RunLintHit `json:"lints" desc:"The lint hits of the run's latest code-check entry: node, lint kind and the reference or Decision; empty when it has none"`
}

// ExitCode is 0 with no hits, else 1.
func (r RunLintCheckResponse) ExitCode() int {
	if r.Passed {
		return 0
	}
	return 1
}

func runRunLintCheck(ctx *Ctx, _ any) (any, error) {
	rr, err := openReviewRun(ctx, "run lint-check", runLintCheckNext)
	if err != nil {
		return nil, err
	}
	resp := RunLintCheckResponse{Run: rr.id, Change: rr.state.Change, Lints: []RunLintHit{}}
	if cc, ok := rr.state.LatestCodeCheck(); ok {
		for _, h := range cc.Lints {
			resp.Lints = append(resp.Lints, RunLintHit{Node: h.Node, Handle: h.Handle, Field: h.Field, Kind: h.Kind, Ref: h.Ref})
		}
	}
	resp.Passed = len(resp.Lints) == 0
	return resp, nil
}

// Human reports the check: one line when clean, else each hit.
func (r RunLintCheckResponse) Human(w io.Writer) {
	if r.Passed {
		fmt.Fprintln(w, "lint-check passed: the code check found no model lint hits")
		return
	}
	fmt.Fprintf(w, "lint-check failed: %d %s\n\n", len(r.Lints), shipPlural(len(r.Lints), "hit", "hits"))
	for _, h := range r.Lints {
		fmt.Fprintf(w, "- %s\n", h.line())
	}
}
