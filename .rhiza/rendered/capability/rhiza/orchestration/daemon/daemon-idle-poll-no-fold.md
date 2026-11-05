# daemon-idle-poll-no-fold

capability · `01M6H72FY8WGHYD44YXB23B2ZA` · parent [daemon](../daemon.md)

## acceptance

A Go test in internal/api builds a repo with a large log (100+ change sets, several in flight), its index (IndexPath set) and an event command set, runs the whole daemon poll closure (daemonPoll1, daemonNotify and daemonEvents) twice with nothing changed and asserts that the second poll, across all three passes, reads the log zero times through the ctx.readLog counter and also runs planRead zero times and folds zero in-flight change sets through a counter hook in planRead/planGather (the index open and the changeset.FoldCanonical of each EachInFlight change set, which do not pass through ctx.readLog); the same test advances the clock past the 15-minute unjournaled hold with no stamp change and asserts the next poll picks the held change set again, and adds a log entry with a lifecycle move and asserts the next poll reads the log and runs planRead again and runs the event command for that move; a Go benchmark reports ns/op for the unchanged poll, and no test asserts a wall-clock threshold

## statement

A daemon poll where the log, git refs and daemon settings are unchanged since the last poll reuses that poll's folded workspace, plan result, candidates, change-set lookups and overlaps, and the lifecycle moves the events pass (daemonEvents) reads from the folded log, instead of re-reading and re-folding the log, while run journals, process liveness, clock-based holds, resume backoff, the disk check, picks, queued work, the notify pass and running the event command for moves not yet handled are still evaluated every poll

## Links

### Out

- decided-by [daemon-poll-caches-fold-not-picks](../../../../decision/daemon-poll-caches-fold-not-picks.md)

### In

- verifies [daemon-idle-poll-no-fold-verified](../../../../evidence/daemon-idle-poll-no-fold-verified.md)
- verifies [daemon-idle-poll-no-fold-qa-1](../../../../evidence/daemon-idle-poll-no-fold-qa-1.md)
- verifies [pr-217-fix-the-daemon-pins-a-full-cpu](../../../../evidence/pr-217-fix-the-daemon-pins-a-full-cpu.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
