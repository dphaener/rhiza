# pr-300-fix-a-run-s-agent-session-that

evidence · `01M6NJGH59TR8KEH936BWQDMW0` · actor [dphaener](../actor/dphaener.md)

## head

f62b24017ec1d38e40e502defc02018eb50212fa

## pr

https://github.com/dphaener/rhiza/pull/300

## result

1 check passed: check SUCCESS

## statement

Pull request #300 (https://github.com/dphaener/rhiza/pull/300) passed CI at f62b240

## status

verified

## Links

### Out

- verifies [run-session-idle-stall](../capability/rhiza/workflows/workflow-runs/run-status/run-session-idle-stall.md)
- verifies [daemon-nudges-stalled](../capability/rhiza/orchestration/daemon/daemon-nudges-stalled.md)
- verifies [run-status-stalled-session](../capability/rhiza/workflows/workflow-runs/run-status/run-status-stalled-session.md)
- verifies [inbox-stalled-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
