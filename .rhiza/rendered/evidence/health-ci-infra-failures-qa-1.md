# health-ci-infra-failures-qa-1

evidence · `01M78XS12C8AG1E69MXVECHQHH` · actor [claude-code-392142a0](../actor/claude-code-392142a0.md)

## result

PASS

## statement

QA at 266a2e2a3 with the built binary in a scratch repo: a record with one TestX flake, one 'scripts (4/4) / UNKNOWN STEP' and one 'CI run 77' -> rhiza health 'ok ci stable: TestX was seen once and ignored; ignored infra failures: CI run 77, scripts (4/4) / UNKNOWN STEP' with detail lines 'flaky TestX: 1 flake in the last 20 runs' and 'infra <job / step>: 1 failure in the last 20 runs; infra failures do not make CI unstable' for each; rhiza health --json carries the same part (state ok, same detail). TestX on 2 heads -> warn naming TestX, its count 2 and 'stops applying after 16 more completed CI runs', plus the no-local-gate warning.

## Links

### Out

- verifies [health-ci-infra-failures](../capability/rhiza/orchestration/shipping/wait-green-pr/health-ci-infra-failures.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
