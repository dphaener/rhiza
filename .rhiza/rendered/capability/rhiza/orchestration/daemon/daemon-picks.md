# daemon-picks

capability · `01M61R7K5W91SXWWX8AGDZ01HC` · parent [daemon](../daemon.md)

## acceptance

A Go test over a fixture model with two accepted milestone change sets with leaves, an accepted leafless milestone change set, an accepted backlog change set with leaves, one running run on a third accepted milestone change set, one interrupted run on a fourth and one failed run on a fifth picks only the first milestone change set under max 2, and none under max 1; with the fifth's failed run abandoned it picks the fifth again; a testscript runs rhiza daemon --once with a fake Rex and sees that one start.

## statement

Each poll, the daemon picks, in each project, the accepted change sets that rhiza ready says can start, that serve a milestone intent or, as daemon-picks-backlog sets out by the repo's [schedule] backlog rule, serve none, have at least one leaf, have no live run (running or interrupted, as run-one-live counts it) and whose latest run is not a failed run still in the inbox (as inbox-list keeps it: not abandoned, not retried, no newer run), in that project's roadmap milestone order (with the backlog placed and ordered as daemon-picks-backlog sets out, not by the roadmap's backlog group), and starts them, round-robin across projects as root-daemon-round-robin sets out, until its running runs (live runs whose process is running) reach the machine-wide maxRuns (default 2), so an interrupted run keeps its change set from being picked but takes no slot and a failed run's change set is never started again by the daemon until the owner clears that item; it never picks a draft, proposed or leafless change set on its own. Runs the owner queued (run-start-through-daemon) are not picks: they skip these filters and only share the maxRuns slots.

## Links

### Out

- depends-on daemon-lock (retired)
- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- decided-by [answer-8cd68bc7-3](../../../../decision/answer-8cd68bc7-3.md)
- decided-by [answer-8cd68bc7-5](../../../../decision/answer-8cd68bc7-5.md)
- decided-by [answer-8cd68bc7-13](../../../../decision/answer-8cd68bc7-13.md)
- decided-by [answer-8cd68bc7-11](../../../../decision/answer-8cd68bc7-11.md)

### In

- depends-on daemon-run-tabs (retired)
- depends-on [daemon-status](daemon-status.md)
- depends-on [inbox-notify](../inbox/inbox-notify.md)
- verifies [daemon-picks-verified](../../../../evidence/daemon-picks-verified.md)
- verifies [daemon-picks-qa-1](../../../../evidence/daemon-picks-qa-1.md)
- verifies [daemon-picks-verified-3](../../../../evidence/daemon-picks-verified-3.md)
- verifies [daemon-picks-qa-3](../../../../evidence/daemon-picks-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)
- verifies [daemon-picks-root-qa-1](../../../../evidence/daemon-picks-root-qa-1.md)
- verifies [daemon-picks-qa-2](../../../../evidence/daemon-picks-qa-2.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)
- verifies [daemon-picks-schedule-qa-1](../../../../evidence/daemon-picks-schedule-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
