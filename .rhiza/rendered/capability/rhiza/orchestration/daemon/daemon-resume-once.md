# daemon-resume-once

capability · `01M61R7MP3HGEB60GJV4M3CWQA` · parent [daemon](../daemon.md)

## acceptance

A testscript with a fixture journal that is unfinished and names a dead process runs rhiza daemon --once twice with a fake Rex: the first opens one tab running rhiza run resume on that run and journals it, the second (with the run interrupted again) opens none; a failed run's journal gets no tab.

## statement

The daemon resumes an interrupted run once, as rhiza run resume <run> in a new Rex tab, and journals that resume in the run's journal; it leaves alone a run interrupted again after that resume, which the inbox lists as a run the daemon gave up on, and every failed run, which the inbox lists as failed.

## Links

### Out

- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- depends-on [daemon-run-tabs](daemon-run-tabs.md)
- decided-by [answer-8cd68bc7-15](../../../../decision/answer-8cd68bc7-15.md)

### In

- depends-on [inbox-list](../inbox/inbox-list.md)
- verifies [daemon-resume-once-verified](../../../../evidence/daemon-resume-once-verified.md)
- verifies [daemon-resume-once-qa-1](../../../../evidence/daemon-resume-once-qa-1.md)
- verifies [daemon-resume-once-verified-r1](../../../../evidence/daemon-resume-once-verified-r1.md)
- verifies [daemon-resume-once-qa-3](../../../../evidence/daemon-resume-once-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
