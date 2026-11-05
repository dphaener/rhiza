package api

// This file holds run.milestone-check, the plan workflow's step between
// the interview and the review: every change set the plan makes must
// serve a milestone, directly or through an intent under one, so a plan
// that serves none fails loudly before it is reviewed or gated.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
)

func init() {
	register(Operation{
		Name:     "run.milestone-check",
		Summary:  "Check that each change set the plan of the run RHIZA_RUN names makes serves a milestone, directly or through an intent under one: pass, or fail naming each that does not",
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
	ChangeSets []string `json:"change_sets" desc:"The plan's change sets checked, by ID, each serving a milestone"`
}

// Human prints the passed check.
func (r RunMilestoneCheckResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Run %s: each of the plan's %d change set(s) serves a milestone\n", r.Run, len(r.ChangeSets))
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
	resp := RunMilestoneCheckResponse{Run: rr.id, ChangeSets: []string{}}
	var bad []string
	for _, c := range plan {
		if !inFlight[c.ID] {
			continue // landed already: it is in no roadmap
		}
		resp.ChangeSets = append(resp.ChangeSets, c.ID.String())
		// Each change set's served intents resolve in its own what-if view,
		// so an intent the plan wrote into it under a milestone counts.
		found, err := milestoneServed(ctx, cmd, c.ID)
		if err != nil {
			return nil, err
		}
		if !found {
			bad = append(bad, fmt.Sprintf("- %s (%s)", c.ID, c.Open.DisplayTitle()))
		}
	}
	if len(bad) > 0 {
		return nil, Errorf("%s: %d change set(s) of the plan of run %s serve no milestone, directly or through an intent under one:\n%s\n"+
			"Fix: link each with `rhiza link <change set> serves <intent> --change <change set>` to an intent under a milestone (or a milestone itself); the plan was started without a milestone, so the interview must name one",
			cmd, len(bad), rr.id, strings.Join(bad, "\n")).WithNext("rhiza run status " + rr.id)
	}
	return resp, nil
}

// milestoneServed reports whether change set cs, in its what-if view,
// serves a milestone directly or through an intent under one.
func milestoneServed(ctx *Ctx, cmd string, cs id.ID) (bool, error) {
	res, in, _, err := planRead(ctx, cmd, viewWhatIfPref+cs.String())
	if err != nil {
		return false, err
	}
	in.intentPath = planAncestorPaths(in.intentPath)
	for _, p := range res.order {
		if p.info.ID != cs {
			continue
		}
		for _, i := range p.serves {
			if _, ok := planMilestoneOf(i, in); ok {
				return true, nil
			}
		}
	}
	return false, nil
}
