# health-ci-infra-failures

capability · `01M78S7501GDG9S1WVKE58HFKW` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run TestHealthCI shows a record with one 'scripts (4/4) / UNKNOWN STEP' infra failure and one single flake of TestX: the part ok (stable), listing the infra failure with count 1 and TestX with count 1 as ignored; a record where TestX flaked on 2 heads: the part warn, naming TestX, its count and when it stops applying; go test -trimpath ./cmd/rhiza -run TestScripts/health-ci passes with the infra line in rhiza health and health --json

## statement

rhiza health's CI part lists, apart from the flaky tests and their counts, the infra failures in the window with their count per job / step, and states plainly that they do not make CI unstable

## Links

### Out

- depends-on [ci-unstable-repeated-flake](ci-unstable-repeated-flake.md)
- depends-on [health-ci-stability](../../daemon/daemon-status/health-command/health-ci-stability.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../decision/local-gate-repeated-named-flakes.md)

### In

- verifies [health-ci-infra-failures-verified](../../../../../evidence/health-ci-infra-failures-verified.md)
- verifies [health-ci-infra-failures-qa-1](../../../../../evidence/health-ci-infra-failures-qa-1.md)
- verifies [pr-509-the-local-ship-gate-switches](../../../../../evidence/pr-509-the-local-ship-gate-switches.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
