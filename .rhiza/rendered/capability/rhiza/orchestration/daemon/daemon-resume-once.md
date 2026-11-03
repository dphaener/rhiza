# daemon-resume-once

capability · `01M61R7MP3HGEB60GJV4M3CWQA` · parent [daemon](../daemon.md)

## acceptance

A testscript with a fixture journal that is unfinished and names a dead process runs rhiza daemon --once twice with a fake Rex: the first starts one detached rhiza run resume --foreground on that run and journals it, the second (with the run interrupted again) starts none; a failed run's journal gets no resume.

## statement

The daemon resumes an interrupted run whose change set is accepted once, as a detached rhiza run resume --foreground <run> process outside Rex (daemon-headless-run-drivers), and journals that resume in the run's journal; it leaves alone an interrupted run whose change set is not accepted, a run interrupted again after that resume, which the inbox lists as a run the daemon gave up on, and every failed run, which the inbox lists as failed.

## Links

### Out

- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- depends-on daemon-run-tabs (retired)
- decided-by [answer-8cd68bc7-15](../../../../decision/answer-8cd68bc7-15.md)

### In

- depends-on [inbox-list](../inbox/inbox-list.md)
- depends-on [daemon-resumes-accepted-only](daemon-resumes-accepted-only.md)
- verifies [daemon-resume-once-verified](../../../../evidence/daemon-resume-once-verified.md)
- verifies [daemon-resume-once-qa-1](../../../../evidence/daemon-resume-once-qa-1.md)
- verifies [daemon-resume-once-verified-r1](../../../../evidence/daemon-resume-once-verified-r1.md)
- verifies [daemon-resume-once-qa-3](../../../../evidence/daemon-resume-once-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
