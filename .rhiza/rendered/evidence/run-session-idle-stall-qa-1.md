# run-session-idle-stall-qa-1

evidence · `01M6NFZ1P972BJT9CVP6DSVYZY` · actor [claude-code-07b49e62](../actor/claude-code-07b49e62.md)

## result

PASS

## statement

QA cycle 1 with the built binary (go build at 1bae8c32) in a scratch repo, fake rex/claude on PATH: two live runs of a one-agent-step workflow, one started with run start --autonomous, one owner-attended; session records' hook reports hand-set, clock moved with RHIZA_TEST_NOW. Both sessions idle at 00:00: rhiza inbox at 00:19 empty; at 00:21 one stalled item, the autonomous run's (the owner-attended idle one gave none); rhiza run status <auto run> --json showed stalled {session, since 2030-01-01T00:00:00Z, step build}. After adding a done mark to the autonomous idle record, inbox at +30m listed no stall for it. Scoped tests: go test ./internal/session ./internal/api -run 'TestStalled|TestDaemonNudge|TestRunStalled' ok.

## Links

### Out

- verifies [run-session-idle-stall](../capability/rhiza/workflows/workflow-runs/run-status/run-session-idle-stall.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
