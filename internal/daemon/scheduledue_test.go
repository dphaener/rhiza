package daemon

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/dispatch"
	"github.com/dphaener/rhiza/internal/schedule"
)

func at(h, m int) time.Time { return time.Date(2026, 11, 8, h, m, 0, 0, time.UTC) }

func hourly(t *testing.T) schedule.Schedule {
	t.Helper()
	e, err := schedule.ParseEvery("1h")
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func sched(catch, overlap string) dispatch.Schedule {
	return dispatch.Schedule{Name: "s", CatchUp: catch, Overlap: overlap}
}

func TestScheduleDueOnTime(t *testing.T) {
	st := ScheduleState{Since: at(8, 30), LastSlot: at(9, 0)}
	d := ScheduleDue(sched("once", "skip"), hourly(t), st, at(10, 0))
	if !d.Fire || !d.Slot.Equal(at(10, 0)) || !d.Next.Equal(at(11, 0)) {
		t.Fatalf("got %+v", d)
	}
	if d := ScheduleDue(sched("once", "skip"), hourly(t), st, at(9, 59)); d.Fire {
		t.Fatalf("fired early: %+v", d)
	}
}

func TestScheduleDueMissedSlots(t *testing.T) {
	st := ScheduleState{Since: at(0, 0), LastSlot: at(1, 0), LastStart: at(1, 0), LastEnd: at(1, 5)}
	now := at(4, 30) // slots 2, 3, 4 missed
	d := ScheduleDue(sched("once", "skip"), hourly(t), st, now)
	if !d.Fire || !d.Slot.Equal(at(4, 0)) {
		t.Fatalf("catch_up once: %+v", d)
	}
	// Fired once: the next decision, with the journal updated, is quiet.
	st.LastSlot, st.LastStart, st.LastEnd = d.Slot, now, now.Add(time.Minute)
	if d := ScheduleDue(sched("once", "skip"), hourly(t), st, at(4, 45)); d.Fire {
		t.Fatalf("fired twice: %+v", d)
	}
	d = ScheduleDue(sched("skip", "skip"), hourly(t), ScheduleState{Since: at(0, 0), LastSlot: at(1, 0), LastEnd: at(1, 5)}, now)
	if d.Fire || !d.Next.Equal(at(5, 0)) {
		t.Fatalf("catch_up skip: %+v", d)
	}
	// Within the grace a skip schedule still fires.
	if d := ScheduleDue(sched("skip", "skip"), hourly(t), ScheduleState{Since: at(0, 0), LastSlot: at(3, 0)}, at(4, 1)); !d.Fire {
		t.Fatalf("on time under skip: %+v", d)
	}
}

func TestScheduleDueNewSchedule(t *testing.T) {
	st := ScheduleState{Since: at(9, 30)}
	if d := ScheduleDue(sched("once", "skip"), hourly(t), st, at(9, 45)); d.Fire {
		t.Fatalf("fired before its first slot: %+v", d)
	}
	if d := ScheduleDue(sched("once", "skip"), hourly(t), st, at(10, 0)); !d.Fire || !d.Slot.Equal(at(10, 0)) {
		t.Fatalf("first slot: %+v", d)
	}
}

func TestScheduleDueOverlap(t *testing.T) {
	// A run for the 9:00 slot is going; 10:00 and 11:00 pass.
	running := ScheduleState{Since: at(0, 0), LastSlot: at(9, 0), LastStart: at(9, 0), Running: true}
	for _, ov := range []string{"skip", "queue"} {
		if d := ScheduleDue(sched("once", ov), hourly(t), running, at(10, 0)); d.Fire {
			t.Fatalf("%s fired during a run: %+v", ov, d)
		}
	}
	ended := running
	ended.Running, ended.LastEnd = false, at(11, 20)
	now := at(11, 30)
	if d := ScheduleDue(sched("once", "skip"), hourly(t), ended, now); d.Fire {
		t.Fatalf("skip kept a dropped slot: %+v", d)
	}
	d := ScheduleDue(sched("once", "queue"), hourly(t), ended, now)
	if !d.Fire || !d.Slot.Equal(at(11, 0)) {
		t.Fatalf("queue: two slots should fire once: %+v", d)
	}
	// The queued slot still fires under catch_up skip, however late.
	if d := ScheduleDue(sched("skip", "queue"), hourly(t), ended, now); !d.Fire {
		t.Fatalf("queue under catch_up skip: %+v", d)
	}
	// After the queued run is journaled, nothing is left.
	after := ScheduleState{Since: at(0, 0), LastSlot: at(11, 0), LastStart: now, LastEnd: at(11, 40)}
	if d := ScheduleDue(sched("once", "queue"), hourly(t), after, at(11, 45)); d.Fire {
		t.Fatalf("queue fired again: %+v", d)
	}
}
