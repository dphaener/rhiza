# ci-unstable-repeated-flake-qa-1

evidence · `01M78XS0F6ZFWP40FN3VV9ZDMW` · actor [claude-code-392142a0](../actor/claude-code-392142a0.md)

## result

PASS

## statement

QA at 266a2e2a3 with the built binary in a scratch repo, hand-written CI stability records: one TestX flake -> rhiza health 'stable: TestX was seen once and ignored' and rhiza ship gate 'Local gate skipped: CI is stable (TestX was seen once and ignored).' exit 0; TestX twice on one head -> stable 'flaked 2 times on one head and was ignored'; TestX at runs 3 and 5 of 8 -> unstable 'TestX flaked 2 times on 2 heads in the last 20 CI runs (stops applying after 16 more completed CI runs)'; with git config rhiza.ship.ciWindow 5 -> '...last 5 CI runs (stops applying after 1 more...)', and with one more clean run appended -> stable; red main -> unstable 'main's latest CI result is red (m1)'; rhiza daemon set rhiza.ship.ciWindow 0 refused, exit 2.

## Links

### Out

- verifies [ci-unstable-repeated-flake](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-unstable-repeated-flake.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
