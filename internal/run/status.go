package run

import (
	"reflect"
	"slices"

	"github.com/dphaener/rhiza/internal/workflow"
)

// The root states of a StatusReport.
const (
	// StatusWorking: the run is going, at Step.
	StatusWorking = "working"
	// StatusIdle: the run stopped (at a wait under NoWait) and can be
	// resumed; Detail says why.
	StatusIdle = "idle"
	// StatusDone: the run finished and succeeded.
	StatusDone = "done"
	// StatusError: the run finished and failed; Detail says why.
	StatusError = "error"
)

// StatusReport is where a run stands, for the caller to show (Engine.Status):
// the run's state at its current step and, while that step is a leaves
// or checks step, its progress and live leaves or checks. It is derived from the journal
// alone (deriveStatus), so this package writes nothing for it.
type StatusReport struct {
	Run    string
	Change string
	Title  string
	// Step is the current step, or the step the run finished or stopped at.
	Step string
	// State is StatusWorking, StatusIdle, StatusDone or StatusError.
	State string
	// Detail is, when finished or stopped, why; "" otherwise.
	Detail string
	// Leaves and Passed are, while Step is a leaves step whose leaves are
	// known, how many leaves it has and how many passed (have Evidence).
	Leaves int
	Passed int
	// Checks is, while Step is a checks step, how many checks it has;
	// Passed is then how many of them passed their latest cycle.
	Checks int
	// Progress is Passed × 100 / Leaves (or / Checks) under the same
	// condition, nil otherwise (and once the run finished or stopped).
	Progress *int
	// Live are the leaves step's open leaves with a session, in leaf
	// order; on a checks step, its checks with a live session, in the
	// step's order, each as a LeafStatus whose Handle is the check's name
	// and Attempt its cycle.
	Live []LeafStatus
	// Ended are the leaves live in the last report and not in this one:
	// the caller clears what it showed for them.
	Ended []LeafStatus
}

// LeafStatus is one live leaf of a StatusReport: the leaf at position N
// (from 1, in leaf order) and its latest session's role on its current
// attempt. Blocked is set while that session waits on the owner, on Kind;
// Arbiter while it waits on the arbiter (an autonomous run's question).
type LeafStatus struct {
	N       int
	Leaf    string
	Handle  string
	Role    string
	Attempt int
	Session string
	Blocked bool
	Kind    string
	Arbiter bool
}

// deriveStatus is the pure report of a run from its folded journal st,
// its definition and, when its current step is a leaves step, that
// step's leaves in leaf order (nil when not known yet).
func deriveStatus(st State, def *workflow.Definition, leaves []Leaf) StatusReport {
	rep := StatusReport{Run: st.ID, Change: st.Change, Title: st.Title, Step: st.Current, State: StatusWorking}
	if rep.Step == "" {
		rep.Step = st.FinishedStep
	}
	if rep.Step == "" && len(st.Steps) > 0 {
		rep.Step = st.Steps[len(st.Steps)-1].Step
	}
	if st.Ended {
		rep.State, rep.Detail = StatusDone, st.Detail
		if st.Result != Succeeded {
			rep.State = StatusError
		}
		return rep
	}
	step, ok := def.Steps[st.Current]
	if ok && step.Kind == workflow.KindChecks {
		return checksStatus(rep, st, step, leaves)
	}
	if !ok || step.Kind != workflow.KindLeaves || leaves == nil {
		return rep
	}
	rep.Leaves = len(leaves)
	for i, lf := range leaves {
		l, ok := st.Leaf(step.Name, lf.ID)
		if !ok {
			continue
		}
		if l.Outcome == LeafPassed {
			rep.Passed++
		}
		if l.Outcome != "" {
			continue
		}
		s, ok := latestSession(l)
		if !ok {
			continue
		}
		ls := LeafStatus{N: i + 1, Leaf: lf.ID, Handle: lf.Handle, Role: s.Role, Attempt: s.Attempt, Session: s.Session}
		if b, ok := st.BlockOf(s.Session); ok {
			if st.Autonomous && b.Block.Kind == "question" {
				ls.Arbiter = true
			} else {
				ls.Blocked, ls.Kind = true, b.Block.Kind
			}
		}
		rep.Live = append(rep.Live, ls)
	}
	p := 100
	if rep.Leaves > 0 {
		p = rep.Passed * 100 / rep.Leaves
	}
	rep.Progress = &p
	return rep
}

// checksStatus is rep with checks step step's progress: its checks, how
// many passed their latest cycle, and each check whose latest session is
// not yet judged; then its live fixers and verifiers (fixStatus). leaves
// are the change set's leaves in order, nil when not read yet.
func checksStatus(rep StatusReport, st State, step workflow.Step, leaves []Leaf) StatusReport {
	rep.Checks = len(step.Checks)
	for i, c := range step.Checks {
		cs, ok := st.Check(step.Name, c.Name)
		if !ok {
			continue
		}
		if o, ok := cs.Latest(); ok && o.Outcome == CheckPassed {
			rep.Passed++
		}
		n := len(cs.Sessions)
		if n == 0 {
			continue
		}
		s := cs.Sessions[n-1]
		if _, judged := cs.Judgement(s.Role, s.Cycle); judged {
			continue
		}
		ls := LeafStatus{N: i + 1, Handle: c.Name, Role: s.Role, Attempt: s.Cycle, Session: s.Session}
		if b, ok := st.BlockOf(s.Session); ok {
			if st.Autonomous && b.Block.Kind == "question" {
				ls.Arbiter = true
			} else {
				ls.Blocked, ls.Kind = true, b.Block.Kind
			}
		}
		rep.Live = append(rep.Live, ls)
	}
	rep.Live = append(rep.Live, fixStatus(st, step, leaves)...)
	p := 100
	if rep.Checks > 0 {
		p = rep.Passed * 100 / rep.Checks
	}
	rep.Progress = &p
	return rep
}

// fixStatus are checks step step's live fix sessions: each leaf's fixer
// or verifier whose latest session is not yet judged, in leaf order (as
// leaves gives it, else the journal's), as a LeafStatus naming its leaf
// and its cycle as Attempt; then the change-set fixer of the latest cycle
// that ran one, until judged, with no Leaf, N 0 and the step's name as
// Handle.
func fixStatus(st State, step workflow.Step, leaves []Leaf) []LeafStatus {
	var out []LeafStatus
	add := func(l LeafState, n int) {
		s, ok := latestSession(l)
		if !ok {
			return
		}
		if _, judged := l.Judged(s.Role, s.Attempt); judged {
			return
		}
		ls := LeafStatus{N: n, Leaf: l.Leaf, Handle: l.Handle, Role: s.Role, Attempt: s.Attempt, Session: s.Session}
		sessionBlock(&ls, st)
		out = append(out, ls)
	}
	n := 0
	for _, lf := range leaves {
		n++
		if l, ok := st.Leaf(step.Name, lf.ID); ok {
			add(l, n)
		}
	}
	for _, l := range st.Leaves {
		if l.Step != step.Name || slices.ContainsFunc(leaves, func(lf Leaf) bool { return lf.ID == l.Leaf }) {
			continue
		}
		n++
		add(l, n)
	}
	for i := len(st.Fixes) - 1; i >= 0; i-- {
		f := st.Fixes[i]
		if f.Step != step.Name {
			continue
		}
		if !f.Judged {
			ls := LeafStatus{Handle: step.Name, Role: workflow.RoleFixer, Attempt: f.Cycle, Session: f.Session}
			sessionBlock(&ls, st)
			out = append(out, ls)
		}
		break
	}
	return out
}

// sessionBlock marks ls blocked on the owner, or waiting on the arbiter,
// when its session has an open block.
func sessionBlock(ls *LeafStatus, st State) {
	if b, ok := st.BlockOf(ls.Session); ok {
		if st.Autonomous && b.Block.Kind == "question" {
			ls.Arbiter = true
		} else {
			ls.Blocked, ls.Kind = true, b.Block.Kind
		}
	}
}

// latestSession is the last started session of l's current attempt.
func latestSession(l LeafState) (LeafSession, bool) {
	attempt := l.Attempt()
	var out LeafSession
	found := false
	for _, s := range l.Sessions {
		if s.Attempt == attempt && (!found || s.Seq > out.Seq) {
			out, found = s, true
		}
	}
	return out, found
}

// sameLive says whether a and b are one record: a leaf, a check or the
// change-set fixer (no leaf, N 0) at the same position.
func sameLive(a, b LeafStatus) bool {
	return a.N == b.N && a.Leaf == b.Leaf && checkRole(a.Role) == checkRole(b.Role)
}

func checkRole(role string) bool { return role == workflow.RoleCheck || role == workflow.RoleCritic }

// endedLeaves are the leaves of prev's Live not in next's.
func endedLeaves(prev, next []LeafStatus) []LeafStatus {
	var out []LeafStatus
	for _, p := range prev {
		live := false
		for _, n := range next {
			live = live || sameLive(n, p)
		}
		if !live {
			out = append(out, p)
		}
	}
	return out
}

// status reports the run to Engine.Status when the report derived from
// the journal differs from the last one reported. A non-empty stopped
// reports the run idle, stopped for that reason.
func (r *running) status(stopped string) {
	if r.e.Status == nil {
		return
	}
	st := r.j.State()
	rep := deriveStatus(st, r.def, r.leafList[st.Current])
	if stopped != "" {
		rep.State, rep.Detail, rep.Progress, rep.Live = StatusIdle, stopped, nil, nil
	}
	var prev []LeafStatus
	if r.lastStatus != nil {
		if reflect.DeepEqual(*r.lastStatus, rep) {
			return
		}
		prev = r.lastStatus.Live
	}
	last := rep
	r.lastStatus = &last
	rep.Ended = endedLeaves(prev, rep.Live)
	r.e.Status(rep)
}
