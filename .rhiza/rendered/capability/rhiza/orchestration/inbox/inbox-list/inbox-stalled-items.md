# inbox-stalled-items

capability · `01M6A0NY4FZ07ZQMR23E700DTQ` · parent [inbox-list](../inbox-list.md)

## acceptance

inbox-stalled.txtar: with a working session record last reported at T and a live run journal, rhiza inbox --json has no stalled item at T+19m and one of kind stalled at T+21m; a blocked record gives none, and an idle one only as run-session-idle-stall says; the inbox app test shows the stalled group

## statement

rhiza inbox lists each stalled session of a live run as a stalled item naming the session and its run, with rhiza session message as the next step, and the inbox app shows them

## Links

### Out

- decided-by [stalled-after-20-minutes](../../../../../decision/stalled-after-20-minutes.md)
- decided-by [idle-run-session-stalls](../../../../../decision/idle-run-session-stalls.md)

### In

- verifies [inbox-stalled-items-verified](../../../../../evidence/inbox-stalled-items-verified.md)
- verifies [inbox-stalled-items-qa-1](../../../../../evidence/inbox-stalled-items-qa-1.md)
- verifies [inbox-stalled-items-verified-2](../../../../../evidence/inbox-stalled-items-verified-2.md)
- verifies [inbox-stalled-items-qa-2](../../../../../evidence/inbox-stalled-items-qa-2.md)
- verifies [pr-158-stalled-sessions-surfaced](../../../../../evidence/pr-158-stalled-sessions-surfaced.md)
- verifies [inbox-stalled-items-qa-01m6n9af-1](../../../../../evidence/inbox-stalled-items-qa-01m6n9af-1.md)
- verifies [pr-300-fix-a-run-s-agent-session-that](../../../../../evidence/pr-300-fix-a-run-s-agent-session-that.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
