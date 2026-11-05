# daemon-nudges-stalled-verified-2

evidence · `01M6NG792DY06EZXCG0VYDK6C8` · actor [claude-code-4fbdf9ea](../actor/claude-code-4fbdf9ea.md)

## result

PASS

## statement

Verifier at 13770358 (attempt 2): `go test -trimpath -count=1 -run 'TestDaemonNudge|Stall' -v ./internal/api/` passes TestDaemonNudgeStalled. Read the test and diff: on a fake clock with a session.Store and a fake sender, an autonomous run A whose session 'auto' now reports hook state idle (Autonomous record, so session.Stalled stalls it after StallAfter=20m) and an owner-attended run O whose session reports working go through daemonNudgeStalled: no send at 19m; exactly one message to 'auto' at 21m, checked to contain 'background shell' and '`rhiza session done`'; a poll at 25m on the same stall sends nothing; runStalledAll still lists both stalled sessions (2) after the nudge; a new report then 21 silent minutes sends a second message; 'owner' is never sent to; a failed send is retried next poll. daemonNudgeStalled skips states whose own run.State.Autonomous is false and keys each stall on session@since; daemonNudge feeds it live, non-ended runs and is called from runDaemon's poll (daemon.go:195). This fixes compliance f1 (idle stall now driven through the nudge).

## Links

### Out

- verifies [daemon-nudges-stalled](../capability/rhiza/orchestration/daemon/daemon-nudges-stalled.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
