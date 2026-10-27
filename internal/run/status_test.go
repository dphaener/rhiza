package run

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// deriveStatus over a leaves step: progress is passed leaves × 100 / all,
// a live leaf is an open one with a session, its role the latest session's
// on its current attempt, blocked on the owner, or on the arbiter in an
// autonomous run's question; a finished run is done or error with no
// progress or leaves.
func TestDeriveStatus(t *testing.T) {
	def, err := workflow.Parse("w", []byte(leavesWF))
	if err != nil {
		t.Fatal(err)
	}
	leaves := []Leaf{{ID: "L1", Handle: "a"}, {ID: "L2", Handle: "b"}, {ID: "L3", Handle: "c"}}
	st := State{
		ID: "R1", Change: "CS1", Title: "T", Current: "build",
		Leaves: []LeafState{
			{Step: "build", Leaf: "L1", Handle: "a", Outcome: LeafPassed,
				Sessions: []LeafSession{{Role: "builder", Attempt: 1, Session: "S1", Seq: 1}}},
			{Step: "build", Leaf: "L2", Handle: "b", Sessions: []LeafSession{
				{Role: "builder", Attempt: 1, Session: "S2", Seq: 2},
				{Role: "verifier", Attempt: 1, Session: "S3", Seq: 3},
			}},
		},
		Blocks: []SessionBlock{{Step: "build", Session: "S3", Block: Block{Kind: "question"}}},
	}
	rep := deriveStatus(st, def, leaves)
	want := StatusReport{
		Run: "R1", Change: "CS1", Title: "T", Step: "build", State: StatusWorking,
		Leaves: 3, Passed: 1, Progress: ptr(33),
		Live: []LeafStatus{{N: 2, Leaf: "L2", Handle: "b", Role: "verifier", Attempt: 1, Session: "S3", Blocked: true, Kind: "question"}},
	}
	if !reflect.DeepEqual(rep, want) {
		t.Errorf("report =\n%+v\nwant\n%+v", rep, want)
	}
	st.Autonomous = true
	if l := deriveStatus(st, def, leaves).Live[0]; l.Blocked || !l.Arbiter {
		t.Errorf("autonomous question: %+v, want waiting on the arbiter", l)
	}
	if rep := deriveStatus(st, def, nil); rep.Progress != nil || rep.Live != nil {
		t.Errorf("leaves unknown: %+v", rep)
	}
	st.Ended, st.Result, st.Detail, st.Current, st.FinishedStep = true, Failed, "boom", "", "build"
	if rep := deriveStatus(st, def, leaves); rep.State != StatusError || rep.Detail != "boom" || rep.Step != "build" ||
		rep.Progress != nil || rep.Live != nil {
		t.Errorf("failed: %+v", rep)
	}
	st.Result = Succeeded
	if rep := deriveStatus(st, def, leaves); rep.State != StatusDone {
		t.Errorf("succeeded: %+v", rep)
	}
}

func ptr(n int) *int { return &n }

// statusLine is a report in one line: state, step, progress and each live
// (+) and ended (-) leaf.
func statusLine(r StatusReport) string {
	s := r.State + " " + r.Step
	if r.Progress != nil {
		s += fmt.Sprintf(" %d%%", *r.Progress)
	}
	for _, l := range r.Live {
		s += fmt.Sprintf(" +%d:%s", l.N, l.Role)
		if l.Blocked {
			s += ":blocked"
		}
	}
	for _, l := range r.Ended {
		s += fmt.Sprintf(" -%d", l.N)
	}
	return s
}

// An engine with Status reports each change once: the step entered, the
// leaves known, each leaf's builder then verifier, a block and its end,
// each leaf ended, and the run done.
func TestEngineStatus(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves
	passAll(r, twoLeaves)
	b := leafSID("leaf-a", "builder", 1)
	r.agentSeq[b] = []AgentStatus{{Blocked: &Block{Kind: "permission"}}, {}}
	var got []string
	e.Status = func(rep StatusReport) { got = append(got, statusLine(rep)) }
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"working build",
		"working build 0%",
		"working build 0% +1:builder",
		"working build 0% +1:builder:blocked",
		"working build 0% +1:builder",
		"working build 0% +1:verifier",
		"working build 50% -1",
		"working build 50% +2:builder",
		"working build 50% +2:verifier",
		"working build 100% -2",
		"working build",
		"done build",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("reports =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// A failed run's last report is error, with the failure's detail.
func TestEngineStatusError(t *testing.T) {
	e, r, _, _ := setup(t, leavesWF)
	r.leaves = twoLeaves[:1]
	var last StatusReport
	e.Status = func(rep StatusReport) { last = rep }
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	if last.State != StatusError || last.Detail == "" || last.Live != nil {
		t.Errorf("last report %+v", last)
	}
}
