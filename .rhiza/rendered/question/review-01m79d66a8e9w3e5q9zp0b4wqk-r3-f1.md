# review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f1

question · `01M79E76ASY4F17XQB4MR28MJS`

## statement

Plan review round 3 left finding f1 (contradiction on queue-wait-journal-lines): The leaf and Decision queue-wait-journal-file treat paused, disk-low, load-high, memory-low, usage-reserve and ship-owner as pre-start reasons that only a not-yet-started change set gets, but the code gives them to change sets with a live run too. daemonCandidates (internal/api/daemon.go:1321) returns every accepted, ready change set, live runs included, and daemonPoll1 (internal/api/daemon.go:1020-1062) skips every candidate for paused, disk-low/machine hold, usage-reserve or ship-owner before daemon.Pick ever tells live/held/shipped apart. So under the leaf's rule (a pre-start reason after started means it waits anew), a load or usage hold or a pause during a live run journals a new wait after started, and the perf analysis sums time when the run was in fact running as queue wait. The leaf must say how it tells these apart, for example by not journaling a hold for a change set whose latest run is live or held, and its acceptance check should cover a hold that comes on while the run is live. The critic confirmed it: Checked the code: daemonCandidates keeps every accepted, ready change set with leaves, live runs included, and live/held/shipped/failed are worked out only inside daemon.Pick (internal/daemon/pick.go), which daemonPoll1 calls only in its final else branch. The paused, disk-low/machine-hold and usage-reserve branches (internal/api/daemon.go ~1040-1062) skip every candidate first, so a change set with a live run gets paused, load-high, memory-low or usage-reserve. Under the leaf's rule (a pre-start reason after started means it waits anew), that journals a new wait while the run is running, and the perf sum counts run time as queue wait. The leaf's acceptance doesn't cover this. Major stands. How should it be settled?

## Links

### Out

- blocks [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

### In

- answers [queue-wait-post-start-first](../decision/queue-wait-post-start-first.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
