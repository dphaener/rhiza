# health-ci-stability

capability · `01M6YCRKNPNZY0F0BTDCXQAP74` · parent [health-command](../health-command.md)

## acceptance

go test -trimpath ./internal/api -run TestHealthCI shows the part ok on a clean record, ok listing TestX with count 1 as ignored on a single flake, warn naming TestX and its count when TestX flaked on 2 distinct heads, warn naming main red, and the no-local-gate warning only when unstable with rhiza.ship.localGate unset; go test -trimpath ./cmd/rhiza -run TestScripts/health-ci prints the part in rhiza health and health --json

## statement

rhiza health shows a CI part: stable or unstable with the reason (a named test flaked on 2 or more distinct heads in the window, main red), each flaky test with its flake count (single flakes marked as ignored), and a warning when CI is unstable and no local gate (rhiza.ship.localGate) is set

## Links

### Out

- depends-on [ci-stability-state](../../../shipping/wait-green-pr/ci-stability-state.md)
- depends-on [ship-local-gate](../../../shipping/wait-green-pr/ship-local-gate.md)
- decided-by [ci-unstable-gate-rules](../../../../../../decision/ci-unstable-gate-rules.md)
- decided-by [local-gate-trigger-identified-flakes](../../../../../../decision/local-gate-trigger-identified-flakes.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../../decision/local-gate-repeated-named-flakes.md)

### In

- depends-on [health-ci-infra-failures](../../../shipping/wait-green-pr/health-ci-infra-failures.md)
- verifies [health-ci-stability-verified](../../../../../../evidence/health-ci-stability-verified.md)
- verifies [health-ci-stability-qa-1](../../../../../../evidence/health-ci-stability-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)
- verifies [pr-509-the-local-ship-gate-switches](../../../../../../evidence/pr-509-the-local-ship-gate-switches.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
