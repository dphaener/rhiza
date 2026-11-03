# daemon-poll-cache-invalidates

capability · `01M6H72NZ2CZJTC8B1VGB4HAH0` · parent [daemon](../daemon.md)

## acceptance

A Go test in internal/api runs a poll, then for each of a new log entry (a change set accepted), a branch ref moved, a daemon setting changed with daemon set, and a run the daemon started, asserts through the ctx.readLog counter that the next poll reads the log again and acts on the change (the newly accepted change set is picked, the new setting is used)

## statement

The daemon's poll cache is dropped, and the next poll folds the log again, when a log entry is added, a git ref moves, a daemon setting changes, or the daemon itself starts or resumes a run

## Links

### Out

- decided-by [daemon-poll-caches-fold-not-picks](../../../../decision/daemon-poll-caches-fold-not-picks.md)

### In

- verifies [daemon-poll-cache-invalidates-verified](../../../../evidence/daemon-poll-cache-invalidates-verified.md)
- verifies [daemon-poll-cache-invalidates-qa-1](../../../../evidence/daemon-poll-cache-invalidates-qa-1.md)
- verifies [pr-217-fix-the-daemon-pins-a-full-cpu](../../../../evidence/pr-217-fix-the-daemon-pins-a-full-cpu.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
