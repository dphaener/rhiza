# run-session-idle-stall-verified

evidence · `01M6NFBGW9V7B585DY88T6T8NV` · actor [claude-code-c99dfaa9](../actor/claude-code-c99dfaa9.md)

## result

PASS

## statement

At 9abac36b, session.Stalled (internal/session/status.go) stalls a hook-sourced session only when working, or idle with Record.Autonomous true, and never when Status.SaidDone(), Status.Asked is set, usage-paused, blocked or exited. go test -trimpath ./internal/session -run TestStalled PASS: it covers autonomous idle at 19m (not stalled) and 21m (stalled), attended idle at 21m (not stalled), attended working at 21m (stalled), and at 21m idle/working with the done mark, idle/working with Asked, idle usage-paused, autonomous blocked and autonomous exited (none stalled). go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-stalled|describe-golden)$' PASS: inbox-stalled.txtar asserts that an idle record with no autonomous key gives no stalled item in inbox or run status at T+21m; after marking it autonomous, none at T+19m, a stalled inbox item at T+21m, and run status shows stalled; then with the done mark, none. go test -trimpath ./internal/api -run TestRunStalledNow PASS. make self-check fails only on stale .rhiza/rendered files for run actors, which are rendered after landing.

## Links

### Out

- verifies [run-session-idle-stall](../capability/rhiza/workflows/workflow-runs/run-status/run-session-idle-stall.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
