package api

// This file holds the plan review ops of a run: run.findings journals a
// reviewer's findings, opening the next review round, and run.verdict
// journals a critic's verdicts on the latest round's findings. Both work
// only inside the run RHIZA_RUN names, which must be unfinished, and both
// check everything before they journal anything. The round count comes
// from the journal: each run.findings call is one round.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.findings",
		Summary:  "Journal a plan reviewer's findings as the next review round of the run RHIZA_RUN names",
		Write:    true,
		Request:  reflect.TypeFor[RunFindingsRequest](),
		Response: reflect.TypeFor[RunFindingsResponse](),
		Input: &Input{
			Summary: `{"findings":[{"node":...,"kind":...,"reason":...},...]}, maybe empty for a clean review; kind is one of ` + strings.Join(run.Kinds, ", "),
			Schema:  findingsFileSchema,
		},
		Run: runRunFindings,
	})
	register(Operation{
		Name:     "run.verdict",
		Summary:  "Journal a critic's verdicts on the findings of the latest review round of the run RHIZA_RUN names",
		Write:    true,
		Request:  reflect.TypeFor[RunVerdictRequest](),
		Response: reflect.TypeFor[RunVerdictResponse](),
		Input: &Input{
			Summary: `{"verdicts":[{"finding":"f1","verdict":"confirmed"|"dropped","reason":...},...]}, at least one; each names a finding of the latest round`,
			Schema:  verdictFileSchema,
		},
		Run: runRunVerdict,
	})
}

// RunFindingsRequest names the file of findings.
type RunFindingsRequest struct {
	File string `json:"file" arg:"0" desc:"A JSON file of findings, or - for stdin (see the input schema)"`
}

// RunVerdictRequest names the file of verdicts.
type RunVerdictRequest struct {
	File string `json:"file" arg:"0" desc:"A JSON file of verdicts, or - for stdin (see the input schema)"`
}

// RunFinding is one finding of a review round, with its latest verdict
// once a critic gave one.
type RunFinding struct {
	ID            string `json:"id" desc:"The finding's ID in its round: f1, f2, ..."`
	Node          string `json:"node" desc:"The node it is about, by ID"`
	Handle        string `json:"handle,omitempty" desc:"The node's handle when the finding was recorded; absent when it had none"`
	Kind          string `json:"kind" desc:"What is wrong: missing-acceptance, missing-provenance, unreal-depends-on, outside-intent or contradiction"`
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

// RunFindingsResponse is the round the findings opened.
type RunFindingsResponse struct {
	Run      string       `json:"run" desc:"The run's full ID"`
	Round    int          `json:"round" desc:"The round the findings opened, from 1"`
	Findings []RunFinding `json:"findings" desc:"The findings recorded, in order; empty for a clean review"`
	DryRun   bool         `json:"dry_run,omitempty" desc:"True if nothing was journalled"`
}

// RunVerdictResponse is the verdicts recorded on a round.
type RunVerdictResponse struct {
	Run      string       `json:"run" desc:"The run's full ID"`
	Round    int          `json:"round" desc:"The round judged: the latest"`
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
	if st.Ended {
		return reviewRun{}, Usagef("%s works inside an unfinished run, and run %s is finished: it %s", cmd, j.ID, st.Result).
			WithNext("rhiza run status " + j.ID)
	}
	cs, err := id.Parse(st.Change)
	if err != nil {
		return reviewRun{}, Errorf("run %s names change set %q: %v", j.ID, st.Change, err)
	}
	return reviewRun{store: rs, id: j.ID, state: st, change: cs}, nil
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
// reason, and a node that resolves in the run's change set's view. Only
// then does it journal them, numbered f1, f2, ..., as the next round.
func runRunFindings(ctx *Ctx, req any) (any, error) {
	r := req.(RunFindingsRequest)
	const cmd = "run findings"
	rr, err := openReviewRun(ctx, cmd, runFindingsNext)
	if err != nil {
		return nil, err
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
	m := w.model(rr.change)
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
		node, err := m.Resolve(ref)
		if err != nil {
			return nil, Usagef("finding %d: node %q: %s, in change set %s's view", n, ref, AsError(err).Message, rr.change).
				WithNext("rhiza ls --change " + rr.change.String())
		}
		found = append(found, run.Finding{
			ID: "f" + strconv.Itoa(n), Node: node.ID.String(), Handle: node.Handle, Kind: in.Kind, Reason: in.Reason,
		})
	}
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
	if len(f.Verdicts) == 0 {
		return nil, Usagef(`%s file: want {"verdicts":[...]}, at least one`, cmd).WithNext(runVerdictNext)
	}
	round, err := rr.latestRound(cmd, runVerdictNext)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for i, v := range f.Verdicts {
		n := i + 1
		if _, ok := round.Finding(v.Finding); !ok {
			return nil, Usagef("verdict %d: round %d holds no finding %q: a critic judges the reviewer's findings and adds none",
				n, round.Number, v.Finding).WithNext("rhiza run status " + rr.id)
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
	resp := RunVerdictResponse{Run: rr.id, Round: round.Number, Verdicts: make([]RunVerdict, 0, len(f.Verdicts)), DryRun: ctx.DryRun}
	for _, v := range f.Verdicts {
		resp.Verdicts = append(resp.Verdicts, RunVerdict{Finding: v.Finding, Verdict: v.Verdict, Reason: v.Reason})
	}
	e := run.Entry{Type: run.TypeVerdict, At: ctx.Now(), Verdicts: f.Verdicts}
	if err := rr.journal(ctx, "verdicts on round "+strconv.Itoa(round.Number), e); err != nil {
		return nil, err
	}
	return resp, nil
}

// runFinding is the response form of finding f, with its latest verdict
// in round rd.
func runFinding(f run.Finding, rd run.Round) RunFinding {
	out := RunFinding{ID: f.ID, Node: f.Node, Handle: f.Handle, Kind: f.Kind, Reason: f.Reason}
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
				"description": "The findings, in order; empty for a clean review. They are numbered f1, f2, ... in the round",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"node":   map[string]any{"type": "string", "description": "The node: a full ID, a unique prefix of 4+ characters or a handle, in the run's change set's view"},
						"kind":   map[string]any{"enum": stringsAny(run.Kinds), "description": "What is wrong with the node"},
						"reason": map[string]any{"type": "string", "minLength": 1, "description": "Why"},
					},
					"required":             []string{"kind", "node", "reason"},
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
				"description": "The verdicts, each on a different finding of the latest round",
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
	if len(r.Findings) == 0 {
		fmt.Fprintf(w, "%s round %d of run %s: no findings\n", verb, r.Round, r.Run)
	} else {
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
	fmt.Fprintf(w, "%s verdicts on round %d of run %s\n\n", verb, r.Round, r.Run)
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
// verdict when it has one. The node shows by handle, else by ID.
func findingLine(f RunFinding) string {
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
