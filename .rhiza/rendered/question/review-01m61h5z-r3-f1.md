# review-01m61h5z-r3-f1

question · `01M61S7J7GTD8FC5Q8MSRPCGJK`

## statement

Plan review round 3 left finding f1 (missing-acceptance on daemon-status): Its acceptance reads daemon status 'after rhiza daemon --once with a fake Rex (the state with the run it started)', but the leaf needs only daemon-lock. Picking and starting a run come from daemon-picks and daemon-run-tabs, so a builder scheduled as soon as daemon-lock lands has no daemon that starts a run, and can't run that part of the check. Either add those needs or drop the started run from the check. The critic confirmed it: daemon-status's acceptance reads the state 'with the run it started' after rhiza daemon --once, but the leaf needs only daemon-lock. Picking a change set and starting its run in a tab come from daemon-picks and daemon-run-tabs. A builder can take daemon-status as soon as daemon-lock lands, and at that point no daemon starts a run, so that part of the check can't run. Either add the needs or drop the started run from the check. How should it be settled?

## Links

### Out

- blocks [daemon-status](../capability/rhiza/orchestration/daemon/daemon-status.md)

### In

- answers [daemon-status-needs-picks](../decision/daemon-status-needs-picks.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
