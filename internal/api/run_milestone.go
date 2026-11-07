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
	"github.com/dphaener/rhiza/internal/policy"
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
