package api

// This file holds the ship operations that finish a change set once its
// pull request's CI has passed: ship.land records the PR and its checks as
// Evidence that verifies the change set's capabilities, lands the change
// set, renders, and commits and pushes that to the PR branch; ship.close-out
// removes the worktree and the branch, locally and on origin, and brings
// main up to date after the merge; ship.command runs the repo's own ship
// command, .rhiza/workflows/ship.sh. They share the state file and helpers
// of ship.go.

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "ship.land",
		Summary:  "Record a change set's passed pull request as Evidence, land the change set, render, and commit and push it to the PR branch",
		Write:    true,
		Request:  reflect.TypeFor[ShipLandRequest](),
		Response: reflect.TypeFor[ShipLandResponse](),
		Run:      runShipLand,
	})
	register(Operation{
		Name:     "ship.close-out",
		Summary:  "After the merge, remove a change set's worktree and branch, locally and on origin, and fast-forward main",
		Request:  reflect.TypeFor[ShipCloseOutRequest](),
		Response: reflect.TypeFor[ShipCloseOutResponse](),
		Run:      runShipCloseOut,
	})
	register(Operation{
		Name:     "ship.command",
		Summary:  "Run the repo's own ship command, .rhiza/workflows/ship.sh, in the main checkout",
		Request:  reflect.TypeFor[ShipCommandRequest](),
		Response: reflect.TypeFor[ShipCommandResponse](),
		Run:      runShipCommand,
	})
}

// shipScript is the repo's ship command, relative to the project root.
const shipScript = ".rhiza/workflows/ship.sh"

// shipLandMessage is the subject of the commit ship.land makes.
const shipLandMessage = "Land the change set and record its CI"

// ShipLandRequest takes nothing: the PR, its head and checks come from the
// state file ship.check wrote.
type ShipLandRequest struct{}

// ShipLandResponse is what ship.land did.
type ShipLandResponse struct {
	Change          string   `json:"change" desc:"The change set"`
	PR              int      `json:"pr" desc:"The pull request recorded"`
	Evidence        string   `json:"evidence" desc:"The Evidence node's handle"`
	EvidenceWritten bool     `json:"evidence_written" desc:"True if this call wrote the Evidence; false if an earlier one had"`
	EvidenceChange  string   `json:"evidence_change,omitempty" desc:"The change set the Evidence was written into: the change set itself, or a record change set when it had already landed"`
	Verifies        []id.ID  `json:"verifies" desc:"The capabilities the Evidence verifies: those the change set creates or changes"`
	Landed          bool     `json:"landed" desc:"True if this call landed the change set; false if it had already landed"`
	Rendered        []string `json:"rendered" desc:"The rendered files written or removed"`
	Commit          string   `json:"commit,omitempty" desc:"The commit made on the PR branch; empty if nothing was staged"`
	Pushed          string   `json:"pushed,omitempty" desc:"The head pushed to the PR branch; empty if nothing needed pushing"`
}

// ShipCloseOutRequest takes nothing: the merge commit comes from the state
// file ship.merge wrote.
type ShipCloseOutRequest struct{}

// ShipCloseOutResponse lists what close-out did.
type ShipCloseOutResponse struct {
	Change      string     `json:"change" desc:"The change set"`
	Branch      string     `json:"branch" desc:"The change set's branch"`
	Worktree    string     `json:"worktree" desc:"The change set's worktree"`
	Main        string     `json:"main" desc:"The main checkout close-out worked from"`
	Steps       []ShipStep `json:"steps" desc:"Each thing close-out does, in order, with how it went"`
	MainUpdated bool       `json:"main_updated" desc:"True if main was fast-forwarded from origin; the command exits 1 otherwise"`
	Detail      string     `json:"detail,omitempty" desc:"Why main was not updated"`
}

// ShipStep is one thing close-out does.
type ShipStep struct {
	Step   string `json:"step" desc:"What: remove-worktree, close-tab, delete-branch, delete-remote-branch or update-main"`
	Result string `json:"result" desc:"done, already (nothing was left to do) or failed"`
	Tab    string `json:"tab,omitempty" desc:"For close-tab, the Rex window of the change set's session"`
	Detail string `json:"detail,omitempty" desc:"Why it failed"`
}

// ShipCommandRequest takes nothing.
type ShipCommandRequest struct{}

// ShipCommandResponse is how the repo's ship command ran.
type ShipCommandResponse struct {
	Script string `json:"script" desc:"The ship command's path"`
	Dir    string `json:"dir" desc:"The folder it ran in"`
	Ran    bool   `json:"ran" desc:"False if the repo has no ship command"`
	Exit   int    `json:"exit" desc:"Its exit code; the command exits with it"`
}

// ExitCode is 0 once main is updated, else 1.
func (r ShipCloseOutResponse) ExitCode() int {
	if r.MainUpdated {
		return 0
	}
	return 1
}

// ExitCode is the ship command's.
func (r ShipCommandResponse) ExitCode() int { return r.Exit }

// refuseShipDryRun refuses --dry-run: the ship steps run programs and
// push, which a dry run cannot predict.
func refuseShipDryRun(ctx *Ctx, name string) error {
	if ctx.DryRun {
		return Usagef("%s does not take --dry-run: it runs git and gh", CommandPath(name)).WithNext("rhiza describe " + name)
	}
	return nil
}

// shipEvidenceHandle is the Evidence node's handle: pr-<n>-<slug>, made a
// valid handle and cut to 64 bytes at most.
func shipEvidenceHandle(pr int, slug string) string {
	return handleSlug(fmt.Sprintf("pr-%d-%s", pr, slug))
}

// shipEvidenceFields are the Evidence node's text fields, as add's
// name=text. Status verified makes the capabilities it verifies show as
// verified in a slice.
func shipEvidenceFields(st shipState) []string {
	parts := make([]string, 0, len(st.Checks))
	for _, c := range st.Checks {
		parts = append(parts, strings.TrimSpace(c.Name+" "+c.Conclusion))
	}
	return []string{
		fmt.Sprintf("statement=Pull request #%d (%s) passed CI at %s", st.PR, st.URL, shortSHA(st.Head)),
		fmt.Sprintf("result=%d %s passed: %s", len(st.Checks), shipPlural(len(st.Checks), "check", "checks"), strings.Join(parts, ", ")),
		"pr=" + st.URL,
		"head=" + st.Head,
		"status=verified",
	}
}

// changedCapabilities are the live capabilities whose nodes the change
// set's ops create, change or move, in op order, each once.
func changedCapabilities(m Model, ops []op.Op) []id.ID {
	out := []id.ID{}
	for _, o := range ops {
		switch o.Type {
		case op.TypeNodeCreate, op.TypeNodeUpdate, op.TypeNodeMove:
		default:
			continue
		}
		n, ok := m.Node(o.Node)
		if !ok || n.Kind != op.KindCapability || n.Retired || slices.Contains(out, n.ID) {
			continue
		}
		out = append(out, n.ID)
	}
	return out
}

// runShipLand records the PR as Evidence, lands the change set, renders,
// and commits and pushes to the PR branch. It runs in the change set's
// worktree (ctx.Root) on its branch. Each part is skipped when an earlier
// call did it, so a resumed run moves on.
func runShipLand(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.land"); err != nil {
		return nil, err
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	st, err := readShipState(t.statePath())
	if err != nil {
		return nil, err
	}
	if st.PR == 0 || st.Head == "" || len(st.Checks) == 0 {
		return nil, Usagef("change set %s has no passed pull request recorded: run rhiza ship check first", t.cs).
			WithNext(shipCheckNext)
	}
	head, err := ctx.Git("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return nil, Errorf("ship land: which branch: %v", err)
	}
	if b := strings.TrimSpace(head); b != t.branch {
		return nil, Conflictf("ship land commits to the change set's branch %s, but %s is on %s", t.branch, ctx.Root, b)
	}

	resp := ShipLandResponse{Change: t.cs.String(), PR: st.PR, Evidence: shipEvidenceHandle(st.PR, t.slug), Rendered: []string{}}
	if err := shipRecordEvidence(ctx, t, st, &resp); err != nil {
		return nil, err
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	if c, _ := w.changeSet(t.cs); c.State != changeset.StateLanded {
		land, _ := Lookup("land")
		lctx := *ctx
		lctx.ChangeSet, lctx.ChangeRef, lctx.Record = t.cs, "", false
		if _, err := land.Run(&lctx, LifecycleRequest{}); err != nil {
			return nil, err
		}
		resp.Landed = true
	}
	r, err := runRender(ctx, RenderRequest{})
	if err != nil {
		return nil, err
	}
	rr := r.(RenderResponse)
	resp.Rendered = append(append(resp.Rendered, rr.Written...), rr.Removed...)

	if _, err := ctx.Git("add", "--", ".rhiza"); err != nil {
		return nil, Errorf("%v", err)
	}
	staged, err := ctx.Git("diff", "--cached", "--name-only", "--", ".rhiza")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	if strings.TrimSpace(staged) != "" {
		if _, err := ctx.Git("commit", "-m", shipLandMessage, "-m", "Change: "+t.cs.String(), "--", ".rhiza"); err != nil {
			return nil, Errorf("%v", err)
		}
		if resp.Commit, err = gitHead(ctx); err != nil {
			return nil, err
		}
	}
	pushed, err := shipPush(ctx, t.branch)
	if err != nil {
		return nil, err
	}
	if pushed {
		if resp.Pushed, err = gitHead(ctx); err != nil {
			return nil, err
		}
		st.Pushed = resp.Pushed
		if err := writeShipState(t.statePath(), st); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// shipRecordEvidence writes the Evidence node and its verifies links in
// one apply: into the change set while it is in flight, else (the agent
// landed it anyway) into a record change set. A live node already holding
// the handle with the same pr field means an earlier call wrote it.
func shipRecordEvidence(ctx *Ctx, t shipTarget, st shipState, resp *ShipLandResponse) error {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return err
	}
	c, ok := w.changeSet(t.cs)
	if !ok {
		return NotFoundf("no change set %s", t.cs).WithNext(changeLsNext)
	}
	record := false
	switch c.State {
	case changeset.StateAccepted:
	case changeset.StateLanded:
		record = true
	default:
		return Conflictf("cannot ship change set %s: it is %s, not accepted", t.cs, c.State).WithNext(legalNext(c))
	}
	m := w.model(t.cs)
	resp.Verifies = changedCapabilities(m, c.Ops)
	if n, taken := m.HandleTaken(resp.Evidence); taken {
		if node, ok := m.Node(n); ok && node.Kind == op.KindEvidence && fieldText(node.Fields["pr"]) == st.URL {
			return nil
		}
		resp.Evidence = uniqueHandle(resp.Evidence, func(h string) bool { _, ok := m.HandleTaken(h); return ok })
	}
	items := []applyItem{{op: "add", build: func(b writeBase) (edit, error) {
		return buildAdd(b, AddRequest{Kind: string(op.KindEvidence), Handle: resp.Evidence, Fields: shipEvidenceFields(st)})
	}}}
	for _, capID := range resp.Verifies {
		req := LinkRequest{From: resp.Evidence, Type: string(op.LinkVerifies), To: capID.String()}
		items = append(items, applyItem{op: "link", build: func(b writeBase) (edit, error) { return buildLink(b, req) }})
	}
	wctx := *ctx
	wctx.ChangeSet, wctx.ChangeRef, wctx.Record = t.cs, "", record
	intent := fmt.Sprintf("Record pull request #%d's CI for change set %s", st.PR, t.cs)
	r, err := runWriteOn(&wctx, "apply", intent, func(w *workspace, cs id.ID) (edit, error) {
		return buildApply(w, cs, ctx.IDs, items)
	})
	if err != nil {
		return err
	}
	resp.EvidenceWritten = true
	resp.EvidenceChange = r.(WriteResponse).ChangeSet.String()
	return nil
}

// gitHead is the commit HEAD names in ctx.Root.
func gitHead(ctx *Ctx) (string, error) {
	out, err := ctx.Git("rev-parse", "HEAD")
	if err != nil {
		return "", Errorf("%v", err)
	}
	return strings.TrimSpace(out), nil
}

// shipPush pushes branch when it has no upstream (setting origin's) or is
// ahead of it, and reports whether it pushed. Hooks run; it never forces.
func shipPush(ctx *Ctx, branch string) (bool, error) {
	if _, err := ctx.Git("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}"); err != nil {
		if _, err := ctx.Git("push", "-u", "origin", branch); err != nil {
			return false, Errorf("%v", err)
		}
		return true, nil
	}
	ahead, err := ctx.Git("rev-list", "--count", "@{upstream}..HEAD")
	if err != nil {
		return false, Errorf("%v", err)
	}
	if n, _ := strconv.Atoi(strings.TrimSpace(ahead)); n == 0 {
		return false, nil
	}
	if _, err := ctx.Git("push"); err != nil {
		return false, Errorf("%v", err)
	}
	return true, nil
}

// gitAt runs git with args in dir through ctx.Exec, for a checkout other
// than ctx.Root.
func gitAt(ctx *Ctx, dir string, args ...string) (string, error) {
	if ctx.Exec == nil {
		return "", Errorf("shipping needs git, and no programs can be run here")
	}
	out, err := ctx.Exec("git", args, dir, nil)
	return string(out), err
}

// runShipCloseOut removes the worktree, closes the Rex tabs of the change
// set's sessions, deletes the branch, locally and on origin, then
// fast-forwards main, all from the main checkout. Each part
// already done is reported as such. A dirty worktree fails git's remove,
// and so the step: nothing is forced.
func runShipCloseOut(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.close-out"); err != nil {
		return nil, err
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	st, err := readShipState(t.statePath())
	if err != nil {
		return nil, err
	}
	if st.MergeCommit == "" {
		return nil, Usagef("change set %s has no merge commit recorded: its pull request is not merged", t.cs).
			WithNext(shipMergeNext)
	}
	resp := ShipCloseOutResponse{Change: t.cs.String(), Branch: t.branch, Worktree: t.worktree, Main: t.main, Steps: []ShipStep{}}
	step := func(name string, done bool) {
		resp.Steps = append(resp.Steps, ShipStep{Step: name, Result: map[bool]string{true: "done", false: "already"}[done]})
	}

	if _, err := os.Stat(t.worktree); err == nil {
		if _, err := gitAt(ctx, t.main, "worktree", "remove", t.worktree); err != nil {
			return nil, Errorf("%v", err)
		}
		step("remove-worktree", true)
	} else if errors.Is(err, fs.ErrNotExist) {
		step("remove-worktree", false)
	} else {
		return nil, Errorf("worktree %s: %v", t.worktree, err)
	}

	recs, err := session.NewStore(session.Dir(t.common)).List()
	if err != nil {
		return nil, Errorf("%v", err)
	}
	for _, rec := range recs {
		if rec.Change != t.cs.String() || rec.Tab.Window == "" {
			continue
		}
		closed, err := closeSessionTab(ctx, rec, t.main)
		if err != nil {
			return nil, err
		}
		resp.Steps = append(resp.Steps, ShipStep{Step: "close-tab", Tab: rec.Tab.Window,
			Result: map[bool]string{true: "done", false: "already"}[closed]})
	}

	if _, err := gitAt(ctx, t.main, "rev-parse", "--verify", "--quiet", "refs/heads/"+t.branch); err == nil {
		if _, err := gitAt(ctx, t.main, "branch", "-D", t.branch); err != nil {
			return nil, Errorf("%v", err)
		}
		step("delete-branch", true)
	} else {
		step("delete-branch", false)
	}

	if _, err := gitAt(ctx, t.main, "push", "origin", "--delete", t.branch); err == nil {
		step("delete-remote-branch", true)
	} else if strings.Contains(err.Error(), "remote ref does not exist") {
		step("delete-remote-branch", false)
	} else {
		return nil, Errorf("%v", err)
	}

	resp.Detail = shipMainBlocked(ctx, t.main)
	if resp.Detail == "" {
		if _, err := gitAt(ctx, t.main, "pull", "--ff-only", "origin", "main"); err != nil {
			resp.Detail = err.Error()
		}
	}
	if resp.Detail != "" {
		resp.Detail = "main was not updated: " + resp.Detail
		resp.Steps = append(resp.Steps, ShipStep{Step: "update-main", Result: "failed", Detail: resp.Detail})
		return resp, nil
	}
	step("update-main", true)
	resp.MainUpdated = true
	return resp, nil
}

// shipMainBlocked says why the main checkout cannot be fast-forwarded, or
// "" if it can: it must be on main with no changes to tracked files
// (untracked files do not stop a fast-forward that does not touch them).
func shipMainBlocked(ctx *Ctx, main string) string {
	b, err := gitAt(ctx, main, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return fmt.Sprintf("%s is not on a branch", main)
	}
	if b = strings.TrimSpace(b); b != "main" {
		return fmt.Sprintf("%s is on %s, not main", main, b)
	}
	status, err := gitAt(ctx, main, "status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return err.Error()
	}
	if strings.TrimSpace(status) != "" {
		return fmt.Sprintf("%s has uncommitted changes", main)
	}
	return ""
}

// runShipCommand runs .rhiza/workflows/ship.sh with sh in the main
// checkout (ctx.Root when there is no git), streaming its output to
// ctx.Stderr, and exits with its exit code. No script is no ship command,
// and exits 0.
func runShipCommand(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.command"); err != nil {
		return nil, err
	}
	dir := shipCommandDir(ctx)
	resp := ShipCommandResponse{Script: filepath.Join(dir, filepath.FromSlash(shipScript)), Dir: dir}
	if fi, err := os.Stat(resp.Script); errors.Is(err, fs.ErrNotExist) {
		return resp, nil
	} else if err != nil {
		return nil, Errorf("ship command: %v", err)
	} else if fi.IsDir() {
		return nil, Errorf("ship command %s is a folder", resp.Script)
	}
	if ctx.Start == nil {
		return nil, Errorf("ship command: no programs can be run here")
	}
	p, err := ctx.Start(Program{Name: "sh", Args: []string{"-c", `sh "$0" 2>&1`, resp.Script}, Dir: dir})
	if err != nil {
		return nil, Errorf("ship command: %v", err)
	}
	var out io.Writer = io.Discard
	if ctx.Stderr != nil {
		out = ctx.Stderr
	}
	if _, err := io.Copy(out, p.Output()); err != nil {
		p.Kill()
	}
	code, err := p.Wait()
	if err != nil {
		return nil, Errorf("ship command: %v", err)
	}
	resp.Ran, resp.Exit = true, code
	return resp, nil
}

// shipCommandDir is the project root in the main checkout: the first entry
// of git worktree list plus ctx.Root's prefix in its checkout. Without
// git it is ctx.Root.
func shipCommandDir(ctx *Ctx) string {
	if ctx.Git == nil {
		return ctx.Root
	}
	out, err := ctx.Git("worktree", "list", "--porcelain")
	if err != nil {
		return ctx.Root
	}
	main, ok := firstWorktree(out)
	if !ok {
		return ctx.Root
	}
	prefix, err := ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return ctx.Root
	}
	return filepath.Join(main, filepath.FromSlash(strings.TrimSpace(prefix)))
}

// Human lists what ship land did.
func (r ShipLandResponse) Human(w io.Writer) {
	if r.EvidenceWritten {
		fmt.Fprintf(w, "Recorded PR #%d as evidence %s, verifying %d %s.\n", r.PR, r.Evidence,
			len(r.Verifies), shipPlural(len(r.Verifies), "capability", "capabilities"))
	} else {
		fmt.Fprintf(w, "Evidence %s was already recorded.\n", r.Evidence)
	}
	if r.Landed {
		fmt.Fprintf(w, "Landed change set %s.\n", r.Change)
	} else {
		fmt.Fprintf(w, "Change set %s was already landed.\n", r.Change)
	}
	fmt.Fprintf(w, "Rendered %d %s.\n", len(r.Rendered), shipPlural(len(r.Rendered), "file", "files"))
	if r.Commit != "" {
		fmt.Fprintf(w, "Committed %s.\n", shortSHA(r.Commit))
	} else {
		fmt.Fprintln(w, "Nothing to commit.")
	}
	if r.Pushed != "" {
		fmt.Fprintf(w, "Pushed %s.\n", shortSHA(r.Pushed))
	} else {
		fmt.Fprintln(w, "Nothing to push.")
	}
}

// shipStepWords says each close-out step for people: done, then already.
var shipStepWords = map[string][2]string{
	"remove-worktree":      {"Removed worktree %s.", "Worktree %s was already gone."},
	"close-tab":            {"Closed Rex tab %s.", "Rex tab %s was already closed."},
	"delete-branch":        {"Deleted branch %s.", "Branch %s was already gone."},
	"delete-remote-branch": {"Deleted origin/%s.", "origin/%s was already gone."},
	"update-main":          {"Updated main in %s.", "main in %s was already current."},
}

// Human lists each close-out step.
func (r ShipCloseOutResponse) Human(w io.Writer) {
	for _, s := range r.Steps {
		if s.Result == "failed" {
			fmt.Fprintf(w, "%s.\n", strings.TrimSuffix(s.Detail, "."))
			continue
		}
		arg := map[string]string{"remove-worktree": r.Worktree, "delete-branch": r.Branch,
			"delete-remote-branch": r.Branch, "update-main": r.Main, "close-tab": s.Tab}[s.Step]
		words := shipStepWords[s.Step]
		f := words[0]
		if s.Result != "done" {
			f = words[1]
		}
		fmt.Fprintf(w, f+"\n", arg)
	}
}

// Human says whether the ship command ran and how it ended.
func (r ShipCommandResponse) Human(w io.Writer) {
	if !r.Ran {
		fmt.Fprintf(w, "no ship command (%s)\n", shipScript)
		return
	}
	if r.Exit != 0 {
		fmt.Fprintf(w, "Ran %s: exit %d.\n", shipScript, r.Exit)
		return
	}
	fmt.Fprintf(w, "Ran %s.\n", shipScript)
}
