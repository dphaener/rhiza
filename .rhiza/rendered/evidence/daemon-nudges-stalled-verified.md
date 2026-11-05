# daemon-nudges-stalled-verified

evidence · `01M6NFQNH1FWC2NJJ0NGHSFDDE` · actor [claude-code-ab6612fb](../actor/claude-code-ab6612fb.md)

## result

PASS

## statement

At 75c23faf ran `go test -trimpath ./internal/api -run 'TestDaemonNudge|Stall' -count=1 -v` (PASS: TestDaemonNudgeStalled, TestRunStalledNow), `go test -trimpath ./internal/session ./internal/daemon` (ok) and go vet (clean). TestDaemonNudgeStalled drives daemonNudgeStalled on a fake clock, a temp session store and a fake sender with one autonomous run (run State Autonomous true) and one owner-attended run: no send at 19m; at 21m silent without done the autonomous session gets exactly one message, asserted to contain 'background shell' and '`rhiza session done`' (attempt 1's finding fixed: the sender now receives the text); a poll at 25m on the same stall sends nothing; after a new report and 21 more silent minutes a second message is sent; the owner-attended session, stalled the same way, is never sent to; runStalledAll still lists both stalled sessions after the nudge; a failed send is retried next poll. Diff read: daemonNudge filters on the run's own st.Autonomous (not the daemon flag), keys stalls as session@report-time, persists Nudged in daemon state across restarts, and is called each poll in runDaemon. The test's stalled sessions report working rather than idle; the nudge path is state-agnostic (it takes runStalledAll's list), and session.Stalled stalls an idle autonomous-launched session the same way, covered by the run-session-idle-stall leaf's tests.

## Links

### Out

- verifies [daemon-nudges-stalled](../capability/rhiza/orchestration/daemon/daemon-nudges-stalled.md)

## History

- `01M6N9AF69ZGMTB1EEG5RK9MPD` Fix A run's agent session that ends its turn while its own background shell runs (e.g. scoped tests started with run_in_background) sits idle at its prompt and is never flagged: no inbox stalled item and no nudge, so the run waits indefinitely. Seen 2026-11-05: run 01M6N133 fixer session 9130a4fa sat at its prompt about an hour after a recap ('Next, once the scoped daemon tests pass, I'll commit...') with 1 shell open, until the orchestrator sent rhiza session message. Stall detection should treat an idle prompt with no tool call for N minutes as stalled even when a background shell exists, and nudge it to check the shell's result.
