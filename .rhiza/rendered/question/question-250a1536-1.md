# question-250a1536-1

question · `01M78S1GBQGMS8H4F359A5SBQS`

## session

250a1536-7da7-4f36-b513-01a1d18ad9d4

## statement

Which landed intent should change set 01M78RM (The local ship gate switches) serve? Its related gate bug 01M78FD0 serves m4-runs-hold-up ('A build, review or ship-change run finishes without hand repair', milestone-4, open); the original CI-stability work served milestone-25, which has no open change sets on the roadmap. / What should a CI failure that names no test (only 'job / step', 'UNKNOWN STEP' or 'CI run <id>') and passes on rerun do once it no longer counts as a flake? / How is 'failed more than once in the window' set for a named flaky test? / How should the local gate keep its full output?

## Links

### In

- answers [local-gate-repeated-named-flakes](../decision/local-gate-repeated-named-flakes.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
