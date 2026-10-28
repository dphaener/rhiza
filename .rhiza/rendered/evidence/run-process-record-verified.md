# run-process-record-verified

evidence · `01M627ETWW2A1HH0H5HBRYWW83` · actor [claude-code-e0b67628](../actor/claude-code-e0b67628.md)

## result

PASS

## statement

At 4963884 (attempt 2), go test -count=1 -run TestRunJournalsProcess ./internal/api passes: run.start with RexSession/RexBlock (the CLI's REX_BLOCK, cmd/rhiza/ctx.go:75) set journals one process entry with this host, os.Getpid(), the block and the window from a stubbed `rex block inspect -s <session> --json <block>`; run.resume from another block appends a second entry, which State() folds as the run's process with Liveness LiveRunning; resumed outside Rex it journals host and pid only and never calls rex. go test -count=1 -run 'TestLiveness|TestStateProcessLatestWins' ./internal/run passes: a fixture's unfinished journal reads LiveInterrupted with a dead pid, LiveRunning with the test's own pid, and a finished journal NotLive. Read the diff: the fix adds only internal/api/run_process_test.go plus actor log/render files, addressing attempt 1's finding f1. make check exits 0.

## Links

### Out

- verifies [run-process-record](../capability/rhiza/workflows/workflow-runs/run-process-record.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
