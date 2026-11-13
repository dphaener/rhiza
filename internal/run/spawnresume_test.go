package run

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/workflow"
)

// resumeRunner is the fake runner with the capabilities a restarted driver
// uses: it watches spawned children and takes hand-backs.
type resumeRunner struct {
	*fakeRunner
	children   map[string][]ChildState // per child, one state per Child call; past the end, the last
	childCalls map[string]int
	handbacks  []string // "<step> <session> <cause step>" per successful Handback
	attempts   int      // Handback calls, successful or not
	crashOn    int      // the Handback call that panics, as a killed driver does; 0 for none
}

func (r *resumeRunner) Child(session string) (ChildState, error) {
	n := r.childCalls[session]
	r.childCalls[session]++
	seq := r.children[session]
	if n >= len(seq) {
		n = len(seq) - 1
	}
	return seq[n], nil
}

func (r *resumeRunner) Handback(step workflow.Step, session string, cause HandbackCause) error {
	r.attempts++
	if r.attempts == r.crashOn {
		panic(errKilled)
	}
	r.handbacks = append(r.handbacks, step.Name+" "+session+" "+cause.Step)
	return nil
}

type killed struct{}

var errKilled = killed{}

// killDriver runs the engine and reports whether a panic (the driver killed
// mid-run) cut it short.
func killDriver(e *Engine) (died bool) {
	defer func() {
		if recover() != nil {
			died = true
		}
	}()
	_, _ = e.Run(runID)
	return false
}

func countEntries(t *testing.T, e *Engine, typ, session string) int {
	t.Helper()
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, en := range j.Entries {
		if en.Type == typ && (session == "" || en.Session == session) {
			n++
		}
	}
	return n
}

const spawnWF = `
[step.build]
kind = "agent"
prompt = "build"
on.pass = "done"
`

func TestSpawnResume(t *testing.T) {
	t.Run("children across a driver restart", func(t *testing.T) {
		e, fr, _, _ := setup(t, spawnWF)
		rr := &resumeRunner{fakeRunner: fr, children: map[string][]ChildState{}, childCalls: map[string]int{}}
		e.Runner = rr
		fr.outcome["S-build"] = workflow.OutcomePass
		// The driver is killed on the build session's first poll, after the
		// session spawned three children.
		fr.onPoll = func(session string, n int) {
			if session != "S-build" || n != 0 {
				return
			}
			for _, c := range []string{"K-live", "K-ended", "K-gone"} {
				if err := e.Store.Append(runID, Entry{Type: TypeSpawn, At: at(2), Step: "build", Parent: "S-build", Session: c, Role: "builder"}); err != nil {
					t.Fatal(err)
				}
			}
			panic(errKilled)
		}
		if !killDriver(e) {
			t.Fatal("driver was not killed")
		}
		// Restarted: K-live still runs for two checks; K-ended ended while the
		// driver was down and its spawn process records it too; K-gone's
		// record is deleted.
		pass := ChildState{Ended: true, Outcome: workflow.OutcomePass, Detail: "exit 0"}
		rr.children["K-live"] = []ChildState{{}, {}, pass}
		rr.children["K-ended"] = []ChildState{pass}
		rr.children["K-gone"] = []ChildState{{Gone: true}}
		fr.agentOut["S-build"] = []bool{true, false, false, false, true} // polls 0-4: S-build stays up meanwhile
		fr.onPoll = func(session string, n int) {
			if session == "S-build" && n == 1 {
				// The spawn process of K-ended wakes and records its own end.
				if _, err := e.Store.AppendSpawnEnded(runID, Entry{At: at(3), Step: "build", Session: "K-ended", Outcome: workflow.OutcomePass, Detail: "exit 0"}); err != nil {
					t.Fatal(err)
				}
			}
		}
		e.Runner = rr
		fr.polls = map[string]int{}
		if _, err := e.Run(runID); err != nil {
			t.Fatal(err)
		}
		j, _ := e.Store.Read(runID)
		st := j.State()
		if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), []string{"build:pass"}) {
			t.Fatalf("result %q steps %v", st.Result, outcomes(st))
		}
		byID := map[string]Spawn{}
		for _, sp := range st.Spawns {
			byID[sp.Session] = sp
		}
		if sp := byID["K-live"]; !sp.Ended || sp.Outcome != workflow.OutcomePass {
			t.Errorf("live child after reattach: %+v", sp)
		}
		if sp := byID["K-ended"]; !sp.Ended || sp.Outcome != workflow.OutcomePass {
			t.Errorf("child that ended while down: %+v", sp)
		}
		if sp := byID["K-gone"]; !sp.Ended || sp.Outcome != workflow.OutcomeFail || !strings.Contains(sp.Detail, "session record is gone") {
			t.Errorf("child with no record: %+v", sp)
		}
		for _, c := range []string{"K-live", "K-ended", "K-gone"} {
			if n := countEntries(t, e, TypeSpawnEnded, c); n != 1 {
				t.Errorf("%s: %d spawn-ended entries, want 1", c, n)
			}
		}
		if rr.childCalls["K-live"] < 3 {
			t.Errorf("live child checked %d times, want it watched until it ended", rr.childCalls["K-live"])
		}
	})

	t.Run("a pending hand-back is delivered once after a restart", func(t *testing.T) {
		e, fr, _, _ := setup(t, handbackWF)
		rr := &resumeRunner{fakeRunner: fr, children: map[string][]ChildState{}, childCalls: map[string]int{}, crashOn: 1}
		e.Runner = rr
		fr.outcome["S-build"] = workflow.OutcomePass
		fr.codes["check"] = []int{1, 0}
		fr.details["check"] = "tests red"
		// Killed while delivering the hand-back to build.
		if !killDriver(e) {
			t.Fatalf("driver was not killed: %v", entryTypes(t, e))
		}
		if n := countEntries(t, e, TypeHandback, ""); n != 1 {
			t.Fatalf("%d hand-back entries before the restart", n)
		}
		j, _ := e.Store.Read(runID)
		if pb := j.State().PendingHandback; pb == nil || pb.Target != "build" {
			t.Fatalf("pending hand-back %+v", pb)
		}
		if _, err := e.Run(runID); err != nil {
			t.Fatal(err)
		}
		j, _ = e.Store.Read(runID)
		st := j.State()
		want := []string{"build:pass", "check:fail", "build:pass", "check:pass"}
		if st.Result != Succeeded || !reflect.DeepEqual(outcomes(st), want) {
			t.Fatalf("result %q steps %v", st.Result, outcomes(st))
		}
		if !reflect.DeepEqual(rr.handbacks, []string{"build S-build check"}) {
			t.Errorf("deliveries %v, want one", rr.handbacks)
		}
		if st.PendingHandback != nil || countEntries(t, e, TypeHandbackDelivered, "") != 1 {
			t.Errorf("pending %+v, delivered entries %d", st.PendingHandback, countEntries(t, e, TypeHandbackDelivered, ""))
		}
		if n := countEntries(t, e, TypeHandback, ""); n != 1 {
			t.Errorf("%d hand-back entries after the restart, want 1", n)
		}
	})
}
