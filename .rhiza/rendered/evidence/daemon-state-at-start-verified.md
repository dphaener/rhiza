# daemon-state-at-start-verified

evidence · `01M6CXNBMA75WR805AM116CD2M` · actor [claude-code-4765b7b5](../actor/claude-code-4765b7b5.md)

## result

PASS

## statement

At 73aaf370 ran go test ./internal/api -run TestDaemon -count=1 -v: TestDaemonStateAtStart passes. It seeds a state file from an earlier image of the same pid (last poll, max runs 3, run RUN1, skip CS2, unjournaled CS3), takes the lock as the new image, writes daemonStartState and calls daemon.status before any poll; it asserts running, the new lock start time, the old last poll, the request's max runs 5 and autonomous, the seeded run, skip and hold, and no stopped line. Read the diff: runDaemon now builds the state with daemonStartState and calls daemon.WriteState before its first poll; Stopped is not carried. Without the start write status would see a holder mismatch and drop the state, so the test is not vacuous. make check passed (exit 0) on a rerun; the first run failed only in the read-budget timing test under load.

## Links

### Out

- verifies [daemon-state-at-start](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary/daemon-state-at-start.md)

## History

- `01M6CA9QMHN6N95AS3F935ZVN2` daemon status keeps the daemon's state through a re-exec, and a terminal current change set comes with the command that clears it
