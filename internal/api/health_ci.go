package api

// This file holds health's CI part: whether the repo's CI is stable, read
// from the CI stability record, each flaky test with its flake count, and a
// warning when CI is unstable and no local gate is set.

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ciFlakeCounts counts each test's flakes over the last window runs of rec.
func ciFlakeCounts(rec CIStability, window int) map[string]int {
	window = max(window, 1)
	runs := rec.Runs
	if len(runs) > window {
		runs = runs[len(runs)-window:]
	}
	counts := map[string]int{}
	for _, r := range runs {
		for _, f := range r.Flakes {
			counts[f]++
		}
	}
	return counts
}

// healthCI is ok while CI is stable and warn while it is unstable (a flake
// in the window, or main's latest result red), with each flaky test and its
// count as detail. Unstable with no rhiza.ship.localGate set, it also warns
// that nothing runs before a push, and points at the setting.
func healthCI(ctx *Ctx) HealthPart {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return healthFailed(HealthCI, err)
	}
	rec, err := readCIStability(ciStabilityPath(common))
	if err != nil {
		return healthFailed(HealthCI, err)
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return healthFailed(HealthCI, err)
	}
	window, err := strconv.Atoi(settingValue(rows, CIWindowKey))
	if err != nil {
		return healthFailed(HealthCI, fmt.Errorf("%s: %v", CIWindowKey, err))
	}
	v := ciVerdict(rec, window)
	if !v.Unstable {
		return HealthPart{Part: HealthCI, State: HealthOK, Summary: "stable"}
	}
	p := HealthPart{Part: HealthCI, State: HealthWarn, Summary: "unstable: " + v.Reason}
	counts := ciFlakeCounts(rec, window)
	tests := make([]string, 0, len(counts))
	for name := range counts {
		tests = append(tests, name)
	}
	sort.Strings(tests)
	for _, name := range tests {
		n := "flakes"
		if counts[name] == 1 {
			n = "flake"
		}
		p.Detail = append(p.Detail, fmt.Sprintf("flaky %s: %d %s in the last %d runs", name, counts[name], n, max(window, 1)))
	}
	if strings.TrimSpace(settingValue(rows, LocalGateKey)) == "" {
		p.Summary += "; no local gate is set, so nothing runs before a push"
		p.Next = "rhiza daemon set " + LocalGateKey + " <command>"
	}
	return p
}
