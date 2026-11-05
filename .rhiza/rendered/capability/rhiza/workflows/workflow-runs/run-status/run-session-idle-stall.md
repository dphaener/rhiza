# run-session-idle-stall

capability · `01M6N9M6BMS0XZBJDD16M2B0Z2` · parent [run-status](../run-status.md)

## acceptance

A Go test of session.Stalled: an idle hook report of an autonomous session at 19 minutes is not stalled and at 21 minutes is; at 21 minutes an idle session with Autonomous false is not stalled while a working one is; at 21 minutes an idle or working session with the done mark, with Status.Asked set, paused at its usage limit, blocked, or exited is not; inbox-stalled.txtar gains an idle record without done and with Autonomous false that gives no stalled item at T+21m in the inbox or run status, then the same record marked autonomous gives a stalled item at T+21m and none at T+19m, and run status shows it stalled; an idle record with the done mark gives none

## statement

rhiza run status and rhiza inbox count a live run-step session as stalled once its hook status has gone unchanged for more than 20 minutes and it has not said `rhiza session done`: working in any run, or idle in a session an autonomous run launched, whether or not it holds a background shell; an idle session of an owner-attended run, a session waiting on the arbiter's answer, paused at its usage limit, blocked or exited never stalls

## Links

### Out

- decided-by [idle-run-session-stalls](../../../../../decision/idle-run-session-stalls.md)

### In

- depends-on [daemon-nudges-stalled](../../../orchestration/daemon/daemon-nudges-stalled.md)
- blocks [review-01m6n9ad-r3-f1](../../../../../question/review-01m6n9ad-r3-f1.md)
- verifies [run-session-idle-stall-verified](../../../../../evidence/run-session-idle-stall-verified.md)
- verifies [run-session-idle-stall-qa-1](../../../../../evidence/run-session-idle-stall-qa-1.md)
- verifies [pr-300-fix-a-run-s-agent-session-that](../../../../../evidence/pr-300-fix-a-run-s-agent-session-that.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
