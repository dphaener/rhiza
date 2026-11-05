# daemon-nudges-stalled

capability · `01M6N9MBXAAWMHC53534PGKQ4H` · parent [daemon](../daemon.md)

## acceptance

A Go test of the daemon poll on a fake clock, session store and message sender, with one autonomous run and one owner-attended run under the same daemon: a session of the autonomous run idle without done for 21 minutes gets one message naming background shells and `rhiza session done`; a second poll at the same stall sends nothing; a new hook report followed by another 21 silent minutes gets a second message; the owner-attended run's session, stalled the same way, gets no message on any poll; the inbox still lists both stalled sessions after the nudge

## statement

The daemon sends each stalled session of its live autonomous runs (the run's own Autonomous, as its journal records it, not the daemon's --autonomous flag) one message per stall, keyed on the report time it stalled at, telling it to check any background shell's result and carry on, or run `rhiza session done` if its work is finished; a session still stalled after the nudge gets no further one and stays in the inbox, and a stalled session of an owner-attended run is never nudged by the daemon

## Links

### Out

- depends-on [run-session-idle-stall](../../workflows/workflow-runs/run-status/run-session-idle-stall.md)
- decided-by [idle-run-session-stalls](../../../../decision/idle-run-session-stalls.md)

### In

- verifies [daemon-nudges-stalled-verified](../../../../evidence/daemon-nudges-stalled-verified.md)
- verifies [daemon-nudges-stalled-qa-1](../../../../evidence/daemon-nudges-stalled-qa-1.md)
- verifies [daemon-nudges-stalled-verified-2](../../../../evidence/daemon-nudges-stalled-verified-2.md)
- verifies [pr-300-fix-a-run-s-agent-session-that](../../../../evidence/pr-300-fix-a-run-s-agent-session-that.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
