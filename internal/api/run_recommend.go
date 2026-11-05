package api

// This file holds run.recommend, which an advisor session runs to report
// its recommendation for a question a normal run's session asked the
// owner. It journals the recommendation in the advisor's run and answers
// nothing: the owner still answers in the asking session's tab, where
// rhiza inbox shows the recommendation beside the question.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.recommend",
		Summary:  "Journal an advisor session's recommendation and reasons for the question it was launched on; it answers nothing",
		Write:    true,
		Request:  reflect.TypeFor[RunRecommendRequest](),
		Response: reflect.TypeFor[RunRecommendResponse](),
		Run:      runRunRecommend,
	})
}

const runRecommendNext = "rhiza guide advisor"

// RunRecommendRequest is the recommendation and why.
type RunRecommendRequest struct {
	Recommendation string `json:"recommendation" desc:"What to answer, in one or two sentences"`
	Reasons        string `json:"reasons" desc:"Why: what was read and weighed"`
	Restatement    string `json:"restatement,omitempty" desc:"What is asked, in plain words with no handles, IDs or finding IDs; for a question a run raised from a template"`
	Session        string `json:"session,omitempty" desc:"The advisor session; default RHIZA_SESSION, which a launched advisor has"`
}

// RunRecommendResponse is the recommendation journalled, or in a dry run
// the one that would be.
type RunRecommendResponse struct {
	Run     string `json:"run" desc:"The run the question was asked in"`
	Asker   string `json:"asker" desc:"The session that asked the question"`
	Session string `json:"session" desc:"The advisor session"`
	DryRun  bool   `json:"dry_run,omitempty" desc:"True if nothing was journalled"`
}

// Human writes the response as text.
func (r RunRecommendResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "Recommended an answer to the question session %s asked in run %s\n", r.Asker, r.Run)
}

// runRunRecommend journals the recommendation against the session the
// advisor was launched on. Only an advisor session (one launched with an
// asking session) may; so may a run that is still going.
func runRunRecommend(ctx *Ctx, req any) (any, error) {
	r := req.(RunRecommendRequest)
	ref := strings.TrimSpace(r.Session)
	if ref == "" && ctx.Getenv != nil {
		ref = strings.TrimSpace(ctx.Getenv(envSession))
	}
	if ref == "" {
		return nil, Usagef("run recommend needs the advisor session: name it with --session, or run it where %s is set, as in a session Rhiza launched", envSession).
			WithNext(runRecommendNext)
	}
	if strings.TrimSpace(r.Recommendation) == "" || strings.TrimSpace(r.Reasons) == "" {
		return nil, Usagef("run recommend needs --recommendation and --reasons, neither blank").WithNext(runRecommendNext)
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, ref)
	if err != nil {
		return nil, err
	}
	if rec.Advises == "" || rec.Run == "" {
		return nil, Conflictf("session %s is not an advisor: only the session launched to advise on a question recommends", rec.ID).
			WithNext(runRecommendNext)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	j, err := readRun(rs, rec.Run)
	if err != nil {
		return nil, err
	}
	if j.State().Ended {
		return nil, Conflictf("run %s is finished: its question is no longer open", j.ID).WithNext("rhiza run status " + j.ID)
	}
	resp := RunRecommendResponse{Run: j.ID, Asker: rec.Advises, Session: rec.ID, DryRun: ctx.DryRun}
	if ctx.DryRun {
		return resp, nil
	}
	e := run.Entry{
		Type: run.TypeRecommended, At: ctx.Now(), Step: rec.Step, Session: rec.ID, Asker: rec.Advises,
		Text: strings.TrimSpace(r.Recommendation), Detail: strings.TrimSpace(r.Reasons),
		Restated: strings.TrimSpace(r.Restatement),
	}
	if err := rs.Append(j.ID, e); err != nil {
		return nil, Errorf("journal the recommendation in run %s: %v", j.ID, err).WithNext("rhiza run status " + j.ID)
	}
	// The recommendation is journalled: the advisor's pane has done its
	// work, so close it (this may end the process that asked). A pane that
	// will not close is left open: the recommendation stands.
	closeSessionPane(ctx, rec, ctx.Root)
	return resp, nil
}
