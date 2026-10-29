# run-process-record

capability · `01M61R7JKQRF2SBF071J10XKRM` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A Go test finds the process, and with REX_BLOCK set the window and block, in the journal of a started run and of a resumed one, and a reader reports a fixture's unfinished journal as live and interrupted when its pid is not running, as live and running when it is the test's own process, and a finished journal as not live.

## statement

A run journals its process (host and pid) and, when it runs in Rex, its Rex window and block, when it starts and each time it is resumed, so a reader tells which of the live runs (unfinished journals, as run-one-live counts them) are running and which are interrupted (an unfinished journal whose process is gone), and where each run's tab is.

## Links

### Out

- decided-by [answer-8cd68bc7-14](../../../../decision/answer-8cd68bc7-14.md)
- decided-by [answer-8cd68bc7-15](../../../../decision/answer-8cd68bc7-15.md)
- decided-by [journal-run-tab](../../../../decision/journal-run-tab.md)

### In

- depends-on [daemon-resume-once](../../orchestration/daemon/daemon-resume-once.md)
- depends-on [daemon-picks](../../orchestration/daemon/daemon-picks.md)
- depends-on [inbox-open](../../orchestration/inbox/inbox-open.md)
- verifies [run-process-record-verified](../../../../evidence/run-process-record-verified.md)
- verifies [run-process-record-qa-1](../../../../evidence/run-process-record-qa-1.md)
- verifies [run-process-record-qa-3](../../../../evidence/run-process-record-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
