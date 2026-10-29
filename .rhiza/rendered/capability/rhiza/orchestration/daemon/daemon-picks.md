# daemon-picks

capability · `01M61R7K5W91SXWWX8AGDZ01HC` · parent [daemon](../daemon.md)

## acceptance

A Go test over a fixture model with two accepted milestone change sets with leaves, an accepted leafless milestone change set, an accepted backlog change set with leaves, one running run on a third accepted milestone change set, one interrupted run on a fourth and one failed run on a fifth picks only the first milestone change set under max 2, and none under max 1; with the fifth's failed run abandoned it picks the fifth again; a testscript runs rhiza daemon --once with a fake Rex and sees that one start.

## statement

Each poll, the daemon picks the accepted change sets that rhiza ready says can start, that serve a milestone intent, have at least one leaf, have no live run (running or interrupted, as run-one-live counts it) and whose latest run is not a failed run still in the inbox (as inbox-list keeps it: not abandoned, not retried, no newer run), in roadmap order, and starts them until its running runs (live runs whose process is running) reach --max-runs (default 2), so an interrupted run keeps its change set from being picked but takes no slot and a failed run's change set is never started again by the daemon until the owner clears that item; it never starts a backlog, draft, proposed or leafless change set.

## Links

### Out

- depends-on [daemon-lock](daemon-lock.md)
- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- decided-by [answer-8cd68bc7-3](../../../../decision/answer-8cd68bc7-3.md)
- decided-by [answer-8cd68bc7-5](../../../../decision/answer-8cd68bc7-5.md)
- decided-by [answer-8cd68bc7-13](../../../../decision/answer-8cd68bc7-13.md)
- decided-by [answer-8cd68bc7-11](../../../../decision/answer-8cd68bc7-11.md)

### In

- depends-on [daemon-run-tabs](daemon-run-tabs.md)
- depends-on [daemon-status](daemon-status.md)
- depends-on [inbox-notify](../inbox/inbox-notify.md)
- verifies [daemon-picks-verified](../../../../evidence/daemon-picks-verified.md)
- verifies [daemon-picks-qa-1](../../../../evidence/daemon-picks-qa-1.md)
- verifies [daemon-picks-verified-3](../../../../evidence/daemon-picks-verified-3.md)
- verifies [daemon-picks-qa-3](../../../../evidence/daemon-picks-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
