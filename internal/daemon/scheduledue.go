package daemon

import (
	"time"

	"github.com/dphaener/rhiza/internal/dispatch"
	"github.com/dphaener/rhiza/internal/schedule"
)

// maxSlotsScanned bounds the walk over missed slots in one decision.
const maxSlotsScanned = 1_000_000

// DueGrace is how late a slot may be and still count as on time under
// catch_up skip: a poll or two of the daemon.
const DueGrace = 2 * time.Minute

// ScheduleState is what the daemon knows about one schedule, rebuilt on every
// poll from the schedule's runs' journals, so it survives restarts.
type ScheduleState struct {
	// Since is when the schedule began to exist for the daemon (the time its
	// entry was first seen). A schedule with no run yet fires no slot at or
	// before Since: it waits for its next slot.
	Since time.Time
	// LastSlot is the slot the last scheduled run fired for; zero if none.
	LastSlot time.Time
	// LastStart is when that run started; zero if none.
	LastStart time.Time
	// LastEnd is when it ended; zero while it is still going or if none.
	LastEnd time.Time
	// Running is whether the schedule's previous run is still going.
	Running bool
}

// DueDecision is the outcome of ScheduleDue.
type DueDecision struct {
	// Fire says a run starts now, once, for Slot.
	Fire bool
	// Slot is the slot the run is for: the latest due slot, however many were
	// missed. Zero when Fire is false.
	Slot time.Time
	// Next is the first slot after now; zero if the schedule has none.
	Next time.Time
}

// ScheduleDue is the pure decision of whether a schedule fires at now.
//
//   - A slot after the baseline (the last fired slot, else Since) and not
//     after now is due; any number of them fire one run, for the latest.
//   - catch_up skip fires only a slot no older than DueGrace, so downtime
//     waits for the next slot; catch_up once fires the late one too.
//   - While the previous run is going nothing fires. Under overlap skip the
//     slots that fell in that run are dropped; under overlap queue they stay
//     pending, so one run starts when it ends, however late.
func ScheduleDue(sc dispatch.Schedule, sch schedule.Schedule, st ScheduleState, now time.Time) DueDecision {
	var d DueDecision
	d.Next, _ = sch.Next(now)
	if st.Running {
		return d
	}
	base := st.LastSlot
	if base.IsZero() {
		base = st.Since
	}
	queue := sc.Overlap == dispatch.OverlapQueue
	if !queue && st.LastEnd.After(base) {
		base = st.LastEnd // slots during the run are dropped
	}
	var first, latest time.Time
	t := base
	for i := 0; i < maxSlotsScanned; i++ {
		n, ok := sch.Next(t)
		if !ok || n.After(now) {
			break
		}
		if first.IsZero() {
			first = n
		}
		latest, t = n, n
	}
	if latest.IsZero() {
		return d
	}
	if sc.CatchUp == dispatch.CatchUpSkip && now.Sub(latest) > DueGrace {
		held := queue && !st.LastStart.IsZero() && first.After(st.LastStart) && !first.After(st.LastEnd)
		if !held {
			return d
		}
	}
	d.Fire, d.Slot = true, latest
	return d
}
