package session

import (
	"testing"
	"time"
)

func TestStartNotSubmitted(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	old := now.Add(-StartGrace - time.Second)
	idle := func(at time.Time) *Report { return &Report{State: StateIdle, At: at} }
	cases := []struct {
		name string
		rec  Record
		want bool
	}{
		{"idle past the grace", Record{Run: "r", Status: Status{Hook: idle(old), Started: idle(old)}}, true},
		{"within the grace", Record{Run: "r", Status: Status{Hook: idle(now), Started: idle(now)}}, false},
		{"prompted", Record{Run: "r", Status: Status{Hook: idle(old), Started: idle(old), Prompted: idle(old)}}, false},
		{"no start reported", Record{Run: "r", Status: Status{Hook: idle(old)}}, false},
		{"not a run's", Record{Status: Status{Hook: idle(old), Started: idle(old)}}, false},
		{"working", Record{Run: "r", Status: Status{Hook: &Report{State: StateWorking, At: old}, Started: idle(old)}}, false},
		{"exited", Record{Run: "r", Status: Status{Hook: idle(old), Started: idle(old), Exit: &Report{State: StateExited, At: old}}}, false},
	}
	for _, c := range cases {
		if got := NeverSubmitted(c.rec, now); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}
