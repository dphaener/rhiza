package model

import (
	"errors"
	"reflect"
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"pgregory.net/rapid"
)

// genSeedLog is genLog with, often, an in-flight change set depending on
// another in flight through a depends-on link of its own.
func genSeedLog(t *rapid.T) []changeset.ChangeSet {
	sets := genLog(t)
	var flight []int
	for i, s := range sets {
		if inFlight(s.State) {
			flight = append(flight, i)
		}
	}
	if len(flight) < 2 || rapid.IntRange(0, 3).Draw(t, "depends?") == 0 {
		return sets
	}
	from := rapid.SampledFrom(flight).Draw(t, "dependent")
	to := rapid.SampledFrom(flight).Draw(t, "dependency")
	if from == to {
		return sets
	}
	sets[from].Ops = append(sets[from].Ops, op.Op{ID: fid(900), Type: op.TypeLinkCreate, Link: ln(50),
		LinkType: op.LinkDependsOn, From: sets[from].ID, To: sets[to].ID})
	return sets
}

// TestAsIsSeedWhatIfFrom: AsIsSeed's model and problems are AsIs's,
// ToBeFrom over its seed is ToBe over the whole log, and WhatIfFrom over it
// is WhatIf for every in-flight change set. ToBeFrom and every WhatIfFrom go
// on from one shared seed, the results are compared only after every call
// has run, and the seed is then still a fresh AsIsSeed's.
func TestAsIsSeedWhatIfFrom(t *testing.T) {
	var deps, skips, handles, views int
	rapid.Check(t, func(t *rapid.T) {
		sets := genSeedLog(t)
		m, ps, seed := AsIsSeed(sets)
		wm, wps := AsIs(sets)
		sameModel(t, "AsIsSeed", m, wm, ps, wps)

		type view struct {
			cs id.ID
			m  Model
			ps []Problem
		}
		tm, tps := ToBeFrom(seed, sets)
		var got []view
		for _, s := range sets {
			if !inFlight(s.State) {
				continue
			}
			m, ps := WhatIfFrom(seed, sets, s.ID)
			got = append(got, view{s.ID, m, ps})
		}
		wtm, wtps := ToBe(sets)
		sameModel(t, "ToBeFrom", tm, wtm, tps, wtps)
		for _, v := range got {
			views++
			if len(changeset.Closure(sets, v.cs)) > 0 {
				deps++
			}
			for _, p := range v.ps {
				switch {
				case errors.Is(p.Err, ErrHandle):
					handles++
				case p.ChangeSet == v.cs && (errors.Is(p.Err, ErrBase) || errors.Is(p.Err, ErrMissing) ||
					errors.Is(p.Err, ErrRetired) || errors.Is(p.Err, ErrExists)):
					skips++
				}
			}
			wm, wps := WhatIf(sets, v.cs)
			sameModel(t, "WhatIfFrom "+v.cs.String(), v.m, wm, v.ps, wps)
		}
		if _, _, fresh := AsIsSeed(sets); !reflect.DeepEqual(seed, fresh) {
			t.Fatalf("the what-ifs changed the shared seed:\n%+v\nvs a fresh one\n%+v", seed, fresh)
		}
	})
	t.Logf("what-ifs %d, with an in-flight dependency %d, skipped ops of the change set %d, duplicate handles %d",
		views, deps, skips, handles)
	if views == 0 || deps == 0 || skips == 0 || handles == 0 {
		t.Errorf("vacuous run: what-ifs %d, dependencies %d, skips %d, duplicate handles %d", views, deps, skips, handles)
	}
}
