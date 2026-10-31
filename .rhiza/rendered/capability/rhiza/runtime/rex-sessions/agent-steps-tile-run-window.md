# agent-steps-tile-run-window

capability · `01M6857JXB8B7WXBN5VHE7MV1K` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A new TestStartAgent in internal/api/run_checks_panes_test.go with paneRunner calls StartAgent for a plan step after a session with block:a1 and asserts one rex split of block:a1 and no window new; with an empty session store it asserts exactly one window new; no test runs real Rex

## statement

An agent step's session, such as a plan run's interview, review or revise, splits the pane of the run's latest session with a pane, and opens a window only when it is the run's first session with a pane

## Links

### Out

- decided-by [m6-overnight-batch](../../../../decision/m6-overnight-batch.md)

### In

- verifies [agent-steps-tile-run-window-verified](../../../../evidence/agent-steps-tile-run-window-verified.md)
- verifies [agent-steps-tile-run-window-qa-1](../../../../evidence/agent-steps-tile-run-window-qa-1.md)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
