# ci-flakes-name-tests

capability · `01M78S74JZND6E67Z618T2TYQM` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipCheckStabilityRecord|TestFlakyTestBugs|TestTrimFailureLog' shows: a rerun-green head whose failed log has '--- FAIL: TestFoo' and '--- FAIL: TestScripts/bar' recorded with flakes TestFoo and TestScripts/bar and no infra entry; a rerun-green head whose log names only 'scripts (4/4) / UNKNOWN STEP' (or no job/step at all) recorded with that job / step (or 'CI run <id>') as infra and no flake; a head still red after its rerun recorded red with neither; an existing record whose flakes hold only 'job / step' entries read as infra; the watchdog queuing a bug for TestFoo and none for an infra-only run

## statement

When rhiza ship check settles a head that went red and then passed on a rerun, it records a flake in the CI stability record only against each failing test the failed log names (go test '--- FAIL: <name>' lines, subtests by their full name), and records a failure that names no test (only 'job / step', 'UNKNOWN STEP' or 'CI run <id>') as an infra failure kept apart from flakes, for which the watchdog files no flaky-test bug

## Links

### Out

- depends-on [ci-stability-record](ci-stability-record.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../decision/local-gate-repeated-named-flakes.md)

### In

- depends-on [ci-unstable-repeated-flake](ci-unstable-repeated-flake.md)
- verifies [ci-flakes-name-tests-verified](../../../../../evidence/ci-flakes-name-tests-verified.md)
- verifies [ci-flakes-name-tests-qa-1](../../../../../evidence/ci-flakes-name-tests-qa-1.md)
- verifies [pr-509-the-local-ship-gate-switches](../../../../../evidence/pr-509-the-local-ship-gate-switches.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
