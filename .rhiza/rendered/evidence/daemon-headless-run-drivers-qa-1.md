# daemon-headless-run-drivers-qa-1

evidence · `01M6G7Q6Q9JSB1PR3Q55QKFY2C` · actor [claude-code-5fcf5cce](../actor/claude-code-5fcf5cce.md)

## result

PASS

## statement

At a5ee2b13 I built rhiza (go build) and ran it by hand in a scratch git repo with RHIZA_* and REX_* unset, a fake rex that logs its calls and runs the agent command, and a fake claude agent (held by a file). rhiza run start <cs> queued the run; rhiza daemon --permission-mode manual --once printed 'Started a run of <cs>' and exited. ps showed the driver as 'rhiza run start --foreground <cs> --permission-mode manual' with PPID 1 in its own session (setsid), cwd the repo. rex.log had only 'session ls' plus the agent's own window new, no window new for the driver. starts.json drained (no new queue entry). .git/rhiza/logs/<cs>.log held the child's output, and rhiza run ls --json showed the run state running, live true, after the daemon had exited. Under rhiza daemon --once --autonomous a picked milestone change set started as 'run start --foreground <cs> --permission-mode manual --autonomous'. A killed driver's run was later resumed by the daemon as a detached 'rhiza run resume --foreground <run>'.

## Links

### Out

- verifies [daemon-headless-run-drivers](../capability/rhiza/orchestration/daemon/daemon-headless-run-drivers.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
