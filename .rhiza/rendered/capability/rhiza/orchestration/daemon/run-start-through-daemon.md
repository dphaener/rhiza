# run-start-through-daemon

capability · `01M6FESFPN50CY4ACP8PEQPKEA` · parent [daemon](../daemon.md)

## acceptance

A testscript: rhiza run start <cs> queues the run and prints that it waits for the daemon; rhiza daemon --once then starts it headless, or leaves it queued while maxRuns runs are live; a run queued on an accepted backlog change set and a plan run queued on a draft change set are both started by rhiza daemon --once; rhiza run start --foreground drives the run in the calling process; a unit test on the roadmap app's fake runner shows its start key queueing the run

## statement

rhiza run start and the roadmap app's start keys queue the run for the daemon, which starts it headless within maxRuns, instead of driving it in the caller's tab; a queued run is the owner's start, not a daemon pick, so the daemon starts it on whatever change set rhiza run start accepts (backlog change sets and plan runs on drafts included) without daemon-picks' milestone, state and leaf filters; with no daemon running they say the run waits for one; rhiza run start --foreground drives the run in place, for debugging and as the daemon's own child

## Links

### Out

- decided-by [m24-only-agent-tabs](../../../../decision/m24-only-agent-tabs.md)

### In

- verifies [run-start-through-daemon-verified](../../../../evidence/run-start-through-daemon-verified.md)
- verifies [run-start-through-daemon-qa-1](../../../../evidence/run-start-through-daemon-qa-1.md)
- verifies [run-start-through-daemon-verified-2](../../../../evidence/run-start-through-daemon-verified-2.md)
- verifies [run-start-through-daemon-qa-2](../../../../evidence/run-start-through-daemon-qa-2.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
