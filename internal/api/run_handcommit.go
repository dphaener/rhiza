package api

// A run's hand commits: the commits on the run's branch since it started
// that the run did not make (Decision hand-commit-by-trailer). The engine
// journals each once as a hand-commit intervention (run.HandCommitScanner).

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// handCommitFormat is git log's format for the scan: the commit, its
// parents and its author's name, then the raw message, the fields split by
// US (0x1f) and the commit ended by RS (0x1e).
const handCommitFormat = "%H%x1f%P%x1f%an%x1f%B%x1e"

// branchHead is the commit branch is at; "" when there is no branch or git
// (ctx.Git, in the main checkout: a branch's ref is shared by every
// worktree) cannot resolve it. A run records it as its Base when it starts.
func branchHead(ctx *Ctx, branch string) string {
	if branch == "" || ctx.Git == nil {
		return ""
	}
	out, err := ctx.Git("rev-parse", "--verify", "-q", "refs/heads/"+branch+"^{commit}")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// HandCommits is run.HandCommitScanner: the hand commits on the run's
// branch since its Base, oldest first. A run with no branch or no recorded
// Base has none. The run's leaves are read from the model, so a model that
// cannot be read is an error and nothing is journalled.
func (r *stepRunner) HandCommits(st run.State) ([]run.Intervention, error) {
	if st.Branch == "" || st.Base == "" || r.ctx.Git == nil {
		return nil, nil
	}
	var leaves []run.Leaf
	var own []string
	if st.Change != "" {
		ls, err := r.Leaves(st.Change)
		if err != nil {
			return nil, err
		}
		leaves = ls
		if own, err = r.processActors(st.Change); err != nil {
			return nil, err
		}
	}
	// Git runs in the main checkout: the run's branch ref is shared by every
	// worktree, and the worktree may be gone before the run finishes.
	return handCommits(r.ctx.Git, st, leaves, r.sessions, own...)
}

// processActors are the handle and ID of the actor the engine's own
// process writes as (RHIZA_ACTOR, else the actor node of git's
// user.email), read only: the actor a run's command steps, such as ship
// rebase's render commit, put in their Actor: trailers, since they inherit
// the process's environment. None when that actor is not in the model yet.
func (r *stepRunner) processActors(change string) ([]string, error) {
	cs, err := id.Parse(change)
	if err != nil {
		return nil, fmt.Errorf("run %s names change set %q: %v", r.run, change, err)
	}
	if _, err := os.Stat(r.worktree); err != nil {
		return nil, nil
	}
	v, err := r.readModel("process-actor:"+cs.String(), cs, func(w *workspace) (any, error) {
		m := w.model(cs)
		var n id.ID
		switch {
		case !r.ctx.Actor.IsZero():
			n = r.ctx.Actor
		case r.ctx.ActorRef != "":
			a, err := resolveActorRef(m, r.ctx.ActorRef)
			if err != nil {
				return []string(nil), nil
			}
			n = a
		default:
			a, ok := m.ActorByEmail(gitValue(r.ctx.GitConfig, "user.email"))
			if !ok {
				return []string(nil), nil
			}
			n = a.ID
		}
		return []string{nodeName(m, n), n.String()}, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]string), nil
}

// handCommits lists, oldest first, the commits on st's branch after
// st.Base, along its first parents only (so the base branch's commits a
// merge brings in are never scanned), that the run did not make. A commit
// is the run's when it is a merge whose second parent is on the base branch
// (origin/HEAD, origin/main or main), or a trailer names the run: Leaf: a
// leaf of leaves (handle or ID), Review: the run, or Actor: the actor of one
// of the run's sessions in sessions or one of own (the engine process's own
// actor, see processActors). Each other commit is a hand commit,
// its actor the commit's Actor: trailer, else its author.
func handCommits(git func(...string) (string, error), st run.State, leaves []run.Leaf, sessions *session.Store, own ...string) ([]run.Intervention, error) {
	out, err := git("log", "--first-parent", "--reverse", "--format="+handCommitFormat, st.Base+"..refs/heads/"+st.Branch, "--")
	if err != nil {
		return nil, fmt.Errorf("list the commits on branch %s since %s: %v", st.Branch, shortRev(st.Base), err)
	}
	if strings.TrimSpace(out) == "" {
		return nil, nil
	}
	ours := map[string]bool{}
	for _, l := range leaves {
		ours[l.Handle], ours[l.ID] = true, true
	}
	delete(ours, "")
	actors := map[string]bool{}
	for _, a := range own {
		if a != "" {
			actors[a] = true
		}
	}
	if sessions != nil {
		recs, err := sessions.List()
		if err != nil {
			return nil, fmt.Errorf("read the run's sessions: %v", err)
		}
		for _, rec := range recs {
			if rec.Run == st.ID && rec.Actor != "" {
				actors[rec.Actor] = true
			}
		}
	}
	var bases []string
	for _, ref := range baseCandidates {
		if _, err := git("rev-parse", "--verify", "-q", ref+"^{commit}"); err == nil {
			bases = append(bases, ref)
		}
	}
	onBase := func(c string) bool {
		for _, ref := range bases {
			if _, err := git("merge-base", "--is-ancestor", c, ref); err == nil {
				return true
			}
		}
		return false
	}
	var found []run.Intervention
	for _, rec := range strings.Split(out, "\x1e") {
		f := strings.SplitN(strings.TrimLeft(rec, "\n"), "\x1f", 4)
		if len(f) < 4 || f[0] == "" {
			continue
		}
		hash, parents, author, msg := f[0], strings.Fields(f[1]), f[2], f[3]
		if len(parents) >= 2 && onBase(parents[1]) {
			continue
		}
		cs := parseTrailerLog(hash + "\x1f" + msg)
		if len(cs) != 1 {
			continue
		}
		c := cs[0]
		if slices.ContainsFunc(c.Leaves, func(l string) bool { return ours[l] }) ||
			slices.Contains(c.Reviews, st.ID) ||
			slices.ContainsFunc(c.Actors, func(a string) bool { return actors[a] }) {
			continue
		}
		actor := author
		if len(c.Actors) > 0 {
			actor = c.Actors[0]
		}
		if strings.TrimSpace(actor) == "" {
			actor = "unknown"
		}
		found = append(found, run.Intervention{
			Kind: run.InterventionHandCommit, Actor: actor, Commit: hash,
			Detail: fmt.Sprintf("commit %s on %s: %s", shortRev(hash), st.Branch, c.Subject),
		})
	}
	return found, nil
}
