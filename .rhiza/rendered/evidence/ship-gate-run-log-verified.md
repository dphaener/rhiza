# ship-gate-run-log-verified

evidence · `01M78WYJN8KEY23XJRE37GM2MT` · actor [claude-code-63b13aa9](../actor/claude-code-63b13aa9.md)

## result

PASS

## statement

Verified at ed1957908: go test -trimpath ./internal/api -run 'TestShipGateRunLog|TestShipRebase' -v passes (13 tests). TestShipGateRunLog with RHIZA_RUN set runs a failing gate printing '--- FAIL: TestX', 20000 filler lines (>64 KB) and make's Error 1 line twice: both logs land as distinct files in run.Dir/<run>/ (gate-<n>.log, O_EXCL so never overwritten), each holds the whole output incl. 'Error 1', the last 20 human lines hold '--- FAIL: TestX' and 'Full log: <path>' within stepTailBytes (the command step's journal detail keeps the last 20 lines / 2048 bytes, so a failing gate step's detail names the log), and a passing gate also writes its log there; TestShipRebaseRunsLocalGate shows the no-run log at <common>/rhiza/gate/cs/ship-it.log. go test -trimpath ./internal/workflow ok; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden ok (golden carries the new log desc). grep -c 'run `rhiza ship gate` to see' on internal/workflow/builtin/ship-change.toml, .rhiza/workflows/ship-change.toml, .rhiza/workflows/docs.toml prints 0 for all three, and each file's two gate-fix prompts tell the fixer to read the full gate log (`Full log:`) instead of rerunning rhiza ship gate. make self-check passes (exit 0) in a clean detached worktree of the commit with RHIZA_ROOT/CHANGE/RUN/LEAF/ACTOR unset; in the run worktree it failed only on the verifier session's own unrendered actor record. Read the diff of internal/api/ship_gate.go: run id resolved via the run store, path-traversal ids refused.

## Links

### Out

- verifies [ship-gate-run-log](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-run-log.md)

## History

- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
