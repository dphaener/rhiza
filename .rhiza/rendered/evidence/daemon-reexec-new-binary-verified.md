# daemon-reexec-new-binary-verified

evidence · `01M6AQH81E84Q98864ZTRKHXVD` · actor [claude-code-264db3bc](../actor/claude-code-264db3bc.md)

## result

PASS

## statement

At a1222ea3 I ran go test -race -count=3 -run 'TestDaemonReexecsRebuiltBinary|TestDaemonLoopReexecWaitsForPending' ./internal/api (PASS) and go test -run TestReexec ./internal/daemon (PASS); gofmt, go vet and go test ./... all pass. make check passes fmt, vet and test; its self-check fails only on stale .rhiza/rendered files from this run's actor writes, not on code. Read the diff of 63a4cafb and a1222ea3. TestDaemonReexecsRebuiltBinary drives runDaemon over 4 polls with daemonReexec and daemonTicker replaced by a fake stat/exec. With a changed size it execs exactly once. The exec happens after the first poll wrote LastPoll. Inside the exec the daemon lock is free for the new image to take (TryLock succeeds), and the lock is held again on later ticks. With an unchanged image it never execs and stats 4 times. TestDaemonLoopReexecWaitsForPending drives daemonLoop with starts pending at polls 1 and 2: the exec runs once, at poll 3, and is not retried over 5 polls. This holds for both a nil and a failing exec. Due and Do in internal/daemon/reexec.go compare size and mod time and release, then retake, the lock around Exec. This addresses finding f1 from attempt 1.

## Links

### Out

- verifies [daemon-reexec-new-binary](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary.md)

## History

- `01M6A0K064DH4TP3JJ0861DFG4` The daemon re-executes itself when its binary is rebuilt, and daemon status shows its version
