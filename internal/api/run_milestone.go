package api

// This file holds run.milestone-check, the plan workflow's step between
// the interview and the review: every change set the plan makes must
// serve the milestone the run was started with, directly or through an
// intent under it, and, when the repo's policy never builds work that
// serves no milestone (backlog = "skip"), serve some milestone. Otherwise a
// plan serving none passes: the daemon builds it by the schedule rule.

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/policy"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.milestone-check",
		Summary:  "Check the change sets the plan of the run RHIZA_RUN names makes: each must serve the milestone the run was started with (RHIZA_MILESTONE), directly or through an intent under it, and some milestone when the repo's .rhiza/policy.toml sets backlog = \"skip\"; otherwise pass. Fails naming each that does not",
		Request:  reflect.TypeFor[RunMilestoneCheckRequest](),
		Response: reflect.TypeFor[RunMilestoneCheckResponse](),
		Run:      runRunMilestoneCheck,
	})
}

// RunMilestoneCheckRequest takes nothing: the run is RHIZA_RUN's.
type RunMilestoneCheckRequest struct{}

// RunMilestoneCheckResponse is the passed check.
type RunMilestoneCheckResponse struct {
	Run        string   `json:"run" desc:"The run's full ID"`
	ChangeSets []string `json:"change_sets" desc:"The plan's change sets checked, by ID, checked against the milestone and policy rules"`
}

// Human prints the passed check.
func (r RunMilestoneCheckResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Run %s: the plan's %d change set(s) pass the milestone check\n", r.Run, len(r.ChangeSets))
}

func runRunMilestoneCheck(ctx *Ctx, _ any) (any, error) {
	const cmd = "run milestone-check"
	rr, err := openReviewRun(ctx, cmd, "rhiza describe run.milestone-check")
	if err != nil {
		return nil, err
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	plan, err := planChangeSets(w.changeSets(), w.model(rr.change), rr.change)
	if err != nil {
		return nil, Errorf("%s: %v", cmd, err).WithNext("rhiza run status " + rr.id)
	}
	res, _, _, err := planRead(ctx, cmd, "")
	if err != nil {
		return nil, err
	}
	inFlight := map[id.ID]bool{}
	for _, p := range res.order {
		inFlight[p.info.ID] = true
	}
	pol, err := policy.Read(os.DirFS(ctx.Root))
	if err != nil {
		return nil, Errorf("%s: %v", cmd, err).WithNext("rhiza check")
	}
	var want id.ID
	if m := strings.TrimSpace(rr.state.Milestone); m != "" {
		if want, err = id.Parse(m); err != nil {
			return nil, Errorf("run %s names milestone %q: %v", rr.id, m, err)
		}
	}
	resp := RunMilestoneCheckResponse{Run: rr.id, ChangeSets: []string{}}
	var bad []string
	var badNum []string
	var landed Model
	if rr.state.MilestoneNumber > 0 {
		landed = newModel(w.diskLog(), id.ID{})
	}
	for _, c := range plan {
		if !inFlight[c.ID] {
			continue // landed already: it is in no roadmap
		}
		resp.ChangeSets = append(resp.ChangeSets, c.ID.String())
		// Each change set's served intents resolve in its own what-if view,
		// so an intent the plan wrote into it under a milestone counts.
		served, err := milestonesServed(ctx, cmd, c.ID)
		if err != nil {
			return nil, err
		}
		if landed != nil {
			badNum = append(badNum, newMilestoneNumberProblems(w.model(c.ID), landed, rr.state.MilestoneNumber)...)
		}
		ok := true
		switch {
		case !want.IsZero():
			ok = slices.Contains(served, want)
		case pol.Backlog == policy.BacklogSkip:
			ok = len(served) > 0
		}
		if !ok {
			bad = append(bad, fmt.Sprintf("- %s (%s)", c.ID, c.Open.DisplayTitle()))
		}
	}
	if len(badNum) > 0 {
		n := rr.state.MilestoneNumber
		slices.Sort(badNum)
		badNum = slices.Compact(badNum)
		return nil, Errorf("%s: the plan of run %s does not use the milestone number %d the run reserved:\n%s\n"+
			"Fix: name the new milestone intent milestone-%d with field milestone=%d and prefix its done points' handles m%d- (RHIZA_MILESTONE_NUMBER is %d)",
			cmd, rr.id, n, strings.Join(badNum, "\n"), n, n, n, n).WithNext("rhiza run status " + rr.id)
	}
	if len(bad) > 0 {
		if !want.IsZero() {
			return nil, Errorf("%s: %d change set(s) of the plan of run %s do not serve milestone %s, directly or through an intent under it:\n%s\n"+
				"Fix: link each with `rhiza link <change set> serves <intent> --change <change set>` to an intent under that milestone (or the milestone itself)",
				cmd, len(bad), rr.id, want, strings.Join(bad, "\n")).WithNext("rhiza run status " + rr.id)
		}
		return nil, Errorf("%s: %d change set(s) of the plan of run %s serve no milestone, directly or through an intent under one:\n%s\n"+
			"Fix: link each with `rhiza link <change set> serves <intent> --change <change set>` to an intent under a milestone (or a milestone itself); the repo's %s sets backlog = %q, so work that serves no milestone is never built",
			cmd, len(bad), rr.id, strings.Join(bad, "\n"), policy.Path, policy.BacklogSkip).WithNext("rhiza run status " + rr.id)
	}
	return resp, nil
}

// milestonesServed is the milestones change set cs, in its what-if view,
// serves directly or through an intent under one.
func milestonesServed(ctx *Ctx, cmd string, cs id.ID) ([]id.ID, error) {
	res, in, _, err := planRead(ctx, cmd, viewWhatIfPref+cs.String())
	if err != nil {
		return nil, err
	}
	var out []id.ID
	in.intentPath = planAncestorPaths(in.intentPath)
	for _, p := range res.order {
		if p.info.ID != cs {
			continue
		}
		for _, i := range p.serves {
			if m, ok := planMilestoneOf(i, in); ok {
				out = append(out, m)
			}
		}
	}
	return out, nil
}

// milestoneNumberLock names the lock a plan run holds from drawing its
// milestone number to journaling it.
const milestoneNumberLock = "milestone-number"

// reserveMilestoneNumber draws the milestone number a plan run for a new
// milestone reserves: one more than the highest of every landed milestone
// intent's milestone field and every number reserved in any run journal of
// rs, ended and failed runs included. The caller holds milestoneNumberLock
// and journals the number before releasing it.
func reserveMilestoneNumber(ctx *Ctx, rs *run.Store) (int, error) {
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return 0, err
	}
	top := 0
	m := newModel(w.diskLog(), id.ID{})
	var walk func(parent id.ID)
	walk = func(parent id.ID) {
		for _, n := range m.Children(op.KindIntent, parent) {
			if planMilestoneField(n.Fields) {
				in := planInputs{milestoneNums: map[id.ID]float64{}}
				planNoteMilestoneNum(in, n.ID, n.Fields)
				if f, ok := in.milestoneNums[n.ID]; ok && int(f) > top {
					top = int(f)
				}
			}
			walk(n.ID)
		}
	}
	walk(id.ID{})
	js, err := rs.List()
	if err != nil {
		return 0, Errorf("read the run journals: %v", err)
	}
	for _, j := range js {
		if n := j.State().MilestoneNumber; n > top {
			top = n
		}
	}
	return top + 1, nil
}

// newMilestoneNumberProblems lists, one line each, how the new milestone
// intents of m (those landed lacks that carry a milestone field) and their
// done points break the reserved number n: the handle must be milestone-n,
// the milestone field n, and each done point's handle prefixed mn-.
func newMilestoneNumberProblems(m, landed Model, n int) []string {
	var out []string
	var walk func(parent id.ID, under *Node)
	walk = func(parent id.ID, under *Node) {
		for _, c := range m.Children(op.KindIntent, parent) {
			if c.Retired {
				continue
			}
			if under != nil {
				if want := fmt.Sprintf("m%d-", n); !strings.HasPrefix(c.Handle, want) {
					out = append(out, fmt.Sprintf("- done point %q of milestone %q lacks the %q prefix (expected milestone number %d)", c.Handle, under.Handle, want, n))
				}
				walk(c.ID, under)
				continue
			}
			if _, ok := landed.Node(c.ID); ok || !planMilestoneField(c.Fields) {
				walk(c.ID, nil)
				continue
			}
			in := planInputs{milestoneNums: map[id.ID]float64{}}
			planNoteMilestoneNum(in, c.ID, c.Fields)
			if h := fmt.Sprintf("milestone-%d", n); c.Handle != h {
				out = append(out, fmt.Sprintf("- milestone intent %q has the wrong handle (expected milestone number %d: %q)", c.Handle, n, h))
			}
			if f, ok := in.milestoneNums[c.ID]; !ok || int(f) != n || f != float64(int(f)) {
				out = append(out, fmt.Sprintf("- milestone intent %q has the wrong milestone field (expected milestone number %d)", c.Handle, n))
			}
			cc := c
			walk(c.ID, &cc)
		}
	}
	walk(id.ID{}, nil)
	return out
}
