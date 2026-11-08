package quiettest

import (
	"testing"
	"time"
)

// TestStopwatchCountsSleeps checks a lap's own time keeps time the process
// spent asleep: only waiting for a CPU comes off the clock, never disk,
// fsync or sleeps, so a read that started sleeping would still fail its
// budget.
func TestStopwatchCountsSleeps(t *testing.T) {
	sw := Start()
	time.Sleep(30 * time.Millisecond)
	l := sw.Stop()
	if l.Wall < 30*time.Millisecond {
		t.Fatalf("wall %v, want at least the 30ms slept", l.Wall)
	}
	if l.Own > l.Wall {
		t.Errorf("own %v over wall %v", l.Own, l.Wall)
	}
	if l.Own < 25*time.Millisecond {
		t.Errorf("own %v left out the sleep (wall %v)", l.Own, l.Wall)
	}
}

func TestSummarize(t *testing.T) {
	ms := func(n int) time.Duration { return time.Duration(n) * time.Millisecond }
	got := Summarize([]Lap{{ms(9), ms(2), ms(9)}, {ms(1), ms(1), 0}, {ms(5), ms(4), 0}, {ms(3), ms(3), 0}})
	want := Timing{Min: ms(1), Median: ms(4), Max: ms(9), Own: ms(2) + ms(1)/2, Pressure: 0.5}
	if got != want {
		t.Errorf("even laps: %+v, want %+v", got, want)
	}
	got = Summarize([]Lap{{ms(7), ms(7), 0}, {ms(2), ms(1), 0}, {ms(4), ms(4), 0}})
	if want := (Timing{Min: ms(2), Median: ms(4), Max: ms(7), Own: ms(4)}); got != want {
		t.Errorf("odd laps: %+v, want %+v", got, want)
	}
	if got := Summarize(nil); got != (Timing{}) {
		t.Errorf("no laps: %+v", got)
	}
}

// TestOver checks a budget fails on own time over it, passes under it, and
// is inconclusive, not failed, when the machine was saturated.
func TestOver(t *testing.T) {
	b := 50 * time.Millisecond
	if over, why := (Timing{Median: 3 * b, Own: b - 1}).Over(b); over || why != "" {
		t.Errorf("own under, clock over: over %v, %q", over, why)
	}
	if over, _ := (Timing{Own: b, Pressure: 0.05}).Over(b); !over {
		t.Error("own over on a quiet machine passed")
	}
	saturated := Timing{Own: 2 * b, Pressure: SaturatedPressure}
	defer func(was bool) { onCI = was }(onCI)
	onCI = false
	if over, why := saturated.Over(b); over || why == "" {
		t.Errorf("own over on a saturated machine: over %v, %q, want inconclusive", over, why)
	}
	onCI = true
	if over, _ := saturated.Over(b); !over {
		t.Error("own over on a saturated CI runner passed: CI judges every budget")
	}
}
