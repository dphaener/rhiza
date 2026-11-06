package api

// This file holds run.ls: the repository's runs newest first, each with
// the little a live view needs. It reads the run journals only and never
// writes.

import (
	"fmt"
	"io"
	"reflect"
	"slices"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "run.ls",
		ReadOnly: true,
		Summary:  "List the repository's runs newest first, each with its ID, change set, title, state, current step, start time and whether it is live",
		Request:  reflect.TypeFor[RunLsRequest](),
		Response: reflect.TypeFor[RunLsResponse](),
		Run:      runRunLs,
	})
}

const (
	runLsDefaultLimit = 20
	runLsMaxLimit     = 1000
)

// RunLsRequest limits the listing.
type RunLsRequest struct {
	Limit int `json:"limit,omitempty" desc:"The most runs to return, 1 to 1000; default 20"`
}

// RunLsItem is one run in the listing.
type RunLsItem struct {
	ID         string `json:"id" desc:"The run ID"`
	Change     string `json:"change" desc:"The change set it runs on"`
	Title      string `json:"title" desc:"The change set's title; empty when it has none"`
	State      string `json:"state" desc:"running (another process runs it now), blocked (its agent step waits on the owner), stopped (resumable), succeeded or failed"`
	Current    string `json:"current,omitempty" desc:"The step being run or to run next; empty once the run is finished"`
	StartedAt  string `json:"started_at" desc:"When it started, RFC 3339 UTC"`
	Live       bool   `json:"live" desc:"True while another process holds the run"`
	Workflow   string `json:"workflow" desc:"The workflow's name"`
	Autonomous bool   `json:"autonomous,omitempty" desc:"True for a run started with --autonomous; absent otherwise"`
}

// RunLsResponse is the listing.
type RunLsResponse struct {
	Runs []RunLsItem `json:"runs" desc:"The runs, newest started first; empty when there are none"`
}

// Human writes the listing as text.
func (r RunLsResponse) Human(w io.Writer) {
	if len(r.Runs) == 0 {
		fmt.Fprintln(w, "No runs")
		return
	}
	for _, it := range r.Runs {
		fmt.Fprintf(w, "%s %s %s", it.ID, it.State, it.Workflow)
		if it.Current != "" {
			fmt.Fprintf(w, " at %s", it.Current)
		}
		fmt.Fprintf(w, " since %s, change set %s", it.StartedAt, it.Change)
		if it.Title != "" {
			fmt.Fprintf(w, ": %s", oneLine(it.Title))
		}
		fmt.Fprintln(w)
	}
}

// runRunLs lists the runs, newest first, up to the limit.
func runRunLs(ctx *Ctx, req any) (any, error) {
	r := req.(RunLsRequest)
	limit := r.Limit
	if limit == 0 {
		limit = runLsDefaultLimit
	}
	if limit < 0 || limit > runLsMaxLimit {
		return nil, Usagef("limit %d is out of range: give 1 to %d", r.Limit, runLsMaxLimit)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	runs, err := rs.List()
	if err != nil {
		return nil, Errorf("read runs: %v", err)
	}
	slices.Reverse(runs)
	if len(runs) > limit {
		runs = runs[:limit]
	}
	sessions := session.NewStore(session.Dir(common))
	out := []RunLsItem{}
	for _, j := range runs {
		st := j.State()
		live := !st.Ended && rs.Live(j.ID)
		info := runInfo(ctx, st, live, sessions)
		out = append(out, RunLsItem{
			ID: st.ID, Change: st.Change, Title: st.Title, State: info.State, Current: st.Current,
			StartedAt: rfc3339(st.StartedAt), Live: live, Workflow: st.Workflow, Autonomous: st.Autonomous,
		})
	}
	return RunLsResponse{Runs: out}, nil
}
