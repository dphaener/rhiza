# pr-509-the-local-ship-gate-switches

evidence · `01M78YMKQESGEG1SKFRKK2YEHD` · actor [dphaener](../actor/dphaener.md)

## head

d0adbcbc76cd3884b452d3edf1535a6f1730abf9

## pr

https://github.com/dphaener/rhiza/pull/509

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #509 (https://github.com/dphaener/rhiza/pull/509) passed CI at d0adbcb

## status

verified

## Links

### Out

- verifies [ci-flakes-name-tests](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-flakes-name-tests.md)
- verifies [ci-unstable-repeated-flake](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-unstable-repeated-flake.md)
- verifies [health-ci-infra-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/health-ci-infra-failures.md)
- verifies [ship-gate-run-log](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-run-log.md)
- verifies [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)
- verifies [health-ci-stability](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-ci-stability.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
