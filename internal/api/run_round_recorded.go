package api

// This file holds run.round-recorded, the check a plan run's review step's
// `until` runs: the reviewer's visit must have journaled a findings round
// (empty for a clean review), or the step does not pass and review-check
// would judge the previous round again. session done runs the same check
// for a reviewer's session (see reviewRoundMissing).

import (
	"fmt"
	"io"
	"reflect"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "run.round-recorded",
		Summary:  "Check that the current visit of the plan review step of the run RHIZA_RUN names journaled a findings round (empty for a clean review); fails naming rhiza run findings if not",
		Request:  reflect.TypeFor[RunRoundRecordedRequest](),
		Response: reflect.TypeFor[RunRoundRecordedResponse](),
		Run:      runRunRoundRecorded,
	})
}

// RunRoundRecordedRequest takes nothing: the run is RHIZA_RUN's.
type RunRoundRecordedRequest struct{}

// RunRoundRecordedResponse is the passed check.
type RunRoundRecordedResponse struct {
	Run   string `json:"run" desc:"The run's full ID"`
	Round int    `json:"round" desc:"The latest plan review round's number, journaled in this visit"`
}

// Human prints the passed check.
func (r RunRoundRecordedResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Run %s: review round %d was journaled in this visit\n", r.Run, r.Round)
}

func runRunRoundRecorded(ctx *Ctx, _ any) (any, error) {
	const cmd = "run round-recorded"
	rr, err := openReviewRun(ctx, cmd, "rhiza describe run.round-recorded")
	if err != nil {
		return nil, err
	}
	if msg := reviewRoundMissing(rr.state); msg != "" {
		return nil, Conflictf("%s: %s", cmd, msg).WithNext("rhiza run findings")
	}
	return RunRoundRecordedResponse{Run: rr.id, Round: len(rr.state.Rounds)}, nil
}

// reviewRoundMissing says why the plan review visit st is in has journaled
// no round yet, or "" when it has one.
func reviewRoundMissing(st run.State) string {
	if st.VisitRound() {
		return ""
	}
	return fmt.Sprintf("run %s's review visit has journaled no findings round: run `rhiza run findings` (with {\"findings\":[]} for a clean review) before saying done", st.ID)
}

// isPlanReviewer reports whether rec is a plan run's review-step session,
// which must journal a round before it says done.
func isPlanReviewer(rec session.Record) bool {
	return rec.Run != "" && rec.Step == "review" && rec.Check == "" && rec.Role == "" && rec.Leaf == ""
}
