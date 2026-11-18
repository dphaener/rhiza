# session-process-tree-reading

capability · `01M7QN2RCKAVANW6QNGVNZ7W1H` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath -p 2 ./internal/proctree passes, with tests over fake ps text (never the host's ps): a table holding two sessions' trees, nested three deep, gives each session its own count and summed rss in bytes (KiB x 1024), the root counted, the other session's processes and unrelated processes left out; a 'rhiza session run' line whose session argument only shares a prefix with the asked id is not its root. (failure case) a session with no root process gives not-found, not a zero reading; malformed and header lines are skipped; a ps that fails returns an error naming ps

## files

internal/proctree/proctree.go, internal/proctree/proctree_test.go

## must_not

Read /proc or use any Linux- or darwin-only call; add a dependency; touch internal/api or the watchdog; run ps in a test

## statement

A new internal/proctree package reads the host's process table portably with one ps call and gives a session's process tree: the 'rhiza session run <session-id>' root, its descendants with pid, ppid, rss and args, the process count and the summed resident memory in bytes

## Links

### Out

- decided-by [runaway-tree-from-session-run](../../../../../decision/runaway-tree-from-session-run.md)

### In

- depends-on [runaway-session-fails-step](runaway-session-fails-step.md)
- verifies [session-process-tree-reading-verified](../../../../../evidence/session-process-tree-reading-verified.md)
- verifies [session-process-tree-reading-qa-1](../../../../../evidence/session-process-tree-reading-qa-1.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
