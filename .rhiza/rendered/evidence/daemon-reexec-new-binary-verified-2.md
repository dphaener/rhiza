# daemon-reexec-new-binary-verified-2

evidence · `01M6AS2D250GYCZMPY1F7P6R06` · actor [claude-code-c81acae7](../actor/claude-code-c81acae7.md)

## result

PASS

## statement

At 8e63b943 (attempt 2), go test -count=1 -run 'Reexec|Daemon(Loop|Keeps)' ./internal/api ./internal/daemon passed: TestDaemonReexecsRebuiltBinary drives runDaemon over 4 polls with an injected stat and exec, a changed executable execs once after the first poll wrote state, with the lock free for the new image (TryLock inside Exec) and held again after, an unchanged one never execs (0 execs, 4 stats); TestDaemonLoopReexecWaitsForPending shows a pending start delays the exec to poll 3, the first with none pending, not retried after a failed exec; TestReexec and TestDaemonKeepsUnjournaledHolds pass. Read the diff: the fix carries unjournaled holds and their run counts through daemon state so a re-exec keeps the 15-minute hold (the review's quality f1). make check passes fmt-check, vet and test; its self-check fails only on stale .rhiza/rendered actor and index files from this run's model writes, not on code.

## Links

### Out

- verifies [daemon-reexec-new-binary](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
