package api

// This file holds ship.open-pr: after a change set's review is journaled,
// it pushes the change set's branch and opens its pull request, titled by
// the change set with the run's review summary as the body, or, when one
// is already open, pushes and replaces that body. It records the PR in
// the ship state file for the steps after it (ship.check, ship.land).

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "ship.open-pr",
		Summary:  "Push a change set's branch and open its pull request, titled by the change set with the run's review summary as the body, or push and replace the body of the one already open",
		Request:  reflect.TypeFor[ShipOpenPRRequest](),
		Response: reflect.TypeFor[ShipOpenPRResponse](),
		Run:      runShipOpenPR,
	})
}

// shipOpenPRNext is the command that opens the pull request.
const shipOpenPRNext = "rhiza ship open-pr"

// ShipOpenPRRequest takes nothing: the run is RHIZA_RUN's and the change
// set its own.
type ShipOpenPRRequest struct{}

// ShipOpenPRResponse is what ship.open-pr did.
type ShipOpenPRResponse struct {
	Change  string `json:"change" desc:"The change set"`
	Branch  string `json:"branch" desc:"The change set's branch"`
	PR      int    `json:"pr" desc:"The pull request's number"`
	URL     string `json:"url" desc:"The pull request's URL"`
	Head    string `json:"head" desc:"The branch's head commit, as pushed"`
	Pushed  bool   `json:"pushed" desc:"True if this call pushed the branch; false if origin already had its head"`
	Created bool   `json:"created" desc:"True if this call opened the pull request; false if it replaced the body of the one already open"`
}

// reviewJournaled reports whether st holds the outcome of some review
// check: the review ran.
func reviewJournaled(st run.State) bool {
	for _, c := range st.Checks {
		if len(c.Outcomes) > 0 {
			return true
		}
	}
	return false
}

// runShipOpenPR pushes the branch and opens or updates the pull request.
// It runs in the change set's worktree (ctx.Root) on its branch, inside an
// unfinished run whose review is journaled. A second call, after a new
// commit, pushes it and edits the body; it never creates a second PR.
func runShipOpenPR(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.open-pr"); err != nil {
		return nil, err
	}
	rr, err := openRun(ctx, "ship open-pr", shipOpenPRNext, false)
	if err != nil {
		return nil, err
	}
	if !reviewJournaled(rr.state) && len(rr.state.Rounds) == 0 {
		return nil, Usagef("run %s has no review journaled: ship open-pr opens the pull request once the review has run", rr.id).
			WithNext("rhiza run status " + rr.id)
	}
	t, err := findShipTarget(ctx)
	if err != nil {
		return nil, err
	}
	if t.cs != rr.change {
		return nil, Conflictf("run %s is on change set %s, not %s", rr.id, rr.change, t.cs)
	}
	head, err := ctx.Git("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return nil, Errorf("ship open-pr: which branch: %v", err)
	}
	if b := strings.TrimSpace(head); b != t.branch {
		return nil, Conflictf("ship open-pr pushes the change set's branch %s, but %s is on %s", t.branch, ctx.Root, b)
	}
	// The branch's PR is looked up first: a PR that is another change set's
	// is refused before anything is fetched into this worktree or pushed
	// onto it.
	dir := t.dir()
	var v prView
	out, err := gh(ctx, dir, "pr", "view", t.branch, "--json", "number,url,state,headRefOid,body,statusCheckRollup")
	exists := true
	switch {
	case noPR(err):
		exists = false
	case err != nil:
		return nil, Errorf("%v", err)
	default:
		if v, err = parsePRView(out); err != nil {
			return nil, Errorf("%v", err)
		}
		if err := requirePROwned(t, v); err != nil {
			return nil, err
		}
	}
	if err := shipFollowPR(ctx, t); err != nil {
		return nil, err
	}
	sha, err := gitHead(ctx)
	if err != nil {
		return nil, err
	}
	resp := ShipOpenPRResponse{Change: t.cs.String(), Branch: t.branch, Head: sha}
	if resp.Pushed, err = shipPush(ctx, t.branch); err != nil {
		return nil, err
	}

	body, err := os.CreateTemp("", "rhiza-pr-body-*.md")
	if err != nil {
		return nil, Errorf("ship open-pr: %v", err)
	}
	defer os.Remove(body.Name())
	summary := reviewSummary(rr.state)
	if len(rr.state.Rounds) > 0 {
		// A plan run: its rounds are the review.
		w, err := openCtxWorkspace(ctx)
		if err != nil {
			return nil, err
		}
		summary = planSummary(rr.state, w.changeSets())
	}
	_, werr := io.WriteString(body, summary+prBodyRef(t.ref))
	if cerr := body.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		return nil, Errorf("ship open-pr: write the body: %v", werr)
	}

	if !exists {
		if _, err := gh(ctx, dir, "pr", "create", "--base", shipMain, "--head", t.branch, "--title", t.title, "--body-file", body.Name()); err != nil {
			return nil, Errorf("%v", err)
		}
		resp.Created = true
		if out, err = gh(ctx, dir, "pr", "view", t.branch, "--json", "number,url,state,headRefOid,body,statusCheckRollup"); err != nil {
			return nil, Errorf("%v", err)
		}
		if v, err = parsePRView(out); err != nil {
			return nil, Errorf("%v", err)
		}
	}
	// A branch whose only pull request was merged or closed before this run
	// (a plan's PR on the change set's branch, say) gets a new one: gh then
	// names the open one for the branch.
	if !strings.EqualFold(v.State, "OPEN") && !resp.Created {
		if _, err := gh(ctx, dir, "pr", "create", "--base", shipMain, "--head", t.branch, "--title", t.title, "--body-file", body.Name()); err != nil {
			return nil, Errorf("%v", err)
		}
		resp.Created = true
		if out, err = gh(ctx, dir, "pr", "view", t.branch, "--json", "number,url,state,headRefOid,body,statusCheckRollup"); err != nil {
			return nil, Errorf("%v", err)
		}
		if v, err = parsePRView(out); err != nil {
			return nil, Errorf("%v", err)
		}
		if err := requirePROwned(t, v); err != nil {
			return nil, err
		}
	}
	if !strings.EqualFold(v.State, "OPEN") {
		return nil, Conflictf("pull request #%d for branch %s is %s, not open", v.Number, t.branch, v.State)
	}
	if !resp.Created {
		if _, err := gh(ctx, dir, "pr", "edit", strconv.Itoa(v.Number), "--body-file", body.Name()); err != nil {
			return nil, Errorf("%v", err)
		}
	}
	resp.PR, resp.URL = v.Number, v.URL

	// Nothing is checked or merged until ship.check records CI on the head
	// just pushed.
	path := t.statePath()
	st, err := readShipState(path)
	if err != nil {
		return nil, err
	}
	st.PR, st.URL, st.Branch, st.Pushed = v.Number, v.URL, t.branch, sha
	st.Run = shipRunID(ctx)
	// A merge commit left by an earlier PR (a plan's) is not this PR's.
	st.Head, st.Checks, st.MergeCommit, st.Main = "", nil, "", nil
	if err := writeShipState(filepath.Clean(path), st); err != nil {
		return nil, err
	}
	return resp, nil
}

// Human is one line naming the PR and what was done to it.
func (r ShipOpenPRResponse) Human(w io.Writer) {
	verb := "Updated"
	if r.Created {
		verb = "Opened"
	}
	pushed := "branch already pushed"
	if r.Pushed {
		pushed = "pushed " + shortSHA(r.Head)
	}
	fmt.Fprintf(w, "%s PR #%d (%s): %s.\n", verb, r.PR, pushed, r.URL)
}

// prBodyRef is the PR body's ref section for a change set's ref: nothing
// when it has none, else one line holding the ref in a single inline code
// span. The ref is untrusted text, so its backticks and line breaks are
// written as escapes (a backtick as \x60, a line break as \n or \r, other
// control characters as \u00XX) and it can neither close the span nor
// start markdown or a mention.
func prBodyRef(ref string) string {
	if ref == "" {
		return ""
	}
	var b strings.Builder
	for _, r := range ref {
		switch {
		case r == '`':
			b.WriteString(`\x60`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case unicode.IsControl(r):
			fmt.Fprintf(&b, `\u%04X`, r)
		default:
			b.WriteRune(r)
		}
	}
	return "\nRef: `" + b.String() + "`\n"
}
