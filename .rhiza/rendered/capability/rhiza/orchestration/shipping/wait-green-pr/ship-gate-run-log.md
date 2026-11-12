# ship-gate-run-log

capability · `01M78S754HVX2TR0XVH15MKYZ4` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestShipGateRunLog|TestShipRebase' shows: with RHIZA_RUN set, a failing gate whose output is over 64 KB with '--- FAIL: TestX' near its start and make's error lines at its end writes a log in the run's directory holding every byte, the response's log field names it, and the last 20 lines of the human output contain '--- FAIL: TestX' and 'Full log: <path>'; a second gate run in the same run writes a second file and leaves the first intact; a passing gate also writes its log; with no RHIZA_RUN the log is rhiza/gate/<branch>.log under the git common dir; a ship-change gate command step that fails keeps a journal detail naming the log path; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after regenerating it for the new field; go test -trimpath ./internal/workflow passes; grep -c 'run `rhiza ship gate` to see' internal/workflow/builtin/ship-change.toml .rhiza/workflows/ship-change.toml .rhiza/workflows/docs.toml prints 0 for all three and each of their gate-fix prompts names the full gate log; make self-check passes

## statement

Each local gate run (rhiza ship gate, and the gate ship rebase runs on a merged head) writes its whole combined output to its own log file in the run's directory, never overwriting an earlier gate log of that run (outside a run, rhiza/gate/<branch>.log under the git common dir); ShipGateResponse's log field and the gate step's journal entry carry the path, a failing gate's human output ends with its failing test lines (go test '--- FAIL:', 'FAIL', 'panic:', at most 10) and 'Full log: <path>', and every gate-fix prompt (internal/workflow/builtin/ship-change.toml, .rhiza/workflows/ship-change.toml and .rhiza/workflows/docs.toml, two each) tells the fixer to read that log instead of rerunning the gate

## Links

### Out

- depends-on [ship-local-gate](ship-local-gate.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../decision/local-gate-repeated-named-flakes.md)
- decided-by [gate-log-kept-whole](../../../../../decision/gate-log-kept-whole.md)
- supersedes ship-gate-full-log (retired)

### In

- verifies [ship-gate-run-log-verified](../../../../../evidence/ship-gate-run-log-verified.md)
- verifies [ship-gate-run-log-qa-1](../../../../../evidence/ship-gate-run-log-qa-1.md)
- verifies [pr-509-the-local-ship-gate-switches](../../../../../evidence/pr-509-the-local-ship-gate-switches.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
