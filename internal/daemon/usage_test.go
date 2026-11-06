package daemon

import (
	"testing"
	"time"
)

// A limit hit recorded by one project's driver, in the root's state folder,
// holds starts in a second project, and lifts for both at the reset time.
func TestGovernorHoldSpansProjectsAndLifts(t *testing.T) {
	dir := RootStateDir(t.TempDir())
	hit := time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)
	reset := hit.Add(5 * time.Hour)
	if err := RecordLimitHitIn(dir, LimitHit{HitAt: hit, ResetAt: reset, Agent: "claude"}); err != nil {
		t.Fatal(err)
	}
	// Project B's daemon poll reads the same folder.
	if h, ok := ActiveHold(dir, hit.Add(time.Hour)); !ok || !h.ResetAt.Equal(reset) {
		t.Fatalf("ActiveHold during = %v, %v; want the hit holding until %v", h, ok, reset)
	}
	if _, ok := ActiveHold(dir, reset.Add(-time.Nanosecond)); !ok {
		t.Error("hold lifted before the reset time")
	}
	if _, ok := ActiveHold(dir, reset); ok {
		t.Error("hold still on at the reset time")
	}
	if _, ok := ActiveHold(t.TempDir(), hit.Add(time.Hour)); ok {
		t.Error("an empty folder holds")
	}
}
