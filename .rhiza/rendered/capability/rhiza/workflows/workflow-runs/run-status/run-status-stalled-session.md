# run-status-stalled-session

capability · `01M6A0NWP7AH4JGMDT37F6X7DG` · parent [run-status](../run-status.md)

## acceptance

A Go test of the stall cutoff: working at 19 minutes is not stalled, at 21 minutes is; blocked, usage-paused and exited sessions never are, and an idle one stalls only as run-session-idle-stall says; a testscript with a hand-written session record and RHIZA_TEST_NOW shows 'stalled' in run status only past the cutoff

## statement

rhiza run status marks a live run-step session stalled as run-session-idle-stall says: once its hook status, working or idle, has gone unchanged for more than 20 minutes and it has not said `rhiza session done`, unless it is waiting on the arbiter's answer, paused for a usage limit, blocked or exited

## Links

### Out

- decided-by [stalled-after-20-minutes](../../../../../decision/stalled-after-20-minutes.md)
- decided-by [idle-run-session-stalls](../../../../../decision/idle-run-session-stalls.md)

### In

- verifies [run-status-stalled-session-verified](../../../../../evidence/run-status-stalled-session-verified.md)
- verifies [run-status-stalled-session-qa-1](../../../../../evidence/run-status-stalled-session-qa-1.md)
- verifies [run-status-stalled-session-verified-2](../../../../../evidence/run-status-stalled-session-verified-2.md)
- verifies [run-status-stalled-session-qa-2](../../../../../evidence/run-status-stalled-session-qa-2.md)
- verifies [pr-158-stalled-sessions-surfaced](../../../../../evidence/pr-158-stalled-sessions-surfaced.md)
- verifies [pr-300-fix-a-run-s-agent-session-that](../../../../../evidence/pr-300-fix-a-run-s-agent-session-that.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
