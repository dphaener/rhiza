package api

// This file holds the ship operations that read and act on a change set's
// pull request: ship.check asks GitHub (through gh) whether the PR, or
// main after the merge, has every CI check passed, and ship.merge squash
// merges the PR pinned to the head whose checks passed. The built-in
// ship-change workflow calls them from its command and wait steps.
//
// Steps are plain commands, so they pass values to each other through a
// state file per change set, <git common dir>/rhiza/ship/<change set>.json:
// ship.check records the PR, its checked head and checks; ship.rebase
// records the head it pushed; ship.merge records the merge commit;
// ship.check --main records main's checks. Like
// a run's journal it is work-layer state, never committed.

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "ship.check",
		Summary:  "Check that a change set's pull request has every CI check passed, or that main after its merge has the tested head's tree; with --settled, rebase a pull request GitHub reports conflicting",
		Write:    true,
		Request:  reflect.TypeFor[ShipCheckRequest](),
		Response: reflect.TypeFor[ShipCheckResponse](),
		Run:      runShipCheck,
	})
	register(Operation{
		Name:     "ship.merge",
		Summary:  "Squash merge a change set's pull request, pinned to the head whose checks passed; merges nothing if origin/main is not in that head",
		Request:  reflect.TypeFor[ShipMergeRequest](),
		Response: reflect.TypeFor[ShipMergeResponse](),
		Run:      runShipMerge,
	})
}

// Next commands of the ship operations.
const (
	shipCheckNext = "rhiza ship check"
	shipMergeNext = "rhiza ship merge"
)

// shipDir is where the state files live, under the git common dir.
const shipDir = "rhiza/ship"

// ShipCheckRequest asks about the change set's PR, or with Main about the
// merge commit on main.
type ShipCheckRequest struct {
	Settled bool `json:"settled,omitempty" flag:"settled" desc:"Exit 0 once every check has finished, passed or not, instead of only when all passed; for a wait step; when GitHub reports the pull request CONFLICTING it rebases the branch onto main as ship rebase does, then reports not settled, and exits 75 with a stop: line on a conflict in code"`
	Main    bool `json:"main,omitempty" flag:"main" desc:"Check that the merge commit rhiza ship merge recorded has the tree of the head whose checks passed, locally with git and no wait on CI, instead of checking the pull request; a differing tree dispatches CI on the merge commit (once) and checks its runs as it does a pull request's, and a stop that cannot be waited out exits 75 with a stop: line"`
	PROpen  bool `json:"pr_open,omitempty" flag:"pr-open" desc:"Exit 0 once a pull request from the change set's branch is open, whatever its CI; records nothing (for an agent step's until)"`
}

// ShipCheckResponse is what ship.check found.
type ShipCheckResponse struct {
	Change    string      `json:"change" desc:"The change set"`
	Branch    string      `json:"branch" desc:"The change set's branch"`
	Main      bool        `json:"main,omitempty" desc:"True if it checked the merge commit on main rather than the pull request"`
	PR        int         `json:"pr,omitempty" desc:"The pull request's number; 0 if none was found"`
	URL       string      `json:"url,omitempty" desc:"The pull request's URL"`
	State     string      `json:"state,omitempty" desc:"The pull request's state, such as OPEN or MERGED"`
	Commit    string      `json:"commit,omitempty" desc:"The commit checked: the PR's head, or the merge commit"`
	Checks    []ShipCheck `json:"checks" desc:"The checks reported for the commit"`
	Settled   bool        `json:"settled" desc:"True if there is an open PR (or a merge commit) and every check on it has finished"`
	Green     bool        `json:"green" desc:"True if settled and every check passed (success, neutral or skipped)"`
	Failed    []string    `json:"failed,omitempty" desc:"The names of the settled checks that did not pass"`
	Passed    bool        `json:"passed" desc:"True if the check asked for holds: open with --pr-open, settled with --settled, else green; the command exits 1 otherwise"`
	Detail    string      `json:"detail,omitempty" desc:"Why it did not pass, such as no open pull request"`
	Written   bool        `json:"written,omitempty" desc:"True if the result was recorded in the state file for the next steps"`
	Rebased   string      `json:"rebased,omitempty" desc:"With --settled on a conflicting pull request: the head pushed after rebasing the branch onto main"`
	Conflicts []string    `json:"conflicts,omitempty" desc:"With --settled on a conflicting pull request: the files whose conflict outside the rendered files stopped the rebase; the command exits 75"`
	Question  string      `json:"question,omitempty" desc:"The handle of the question raised about those conflicts"`
	Stop      string      `json:"stop,omitempty" desc:"With --main: why main cannot be confirmed, such as no workflow to dispatch CI on a merge commit whose tree differs from the tested head's, or CI on it failing again after a rerun; the command exits 75"`
	Reran     []string    `json:"reran,omitempty" desc:"The GitHub Actions run IDs whose failed jobs this check reran, once per head, instead of failing; the check then reports not settled"`
}

// ShipCheck is one CI check on a commit, normalised from gh's check runs,
// status contexts and workflow runs.
type ShipCheck struct {
	Name       string `json:"name" desc:"The check's name"`
	Status     string `json:"status" desc:"COMPLETED once it has finished, else its status, such as IN_PROGRESS or PENDING"`
	Conclusion string `json:"conclusion,omitempty" desc:"How it finished, such as SUCCESS or FAILURE; empty until it has"`
	// RunID is the GitHub Actions run the check belongs to, when known
	// (from the check's detailsUrl, or a workflow run's databaseId). It
	// is what a rerun needs; checks without one are never rerun.
	RunID string `json:"-"`
}

// ShipMergeRequest takes nothing: the PR and its head come from the state
// file ship.check wrote.
type ShipMergeRequest struct{}

// ShipMergeResponse is the merged PR, or the refusal to merge a head that
// lacks origin/main.
type ShipMergeResponse struct {
	Change        string `json:"change" desc:"The change set"`
	Moved         bool   `json:"moved,omitempty" desc:"True if nothing was merged because origin/main is not in the checked head: another change set merged after its checks; the command exits 1 and the run goes back to rebase and check again"`
	PR            int    `json:"pr" desc:"The pull request merged"`
	Head          string `json:"head" desc:"The head commit the merge was pinned to"`
	MergeCommit   string `json:"merge_commit" desc:"The squash commit on main"`
	AlreadyMerged bool   `json:"already_merged,omitempty" desc:"True if an earlier ship merge had merged it and recorded the merge commit"`
}

// ExitCode is 1 when ship merge merged nothing because main moved.
func (r ShipMergeResponse) ExitCode() int {
	if r.Moved {
		return 1
	}
	return 0
}

// ExitCode is 0 when the check asked for holds, else 1.
func (r ShipCheckResponse) ExitCode() int {
	if len(r.Conflicts) > 0 || r.Stop != "" {
		return run.StopExit
	}
	if r.Passed {
		return 0
	}
	return 1
}

// shipTarget is what the ship operations know about a change set: where
// its branch, worktree and main checkout are, and its state file.
type shipTarget struct {
	cs       id.ID
	title    string
	slug     string
	branch   string // cs/<slug>
	worktree string // the change set's worktree; it may not exist (yet, or any more)
	main     string // the main checkout: the first entry of git worktree list
	common   string // the absolute git common dir
}

// findShipTarget resolves the change set ctx names (ctx.ChangeSet, else
// --change or RHIZA_CHANGE, else the current one) in any state, landed
// included, and derives its branch and worktree as session.start does.
// It writes nothing.
func findShipTarget(ctx *Ctx) (shipTarget, error) {
	var t shipTarget
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return t, err
	}
	cs, err := w.effectiveChangeSet(ctx)
	if err != nil {
		return t, err
	}
	if cs.IsZero() {
		return t, Usagef("no change set: name one with --change or make one current").WithNext(changeLsNext)
	}
	c, _ := w.changeSet(cs)
	common, err := gitCommonDir(ctx)
	if err != nil {
		return t, err
	}
	out, err := ctx.Git("worktree", "list", "--porcelain")
	if err != nil {
		return t, Errorf("list worktrees: %v", err)
	}
	main, ok := firstWorktree(out)
	if !ok {
		return t, Errorf("git worktree list names no main checkout")
	}
	t = shipTarget{cs: cs, title: c.Open.DisplayTitle(), main: main, common: common}
	t.slug = session.Slug(t.title)
	t.branch = session.Branch(t.slug)
	mainRoot := filepath.Dir(common)
	t.worktree = session.Worktree(mainRoot, filepath.Base(mainRoot), t.slug)
	return t, nil
}

// firstWorktree is the path of the first entry of `git worktree list
// --porcelain`: the main checkout.
func firstWorktree(porcelain string) (string, bool) {
	for l := range strings.SplitSeq(porcelain, "\n") {
		if p, ok := strings.CutPrefix(l, "worktree "); ok && p != "" {
			return p, true
		}
	}
	return "", false
}

// dir is where commands about the change set run: its worktree while it
// exists, else the main checkout (close-out removes the worktree).
func (t shipTarget) dir() string {
	if fi, err := os.Stat(t.worktree); err == nil && fi.IsDir() {
		return t.worktree
	}
	return t.main
}

// statePath is the change set's state file.
func (t shipTarget) statePath() string {
	return filepath.Join(t.common, filepath.FromSlash(shipDir), t.cs.String()+".json")
}

// shipState is what one ship step leaves for the next, in the state file.
type shipState struct {
	PR          int         `json:"pr,omitempty"`
	URL         string      `json:"url,omitempty"`
	Branch      string      `json:"branch,omitempty"`
	Head        string      `json:"head,omitempty"`
	Checks      []ShipCheck `json:"checks,omitempty"`
	MergeCommit string      `json:"merge_commit,omitempty"`
	Main        []ShipCheck `json:"main,omitempty"`
	// Pushed is the head ship.land or ship.rebase pushed to the PR
	// branch. Until GitHub reports it as the PR's head, ship.check counts
	// the PR as not settled, so a wait never passes on the checks of the
	// head before. ship.rebase also clears Head and Checks, so nothing
	// lands or merges until ship.check records the rebased head's CI.
	Pushed string `json:"pushed,omitempty"`
	// RerunHead is the head (or merge commit) whose failed checks were
	// rerun, and Reruns the Actions run IDs rerun for it: a head is rerun
	// once, so a second failure on it fails.
	RerunHead string   `json:"rerun_head,omitempty"`
	Reruns    []string `json:"reruns,omitempty"`
	// DispatchCommit is the merge commit CI was dispatched on, and
	// Dispatched the workflow files dispatched for it, once each.
	DispatchCommit string   `json:"dispatch_commit,omitempty"`
	Dispatched     []string `json:"dispatched,omitempty"`
}

// readShipState reads the state file at path; a missing file is the zero
// state.
func readShipState(path string) (shipState, error) {
	var st shipState
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return st, nil
	}
	if err != nil {
		return st, Errorf("read ship state: %v", err)
	}
	if err := json.Unmarshal(b, &st); err != nil {
		return st, Errorf("ship state %s: %v", path, err)
	}
	return st, nil
}

// writeShipState writes st to path through a temporary file and a rename,
// creating the folders.
func writeShipState(path string, st shipState) error {
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return Errorf("ship state: %v", err)
	}
	if err := writeAtomic(path, append(b, '\n')); err != nil {
		return Errorf("write ship state: %v", err)
	}
	return nil
}

// gh runs gh with args in dir and returns its standard output. A failure
// carries gh's standard error (see Ctx.Exec).
func gh(ctx *Ctx, dir string, args ...string) ([]byte, error) {
	if ctx.Exec == nil {
		return nil, Errorf("shipping needs gh, and no programs can be run here")
	}
	return ctx.Exec("gh", args, dir, nil)
}

// noPR reports whether a gh failure says there is no pull request.
func noPR(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "no pull requests found")
}

// prView is the part of `gh pr view --json
// number,url,state,headRefOid,mergeable,statusCheckRollup` ship.check uses.
type prView struct {
	Number int
	URL    string
	State  string
	Head   string
	// Mergeable is GitHub's verdict: MERGEABLE, CONFLICTING or UNKNOWN.
	Mergeable string
	Checks    []ShipCheck
}

// parsePRView parses gh pr view's JSON. Rollup items are check runs
// ({name, status, conclusion}) or status contexts ({context, state});
// both become ShipChecks. A status context in state PENDING or EXPECTED
// has not finished; any other state is its conclusion.
func parsePRView(b []byte) (prView, error) {
	var raw struct {
		Number     int    `json:"number"`
		URL        string `json:"url"`
		State      string `json:"state"`
		HeadRefOid string `json:"headRefOid"`
		Mergeable  string `json:"mergeable"`
		Rollup     []struct {
			Name       string `json:"name"`
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
			Context    string `json:"context"`
			State      string `json:"state"`
			DetailsURL string `json:"detailsUrl"`
		} `json:"statusCheckRollup"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return prView{}, fmt.Errorf("gh pr view: %v", err)
	}
	v := prView{Number: raw.Number, URL: raw.URL, State: raw.State, Head: raw.HeadRefOid, Mergeable: raw.Mergeable, Checks: []ShipCheck{}}
	for _, r := range raw.Rollup {
		if r.Context != "" || (r.Status == "" && r.State != "") {
			c := ShipCheck{Name: r.Context, Status: "COMPLETED", Conclusion: r.State}
			if s := strings.ToUpper(r.State); s == "PENDING" || s == "EXPECTED" || s == "" {
				c.Status, c.Conclusion = cmp.Or(s, "PENDING"), ""
			}
			v.Checks = append(v.Checks, c)
			continue
		}
		v.Checks = append(v.Checks, ShipCheck{Name: r.Name, Status: r.Status, Conclusion: r.Conclusion, RunID: actionsRunID(r.DetailsURL)})
	}
	return v, nil
}

// actionsRunRe finds the run ID in an Actions details URL such as
// https://github.com/o/r/actions/runs/42/job/9.
var actionsRunRe = regexp.MustCompile(`/actions/runs/([0-9]+)`)

// actionsRunID is the GitHub Actions run ID in a check's details URL, or
// "" when the URL is not an Actions one.
func actionsRunID(url string) string {
	if m := actionsRunRe.FindStringSubmatch(url); m != nil {
		return m[1]
	}
	return ""
}

// parseRunList parses `gh run list --commit <sha> --json
// databaseId,name,status,conclusion`.
func parseRunList(b []byte) ([]ShipCheck, error) {
	var raw []struct {
		DatabaseID int64  `json:"databaseId"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("gh run list: %v", err)
	}
	out := make([]ShipCheck, 0, len(raw))
	for _, r := range raw {
		c := ShipCheck{Name: r.Name, Status: r.Status, Conclusion: r.Conclusion}
		if r.DatabaseID != 0 {
			c.RunID = strconv.FormatInt(r.DatabaseID, 10)
		}
		out = append(out, c)
	}
	return out, nil
}

// passingConclusions are the conclusions that count as passed, upper-cased.
var passingConclusions = map[string]bool{"SUCCESS": true, "NEUTRAL": true, "SKIPPED": true}

// checksVerdict classifies checks. No checks at all is not settled (CI
// has not reported yet, or there is none: a wait times out rather than
// merging unchecked); neither is any check whose status is not completed.
// Settled checks are green when every conclusion passed; failed names the
// others, in order.
func checksVerdict(checks []ShipCheck) (settled, green bool, failed []string) {
	if len(checks) == 0 {
		return false, false, nil
	}
	for _, c := range checks {
		if !strings.EqualFold(c.Status, "completed") {
			return false, false, nil
		}
	}
	for _, c := range checks {
		if !passingConclusions[strings.ToUpper(c.Conclusion)] {
			failed = append(failed, c.Name)
		}
	}
	return true, len(failed) == 0, failed
}

// runShipCheck checks the PR (or with Main the merge commit) and, when it
// is green and Settled was not asked for, records it in the state file.
// Nothing is recorded on failure. With PROpen it asks only whether the PR
// is open.
func runShipCheck(ctx *Ctx, req any) (any, error) {
	r := req.(ShipCheckRequest)
	if r.PROpen && (r.Main || r.Settled) {
		return nil, Usagef("--pr-open asks only whether the pull request is open; it takes neither --main nor --settled")
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	resp := ShipCheckResponse{Change: t.cs.String(), Branch: t.branch, Main: r.Main, Checks: []ShipCheck{}}
	path := t.statePath()
	st, err := readShipState(path)
	if err != nil {
		return nil, err
	}
	if r.Main {
		if st.MergeCommit == "" {
			return nil, Usagef("change set %s has no merge commit recorded: its pull request is not merged", t.cs).
				WithNext(shipMergeNext)
		}
		resp.PR, resp.URL, resp.Commit = st.PR, st.URL, st.MergeCommit
		return shipCheckMainTree(ctx, t, st, resp, path, r.Settled)
	} else {
		out, err := gh(ctx, t.dir(), "pr", "view", t.branch, "--json", "number,url,state,headRefOid,mergeable,statusCheckRollup")
		if noPR(err) {
			resp.Detail = fmt.Sprintf("no open pull request for branch %s", t.branch)
			return resp, nil
		}
		if err != nil {
			return nil, Errorf("%v", err)
		}
		v, err := parsePRView(out)
		if err != nil {
			return nil, Errorf("%v", err)
		}
		resp.PR, resp.URL, resp.State, resp.Commit, resp.Checks = v.Number, v.URL, v.State, v.Head, v.Checks
		if !strings.EqualFold(v.State, "OPEN") {
			resp.Detail = fmt.Sprintf("pull request #%d for branch %s is %s, not open", v.Number, t.branch, v.State)
			return resp, nil
		}
		if r.PROpen {
			resp.Passed = true
			return resp, nil
		}
		if r.Settled && strings.EqualFold(v.Mergeable, "CONFLICTING") {
			return shipSettleConflict(ctx, resp)
		}
		if st.Pushed != "" && v.Head != st.Pushed {
			// Someone else may have pushed on top of the run's head:
			// follow a head that contains the pushed one.
			if shipHeadHolds(ctx, t, v.Head, st.Pushed) {
				st.Pushed = v.Head
			} else if shipHeadHolds(ctx, t, st.Pushed, v.Head) {
				resp.Detail = fmt.Sprintf("not settled: the pull request's head is %s, not the pushed %s yet", shortSHA(v.Head), shortSHA(st.Pushed))
				return resp, nil
			} else {
				resp.Detail = fmt.Sprintf("not settled: the pull request's head %s drops the pushed %s", shortSHA(v.Head), shortSHA(st.Pushed))
				return resp, nil
			}
		}
	}
	resp.Settled, resp.Green, resp.Failed = checksVerdict(resp.Checks)
	switch {
	case r.Settled:
		resp.Passed = resp.Settled
	default:
		resp.Passed = resp.Green
	}
	if resp.Settled && !resp.Green {
		if err := shipRerunFailed(ctx, t, &st, &resp, path); err != nil {
			return nil, err
		}
	}
	if !resp.Settled {
		if resp.Detail == "" {
			resp.Detail = notSettled(resp.Checks)
		}
	} else if !resp.Green {
		resp.Detail = "failed: " + strings.Join(resp.Failed, ", ")
	}
	if !resp.Passed || r.Settled {
		return resp, nil
	}
	if r.Main {
		st.Main = resp.Checks
	} else {
		st.PR, st.URL, st.Branch, st.Head, st.Checks = resp.PR, resp.URL, t.branch, resp.Commit, resp.Checks
	}
	if err := writeShipState(path, st); err != nil {
		return nil, err
	}
	resp.Written = true
	return resp, nil
}

// shipCheckMainTree confirms main locally, with git and no wait on CI: the
// merge commit's tree must equal the tree of the PR head whose checks
// passed. ship.merge only merges a head that holds origin/main, so a squash
// of it has the head's tree; a differing tree means main changed between
// that check and the merge. Then main is not what the checks passed on, and
// shipCheckMainDispatch runs CI on the merge commit instead.
func shipCheckMainTree(ctx *Ctx, t shipTarget, st shipState, resp ShipCheckResponse, path string, settled bool) (any, error) {
	if st.Head == "" {
		return nil, Usagef("change set %s has no checked pull request head recorded: run rhiza ship check first", t.cs).
			WithNext(shipCheckNext)
	}
	dir := t.dir()
	if _, err := gitAt(ctx, dir, "fetch", "-q", "origin"); err != nil {
		return nil, Errorf("%v", err)
	}
	tree := func(commit string) (string, error) {
		out, err := gitAt(ctx, dir, "rev-parse", "--verify", "--quiet", commit+"^{tree}")
		if err != nil {
			return "", Errorf("tree of %s: %v", shortSHA(commit), err)
		}
		return strings.TrimSpace(out), nil
	}
	mt, err := tree(st.MergeCommit)
	if err != nil {
		return nil, err
	}
	ht, err := tree(st.Head)
	if err != nil {
		return nil, err
	}
	if mt != ht {
		return shipCheckMainDispatch(ctx, t, st, resp, path, settled)
	}
	resp.Settled, resp.Green, resp.Passed = true, true, true
	return resp, nil
}

// shipRerunFailed reruns, once per head, the failed jobs of the Actions
// runs of the settled checks that failed: it records the runs in the state
// file and turns resp into not settled, so a wait keeps waiting for the
// rerun. A head already rerun, and failed checks with no Actions run ID,
// are left to fail as before.
func shipRerunFailed(ctx *Ctx, t shipTarget, st *shipState, resp *ShipCheckResponse, path string) error {
	if st.RerunHead == resp.Commit && len(st.Reruns) > 0 {
		return nil
	}
	var runs []string
	seen := map[string]bool{}
	for _, c := range resp.Checks {
		if c.RunID == "" || seen[c.RunID] || passingConclusions[strings.ToUpper(c.Conclusion)] {
			continue
		}
		seen[c.RunID] = true
		runs = append(runs, c.RunID)
	}
	if len(runs) == 0 {
		return nil
	}
	// Record before rerunning: a crash between the two then costs a
	// missed rerun, never a second one.
	st.RerunHead, st.Reruns = resp.Commit, runs
	if err := writeShipState(path, *st); err != nil {
		return err
	}
	for _, id := range runs {
		if _, err := gh(ctx, t.dir(), "run", "rerun", id, "--failed"); err != nil {
			return Errorf("%v", err)
		}
	}
	resp.Reran = runs
	resp.Settled, resp.Green, resp.Passed = false, false, false
	resp.Detail = fmt.Sprintf("not settled: reran the failed jobs of %s %s after %s failed", shipPlural(len(runs), "run", "runs"), strings.Join(runs, ", "), strings.Join(resp.Failed, ", "))
	resp.Failed = nil
	return nil
}

// shipSettleConflict runs the ship rebase logic for a pull request GitHub
// reports CONFLICTING and reports it not settled: after a push the next
// check sees the new head; a conflict in code raises the rebase-conflict
// question and stops the run (exit 75).
func shipSettleConflict(ctx *Ctx, resp ShipCheckResponse) (any, error) {
	out, err := runShipRebase(ctx, ShipRebaseRequest{})
	if err != nil {
		return nil, err
	}
	rb := out.(ShipRebaseResponse)
	switch {
	case len(rb.Conflicts) > 0:
		resp.Conflicts, resp.Question, resp.Detail = rb.Conflicts, rb.Question, rb.Detail
	case rb.Pushed != "":
		resp.Rebased = rb.Pushed
		resp.Detail = fmt.Sprintf("not settled: the pull request conflicted; rebased onto main and pushed %s", shortSHA(rb.Pushed))
	default:
		resp.Detail = "not settled: GitHub reports the pull request conflicting, though the branch holds main"
	}
	return resp, nil
}

// shipHeadHolds reports whether commit anc is an ancestor of (or equal to)
// desc, after fetching the branch so a head pushed from elsewhere is
// known locally. Any git failure counts as not holding.
func shipHeadHolds(ctx *Ctx, t shipTarget, desc, anc string) bool {
	if _, err := gitAt(ctx, t.dir(), "fetch", "-q", "origin", t.branch); err != nil {
		return false
	}
	_, err := gitAt(ctx, t.dir(), "merge-base", "--is-ancestor", anc, desc)
	return err == nil
}

// notSettled says why checks are not settled.
func notSettled(checks []ShipCheck) string {
	if len(checks) == 0 {
		return "not settled: no checks reported yet"
	}
	done := 0
	for _, c := range checks {
		if strings.EqualFold(c.Status, "completed") {
			done++
		}
	}
	return fmt.Sprintf("not settled: %d of %d checks finished", done, len(checks))
}

// runShipMerge squash merges the PR ship.check recorded, pinned to the
// head whose checks passed, then records the merge commit. It never passes
// --admin, --auto or --delete-branch: the merge waits for nobody and
// close-out deletes the branch. Once a merge commit is recorded it merges
// nothing again, so a resumed run moves on.
func runShipMerge(ctx *Ctx, _ any) (any, error) {
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	path := t.statePath()
	st, err := readShipState(path)
	if err != nil {
		return nil, err
	}
	if st.PR == 0 || st.Head == "" {
		return nil, Usagef("change set %s has no checked pull request head recorded: run rhiza ship check first", t.cs).
			WithNext(shipCheckNext)
	}
	resp := ShipMergeResponse{Change: t.cs.String(), PR: st.PR, Head: st.Head}
	if st.MergeCommit != "" {
		resp.MergeCommit, resp.AlreadyMerged = st.MergeCommit, true
		return resp, nil
	}
	pr := fmt.Sprint(st.PR)
	dir := t.dir()
	// A head that lacks origin/main would squash to a tree nobody tested:
	// merge nothing and let the run rebase and wait on the new head.
	if _, err := gitAt(ctx, dir, "fetch", "-q", "origin"); err != nil {
		return nil, Errorf("%v", err)
	}
	if _, err := gitAt(ctx, dir, "rev-parse", "--verify", "--quiet", "refs/remotes/origin/"+shipMain); err != nil {
		return nil, Errorf("ship merge: origin has no %s branch: %v", shipMain, err)
	}
	if _, err := gitAt(ctx, dir, "merge-base", "--is-ancestor", "refs/remotes/origin/"+shipMain, st.Head); err != nil {
		resp.Moved = true
		return resp, nil
	}
	if _, err := gh(ctx, dir, "pr", "merge", pr, "--squash", "--match-head-commit", st.Head); err != nil {
		return nil, Errorf("%v", err)
	}
	out, err := gh(ctx, dir, "pr", "view", pr, "--json", "mergeCommit,state")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	var v struct {
		State       string `json:"state"`
		MergeCommit *struct {
			Oid string `json:"oid"`
		} `json:"mergeCommit"`
	}
	if err := json.Unmarshal(out, &v); err != nil {
		return nil, Errorf("gh pr view: %v", err)
	}
	if !strings.EqualFold(v.State, "MERGED") || v.MergeCommit == nil || v.MergeCommit.Oid == "" {
		return nil, Errorf("pull request #%d is %s after the merge, with no merge commit", st.PR, cmp.Or(v.State, "in an unknown state"))
	}
	st.MergeCommit = v.MergeCommit.Oid
	if err := writeShipState(path, st); err != nil {
		return nil, err
	}
	resp.MergeCommit = st.MergeCommit
	return resp, nil
}

// shortSHA is the first seven characters of a commit ID.
func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

// Human is one line: what was checked and how it stands.
func (r ShipCheckResponse) Human(w io.Writer) {
	what := fmt.Sprintf("PR #%d at %s", r.PR, shortSHA(r.Commit))
	if r.Main {
		what = fmt.Sprintf("main at %s", shortSHA(r.Commit))
	}
	switch {
	case r.Stop != "":
		fmt.Fprintf(w, "stop: %s\n", r.Stop)
	case len(r.Conflicts) > 0:
		fmt.Fprintf(w, "%s: %s\n", what, r.Detail)
		fmt.Fprintf(w, "stop: %s; answer question %s, then resume the run\n", r.Detail, r.Question)
	case r.PR == 0 && !r.Main:
		fmt.Fprintf(w, "%s.\n", r.Detail)
	case r.State != "" && !strings.EqualFold(r.State, "OPEN"):
		fmt.Fprintf(w, "%s.\n", r.Detail)
	case r.Passed && !r.Settled:
		// Only --pr-open passes without settled checks.
		fmt.Fprintf(w, "%s: open from %s\n", what, r.Branch)
	case r.Green && r.Main && len(r.Checks) == 0:
		fmt.Fprintf(w, "%s: the merge commit's tree is the tested head's\n", what)
	case r.Green:
		fmt.Fprintf(w, "%s: %d %s passed\n", what, len(r.Checks), shipPlural(len(r.Checks), "check", "checks"))
	default:
		fmt.Fprintf(w, "%s: %s\n", what, r.Detail)
	}
}

// shipPlural is one when n is 1, else many.
func shipPlural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

// Human is one line naming the PR, its head and the merge commit.
func (r ShipMergeResponse) Human(w io.Writer) {
	if r.Moved {
		fmt.Fprintf(w, "Not merged: origin/%s has moved past PR #%d's checked head %s; rebase and check again.\n", shipMain, r.PR, shortSHA(r.Head))
		return
	}
	verb := "Merged"
	if r.AlreadyMerged {
		verb = "Already merged"
	}
	fmt.Fprintf(w, "%s PR #%d at %s as %s.\n", verb, r.PR, shortSHA(r.Head), shortSHA(r.MergeCommit))
}
