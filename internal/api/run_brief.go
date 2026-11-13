package api

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/slice"
	"github.com/dphaener/rhiza/internal/workflow"
)

// The files a run's session finds in its directory, and the environment
// variables that name them.
const (
	briefFile    = "brief.md"
	sliceFile    = "slice.json"
	handbackFile = "handback.md"

	envBrief    = "RHIZA_BRIEF"
	envSlice    = "RHIZA_SLICE"
	envHandback = "RHIZA_HANDBACK"
)

// prepareSession fills the directory of the session named session of step
// (in its current visit) with brief.md, built from the model as the run's
// worktree holds it, and slice.json, the bytes `rhiza slice --json` prints
// for the run's change set now, and records the directory on req so the
// agent gets RHIZA_BRIEF, RHIZA_SLICE and, when handback.md is there,
// RHIZA_HANDBACK. A dry run writes nothing. A directory or file that cannot
// be written fails the start with its path.
func (r *stepRunner) prepareSession(req *SessionStartRequest, step, session string) error {
	common, err := gitCommonDir(r.ctx)
	if err != nil {
		return err
	}
	dir := run.SessionDir(common, r.run, step, r.runState().CurrentVisit, session)
	req.briefDir = dir
	return r.fillSession(dir)
}

// prepareKeyed is prepareSession for a session that is no step visit, an
// arbiter or an advisor: it sets req to fill <runs>/<run>/<kind>/<session
// id> once runSessionStart has allocated the session's ID, before the
// launch.
func (r *stepRunner) prepareKeyed(req *SessionStartRequest, kind string) {
	req.briefFor = func(sid string) (string, error) {
		common, err := gitCommonDir(r.ctx)
		if err != nil {
			return "", err
		}
		dir := run.KeyedSessionDir(common, r.run, kind, sid)
		return dir, r.fillSession(dir)
	}
}

// fillSession writes brief.md and slice.json into dir for the run's change
// set. A dry run writes nothing.
func (r *stepRunner) fillSession(dir string) error {
	if r.ctx.DryRun {
		return nil
	}
	cs, err := id.Parse(r.change)
	if err != nil {
		return fmt.Errorf("run %s names change set %q: %v", r.run, r.change, err)
	}
	// The model as the run's worktree holds it, else the caller's checkout
	// (the worktree does not hold the change set yet).
	var c *Ctx
	var sl SliceResponse
	for _, wt := range []bool{true, false} {
		cc := *r.ctx
		cc.DryRun, cc.ChangeSet, cc.ChangeRef = false, cs, ""
		if wt {
			root, rerr := r.worktreeRoot()
			if _, serr := os.Stat(r.worktree); rerr != nil || serr != nil {
				continue
			}
			cc.Root = root
			cc.Git = func(args ...string) (string, error) {
				out, err := r.ctx.Exec("git", append([]string{"-C", root}, args...), root, nil)
				return string(out), err
			}
		}
		resp, serr := runSlice(&cc, SliceRequest{})
		if serr != nil {
			err = serr
			continue
		}
		c, sl, err = &cc, resp.(SliceResponse), nil
		break
	}
	if c == nil {
		if err == nil {
			err = fmt.Errorf("no checkout holds change set %s", cs)
		}
		return err
	}
	js, err := canon.Marshal(sl)
	if err != nil {
		return fmt.Errorf("encode slice: %w", err)
	}
	// brief.md lists every decision, so it reads a slice nothing is cut
	// from; slice.json stays what `rhiza slice --json` prints.
	full, err := runSlice(c, SliceRequest{Budget: unboundedSliceBudget})
	if err != nil {
		return err
	}
	in, err := r.briefInput(c, cs, full.(SliceResponse))
	if err != nil {
		return err
	}
	return run.WriteSessionFiles(dir, run.RenderBrief(in), js)
}

// unboundedSliceBudget is a token budget no slice reaches, so no section
// is cut.
const unboundedSliceBudget = SliceBudget(math.MaxInt32)

// briefInput is what brief.md is rendered from, over the uncut slice sl of
// change set cs: its decisions are the change set's own (focus) Decisions,
// then the other decisions that apply.
func (r *stepRunner) briefInput(c *Ctx, cs id.ID, sl SliceResponse) (run.BriefInput, error) {
	in := run.BriefInput{ChangeID: cs.String()}
	if sl.Change != nil {
		in.Intent = sl.Change.Intent
	}
	w, err := openWorkspace(c.Root, c.Now, c.IDs, true)
	if err != nil {
		return in, err
	}
	ch, ok := w.changeSet(cs)
	if !ok {
		return in, fmt.Errorf("change set %s is not in the log", cs)
	}
	in.ChangeTitle = ch.Open.DisplayTitle()
	m := w.model(cs)
	for _, l := range changeLeaves(m, ch.Ops) {
		n, ok := m.Node(l.ID)
		if !ok {
			continue
		}
		in.Leaves = append(in.Leaves, run.BriefLeaf{
			Handle: l.Handle, Statement: nodeField(n.Fields, "statement"), Acceptance: nodeField(n.Fields, "acceptance"),
		})
	}
	in.Decisions = decisionStatements(sl)
	in.Constraints = itemStatements(sl.Constraints)
	in.PrePushCheck = r.configValue(PrePushCheckKey)
	in.LocalGate = r.configValue(LocalGateKey)
	return in, nil
}

// nodeField is the string field name of a node's fields, "" when absent.
func nodeField(fields map[string]json.RawMessage, name string) string {
	var s string
	if raw, ok := fields[name]; ok {
		_ = json.Unmarshal(raw, &s)
	}
	return s
}

// decisionStatements are the statements of sl's live focus decisions, then
// of its decisions not among them.
func decisionStatements(sl SliceResponse) []string {
	var out []string
	seen := map[id.ID]bool{}
	for _, f := range sl.Focus {
		if f.Kind != op.KindDecision || f.Retired {
			continue
		}
		seen[f.ID] = true
		out = append(out, f.Statement)
	}
	for _, it := range sl.Decisions {
		if !seen[it.ID] {
			out = append(out, itemStatements([]slice.DocItem{it})...)
		}
	}
	return out
}

// itemStatements are the statements of items, the handle where one has none.
func itemStatements(items []slice.DocItem) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		s := it.Statement
		if s == "" {
			s = it.Handle
		}
		out = append(out, s)
	}
	return out
}

// configValue is git config key as the caller's checkout has it, "" unset.
func (r *stepRunner) configValue(key string) string {
	out, err := r.ctx.Git("config", "--get", key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// briefEnv are the variables naming the files in dir: brief.md and
// slice.json, handback.md only when it exists; nil for no directory.
func briefEnv(dir string) (set, unset []string) {
	if dir == "" {
		return nil, []string{envBrief, envSlice, envHandback}
	}
	set = []string{envBrief + "=" + filepath.Join(dir, briefFile), envSlice + "=" + filepath.Join(dir, sliceFile)}
	if _, err := os.Stat(filepath.Join(dir, handbackFile)); err == nil {
		return append(set, envHandback+"="+filepath.Join(dir, handbackFile)), nil
	}
	return set, []string{envHandback}
}

// promptVars are the values a session's invoke and handback templates
// take: the change set, run and step, the paths of the session's brief.md,
// slice.json and handback.md (named whether or not the last exists yet)
// and the run's worktree.
func (r *stepRunner) promptVars(t changeWorktree, step, session string) workflow.PromptVars {
	v := workflow.PromptVars{Change: t.cs.String(), Slug: t.slug, Title: t.title, Run: r.run, Step: step, Worktree: t.worktree}
	if common, err := gitCommonDir(r.ctx); err == nil {
		dir := run.SessionDir(common, r.run, step, r.runState().CurrentVisit, session)
		v.Brief, v.Slice, v.Handback = filepath.Join(dir, briefFile), filepath.Join(dir, sliceFile), filepath.Join(dir, handbackFile)
	}
	return v
}
