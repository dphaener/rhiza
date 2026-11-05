# daemon-cpu-plan-autonomous

source · `01M6H71GW7QB7VGB1HXMXBV560`

## quotes

```json
[]
```

## raw_sha256

dd22db6f1cd15dc8fc00b44fc1787c210cc7352472221fcbc4478aa2f1dfbf78

## session

37d69292-5976-44cc-847c-01a1a26d0332

## summary

Autonomous plan of change set 01M6H6S (the daemon pins a full CPU core on the 4-core VM because every 5 s poll re-reads and re-folds the whole log). No owner present. The planner read daemonPoll1, planRead, openCtxWorkspace and inboxOverlaps in internal/api and raised one question to the arbiter: what an idle poll may skip, and how a poll's cost is measured. The arbiter decided (daemon-poll-caches-fold-not-picks): cache only the fold-derived work (plan result, workspace, candidates, change-set lookups) keyed on log, git ref and settings stamps; keep run journals, liveness, clock-based holds, picks and queued work evaluated every poll; the daemon's own writes invalidate the cache; measure with a deterministic fold-count test plus an optional benchmark, no wall-clock threshold. The planner put the change set under milestone 24's on-the-VM intent, as the fault shows on the VM, and planned two leaves under the daemon capability.

## title

Daemon CPU fix plan, autonomous run

## Links

### In

- derived-from [daemon-poll-caches-fold-not-picks](../decision/daemon-poll-caches-fold-not-picks.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
