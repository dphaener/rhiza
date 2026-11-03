# run-continue-through-daemon-qa-1

evidence · `01M6G7RSNS21PFYA7RSYW37RRN` · actor [claude-code-5fcf5cce](../actor/claude-code-5fcf5cce.md)

## result

PASS

## statement

At a5ee2b13 I built rhiza (go build) and ran it by hand in a scratch git repo with RHIZA_* and REX_* unset, a fake rex that logs its calls and runs the agent command, and a fake claude agent (held by a file). Three review runs waited at sign-off. rhiza run request-changes - <r1> (findings on stdin), rhiza run approve <r2> and rhiza run resume <r3> each journaled their action and printed 'Queued the continuation of run ... waits for one'. A second approve of r3 exited 2 with 'already has a continuation queued'. rhiza daemon --max-runs 1 --once printed 'Resumed run <r1>' and left two entries queued. The child was 'rhiza run resume --foreground <r1>' with PPID 1, and no window new carried 'run resume'. run ls showed r1 running and live after the daemon exited, and its log showed signoff changes leading to review visit 2 with a fixer. The next daemon --once continued r2 to succeeded, with continues.json drained. go test ./internal/tui/inboxapp passes (fake runner). I did not drive the TUI itself.

## Links

### Out

- verifies [run-continue-through-daemon](../capability/rhiza/orchestration/daemon/run-continue-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
