# run-continue-through-daemon

capability · `01M6FFE7252R1S1V4AQTV2RWVN` · parent [daemon](../daemon.md)

## acceptance

A testscript: rhiza run resume <run>, rhiza run approve <run> and rhiza run request-changes - <run> each record the action, queue the continuation and print that it waits for the daemon; rhiza daemon --once then continues the run as a detached rhiza run resume --foreground child with no rex window new for it, and rhiza run ls shows it live after the daemon has exited; each with --foreground drives the run in the calling process; a unit test on the inbox app's fake runner shows its retry, approve and request-changes keys calling rhiza with no rex window new

## statement

rhiza run resume [--retry], rhiza run approve and rhiza run request-changes, and the inbox app's retry/resume, approve and request-changes keys, record the owner's action and queue the run's continuation for the daemon, which drives it headless within maxRuns as a detached child (as daemon-headless-run-drivers resumes a run), instead of continuing the run's driver in the caller's process or in a new Rex tab; like a queued start it is the owner's action, not a daemon pick; with no daemon running they say the run waits for one; --foreground on each drives the continuation in place, for debugging and as the daemon's own child

## Links

### Out

- decided-by [m24-only-agent-tabs](../../../../decision/m24-only-agent-tabs.md)

### In

- verifies [run-continue-through-daemon-verified](../../../../evidence/run-continue-through-daemon-verified.md)
- verifies [run-continue-through-daemon-qa-1](../../../../evidence/run-continue-through-daemon-qa-1.md)
- verifies [run-continue-through-daemon-verified-2](../../../../evidence/run-continue-through-daemon-verified-2.md)
- verifies [run-continue-through-daemon-qa-2](../../../../evidence/run-continue-through-daemon-qa-2.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
