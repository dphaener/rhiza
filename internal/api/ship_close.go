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
	"github.com/dphaener/rhiza/internal/daemon"
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

// ShipCloseOutRequest takes the merge commit from the state file ship.merge
// wrote.
type ShipCloseOutRequest struct {
	Forget bool `json:"forget,omitempty" flag:"forget" desc:"Delete the change set's ship state file once closed out, so a later ship run opens and merges its own pull request (for a plan run, whose merge does not land the change set)"`
}

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
// worktree (ctx.Root) on its branch. In Rhiza's own repo the land and
// render run with rhiza built from the PR head (see ship_build.go). Each
// part is skipped when an earlier call did it, so a resumed run moves on.
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
	if err := shipFollowPR(ctx, t); err != nil {
		return nil, err
	}

	// In Rhiza's own repo, build the PR head first: a branch that changes
	// rendering lands with its own, and a failed build writes nothing.
	bin, err := shipHeadBuild(ctx, t)
	if err != nil {
		return nil, err
	}

	resp := ShipLandResponse{Change: t.cs.String(), PR: st.PR, Evidence: shipEvidenceHandle(st.PR, t.slug), Rendered: []string{}}
	if err := shipRecordEvidence(ctx, t, st, &resp); err != nil {
		return nil, err
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	actor, err := shipActorHandle(w, t.cs, ctx)
	if err != nil {
		return nil, err
	}
	if c, _ := w.changeSet(t.cs); c.State != changeset.StateLanded {
		if err := shipLandWith(ctx, t, bin); err != nil {
			return nil, err
		}
		resp.Landed = true
	}
	rendered, err := shipRenderWith(ctx, bin)
	if err != nil {
		return nil, err
	}
	resp.Rendered = append(resp.Rendered, rendered...)

	if _, err := ctx.Git("add", "--", ".rhiza"); err != nil {
		return nil, Errorf("%v", err)
	}
	staged, err := ctx.Git("diff", "--cached", "--name-only", "--", ".rhiza")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	if strings.TrimSpace(staged) != "" {
		if _, err := ctx.Git("commit", "-m", shipLandMessage, "-m", "Change: "+t.cs.String()+"\nActor: "+actor, "--", ".rhiza"); err != nil {
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
		if !shipConflictUnresolved(ctx, t, st.MergeConflict) {
			st.MergeConflict = nil
		}
		if err := writeShipState(t.statePath(), st); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// shipActorHandle is the handle of the actor making ship's land commit,
// resolved as every write resolves it (RHIZA_ACTOR, else git's
// user.email), for the commit's Actor: trailer. A node with no handle
// gives its ID.
func shipActorHandle(w *workspace, cs id.ID, ctx *Ctx) (string, error) {
	actor, _, err := w.resolveActor(cs, ctx)
	if err != nil {
		return "", err
	}
	if n, ok := w.model(cs).Node(actor); ok && n.Handle != "" {
		return n.Handle, nil
	}
	return actor.String(), nil
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
		return buildAdd(b, AddRequest{Kind: string(op.KindEvidence), Handle: resp.Evidence, Fields: shipEvidenceFields(st), lenient: true})
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

// shipFollowPR brings the worktree up to origin's copy of the change set's
// branch before ship commits or pushes, for a branch someone else (the
// orchestrator fixing a conflict, say) pushed to. When origin's branch
// holds HEAD, or origin has no such branch, nothing changes. When it
// contains HEAD and is ahead, a clean worktree is fast-forwarded to it; a
// worktree with uncommitted changes is refused, naming them. When the two
// have diverged, local commits origin lacks are refused, naming them.
func shipFollowPR(ctx *Ctx, t shipTarget) error {
	// A failed fetch leaves what is known; the push then says what is wrong.
	ctx.Git("fetch", "-q", "origin")
	out, err := ctx.Git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+t.branch)
	if err != nil {
		return nil
	}
	remote := strings.TrimSpace(out)
	if remote == "" {
		return nil
	}
	if _, err := ctx.Git("merge-base", "--is-ancestor", remote, "HEAD"); err == nil {
		return nil
	}
	if _, err := ctx.Git("merge-base", "--is-ancestor", "HEAD", remote); err != nil {
		local, _ := ctx.Git("rev-list", "--oneline", remote+"..HEAD")
		return Conflictf("origin's %s has commits that %s does not, and %s has local commits that origin's lacks: %s; reconcile them before shipping",
			t.branch, ctx.Root, ctx.Root, strings.Join(gitLines(local), "; "))
	}
	dirty, err := ctx.Git("status", "--porcelain", "--untracked-files=no")
	if err != nil {
		return Errorf("%v", err)
	}
	if files := gitLines(dirty); len(files) > 0 {
		return Conflictf("origin's %s has moved to %s, but %s has uncommitted changes that stop it following: %s; commit or discard them",
			t.branch, shortSHA(remote), ctx.Root, strings.Join(files, "; "))
	}
	if _, err := ctx.Git("merge", "--ff-only", remote); err != nil {
		return Errorf("%v", err)
	}
	return nil
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
func runShipCloseOut(ctx *Ctx, req any) (any, error) {
	r := req.(ShipCloseOutRequest)
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
		dropIncomingLogCopies(ctx, t.main)
		if _, err := gitAt(ctx, t.main, "pull", "--ff-only", "origin", "main"); err != nil {
			resp.Detail = err.Error()
		}
	}
	if resp.Detail != "" {
		resp.Detail = "main was not updated: " + resp.Detail
		resp.Steps = append(resp.Steps, ShipStep{Step: "update-main", Result: "failed", Detail: resp.Detail})
		// Record it where rhiza health and the inbox read it, not only here.
		if err := daemon.WriteMainStale(t.common, daemon.MainStale{Since: ctx.Now(), Change: t.cs.String(), Detail: resp.Detail}); err != nil {
			return nil, Errorf("record the failed pull: %v", err)
		}
		return resp, nil
	}
	if err := daemon.ClearMainStale(t.common); err != nil {
		return nil, Errorf("clear the failed pull: %v", err)
	}
	step("update-main", true)
	resp.MainUpdated = true
	// Forget only now: a close-out whose pull failed keeps its state, so a
	// retry finds the merge commit and finishes.
	if r.Forget {
		if err := os.Remove(t.statePath()); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, Errorf("forget ship state: %v", err)
		}
	}
	return resp, nil
}

// dropIncomingLogCopies removes the untracked .rhiza/log files in main that
// origin/main holds byte for byte, so the fast-forward may write them as
// tracked files. A change set opened in the main checkout leaves its first
// entries there. Files that differ, or that origin/main lacks, stay; the
// pull then fails and names them. Best effort: any git failure drops nothing.
func dropIncomingLogCopies(ctx *Ctx, main string) {
	if _, err := gitAt(ctx, main, "fetch", "origin", "main"); err != nil {
		return
	}
	out, err := gitAt(ctx, main, "ls-files", "--others", "--exclude-standard", "-z", "--", ".rhiza/log")
	if err != nil {
		return
	}
	for _, f := range strings.Split(out, "\x00") {
		if f == "" {
			continue
		}
		ours, err := gitAt(ctx, main, "hash-object", "--", f)
		if err != nil {
			continue
		}
		theirs, err := gitAt(ctx, main, "rev-parse", "--verify", "--quiet", "origin/main:"+f)
		if err != nil || strings.TrimSpace(ours) != strings.TrimSpace(theirs) {
			continue
		}
		_ = os.Remove(filepath.Join(main, f))
	}
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
