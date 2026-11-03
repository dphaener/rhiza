# daemon-headless-run-drivers

capability · `01M6FESDV933JY79RW9C3YJCCQ` · parent [daemon](../daemon.md)

## acceptance

A testscript runs rhiza daemon --once against a fake Rex: the run's driver starts as a rhiza run start --foreground process with the change set's ID (and --autonomous only under rhiza daemon --once --autonomous), with no rex window new for it and no new queue entry, a line the child itself prints (step build/build: entered) lands in the run's log file, and rhiza run ls shows the run live after the daemon has exited; testscripts for plan requests and for resume see the daemon start a request's plan run as a detached rhiza run start --foreground --workflow plan --new <intent> --request <request> process and an interrupted run's driver as a detached rhiza run resume --foreground <run> process

## statement

The daemon starts each run as a detached rhiza run start --foreground <change set> process, starts a request's plan run as a detached rhiza run start --foreground --workflow plan --new <intent> --request <request> process (with --milestone for a milestone request), and resumes or continues one as a detached rhiza run resume --foreground <run> process, in the repo and outside Rex, adding --autonomous only when the daemon was started with --autonomous, with output in the run's log file under .git/rhiza, so no Rex tab holds a run's driver, the child drives the run rather than queueing it again, and a run still outlives the daemon; the daemon reaps each such child when it exits, so a driver that dies shows as interrupted, not live; it replaces daemon-run-tabs, which it retires

## Links

### Out

- decided-by [m24-only-agent-tabs](../../../../decision/m24-only-agent-tabs.md)

### In

- blocks [review-01m6fdsp-r3-f1](../../../../question/review-01m6fdsp-r3-f1.md)
- verifies [daemon-headless-run-drivers-verified](../../../../evidence/daemon-headless-run-drivers-verified.md)
- verifies [daemon-headless-run-drivers-qa-1](../../../../evidence/daemon-headless-run-drivers-qa-1.md)
- verifies [daemon-headless-run-drivers-verified-2](../../../../evidence/daemon-headless-run-drivers-verified-2.md)
- verifies [daemon-headless-run-drivers-qa-2](../../../../evidence/daemon-headless-run-drivers-qa-2.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
