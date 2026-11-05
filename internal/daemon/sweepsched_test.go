package daemon

import (
	"testing"
	"time"
)

func TestSweepSchedule(t *testing.T) {
	const gib = 1 << 30
	var s SweepSchedule
	t0 := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	if !s.Due(t0, 100*gib, 10*gib, true) {
		t.Fatal("no sweep due at start")
	}
	s.Done(t0)
	if s.Due(t0.Add(5*time.Minute), 100*gib, 10*gib, true) {
		t.Fatal("sweep due 5 minutes after one with plenty of space")
	}
	if s.Due(t0.Add(59*time.Minute), 100*gib, 10*gib, true) {
		t.Fatal("sweep due before the hour with plenty of space")
	}
	if !s.Due(t0.Add(time.Hour), 100*gib, 10*gib, true) {
		t.Fatal("no sweep due after an hour")
	}
	// Under 2x minFree: at once, but not within 10 minutes of the last.
	if s.Due(t0.Add(9*time.Minute), 19*gib, 10*gib, true) {
		t.Fatal("urgent sweep within 10 minutes of the last")
	}
	if !s.Due(t0.Add(10*time.Minute), 19*gib, 10*gib, true) {
		t.Fatal("no urgent sweep under 2x minFree after 10 minutes")
	}
	if s.Due(t0.Add(10*time.Minute), 20*gib, 10*gib, true) {
		t.Fatal("urgent sweep at exactly 2x minFree")
	}
	if s.Due(t0.Add(10*time.Minute), 1*gib, 0, true) {
		t.Fatal("urgent sweep with the guard off")
	}
	if s.Due(t0.Add(10*time.Minute), 0, 10*gib, false) {
		t.Fatal("urgent sweep with unknown free space")
	}
	s.Done(t0.Add(10 * time.Minute))
	if s.Due(t0.Add(19*time.Minute), 19*gib, 10*gib, true) {
		t.Fatal("urgent sweep within 10 minutes of the urgent one")
	}
}
