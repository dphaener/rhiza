package api

// This file holds run.open: focus the Rex tab of a change set's live run.

import (
	"fmt"
	"io"
	"reflect"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.open",
		Summary:  "Focus the Rex tab of a change set's live run; says so when it has no live run or its tab is gone",
		Request:  reflect.TypeFor[RunOpenRequest](),
		Response: reflect.TypeFor[RunOpenResponse](),
		Run:      runRunOpen,
	})
}

// RunOpenRequest names the change set whose live run to open.
type RunOpenRequest struct {
	Change string `json:"change" arg:"0" desc:"The change set: its ID or a unique prefix of it"`
}

// RunOpenResponse says which tab was focused, or why none was.
type RunOpenResponse struct {
	Change  string `json:"change" desc:"The change set's ID"`
	Live    bool   `json:"live" desc:"True when the change set has a live run"`
	Run     string `json:"run,omitempty" desc:"The live run, when there is one"`
	Opened  bool   `json:"opened" desc:"True when the tab was focused"`
	Tab     string `json:"tab,omitempty" desc:"The Rex block of the tab, when one is recorded"`
	TabGone bool   `json:"tab_gone,omitempty" desc:"True when the live run's tab is no longer in Rex, or none was recorded"`
}

// Human prints the response as a line.
func (r RunOpenResponse) Human(w io.Writer) {
	switch {
	case r.Opened:
		fmt.Fprintf(w, "Focused %s, the tab of run %s\n", r.Tab, r.Run)
	case !r.Live:
		fmt.Fprintf(w, "%s has no live run\n", r.Change)
	default:
		fmt.Fprintf(w, "run %s: tab gone\n", r.Run)
	}
}

func runRunOpen(ctx *Ctx, req any) (any, error) {
	r := req.(RunOpenRequest)
	if r.Change == "" {
		return nil, Usagef("run open needs a change set")
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	cs, err := w.resolveChangeSet(r.Change)
	if err != nil {
		return nil, err
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
	resp := RunOpenResponse{Change: cs.String()}
	var st run.State
	for _, j := range runs {
		s := j.State()
		if s.Change != resp.Change || !rs.Live(j.ID) {
			continue
		}
		if !resp.Live || s.StartedAt.After(st.StartedAt) {
			st, resp.Live, resp.Run = s, true, j.ID
		}
	}
	if !resp.Live {
		return resp, nil
	}
	var session, block string
	if st.Process != nil {
		session, block = st.Process.RexSession, st.Process.RexBlock
	}
	resp.Tab = block
	if block == "" || ctx.Exec == nil {
		resp.TabGone = true
		return resp, nil
	}
	scope := []string{}
	if session != "" {
		scope = []string{"-s", session}
	}
	inspect := append(append([]string{"block", "inspect"}, scope...), "--json", block)
	if _, err := ctx.Exec("rex", inspect, "", nil); err != nil {
		resp.TabGone = true
		return resp, nil
	}
	if _, err := ctx.Exec("rex", append(append([]string{"focus"}, scope...), block), "", nil); err != nil {
		return nil, Errorf("rex could not focus %s: %v", block, err)
	}
	resp.Opened = true
	return resp, nil
}
