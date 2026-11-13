# slow-poll-plan-2026-11-13

source · `01M79ZZXEXW704Y9Y51WSANGQJ`

## quotes

```json
[]
```

## raw_sha256

0e87e334e2d5cf5e495088ecf5bcbbc04189f3e9a2a219db2e015db95e70b21a

## session

a6499807-5d6b-448c-88f4-01a1d3fb37ca

## summary

Autonomous plan run for bug request 01M79YDVQC0XXJ652A36ST5BDJ: no owner was asked, so nothing here is quoted. The planner's research found the following. (1) Every multi-minute poll in the daemon journal on 2026-11-13 (3m3s to 7m9s, ten of them between 01:59Z and 07:53Z) ends 5 to 20 s after a 'leftovers sweep freed' line, and no poll near a sweep was short; polls without a sweep stayed under about 55 s. The hourly leftovers sweep (daemonSweepLeftovers, called from the home project's prepare step inside the poll) runs leftovers.Scratch, which calls filepath.Glob('<TMPDIR>/claude-*/*/<session id>') once per recorded session: the rhiza project records 4136 sessions and ~/tmp/claude-1000 holds 2657 folders, so one sweep reads that folder about 4k times and makes about 11M lstat calls, under a load of 10-28. That is the step that blocks for minutes. (2) The four defunct [git] children of rhiza:935303 started at 10:33:36, 13:02:43, 20:54:19 and 23:33:28 on 2026-11-12, each in the same second as a 'binary changed, re-executing' line: git programs still running when the daemon re-executed itself, whose waiting goroutines died with the old image. The driver list reaps only detached run drivers, so nothing waits on them. (3) The daemon's per-project git calls (projectCtx's Git and the home ctx's gitIn) and most ctx.Exec calls in a poll run with no time limit; only the model-branch network work (30 s), the self-update fetch (30 s), hooks (10 s) and the failure classifier (60 s) are bounded. (4) The poll-behind watchdog line says only how long the poll took, and it is written after the poll ends. The planner serves the landed intent m24-on-the-vm, which the earlier daemon CPU fix 01M6H6SGKR served. The decisions below are the planner's, made under autonomy from this research: make the sweep one listing, take housekeeping off the poll, bound each git program a poll runs at 30 s, name the slow steps in the existing log line and watchdog record and add one line while a poll runs over its limit, and have each new daemon image adopt its inherited children onto the driver list.

## title

Plan of the slow root daemon poll fix (autonomous plan run)

## Links

### In

- derived-from [decision-scratch-sweep-one-listing](../decision/decision-scratch-sweep-one-listing.md)
- derived-from [decision-housekeeping-off-poll](../decision/decision-housekeeping-off-poll.md)
- derived-from [decision-poll-git-timeout](../decision/decision-poll-git-timeout.md)
- derived-from [decision-slow-step-in-log](../decision/decision-slow-step-in-log.md)
- derived-from [decision-adopt-inherited-children](../decision/decision-adopt-inherited-children.md)

## History

- `01M79ZNKT8C95GCF3F1X03G97F` Fix The root daemon's poll sometimes takes minutes, not seconds: on 2026-11-13 the watchdog logged polls of 6m28s (07:53:03Z), 4m3s (05:41:07Z), 7m5s (05:26:58Z), 7m9s (05:09:37Z), 3m3s, 3m41s, 3m50s and 3m15s (01:59-04:10Z), against a 15s limit. While a poll runs, Last poll in rhiza daemon status stays put (07:46:32Z until 07:53:26Z), so no run starts, no PR merges and the watchdog only reports afterwards. Load was 10-28 at the time (3 live runs with go test). The daemon (process rhiza:935303) also holds 4 defunct [git] children it never reaped (pids 356889, 358417, 2171743, 2340896, 8-21h old). Nothing was done: each long poll ended by itself. Wanted: find what in a poll can block for minutes (a git call under load, a lock, a slow project) and bound it with a timeout so one slow step can't stop the loop; reap git children; journal which step was slow when the watchdog fires. Related but different: 01M6G0KFRR and 01M78QS887 (CPU while idle).
