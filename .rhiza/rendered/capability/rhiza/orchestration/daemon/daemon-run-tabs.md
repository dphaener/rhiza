# daemon-run-tabs

capability · `01M61R7M41QVVEDQ5GWRXE7D7M` · parent [daemon](../daemon.md)

## acceptance

A testscript runs rhiza daemon --once with a fake Rex, which records one new tab whose command is rhiza run start with the change set's ID and no --autonomous; under rhiza daemon --once --autonomous the command has --autonomous.

## statement

The daemon starts each run as rhiza run start <change set> in a new Rex tab in the repo, adding --autonomous only when the daemon was started with --autonomous, so a run outlives the daemon.

## Links

### Out

- depends-on [daemon-picks](daemon-picks.md)
- decided-by [answer-8cd68bc7-14](../../../../decision/answer-8cd68bc7-14.md)

### In

- depends-on [daemon-resume-once](daemon-resume-once.md)
- depends-on [daemon-status](daemon-status.md)
- verifies [daemon-run-tabs-verified](../../../../evidence/daemon-run-tabs-verified.md)
- verifies [daemon-run-tabs-qa-1](../../../../evidence/daemon-run-tabs-qa-1.md)
- verifies [daemon-run-tabs-verified-2](../../../../evidence/daemon-run-tabs-verified-2.md)
- verifies [daemon-run-tabs-verified-3](../../../../evidence/daemon-run-tabs-verified-3.md)
- verifies [daemon-run-tabs-qa-3](../../../../evidence/daemon-run-tabs-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
