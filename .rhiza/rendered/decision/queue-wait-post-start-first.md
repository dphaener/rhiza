# queue-wait-post-start-first

decision · `01M79E8SX5T7HQTS913N1E04XG`

## reasoning

The finding is right: daemonPoll1 (internal/api/daemon.go ~1040-1062) gives every candidate the paused, disk-low/machine-hold or usage-reserve reason before daemon.Pick works out live/held/shipped/failed, and daemonCandidates keeps change sets with live runs, so the leaf's rule (a pre-start reason after started means a new wait) would count run time as queue wait and skew the perf sum the intent exists to produce. The run journals the poll already holds (runs, pending, host, alive, ship) are exactly what decide uses to tell post-start states apart, so applying that classification before journaling needs no new data and no change to the poll's own skip reasons in daemon.state.json, which other surfaces show. Sharing one function with decide keeps the two from drifting. This keeps Decision queue-wait-journal-file's shape (lines only on change, pre-start reasons, started) and only pins how pre-start is told apart, which is what the finding asks for; the extra acceptance case covers the hold-during-live-run path the critic found uncovered.

## statement

Before journaling, the daemon classifies each skipped change set by its run state first, the same way daemon.Pick's decide does (internal/daemon/pick.go): a change set that is live (its latest run is running or interrupted), held (pending, started and not yet journaled), unjournaled, shipped (its latest ship run succeeded after it was last accepted or reopened) or failed (its latest run failed and was neither abandoned nor retried) is in a post-start state, and no line is journaled for it whatever reason the poll gave it, so a pause, disk-low, load-high, memory-low, usage-reserve or ship-owner skip that comes on while its run is live or held writes nothing. A pre-start reason is journaled only for a change set with no post-start state: one with no run yet, or whose latest run was abandoned or retried, or that was reaccepted or reopened after shipping (a change set waiting anew). The classification is one exported function in internal/daemon that decide shares, returning the post-start reason for a change set or none, so the poll and the journal cannot disagree. queue-wait-journal-lines's acceptance gains a case: a change set started and then live, with a usage-reserve hold, a machine (load-high) hold and a pause each coming on during the live run, journals nothing after its started line.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f1](../question/review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f1.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
