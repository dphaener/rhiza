# ci-unstable-repeated-flake

capability · `01M78S74SS6J69T9H1YPMK7YSR` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestCIStability|TestShipLocalGate' shows: one flake of TestX in the window stable, with a reason saying TestX was seen once and ignored, and ship gate skipping with that reason; TestX flaked on the same head twice stable; TestX and TestY flaked once each stable; TestX flaked on 2 distinct heads unstable, the reason naming TestX, its count 2 and the right N for the default window of 20 and for rhiza.ship.ciWindow=5, and stable again after N more clean runs; a window holding only infra failures stable, its reason naming them as ignored infra; a red main result unstable with no flake, and stable after a green ship check --main

## statement

Rhiza judges a repo's CI unstable only while one named test was recorded as a flake on at least 2 distinct heads (main or run PR) within the last rhiza.ship.ciWindow completed CI runs, or while main's latest result is red; the unstable reason names each such test, its flake count in the window and when it stops applying (after N more completed CI runs, N being the runs until the older of its two latest flakes leaves the window), and a stable verdict whose window holds infra failures or single flakes says they were seen and ignored, as does ship gate's skip reason

## Links

### Out

- depends-on [ci-flakes-name-tests](ci-flakes-name-tests.md)
- depends-on [ci-stability-state](ci-stability-state.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../decision/local-gate-repeated-named-flakes.md)

### In

- depends-on [health-ci-infra-failures](health-ci-infra-failures.md)
- verifies [ci-unstable-repeated-flake-verified](../../../../../evidence/ci-unstable-repeated-flake-verified.md)
- verifies [ci-unstable-repeated-flake-qa-1](../../../../../evidence/ci-unstable-repeated-flake-qa-1.md)
- verifies [pr-509-the-local-ship-gate-switches](../../../../../evidence/pr-509-the-local-ship-gate-switches.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
