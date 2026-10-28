package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

// A check or critic session's RHIZA_CYCLE names its round's cycle when it
// is one of the step's cycles that runs the check; anything else is 0, so
// the cycle is inferred.
func TestEnvCheckCycle(t *testing.T) {
	st := run.State{Cycles: []run.CheckCycle{
		{Step: "review", Cycle: 1, Checks: []string{"alpha", "beta"}},
		{Step: "review", Cycle: 2, Checks: []string{"beta"}},
	}}
	for _, tt := range []struct {
		env, check string
		want       int
	}{
		{"1", "alpha", 1},
		{" 2 ", "beta", 2},
		{"2", "alpha", 0},
		{"3", "beta", 0},
		{"0", "alpha", 0},
		{"x", "alpha", 0},
		{"", "alpha", 0},
	} {
		ctx := &Ctx{Getenv: func(k string) string {
			if k == envCycle {
				return tt.env
			}
			return ""
		}}
		if got := envCheckCycle(ctx, st, "review", tt.check); got != tt.want {
			t.Errorf("RHIZA_CYCLE=%q for %s: %d, want %d", tt.env, tt.check, got, tt.want)
		}
	}
}
