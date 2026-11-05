# daemon-nudges-stalled-qa-1

evidence · `01M6NFYX5NCNY1GEZGZQNWEHK6` · actor [claude-code-07b49e62](../actor/claude-code-07b49e62.md)

## result

PASS

## statement

QA cycle 1 with the built binary (go build at 1bae8c32) in a scratch repo, fake rex/claude on PATH: two live runs of a one-agent-step workflow, one started with run start --autonomous, one owner-attended; session records' hook reports hand-set, clock moved with RHIZA_TEST_NOW. Autonomous session idle since 00:00, owner session working since 00:00: rhiza daemon --once at 00:19 sent nothing; at 00:21 rex.log showed one send to the autonomous session's block pointing at messages/<sid>/0001.md, whose text names a background shell and `rhiza session done`; daemon.state.json nudged=[<sid>@2030-01-01T00:00:00Z]; poll at 00:25 sent nothing and rhiza inbox still listed both runs stalled; a new hook report at 00:30 then polls at 00:31 (nothing) and 00:51 sent 0002.md; 00:55 nothing. The owner-attended run's stalled session was never sent to, including under rhiza daemon --once --autonomous. An autonomous session marked done got no nudge.

## Links

### Out

- verifies [daemon-nudges-stalled](../capability/rhiza/orchestration/daemon/daemon-nudges-stalled.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
