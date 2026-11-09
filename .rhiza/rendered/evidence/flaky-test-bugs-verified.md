# flaky-test-bugs-verified

evidence · `01M6YMQHSDR7WW026AF5SKH1GE` · actor [claude-code-1a0f5132](../actor/claude-code-1a0f5132.md)

## result

PASS

## statement

At 21c1eb26 I ran go test -trimpath -count=1 ./internal/api -run TestFlakyTestBugs -v (PASS, and PASS again with -count=3). It covers: one bug request for a newly flaky test naming the test, PR #7 and head aaa111 and asking to reproduce locally and fix or quarantine; no second request on a repeat poll or on a second flake of the same test while the first is open; TestB and TestAB each queued once; no request with fileBugs off. On attempt 1 the finding was that a planned request let the same flake be filed again. The new flakyBugFiled treats queued and planning requests as open. A planned or failed request covers the runs its text names and any flake recorded before its run finished, so only a later flake files the test again. The new test step plans request R1 through a finished journal, sees nothing queued over two polls, then sees exactly one new TestA request for the later flake ddd444. go vet and gofmt are clean on internal/api. A full go test -trimpath -count=1 ./internal/api ran every test with no --- FAIL. The package was failed only by the harness guard that the enclosing ~/code root's .rhiza-root/usage.live.json changed. That file is the live root daemon's and this diff (watchdog_flaky.go and its test) does not touch it. An earlier full run hit go's 10m default timeout under load average ~8; the rerun took 740s.

## Links

### Out

- verifies [flaky-test-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/flaky-test-bugs.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
