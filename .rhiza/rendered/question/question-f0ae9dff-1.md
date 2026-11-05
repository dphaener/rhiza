# question-f0ae9dff-1

question · `01M6N9DXRSKHJWY3GARDJ68ETG`

## session

f0ae9dff-3406-48f0-9c7d-01a1aa956836

## statement

Change set 01M6N9AF (idle session with a background shell never flagged): today session.Stalled counts only hook-status 'working' sessions silent >20 min; an idle session that never said `rhiza session done` waits forever. What should count as stalled? / Once such a session counts as stalled, who nudges it?

## Links

### In

- answers [idle-run-session-stalls](../decision/idle-run-session-stalls.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
