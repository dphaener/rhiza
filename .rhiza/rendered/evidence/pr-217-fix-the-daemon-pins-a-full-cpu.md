# pr-217-fix-the-daemon-pins-a-full-cpu

evidence · `01M6HGE0QT0YXXJ37RJ2541EZM` · actor [dphaener](../actor/dphaener.md)

## head

9f8b2d1a3379b09a18c9b5d271953ec689ce95b7

## pr

https://github.com/dphaener/rhiza/pull/217

## result

1 check passed: check SUCCESS

## statement

Pull request #217 (https://github.com/dphaener/rhiza/pull/217) passed CI at 9f8b2d1

## status

verified

## Links

### Out

- verifies [daemon-idle-poll-no-fold](../capability/rhiza/orchestration/daemon/daemon-idle-poll-no-fold.md)
- verifies [daemon-poll-cache-invalidates](../capability/rhiza/orchestration/daemon/daemon-poll-cache-invalidates.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
