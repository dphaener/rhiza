# ship-gate-run-log-qa-1

evidence · `01M78XS1FVYGZ7N3057JH9E12C` · actor [claude-code-392142a0](../actor/claude-code-392142a0.md)

## result

PASS

## statement

QA at 266a2e2a3 with the built binary in a scratch repo, CI unstable (TestX on 2 heads), rhiza.ship.localGate set to a script printing '--- FAIL: TestX', 20000 filler lines (~400 KB) and make's Error 1, exit 2. Outside a run: rhiza ship gate exit 1, log at .git/rhiza/gate/cs/qa-gate-test.log with all 20004 lines; last lines of output were '--- FAIL: TestX', 'FAIL ...', 'Full log: <path>'. With RHIZA_RUN set: two failing gates wrote runs/<run>/gate-1.log and gate-2.log (both 408997 bytes, first intact), a passing gate gate-3.log; RHIZA_RUN=../evil wrote no log. A real foreground run (rhiza run start --workflow <file> --foreground, one command step 'rhiza ship gate') on a new head: run status --json step detail kept '--- FAIL: TestX' and 'Full log: .../runs/01M78XMMZH1P5E6BGY40B33SDS/gate-1.log'. grep -c 'run `rhiza ship gate` to see' prints 0 for internal/workflow/builtin/ship-change.toml, .rhiza/workflows/ship-change.toml, .rhiza/workflows/docs.toml; each has 2 'Full log' prompts.

## Links

### Out

- verifies [ship-gate-run-log](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-run-log.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
