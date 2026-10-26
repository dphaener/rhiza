# status-in-genesis-files-test

evidence · `01M5WAZ9XHM6Y0QGGMC4RW2VH2` · actor [claude-code-dd482246](../actor/claude-code-dd482246.md)

## statement

internal/session/store_test.go TestStoreWriteReadAtomic passes; internal/session/race_test.go TestStoreUpdateConcurrentWriters and TestStoreUpdateConcurrentProcesses pass under -race (a hook writer and a Rex writer make 200 updates each at once, as goroutines and as separate processes, and neither loses a report, since each update holds a per-record flock; without the lock both fail); internal/api/session_test.go TestSessionLsShow passes; cmd/rhiza/testdata/script/session-run.txtar shows the recorded status with rex removed

## Links

### Out

- verifies [status-in-genesis-files](../capability/rhiza/runtime/agent-status/status-in-genesis-files.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
