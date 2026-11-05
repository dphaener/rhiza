package api

// This file holds the plan review ops of a run: run.findings journals a
// reviewer's findings, opening the next review round, run.verdict
// journals a critic's verdicts on the latest round's findings,
// run.review-check decides whether the review may go on to the gate, and
// run.gate ends the plan: it proposes the plan's change sets and, in an
// autonomous run whose review passed, accepts them. All work only inside
// the run RHIZA_RUN names, which must be unfinished. The first three check
// everything before they write anything. The round count comes from the
// journal: each run.findings call is one round.

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "run.findings",
		Summary:  "Journal a plan reviewer's findings as the next review round of the run RHIZA_RUN names, with --check a review check's findings as that check's next round, or with --leaf a verifier's findings on a leaf's current attempt",
		Write:    true,
		Request:  reflect.TypeFor[RunFindingsRequest](),
		Response: reflect.TypeFor[RunFindingsResponse](),
		Input: &Input{
			Summary: `{"findings":[{"node":...,"kind":...,"reason":...},...]}, maybe empty for a clean review; kind is one of ` + strings.Join(run.Kinds, ", ") +
				`. With --check: {"findings":[{"severity":...,"leaf":...,"place":...,"reason":...},...]}, maybe empty; severity is one of ` + strings.Join(run.Severities, ", ") +
				`, leaf optional. With --leaf: {"findings":[{"reason":...},...]}, at least one`,
			Schema: findingsFileSchema,
		},
		Run: runRunFindings,
	})
	register(Operation{
		Name:     "run.verdict",
		Summary:  "Journal a critic's verdicts on the findings of the latest review round of the run RHIZA_RUN names, or with --check of that review check's latest round",
		Write:    true,
		Request:  reflect.TypeFor[RunVerdictRequest](),
		Response: reflect.TypeFor[RunVerdictResponse](),
		Input: &Input{
			Summary: `{"verdicts":[{"finding":"f1","verdict":"confirmed"|"dropped","reason":...},...]}, at least one; each names a finding of the latest round (with --check, of the check's latest round)`,
			Schema:  verdictFileSchema,
		},
		Run: runRunVerdict,
	})
	register(Operation{
		Name:     "run.review-check",
		Summary:  "Check the latest review round of the run RHIZA_RUN names: pass with no confirmed finding, fail with one before round 3, else leave each as a question and pass",
		Write:    true,
		Request:  reflect.TypeFor[RunReviewCheckRequest](),
		Response: reflect.TypeFor[RunReviewCheckResponse](),
		Run:      runRunReviewCheck,
	})
	register(Operation{
		Name:     "run.gate",
		Summary:  "End the plan of the run RHIZA_RUN names: propose each draft among its change sets; an autonomous run whose review passed accepts each under rule autonomous",
		Write:    true,
		Request:  reflect.TypeFor[RunGateRequest](),
		Response: reflect.TypeFor[RunGateResponse](),
		Run:      runRunGate,
	})
}

// RunFindingsRequest names the file of findings.
type RunFindingsRequest struct {
	File  string `json:"file" arg:"0" desc:"A JSON file of findings, or - for stdin (see the input schema)"`
	Leaf  string `json:"leaf,omitempty" flag:"leaf" desc:"A leaf of the run's change set, by handle or ID: journal a verifier's findings, each just a reason, on the leaf's current attempt in the run's leaves step instead of a plan review round"`
	Check string `json:"check,omitempty" flag:"check" desc:"A check of the checks step the run is in: journal that review check's findings, each with a severity, place and reason, as its next round; required in a review run, refused outside one"`
}

// RunVerdictRequest names the file of verdicts.
type RunVerdictRequest struct {
	File  string `json:"file" arg:"0" desc:"A JSON file of verdicts, or - for stdin (see the input schema)"`
	Check string `json:"check,omitempty" flag:"check" desc:"A check of the checks step the run is in: judge that review check's latest round; required in a review run, refused outside one"`
}

// RunFinding is one finding of a review round, with its latest verdict
// once a critic gave one.
type RunFinding struct {
	ID            string `json:"id" desc:"The finding's ID in its round: f1, f2, ..."`
	Node          string `json:"node,omitempty" desc:"The node it is about, by ID; absent for a leaf's finding"`
	Handle        string `json:"handle,omitempty" desc:"The node's handle when the finding was recorded; absent when it had none"`
	Kind          string `json:"kind,omitempty" desc:"What is wrong: missing-acceptance, missing-provenance, unreal-depends-on, outside-intent or contradiction; absent for a leaf's or a review check's finding"`
	Check         string `json:"check,omitempty" desc:"A review check's finding: the check that found it"`
	Severity      string `json:"severity,omitempty" desc:"A review check's finding: blocker, major, minor or nit"`
	Leaf          string `json:"leaf,omitempty" desc:"A review check's finding: the leaf it is about, by ID; absent when it names none"`
	LeafHandle    string `json:"leaf_handle,omitempty" desc:"That leaf's handle when the finding was recorded"`
	Place         string `json:"place,omitempty" desc:"A review check's finding: where, as file:line or a node"`
	Reason        string `json:"reason" desc:"Why the reviewer found it"`
	Verdict       string `json:"verdict,omitempty" desc:"The critic's latest verdict, confirmed or dropped; absent before one"`
	VerdictReason string `json:"verdict_reason,omitempty" desc:"Why the critic gave it; absent before a verdict"`
}

// RunVerdict is one verdict a critic gave.
type RunVerdict struct {
	Finding string `json:"finding" desc:"The finding's ID in the round"`
	Verdict string `json:"verdict" desc:"confirmed or dropped"`
	Reason  string `json:"reason" desc:"Why"`
}

// RunFindingsResponse is the round the findings opened, or the leaf
// attempt they were journalled on.
type RunFindingsResponse struct {
	Run      string       `json:"run" desc:"The run's full ID"`
	Round    int          `json:"round,omitempty" desc:"The round the findings opened, from 1 (with --check, in the check's rounds); absent for a leaf's findings"`
	Check    string       `json:"check,omitempty" desc:"With --check: the check"`
	Cycle    int          `json:"cycle,omitempty" desc:"With --check: the check's cycle the round is in, from 1"`
	Leaf     string       `json:"leaf,omitempty" desc:"With --leaf: the leaf's handle"`
	Attempt  int          `json:"attempt,omitempty" desc:"With --leaf: the leaf's attempt the findings are on, from 1"`
	Findings []RunFinding `json:"findings" desc:"The findings recorded, in order; empty for a clean review"`
	DryRun   bool         `json:"dry_run,omitempty" desc:"True if nothing was journalled"`
}

// RunVerdictResponse is the verdicts recorded on a round.
type RunVerdictResponse struct {
	Run      string       `json:"run" desc:"The run's full ID"`
	Round    int          `json:"round" desc:"The round judged: the latest (with --check, the check's latest)"`
	Check    string       `json:"check,omitempty" desc:"With --check: the check"`
	Verdicts []RunVerdict `json:"verdicts" desc:"The verdicts recorded, in order"`
	DryRun   bool         `json:"dry_run,omitempty" desc:"True if nothing was journalled"`
}

// reviewRun is the run RHIZA_RUN names, unfinished, with its store,
// journal folded and change set: what every review op works in.
type reviewRun struct {
	store  *run.Store
	id     string
	state  run.State
	change id.ID
}

// openReviewRun finds the run RHIZA_RUN names and refuses, as usage,
// unless it is an unfinished run. cmd names the command for messages,
// e.g. "run findings"; next is where a caller learns more.
func openReviewRun(ctx *Ctx, cmd, next string) (reviewRun, error) {
	return openRun(ctx, cmd, next, false)
}

// openRun is openReviewRun; with finished it takes a finished run too.
func openRun(ctx *Ctx, cmd, next string, finished bool) (reviewRun, error) {
	ref := strings.TrimSpace(ctxGetenv(ctx)(envRun))
	if ref == "" {
		return reviewRun{}, Usagef("%s works inside a run, and %s is not set: run it from a step of the run", cmd, envRun).WithNext(next)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return reviewRun{}, err
	}
	rs := run.NewStore(run.Dir(common))
	j, err := rs.Read(ref)
	var nf *run.NotFoundError
	var amb *run.AmbiguousError
	switch {
	case errors.As(err, &nf), errors.As(err, &amb):
		return reviewRun{}, Usagef("%s works inside a run, and %s=%s names no run: %v", cmd, envRun, ref, err).WithNext(next)
	case err != nil:
		return reviewRun{}, Errorf("read run %s: %v", ref, err)
	}
	st := j.State()
	if st.Ended && !finished {
		return reviewRun{}, Usagef("%s works inside an unfinished run, and run %s is finished: it %s", cmd, j.ID, st.Result).
			WithNext("rhiza run status " + j.ID)
	}
	cs, err := id.Parse(st.Change)
	if err != nil {
		return reviewRun{}, Errorf("run %s names change set %q: %v", j.ID, st.Change, err)
	}
	return reviewRun{store: rs, id: j.ID, state: st, change: cs}, nil
}

// reviewStep is the checks step the run is in, if it is in one: then it
// is a review run, and its findings and verdicts belong to a check.
func (r reviewRun) reviewStep() (workflow.Step, bool, error) {
	if r.state.Current == "" {
		return workflow.Step{}, false, nil
	}
	def, err := r.state.ParseDefinition()
	if err != nil {
		return workflow.Step{}, false, Errorf("run %s: its workflow no longer parses: %v", r.id, err)
	}
	s, ok := def.Steps[r.state.Current]
	if !ok || s.Kind != workflow.KindChecks {
		return workflow.Step{}, false, nil
	}
	return s, true, nil
}

// reviewCheck checks the --check of cmd against the run: in a review run
// it must name a check of its step, and outside one it must be empty.
// review is false outside a review run.
func (r reviewRun) reviewCheck(cmd, check, next string) (step workflow.Step, review bool, err error) {
	step, review, err = r.reviewStep()
	if err != nil {
		return workflow.Step{}, false, err
	}
	check = strings.TrimSpace(check)
	switch {
	case !review && check != "":
		return workflow.Step{}, false, Usagef("%s --check: run %s is not in a checks step: --check works only in a review run", cmd, r.id).WithNext(next)
	case !review:
		return workflow.Step{}, false, nil
	case check == "":
		return workflow.Step{}, false, Usagef("%s works in a review run with --check: run %s is in checks step %s, so name one of its checks: %s",
			cmd, r.id, step.Name, strings.Join(checkNames(step), ", ")).WithNext(next)
	}
	if _, ok := step.Check(check); !ok {
		return workflow.Step{}, false, Usagef("%s: check %q is not a check of step %s: its checks are %s",
			cmd, check, step.Name, strings.Join(checkNames(step), ", ")).WithNext(next)
	}
	return step, true, nil
}

// checkNames are step's check names, in order.
func checkNames(step workflow.Step) []string {
	out := make([]string, 0, len(step.Checks))
	for _, c := range step.Checks {
		out = append(out, c.Name)
	}
	return out
}

// latestRound is the run's latest review round; with none it is usage.
func (r reviewRun) latestRound(cmd, next string) (run.Round, error) {
	rd, ok := r.state.LatestRound()
	if !ok {
		return run.Round{}, Usagef("%s: run %s has no review round yet: rhiza run findings opens one", cmd, r.id).WithNext(next)
	}
	return rd, nil
}

// journal appends e to the run's journal, unless ctx is a dry run.
func (r reviewRun) journal(ctx *Ctx, what string, e run.Entry) error {
	if ctx.DryRun {
		return nil
	}
	if err := r.store.Append(r.id, e); err != nil {
		return Errorf("journal the %s in run %s: %v", what, r.id, err).WithNext("rhiza run status " + r.id)
	}
	return nil
}

// lockJournal takes the run's journal lock and re-reads its state, so a
// caller numbers findings from the journal as it is when it appends. A dry
// run takes no lock and keeps the state it has. unlock is always callable.
func (r reviewRun) lockJournal(ctx *Ctx) (reviewRun, func(), error) {
	if ctx.DryRun {
		return r, func() {}, nil
	}
	unlock, err := r.store.LockJournal(r.id)
	if err != nil {
		return r, func() {}, Errorf("lock run %s's journal: %v", r.id, err).WithNext("rhiza run status " + r.id)
	}
	j, err := r.store.Read(r.id)
	if err != nil {
		unlock()
		return r, func() {}, Errorf("read run %s: %v", r.id, err)
	}
	r.state = j.State()
	return r, unlock, nil
}

// readReviewFile reads file, or stdin for "-", and decodes it strictly
// into v: one JSON object with no unknown fields.
func readReviewFile(ctx *Ctx, cmd, next, file string, v any) error {
	var data []byte
	var err error
	switch file = strings.TrimSpace(file); file {
	case "":
		return Usagef("%s needs a file, or - for stdin", cmd).WithNext(next)
	case "-":
		if ctx.Stdin == nil {
			return Usagef("%s -: no stdin", cmd).WithNext(next)
		}
		if data, err = io.ReadAll(ctx.Stdin); err != nil {
			return Usagef("%s -: cannot read stdin: %v", cmd, err).WithNext(next)
		}
	default:
		if data, err = os.ReadFile(file); err != nil {
			return Usagef("%s: cannot read %s: %v", cmd, file, errors.Unwrap(err)).WithNext(next)
		}
	}
	dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(data)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return Usagef("%s file: %v", cmd, err).WithNext(next)
	}
	if dec.More() {
		return Usagef("%s file: more than one JSON value", cmd).WithNext(next)
	}
	return nil
}

const (
	runFindingsNext = "rhiza describe run.findings"
	runVerdictNext  = "rhiza describe run.verdict"
)

// findingsFile is run.findings' input.
type findingsFile struct {
	Findings *[]struct {
		Node   string `json:"node"`
		Kind   string `json:"kind"`
		Reason string `json:"reason"`
	} `json:"findings"`
}

// verdictFile is run.verdict's input.
type verdictFile struct {
	Verdicts []run.Verdict `json:"verdicts"`
}

// runRunFindings checks the run, then each finding: a known kind, a
// reason, and a node that resolves in the run's change set's view or,
// failing that, in its plan's change sets (planResolver). Only
// then does it journal them, numbered f1, f2, ..., as the next round.
func runRunFindings(ctx *Ctx, req any) (any, error) {
	r := req.(RunFindingsRequest)
	const cmd = "run findings"
	rr, err := openReviewRun(ctx, cmd, runFindingsNext)
	if err != nil {
		return nil, err
	}
	if r.Leaf != "" {
		if strings.TrimSpace(r.Check) != "" {
			return nil, Usagef("%s: give --leaf or --check, not both", cmd).WithNext(runFindingsNext)
		}
		return runLeafFindings(ctx, rr, r)
	}
	step, review, err := rr.reviewCheck(cmd, r.Check, runFindingsNext)
	if err != nil {
		return nil, err
	}
	if review {
		return runCheckFindings(ctx, rr, step, strings.TrimSpace(r.Check), r.File)
	}
	var f findingsFile
	if err := readReviewFile(ctx, cmd, runFindingsNext, r.File, &f); err != nil {
		return nil, err
	}
	if f.Findings == nil {
		return nil, Usagef(`%s file: want {"findings":[...]}; a clean review is {"findings":[]}`, cmd).WithNext(runFindingsNext)
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	r2 := planResolver{w: w, run: rr.change}
	found := make([]run.Finding, 0, len(*f.Findings))
	for i, in := range *f.Findings {
		n := i + 1
		if !run.ValidKind(in.Kind) {
			return nil, Usagef("finding %d: kind %q is not one of %s", n, in.Kind, strings.Join(run.Kinds, ", ")).WithNext(runFindingsNext)
		}
		if strings.TrimSpace(in.Reason) == "" {
			return nil, Usagef("finding %d: no reason", n).WithNext(runFindingsNext)
		}
		ref := strings.TrimSpace(in.Node)
		if ref == "" {
			return nil, Usagef("finding %d: no node", n).WithNext(runFindingsNext)
		}
		node, err := r2.resolve(ref)
		if err != nil {
			return nil, Usagef("finding %d: node %q: %s", n, ref, AsError(err).Message).
				WithNext("rhiza ls --change " + rr.change.String())
		}
		found = append(found, run.Finding{
			ID: "f" + strconv.Itoa(n), Node: node.ID.String(), Handle: node.Handle, Kind: in.Kind, Reason: in.Reason,
		})
	}
	rr, unlock, err := rr.lockJournal(ctx)
	if err != nil {
		return nil, err
	}
	defer unlock()
	resp := RunFindingsResponse{Run: rr.id, Round: len(rr.state.Rounds) + 1, Findings: make([]RunFinding, 0, len(found)), DryRun: ctx.DryRun}
	for _, fd := range found {
		resp.Findings = append(resp.Findings, runFinding(fd, run.Round{}))
	}
	e := run.Entry{Type: run.TypeFindings, At: ctx.Now(), Findings: found}
	if err := rr.journal(ctx, "findings of round "+strconv.Itoa(resp.Round), e); err != nil {
		return nil, err
	}
	return resp, nil
}

// checkFindingsFile is run.findings' input with --check.
type checkFindingsFile struct {
	Findings *[]struct {
		Severity string `json:"severity"`
		Leaf     string `json:"leaf"`
		Place    string `json:"place"`
		Reason   string `json:"reason"`
	} `json:"findings"`
}

// runCheckFindings is run.findings with --check in a review run: it checks
// each finding, a known severity, a place, a reason and, when it names
// one, a leaf of the run's change set. Only then does it journal them,
// numbered f1, f2, ..., as check's next round, in the cycle RHIZA_CYCLE
// names when it is one of the step's that runs check, else in the step's
// latest cycle when that cycle runs check (its session may write before the
// engine journals it), else the cycle of the check's latest session, else
// the step's latest cycle. An empty list is a round with no findings.
func runCheckFindings(ctx *Ctx, rr reviewRun, step workflow.Step, check, file string) (any, error) {
	const cmd = "run findings --check"
	var f checkFindingsFile
	if err := readReviewFile(ctx, cmd, runFindingsNext, file, &f); err != nil {
		return nil, err
	}
	if f.Findings == nil {
		return nil, Usagef(`%s file: want {"findings":[...]}; a clean check is {"findings":[]}`, cmd).WithNext(runFindingsNext)
	}
	leafOf := changeLeafResolver(ctx, rr)
	found := make([]run.Finding, 0, len(*f.Findings))
	for i, in := range *f.Findings {
		n := i + 1
		if !run.ValidSeverity(in.Severity) {
			return nil, Usagef("finding %d: severity %q is not one of %s", n, in.Severity, strings.Join(run.Severities, ", ")).WithNext(runFindingsNext)
		}
		if strings.TrimSpace(in.Place) == "" {
			return nil, Usagef("finding %d: no place: give file:line or a node", n).WithNext(runFindingsNext)
		}
		if strings.TrimSpace(in.Reason) == "" {
			return nil, Usagef("finding %d: no reason", n).WithNext(runFindingsNext)
		}
		fd := run.Finding{ID: "f" + strconv.Itoa(n), Check: check, Severity: in.Severity, Place: in.Place, Reason: in.Reason}
		if ref := strings.TrimSpace(in.Leaf); ref != "" {
			node, err := leafOf(ref)
			if err != nil {
				if e := AsError(err); e.Code != CodeUsage {
					return nil, err
				}
				return nil, Usagef("finding %d: leaf %q: %s", n, ref, AsError(err).Message).WithNext("rhiza change show " + rr.change.String())
			}
			fd.Leaf, fd.LeafHandle = node.ID.String(), node.Handle
		}
		found = append(found, fd)
	}
	rr, unlock, err := rr.lockJournal(ctx)
	if err != nil {
		return nil, err
	}
	defer unlock()
	cs, _ := rr.state.Check(step.Name, check)
	cycle := envCheckCycle(ctx, rr.state, step.Name, check)
	cy, latest := rr.state.LatestCycle(step.Name)
	switch n := len(cs.Sessions); {
	case cycle > 0:
	case n > 0 && !(latest && slices.Contains(cy.Checks, check)):
		cycle = cs.Sessions[n-1].Cycle
	case latest:
		cycle = cy.Cycle
	}
	resp := RunFindingsResponse{
		Run: rr.id, Round: len(cs.Rounds) + 1, Check: check, Cycle: cycle, Findings: make([]RunFinding, 0, len(found)), DryRun: ctx.DryRun,
	}
	for _, fd := range found {
		resp.Findings = append(resp.Findings, runFinding(fd, run.Round{}))
	}
	e := run.Entry{Type: run.TypeFindings, At: ctx.Now(), Step: step.Name, Check: check, Cycle: cycle, Findings: found}
	if err := rr.journal(ctx, fmt.Sprintf("findings of check %s, round %d,", check, resp.Round), e); err != nil {
		return nil, err
	}
	return resp, nil
}

// changeLeafResolver resolves a finding's leaf ref among the leaves of
// rr's change set, reading the model only on its first call. A ref that
// names nothing, or a node that is not one of the leaves, is usage.
func changeLeafResolver(ctx *Ctx, rr reviewRun) func(ref string) (Node, error) {
	var m Model
	var leaves []ChangeLeaf
	return func(ref string) (Node, error) {
		if leaves == nil {
			w, err := openCtxWorkspace(ctx)
			if err != nil {
				return Node{}, err
			}
			c, ok := w.changeSet(rr.change)
			if !ok {
				return Node{}, Errorf("run %s's change set %s is not in the log", rr.id, rr.change)
			}
			m = w.model(rr.change)
			leaves = changeLeaves(m, c.Ops)
		}
		node, err := m.Resolve(ref)
		if err != nil {
			return Node{}, Usagef("%s", AsError(err).Message)
		}
		if !slices.ContainsFunc(leaves, func(l ChangeLeaf) bool { return l.ID == node.ID }) {
			return Node{}, Usagef("%s is not one of change set %s's leaves", ref, rr.change)
		}
		return node, nil
	}
}

// envCheckCycle is the cycle the caller's RHIZA_CYCLE names, which the
// engine sets on every check and critic session, when it is a cycle of
// step that runs check; else 0, and the cycle is inferred.
func envCheckCycle(ctx *Ctx, st run.State, step, check string) int {
	n, err := strconv.Atoi(strings.TrimSpace(ctxGetenv(ctx)(envCycle)))
	if err != nil || n < 1 {
		return 0
	}
	for _, cy := range st.Cycles {
		if cy.Step == step && cy.Cycle == n && slices.Contains(cy.Checks, check) {
			return n
		}
	}
	return 0
}

// leafFindingsFile is run.findings' input with --leaf: each finding is
// just a reason.
type leafFindingsFile struct {
	Findings []struct {
		Reason string `json:"reason"`
	} `json:"findings"`
}

// runLeafFindings is run.findings with --leaf: it checks the leaf is one
// of the run change set's leaves, open with a current attempt in the leaves
// step the run is in, then the findings, at least one, each with a
// reason. Only then does it journal one findings entry carrying the leaf
// and attempt, its findings numbered on from those the attempt has.
func runLeafFindings(ctx *Ctx, rr reviewRun, r RunFindingsRequest) (any, error) {
	const cmd = "run findings --leaf"
	statusNext := "rhiza run status " + rr.id
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	c, ok := w.changeSet(rr.change)
	if !ok {
		return nil, Errorf("run %s's change set %s is not in the log", rr.id, rr.change)
	}
	m := w.model(rr.change)
	ref := strings.TrimSpace(r.Leaf)
	node, err := m.Resolve(ref)
	if err == nil && !slices.ContainsFunc(changeLeaves(m, c.Ops), func(l ChangeLeaf) bool { return l.ID == node.ID }) {
		err = Usagef("%s is not one of change set %s's leaves", ref, rr.change)
	}
	if err != nil {
		return nil, Usagef("%s: leaf %q: %s", cmd, ref, AsError(err).Message).WithNext("rhiza change show " + rr.change.String())
	}
	step := rr.state.Current
	l, ok := rr.state.Leaf(step, node.ID.String())
	switch {
	case !ok || l.Attempt() == 0:
		return nil, Usagef("%s: leaf %s has no attempt in a leaves step of run %s: a verifier journals findings while its leaf's attempt runs",
			cmd, node.Handle, rr.id).WithNext(statusNext)
	case l.Outcome != "":
		return nil, Usagef("%s: leaf %s already %s in run %s", cmd, node.Handle, l.Outcome, rr.id).WithNext(statusNext)
	}
	var f leafFindingsFile
	if err := readReviewFile(ctx, cmd, runFindingsNext, r.File, &f); err != nil {
		return nil, err
	}
	if len(f.Findings) == 0 {
		return nil, Usagef(`%s file: want {"findings":[{"reason":...},...]}, at least one`, cmd).WithNext(runFindingsNext)
	}
	for i, in := range f.Findings {
		if strings.TrimSpace(in.Reason) == "" {
			return nil, Usagef("finding %d: no reason", i+1).WithNext(runFindingsNext)
		}
	}
	rr, unlock, err := rr.lockJournal(ctx)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if l, _ = rr.state.Leaf(step, node.ID.String()); l.Attempt() == 0 || l.Outcome != "" {
		return nil, Usagef("%s: leaf %s no longer has an attempt running in run %s", cmd, node.Handle, rr.id).WithNext(statusNext)
	}
	attempt := l.Attempt()
	have := len(l.AttemptFindings(attempt))
	found := make([]run.Finding, 0, len(f.Findings))
	for i, in := range f.Findings {
		found = append(found, run.Finding{ID: "f" + strconv.Itoa(have+i+1), Reason: in.Reason})
	}
	resp := RunFindingsResponse{Run: rr.id, Leaf: node.Handle, Attempt: attempt, Findings: make([]RunFinding, 0, len(found)), DryRun: ctx.DryRun}
	for _, fd := range found {
		resp.Findings = append(resp.Findings, runFinding(fd, run.Round{}))
	}
	e := run.Entry{
		Type: run.TypeFindings, At: ctx.Now(), Step: step, Leaf: node.ID.String(), Handle: node.Handle, Attempt: attempt, Findings: found,
	}
	if err := rr.journal(ctx, fmt.Sprintf("findings on leaf %s, attempt %d,", node.Handle, attempt), e); err != nil {
		return nil, err
	}
	return resp, nil
}

// planResolver resolves a finding's node across a run's plan: the run's
// change set's view first, then the view of each of its plan's change sets
// (planChangeSets, sorted by ID), read once on first need.
type planResolver struct {
	w     *workspace
	run   id.ID
	views []Model
	read  bool
}

// resolve is the node ref names in the run's change set's view; failing
// that, the node it names in the plan's change sets' views, where it must
// name the same node in each view it resolves in. Not found anywhere it is
// the run's view's error, naming the plan too.
func (p *planResolver) resolve(ref string) (Node, error) {
	runView := p.w.model(p.run)
	node, err := runView.Resolve(ref)
	if err == nil {
		return node, nil
	}
	if !p.read {
		p.read = true
		sets := p.w.changeSets()
		plan, perr := planChangeSets(sets, runView, p.run)
		if perr != nil {
			return Node{}, Errorf("%v", perr)
		}
		for _, c := range plan {
			if c.ID != p.run {
				p.views = append(p.views, p.w.model(c.ID))
			}
		}
	}
	var hit *Node
	var ids []string
	for _, v := range p.views {
		n, verr := v.Resolve(ref)
		if verr != nil {
			continue
		}
		if hit == nil {
			hit = &n
		}
		if !slices.Contains(ids, n.ID.String()) {
			ids = append(ids, n.ID.String())
		}
	}
	switch {
	case len(ids) > 1:
		return Node{}, Usagef("names different nodes in the plan's change sets: %s", strings.Join(ids, ", "))
	case hit != nil:
		return *hit, nil
	}
	return Node{}, Usagef("%s, in change set %s's view or any of its plan's change sets", AsError(err).Message, p.run)
}

// runRunVerdict checks the run, its latest round, then each verdict: a
// finding the round holds, named once, confirmed or dropped, with a
// reason. Only then does it journal them.
func runRunVerdict(ctx *Ctx, req any) (any, error) {
	r := req.(RunVerdictRequest)
	const cmd = "run verdict"
	rr, err := openReviewRun(ctx, cmd, runVerdictNext)
	if err != nil {
		return nil, err
	}
	var f verdictFile
	if err := readReviewFile(ctx, cmd, runVerdictNext, r.File, &f); err != nil {
		return nil, err
	}
	step, review, err := rr.reviewCheck(cmd, r.Check, runVerdictNext)
	if err != nil {
		return nil, err
	}
	check := strings.TrimSpace(r.Check)
	if len(f.Verdicts) == 0 {
		return nil, Usagef(`%s file: want {"verdicts":[...]}, at least one`, cmd).WithNext(runVerdictNext)
	}
	var round run.Round
	what := ""
	if review {
		cs, _ := rr.state.Check(step.Name, check)
		rd, ok := cs.LatestRound()
		if !ok {
			return nil, Usagef("%s: check %s of run %s has no round yet: rhiza run findings --check %s opens one", cmd, check, rr.id, check).WithNext(runVerdictNext)
		}
		round, what = rd, "check "+check+"'s round "+strconv.Itoa(rd.Number)
	} else {
		if round, err = rr.latestRound(cmd, runVerdictNext); err != nil {
			return nil, err
		}
		what = "round " + strconv.Itoa(round.Number)
	}
	seen := map[string]bool{}
	for i, v := range f.Verdicts {
		n := i + 1
		if _, ok := round.Finding(v.Finding); !ok {
			return nil, Usagef("verdict %d: %s holds no finding %q: a critic judges the reviewer's findings and adds none",
				n, what, v.Finding).WithNext("rhiza run status " + rr.id)
		}
		if seen[v.Finding] {
			return nil, Usagef("verdict %d: finding %s is judged twice", n, v.Finding).WithNext(runVerdictNext)
		}
		seen[v.Finding] = true
		if v.Verdict != run.Confirmed && v.Verdict != run.Dropped {
			return nil, Usagef("verdict %d: verdict %q is not confirmed or dropped", n, v.Verdict).WithNext(runVerdictNext)
		}
		if strings.TrimSpace(v.Reason) == "" {
			return nil, Usagef("verdict %d: no reason", n).WithNext(runVerdictNext)
		}
	}
	resp := RunVerdictResponse{Run: rr.id, Round: round.Number, Check: check, Verdicts: make([]RunVerdict, 0, len(f.Verdicts)), DryRun: ctx.DryRun}
	for _, v := range f.Verdicts {
		resp.Verdicts = append(resp.Verdicts, RunVerdict{Finding: v.Finding, Verdict: v.Verdict, Reason: v.Reason})
	}
	e := run.Entry{Type: run.TypeVerdict, At: ctx.Now(), Verdicts: f.Verdicts}
	if review {
		e.Step, e.Check = step.Name, check
	}
	if err := rr.journal(ctx, "verdicts on "+what, e); err != nil {
		return nil, err
	}
	return resp, nil
}

// runFinding is the response form of finding f, with its latest verdict
// in round rd.
func runFinding(f run.Finding, rd run.Round) RunFinding {
	out := RunFinding{
		ID: f.ID, Node: f.Node, Handle: f.Handle, Kind: f.Kind, Check: f.Check, Severity: f.Severity,
		Leaf: f.Leaf, LeafHandle: f.LeafHandle, Place: f.Place, Reason: f.Reason,
	}
	if v, ok := rd.Verdict(f.ID); ok {
		out.Verdict, out.VerdictReason = v.Verdict, v.Reason
	}
	return out
}

// findingsFileSchema is the JSON Schema of run.findings' input.
func findingsFileSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"findings": map[string]any{
				"type":        "array",
				"description": "The findings, in order; empty for a clean review or check, but at least one with --leaf. They are numbered f1, f2, ... in the round, or in the leaf's attempt",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"node":     map[string]any{"type": "string", "description": "The node: a full ID, a unique prefix of 4+ characters or a handle, in the run's change set's view. Required, except with --check or --leaf, which refuse it"},
						"kind":     map[string]any{"enum": stringsAny(run.Kinds), "description": "What is wrong with the node. Required, except with --check or --leaf, which refuse it"},
						"severity": map[string]any{"enum": stringsAny(run.Severities), "description": "How serious it is. Required with --check, refused otherwise"},
						"leaf":     map[string]any{"type": "string", "description": "The leaf of the run's change set it is about, by handle or ID. Optional with --check, refused otherwise"},
						"place":    map[string]any{"type": "string", "minLength": 1, "description": "Where: file:line or a node. Required with --check, refused otherwise"},
						"reason":   map[string]any{"type": "string", "minLength": 1, "description": "Why"},
					},
					"required":             []string{"reason"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"findings"},
		"additionalProperties": false,
	}
}

// verdictFileSchema is the JSON Schema of run.verdict's input.
func verdictFileSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"verdicts": map[string]any{
				"type":        "array",
				"minItems":    1,
				"description": "The verdicts, each on a different finding of the latest round (with --check, of the check's latest round)",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"finding": map[string]any{"type": "string", "description": "The finding's ID in the latest round: f1, f2, ..."},
						"verdict": map[string]any{"enum": []any{run.Confirmed, run.Dropped}, "description": "confirmed or dropped"},
						"reason":  map[string]any{"type": "string", "minLength": 1, "description": "Why"},
					},
					"required":             []string{"finding", "reason", "verdict"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"verdicts"},
		"additionalProperties": false,
	}
}

// stringsAny is ss as a []any, for a schema's enum.
func stringsAny(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

// Human reports the round the findings opened.
func (r RunFindingsResponse) Human(w io.Writer) {
	verb := "Recorded"
	if r.DryRun {
		verb = "Would record"
	}
	switch {
	case r.Leaf != "":
		fmt.Fprintf(w, "%s %d finding(s) on leaf %s, attempt %d, of run %s\n\n", verb, len(r.Findings), r.Leaf, r.Attempt, r.Run)
		for _, f := range r.Findings {
			fmt.Fprintf(w, "- %s: %s\n", f.ID, oneLine(f.Reason))
		}
	case r.Check != "" && len(r.Findings) == 0:
		fmt.Fprintf(w, "%s round %d of check %s (cycle %d) of run %s: no findings\n", verb, r.Round, r.Check, r.Cycle, r.Run)
	case r.Check != "":
		fmt.Fprintf(w, "%s round %d of check %s (cycle %d) of run %s\n\n", verb, r.Round, r.Check, r.Cycle, r.Run)
		for _, f := range r.Findings {
			fmt.Fprintf(w, "- %s\n", findingLine(f))
		}
	case len(r.Findings) == 0:
		fmt.Fprintf(w, "%s round %d of run %s: no findings\n", verb, r.Round, r.Run)
	default:
		fmt.Fprintf(w, "%s round %d of run %s\n\n", verb, r.Round, r.Run)
		for _, f := range r.Findings {
			fmt.Fprintf(w, "- %s\n", findingLine(f))
		}
	}
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was journalled.\n")
		return
	}
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}

// Human reports the verdicts recorded.
func (r RunVerdictResponse) Human(w io.Writer) {
	verb := "Recorded"
	if r.DryRun {
		verb = "Would record"
	}
	if r.Check != "" {
		fmt.Fprintf(w, "%s verdicts on round %d of check %s of run %s\n\n", verb, r.Round, r.Check, r.Run)
	} else {
		fmt.Fprintf(w, "%s verdicts on round %d of run %s\n\n", verb, r.Round, r.Run)
	}
	for _, v := range r.Verdicts {
		fmt.Fprintf(w, "- %s: %s: %s\n", v.Finding, v.Verdict, oneLine(v.Reason))
	}
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was journalled.\n")
		return
	}
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}

// findingLine is a finding as "f1 kind on handle: reason", with its
// verdict when it has one. The node shows by handle, else by ID. A review
// check's finding is "f1 severity at place on leaf: reason", "on leaf"
// only when it names one.
func findingLine(f RunFinding) string {
	if f.Severity != "" {
		s := f.ID + " " + f.Severity + " at " + oneLine(f.Place)
		if leaf := cmp.Or(f.LeafHandle, f.Leaf); leaf != "" {
			s += " on " + leaf
		}
		s += ": " + oneLine(f.Reason)
		if f.Verdict != "" {
			s += " — " + f.Verdict + ": " + oneLine(f.VerdictReason)
		}
		return s
	}
	node := f.Handle
	if node == "" {
		node = f.Node
	}
	s := f.ID + " " + f.Kind + " on " + node + ": " + oneLine(f.Reason)
	if f.Verdict != "" {
		s += " — " + f.Verdict + ": " + oneLine(f.VerdictReason)
	}
	return s
}

// reviewRoundCap is the review round from which confirmed findings no
// longer fail the check but are left as questions for the gate.
const reviewRoundCap = 3

// reviewQuestionPrefix starts the handle of every question node
// run.review-check adds: review-<run8>-r<round>-<finding>.
const reviewQuestionPrefix = "review-"

const runReviewCheckNext = "rhiza describe run.review-check"

// RunReviewCheckRequest takes nothing: the run is RHIZA_RUN's.
type RunReviewCheckRequest struct{}

// RunReviewQuestion is a question node run.review-check left for a
// confirmed finding.
type RunReviewQuestion struct {
	ID      string `json:"id" desc:"The question node's ID"`
	Handle  string `json:"handle" desc:"Its handle: review-<run8>-r<round>-<finding>"`
	Finding string `json:"finding" desc:"The finding it holds, by ID in the round"`
	Blocks  string `json:"blocks,omitempty" desc:"The capability it blocks, by ID; absent when the finding names another kind of node"`
	Added   bool   `json:"added" desc:"True if this check added it; false if an earlier check had"`
}

// RunReviewCheckResponse is the check's outcome. A failing check is an
// error (exit 1), not a response.
type RunReviewCheckResponse struct {
	Run       string              `json:"run" desc:"The run's full ID"`
	Round     int                 `json:"round" desc:"The latest review round, from 1"`
	Passed    bool                `json:"passed" desc:"True: the review may go on to the gate"`
	Confirmed []RunFinding        `json:"confirmed" desc:"The round's findings the critic confirmed, in order"`
	Questions []RunReviewQuestion `json:"questions,omitempty" desc:"From round 3, the question nodes left for the confirmed findings, in order"`
	DryRun    bool                `json:"dry_run,omitempty" desc:"True if nothing was written or journalled"`
}

// reviewQuestionHandle is the handle of the question run rid leaves for
// finding fid of round n.
func reviewQuestionHandle(rid string, n int, fid string) string {
	return reviewQuestionPrefix + strings.ToLower(rid[:8]) + "-r" + strconv.Itoa(n) + "-" + fid
}

// reviewQuestionText is the statement of the question left for finding
// f, confirmed for reason, in round n.
func reviewQuestionText(n int, f run.Finding, reason string) string {
	node := f.Handle
	if node == "" {
		node = f.Node
	}
	return fmt.Sprintf("Plan review round %d left finding %s (%s on %s): %s. The critic confirmed it: %s. How should it be settled?",
		n, f.ID, f.Kind, node, strings.TrimRight(oneLine(f.Reason), "."), strings.TrimRight(oneLine(reason), "."))
}

// runRunReviewCheck reads the run's latest review round. With no round
// it fails; with no confirmed finding it passes. Before round 3 a
// confirmed finding fails it (exit 1), writing nothing. From round 3 it
// adds a question node per confirmed finding, linked blocks to the
// finding's node when that is a live capability where the question goes
// (reviewQuestionTarget): the run's change set while it is open; once it
// has landed, as a milestone plan's record does after the interview, the
// open plan change set that created the node, else a record write. It
// writes one entry per target as the caller, journals an asked entry for
// each question (Step RHIZA_STEP, no Session), so an autonomous run's
// arbiter decides it, and passes. A question already created in any open
// or landed change set, or already journalled, is not again.
func runRunReviewCheck(ctx *Ctx, _ any) (any, error) {
	const cmd = "run review-check"
	rr, err := openReviewRun(ctx, cmd, runReviewCheckNext)
	if err != nil {
		return nil, err
	}
	rd, ok := rr.state.LatestRound()
	if !ok {
		return nil, Errorf("%s: run %s has no review round yet, so the plan is unreviewed: rhiza run findings opens one", cmd, rr.id).
			WithNext("rhiza describe run.findings")
	}
	confirmed := rd.Confirmed()
	resp := RunReviewCheckResponse{Run: rr.id, Round: rd.Number, Confirmed: make([]RunFinding, 0, len(confirmed)), DryRun: ctx.DryRun}
	for _, f := range confirmed {
		resp.Confirmed = append(resp.Confirmed, runFinding(f, rd))
	}
	if len(confirmed) == 0 {
		resp.Passed = true
		return resp, nil
	}
	if rd.Number < reviewRoundCap {
		lines := make([]string, 0, len(resp.Confirmed))
		for _, f := range resp.Confirmed {
			lines = append(lines, "- "+findingLine(f))
		}
		return nil, Errorf("%s: review round %d of run %s has %d confirmed finding(s); revise the plan and review it again:\n%s",
			cmd, rd.Number, rr.id, len(confirmed), strings.Join(lines, "\n")).WithNext("rhiza run status " + rr.id)
	}

	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	c, ok := w.changeSet(rr.change)
	if !ok {
		return nil, NotFoundf("run %s's change set %s is not in the log", rr.id, rr.change).WithNext("rhiza run status " + rr.id)
	}
	sets := w.changeSets()
	var plan []changeset.ChangeSet
	if c.State.Terminal() {
		if plan, err = planChangeSets(sets, w.model(rr.change), rr.change); err != nil {
			return nil, Errorf("%s: %v", cmd, err).WithNext("rhiza run status " + rr.id)
		}
	}
	asked := map[string]bool{}
	for _, a := range rr.state.Asked {
		asked[a.Handle] = true
	}
	var targets []reviewTarget
	items := map[reviewTarget][]applyItem{}
	for _, f := range confirmed {
		v, _ := rd.Verdict(f.ID)
		q := RunReviewQuestion{Handle: reviewQuestionHandle(rr.id, rd.Number, f.ID), Finding: f.ID}
		t := reviewQuestionTarget(c, plan, mustParseID(f.Node))
		if n, ok := w.model(t.cs).Node(mustParseID(f.Node)); ok && n.Kind == op.KindCapability && !n.Retired {
			q.Blocks = n.ID.String()
		}
		if qid, ok := reviewQuestionNode(sets, q.Handle); ok {
			q.ID = qid.String()
		} else {
			q.Added = true
			if _, ok := items[t]; !ok {
				targets = append(targets, t)
			}
			add := AddRequest{Kind: string(op.KindQuestion), Handle: q.Handle,
				Fields: []string{"statement=" + reviewQuestionText(rd.Number, f, v.Reason)}, lenient: true}
			items[t] = append(items[t], applyItem{op: "add", build: func(b writeBase) (edit, error) { return buildAdd(b, add) }})
			if q.Blocks != "" {
				link := LinkRequest{From: q.Handle, Type: string(op.LinkBlocks), To: q.Blocks}
				items[t] = append(items[t], applyItem{op: "link", build: func(b writeBase) (edit, error) { return buildLink(b, link) }})
			}
		}
		resp.Questions = append(resp.Questions, q)
	}
	created := map[string]string{}
	for _, t := range targets {
		wctx := *ctx
		wctx.ChangeSet, wctx.ChangeRef, wctx.Record = t.cs, "", t.record
		intent := fmt.Sprintf("Leave the questions of plan review round %d of run %s", rd.Number, rr.id)
		r, err := runWriteOn(&wctx, "add", intent, func(w *workspace, cs id.ID) (edit, error) {
			return buildApply(w, cs, ctx.IDs, items[t])
		})
		if err != nil {
			return nil, err
		}
		for n, h := range r.(WriteResponse).handles {
			created[h] = n.String()
		}
	}
	step := strings.TrimSpace(ctxGetenv(ctx)("RHIZA_STEP"))
	for i, q := range resp.Questions {
		if q.Added {
			resp.Questions[i].ID = created[q.Handle]
		}
		if resp.Questions[i].ID == "" {
			return nil, Errorf("question %s has no node", q.Handle)
		}
		if asked[q.Handle] {
			continue
		}
		f, _ := rd.Finding(q.Finding)
		v, _ := rd.Verdict(q.Finding)
		e := run.Entry{Type: run.TypeAsked, At: ctx.Now(), Step: step, Question: resp.Questions[i].ID,
			Handle: q.Handle, Detail: reviewQuestionText(rd.Number, f, v.Reason)}
		if err := rr.journal(ctx, "question "+q.Handle, e); err != nil {
			return nil, err
		}
	}
	resp.Passed = true
	return resp, nil
}

// reviewTarget is where run.review-check writes a question: change set
// cs, or with record a record change set it opens and lands.
type reviewTarget struct {
	cs     id.ID
	record bool
}

// reviewQuestionTarget is where the question on node goes: the run's
// change set rc while it is open; else the draft, proposed or accepted
// change set among plan whose ops create node; else a record write
// (decision review-questions-follow-the-node).
func reviewQuestionTarget(rc changeset.ChangeSet, plan []changeset.ChangeSet, node id.ID) reviewTarget {
	if !rc.State.Terminal() {
		return reviewTarget{cs: rc.ID}
	}
	for _, c := range plan {
		if !planInFlight(c.State) {
			continue
		}
		if slices.ContainsFunc(c.Ops, func(o op.Op) bool { return o.Type == op.TypeNodeCreate && o.Node == node }) {
			return reviewTarget{cs: c.ID}
		}
	}
	return reviewTarget{record: true}
}

// reviewQuestionNode is the node an earlier check created for the question
// with handle h: a node.create op with that handle in any change set that
// is open or landed.
func reviewQuestionNode(sets []changeset.ChangeSet, h string) (id.ID, bool) {
	for _, c := range sets {
		if c.State.Terminal() && c.State != changeset.StateLanded {
			continue
		}
		for _, o := range c.Ops {
			if o.Type == op.TypeNodeCreate && o.Handle == h {
				return o.Node, true
			}
		}
	}
	return id.ID{}, false
}

// mustParseID parses s, or is the zero ID when s is not one.
func mustParseID(s string) id.ID {
	n, _ := id.Parse(s)
	return n
}

// Human reports the check's outcome.
func (r RunReviewCheckResponse) Human(w io.Writer) {
	if len(r.Confirmed) == 0 {
		fmt.Fprintf(w, "Review round %d of run %s passed: no confirmed finding\n", r.Round, r.Run)
		return
	}
	verb := "left"
	if r.DryRun {
		verb = "would leave"
	}
	fmt.Fprintf(w, "Review round %d of run %s passed: round %d is the last, so it %s %d question(s) for the gate\n\n",
		r.Round, r.Run, r.Round, verb, len(r.Questions))
	for _, q := range r.Questions {
		s := "- " + q.Handle + " for " + q.Finding
		if q.Blocks != "" {
			s += ", blocks " + q.Blocks
		}
		if !q.Added {
			s += " (already there)"
		}
		fmt.Fprintln(w, s)
	}
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written or journalled.\n")
		return
	}
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}

const runGateNext = "rhiza describe run.gate"

// RunGateRequest takes nothing: the run is RHIZA_RUN's.
type RunGateRequest struct{}

// RunGateChangeSet is one of the plan's change sets and what the gate
// did with it.
type RunGateChangeSet struct {
	ID    string          `json:"id" desc:"The change set's ID"`
	Title string          `json:"title" desc:"Its title"`
	From  changeset.State `json:"from" desc:"Its state before the gate"`
	State changeset.State `json:"state" desc:"Its state after the gate (in a dry run, the state it would have)"`
	Moves []string        `json:"moves,omitempty" desc:"What the gate did, in order: proposed, accepted; absent when it left the change set as it was"`
	Next  string          `json:"next,omitempty" desc:"The owner's command when it waits proposed in a normal run; absent otherwise"`
}

// RunGateResponse is the plan's change sets after the gate. An
// autonomous run whose review has not passed is an error (exit 1), not a
// response; one whose review left a question open stops (ExitCode).
type RunGateResponse struct {
	Run        string             `json:"run" desc:"The run's full ID"`
	Autonomous bool               `json:"autonomous" desc:"True if the run is autonomous, so the gate accepted in the owner's place"`
	ChangeSets []RunGateChangeSet `json:"change_sets" desc:"The plan's change sets: the run's change set and every change set serving an intent it adds, by ID"`
	Open       []RunGateQuestion  `json:"open" desc:"Every question the plan left with no live answering Decision: the review's questions and each question node created in the plan's change sets, in order; empty when none"`
	Commit     string             `json:"commit,omitempty" desc:"The commit that renders and commits .rhiza on the run's branch; absent when there was nothing to commit, or in a dry run"`
	DryRun     bool               `json:"dry_run,omitempty" desc:"True if nothing was written or journalled"`
}

// ExitCode is run.StopExit while questions are open (a plan run then
// stops for the owner), else 0.
func (r RunGateResponse) ExitCode() int {
	if len(r.Open) > 0 && !r.DryRun {
		return run.StopExit
	}
	return 0
}

// stopDetail is why a plan run stops at the gate: the open questions, by
// handle, and the command that resumes it.
func (r RunGateResponse) stopDetail() string {
	names := make([]string, 0, len(r.Open))
	for _, q := range r.Open {
		if q.Handle != "" {
			names = append(names, q.Handle)
		} else {
			names = append(names, q.ID)
		}
	}
	return run.OpenStop + strings.Join(names, ", ") + "; answer each with a Decision, then rhiza run resume " + r.Run
}

// RunGateQuestion is one question the plan left unanswered.
type RunGateQuestion struct {
	ID     string `json:"id" desc:"The question node's full ID"`
	Handle string `json:"handle,omitempty" desc:"Its handle, when it has one"`
}

// gateOpenQuestions are the questions the plan left unanswered: the
// run's review questions (reviewQuestions) and every question node a
// plan change set's ops create, each with no live answering Decision
// (reviewAnswered) and not retired, once each, in discovery order.
func gateOpenQuestions(sets []changeset.ChangeSet, plan []changeset.ChangeSet, m Model, st run.State) []RunGateQuestion {
	out := []RunGateQuestion{}
	seen := map[id.ID]bool{}
	consider := func(q id.ID, handle string) {
		if seen[q] {
			return
		}
		seen[q] = true
		if !questionOpen(sets, m, q) {
			return
		}
		if n, ok := m.Node(q); ok && handle == "" {
			handle = n.Handle
		}
		out = append(out, RunGateQuestion{ID: q.String(), Handle: handle})
	}
	for _, a := range reviewQuestions(st) {
		consider(mustParseID(a.Question), a.Handle)
	}
	for _, c := range plan {
		for _, o := range c.Ops {
			if o.Type == op.TypeNodeCreate && o.Kind == op.KindQuestion {
				consider(o.Node, o.Handle)
			}
		}
	}
	return out
}

// gateCommitMessage is the message of the commit that closes a plan run's
// branch.
const gateCommitMessage = "Render the plan"

// gateRenderCommit renders .rhiza/rendered and commits it with every
// .rhiza/log entry the run wrote, with the Change and Actor trailers, so
// the branch passes rhiza check. Nothing to commit commits nothing.
func gateRenderCommit(ctx *Ctx, cs id.ID) (string, error) {
	if _, err := runRender(ctx, RenderRequest{}); err != nil {
		return "", err
	}
	if _, err := ctx.Git("add", "-A", "--", ".rhiza"); err != nil {
		return "", Errorf("%v", err)
	}
	return shipCommitStaged(ctx, shipTarget{cs: cs}, ".rhiza", gateCommitMessage)
}

// runRunGate ends the plan of the run RHIZA_RUN names. It proposes each
// draft among the plan's change sets (planChangeSets) and leaves the
// others as they are, then journals a gated entry naming them all. A
// run then checks the review (gateReview). An autonomous run whose review
// has not passed fails (exit 1) having accepted nothing, the proposals
// kept; an owner-run one just leaves them, each with rhiza accept --rule
// owner next, as it does when rhiza.autonomy.accept is owner. Else it
// accepts each proposed change set, under rule autonomous (autonomous
// run) or autonomy (owner-run), which journals an accepted entry per
// change set. Running it again moves nothing already moved.
func runRunGate(ctx *Ctx, _ any) (any, error) {
	const cmd = "run gate"
	rr, err := openReviewRun(ctx, cmd, runGateNext)
	if err != nil {
		return nil, err
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	m := w.model(rr.change)
	sets := w.changeSets()
	plan, err := planChangeSets(sets, m, rr.change)
	if err != nil {
		return nil, Errorf("%s: %v", cmd, err).WithNext("rhiza run status " + rr.id)
	}
	resp := RunGateResponse{Run: rr.id, Autonomous: rr.state.Autonomous, ChangeSets: make([]RunGateChangeSet, 0, len(plan)), Open: gateOpenQuestions(sets, plan, m, rr.state), DryRun: ctx.DryRun}
	propose, _ := Lookup("propose")
	ids := make([]string, 0, len(plan))
	for _, c := range plan {
		g := RunGateChangeSet{ID: c.ID.String(), Title: c.Open.DisplayTitle(), From: c.State, State: c.State}
		if c.State == changeset.StateDraft {
			if _, err := propose.Run(gateCtx(ctx, c.ID), LifecycleRequest{}); err != nil {
				return nil, err
			}
			g.State, g.Moves = changeset.StateProposed, append(g.Moves, "proposed")
		}
		resp.ChangeSets = append(resp.ChangeSets, g)
		ids = append(ids, g.ID)
	}
	if !slices.Equal(ids, rr.state.Gated) {
		if err := rr.journal(ctx, "gate", run.Entry{Type: run.TypeGated, At: ctx.Now(), Changes: ids}); err != nil {
			return nil, err
		}
	}
	leave := func() (any, error) {
		for i, g := range resp.ChangeSets {
			if g.State == changeset.StateProposed {
				resp.ChangeSets[i].Next = ownerAcceptNext(g.ID)
			}
		}
		if err := gateCommitResp(ctx, &resp, rr.change); err != nil {
			return nil, err
		}
		return resp, nil
	}
	reviewOpen, err := gateReview(sets, m, rr)
	if err != nil {
		if rr.state.Autonomous {
			return nil, err
		}
		// An owner-run plan whose review has not passed accepts nothing.
		return leave()
	}
	if reviewOpen && !rr.state.Autonomous {
		return leave()
	}
	if reviewOpen {
		// A review question still open stops the run for the owner, as a
		// question the planner parked does: it accepts nothing (exit
		// run.StopExit, resp.Open names each).
		if err := gateCommitResp(ctx, &resp, rr.change); err != nil {
			return nil, err
		}
		return resp, nil
	}
	acceptMode, err := autonomyStage(ctx, "accept")
	if err != nil {
		return nil, err
	}
	if acceptMode == AutonomyOwner {
		// rhiza.autonomy.accept is owner: the proposals stay for the owner,
		// and the inbox lists each as waiting for acceptance.
		return leave()
	}
	rule := autonomousRule
	if !rr.state.Autonomous {
		rule = autonomyRule
	}
	accept, _ := Lookup("accept")
	for i, g := range resp.ChangeSets {
		if g.State != changeset.StateProposed {
			continue
		}
		if !ctx.DryRun {
			if _, err := accept.Run(gateCtx(ctx, mustParseID(g.ID)), DecisionRequest{Rule: rule}); err != nil {
				return nil, err
			}
		}
		resp.ChangeSets[i].State, resp.ChangeSets[i].Moves = changeset.StateAccepted, append(g.Moves, "accepted")
	}
	if err := gateCommitResp(ctx, &resp, rr.change); err != nil {
		return nil, err
	}
	return resp, nil
}

// gateCommitResp renders and commits the branch after the gate's moves
// (not in a dry run) and records the commit in resp.
func gateCommitResp(ctx *Ctx, resp *RunGateResponse, cs id.ID) error {
	if ctx.DryRun {
		return nil
	}
	c, err := gateRenderCommit(ctx, cs)
	resp.Commit = c
	return err
}

// gateCtx is ctx acting on change set cs, as the gate's propose and
// accept do.
func gateCtx(ctx *Ctx, cs id.ID) *Ctx {
	c := *ctx
	c.ChangeSet, c.ChangeRef, c.Record = cs, "", false
	return &c
}

// ownerAcceptNext is the owner's command to accept change set cs.
func ownerAcceptNext(cs string) string {
	return "rhiza accept --rule owner --change " + cs
}

// planChangeSets is the plan of change set cs, sorted by ID: cs itself and
// every change set with a serves link to an intent cs's ops create, found
// in the change sets' own link.create ops (where a change set's serves
// links are written) and in m's live links. Rejected, withdrawn and
// superseded ones are left out; landed ones stay in.
func planChangeSets(sets []changeset.ChangeSet, m Model, cs id.ID) ([]changeset.ChangeSet, error) {
	byID := map[id.ID]changeset.ChangeSet{}
	for _, c := range sets {
		byID[c.ID] = c
	}
	own, ok := byID[cs]
	if !ok {
		return nil, fmt.Errorf("change set %s is not in the log", cs)
	}
	intents := map[id.ID]bool{}
	for _, o := range own.Ops {
		if o.Type == op.TypeNodeCreate && o.Kind == op.KindIntent {
			intents[o.Node] = true
		}
	}
	in := map[id.ID]bool{cs: true}
	for i := range intents {
		for _, l := range m.LinksTo(i, op.LinkServes) {
			in[l.From] = true
		}
	}
	for _, c := range sets {
		for _, o := range c.Ops {
			if o.Type == op.TypeLinkCreate && o.LinkType == op.LinkServes && intents[o.To] {
				in[o.From] = true
			}
		}
	}
	var out []changeset.ChangeSet
	for i := range in {
		c, ok := byID[i]
		if !ok || (c.State.Terminal() && c.State != changeset.StateLanded) {
			continue
		}
		out = append(out, c)
	}
	slices.SortFunc(out, func(a, b changeset.ChangeSet) int { return id.Compare(a.ID, b.ID) })
	return out, nil
}

// reviewQuestions are the asked entries run.review-check journalled in
// run st: those with no Session whose handle starts review-.
func reviewQuestions(st run.State) []run.Asked {
	var out []run.Asked
	for _, a := range st.Asked {
		if a.Session == "" && strings.HasPrefix(a.Handle, reviewQuestionPrefix) {
			out = append(out, a)
		}
	}
	return out
}

// reviewAnswered reports whether question q is answered: m (the run's
// change set's view) has a live answers link to it from a live Decision,
// or a draft, proposed or accepted change set's ops create an answers
// link to it.
func reviewAnswered(sets []changeset.ChangeSet, m Model, q id.ID) bool {
	for _, l := range m.LinksTo(q, op.LinkAnswers) {
		if n, ok := m.Node(l.From); ok && n.Kind == op.KindDecision && !n.Retired {
			return true
		}
	}
	for _, c := range sets {
		if !planInFlight(c.State) {
			continue
		}
		for _, o := range c.Ops {
			if o.Type == op.TypeLinkCreate && o.LinkType == op.LinkAnswers && o.To == q {
				return true
			}
		}
	}
	return false
}

// questionOpen reports whether question q is still open: its node is not
// retired and it is not answered (reviewAnswered). gateOpenQuestions and
// gateReview both use it, so the gate's stop and resp.Open agree.
func questionOpen(sets []changeset.ChangeSet, m Model, q id.ID) bool {
	if n, ok := m.Node(q); ok && n.Retired {
		return false
	}
	return !reviewAnswered(sets, m, q)
}

// gateReview decides whether an autonomous run's review passed, so the
// gate may accept. With no review round it has not. The latest round
// passes with no confirmed finding, or with a question left for each
// confirmed one. It reports whether a question the review left
// (reviewQuestions) is still open (questionOpen): the gate then stops for
// the owner and accepts nothing.
func gateReview(sets []changeset.ChangeSet, m Model, rr reviewRun) (bool, error) {
	const cmd = "run gate"
	rd, ok := rr.state.LatestRound()
	if !ok {
		return false, Errorf("%s: run %s is autonomous and has no review round, so the plan is unreviewed and it accepts nothing: rhiza run findings opens one",
			cmd, rr.id).WithNext("rhiza describe run.findings")
	}
	asked := map[string]bool{}
	open := false
	for _, a := range reviewQuestions(rr.state) {
		asked[a.Handle] = true
		if questionOpen(sets, m, mustParseID(a.Question)) {
			open = true
		}
	}
	for _, f := range rd.Confirmed() {
		if !asked[reviewQuestionHandle(rr.id, rd.Number, f.ID)] {
			return false, Errorf("%s: run %s is autonomous and review round %d has confirmed finding %s with no question left for it, so the review has not passed and it accepts nothing: run rhiza run review-check",
				cmd, rr.id, rd.Number, f.ID).WithNext(runReviewCheckNext)
		}
	}
	return open, nil
}

// Human reports the plan's change sets after the gate.
func (r RunGateResponse) Human(w io.Writer) {
	mode := "normal"
	if r.Autonomous {
		mode = "autonomous"
	}
	verb := "Gated"
	if r.DryRun {
		verb = "Would gate"
	}
	fmt.Fprintf(w, "%s the plan of %s run %s\n\n", verb, mode, r.Run)
	var next []string
	for _, c := range r.ChangeSets {
		s := "- " + c.ID + " " + c.Title + ": " + string(c.State)
		if len(c.Moves) > 0 {
			s += " (by the gate: " + strings.Join(c.Moves, ", ") + ")"
		} else {
			s += " (left as it was)"
		}
		fmt.Fprintln(w, s)
		if c.Next != "" {
			next = append(next, c.Next)
		}
	}
	if len(r.Open) > 0 {
		fmt.Fprintf(w, "\n%d question(s) the plan left are open:\n", len(r.Open))
		for _, q := range r.Open {
			s := "- " + q.ID
			if q.Handle != "" {
				s = "- " + q.Handle + " (" + q.ID + ")"
			}
			fmt.Fprintln(w, s)
		}
	}
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written or journalled.\n")
		return
	}
	if len(r.Open) > 0 {
		fmt.Fprintf(w, "\nstop: %s\n", r.stopDetail())
		fmt.Fprintf(w, "next: rhiza run resume %s\n", r.Run)
		return
	}
	if len(next) > 0 {
		fmt.Fprintln(w)
		for _, n := range next {
			fmt.Fprintf(w, "next: %s\n", n)
		}
		return
	}
	fmt.Fprintf(w, "\nnext: rhiza run status %s\n", r.Run)
}
