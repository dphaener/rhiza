# headless-run-agent-tabs

capability · `01M6FESHESZGCR78N9DK7ZQJVY` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A testscript with a fake Rex drives a headless run: each agent session calls rex window new -s rhiza-<repo> with a title naming the run and the step, no rex split is anchored on a driver block, and closing a finished session closes only its own tab

## statement

Each agent session of a run (agent steps, checks, fixers and question sessions) opens its own Rex tab in the repo's session, titled for the run and step, instead of splitting a run window anchored on the driver's pane; every run is headless, so it replaces and retires the run-window tiling leaves run-sessions-one-window, agent-steps-tile-run-window, question-sessions-tile-run-window, checks-and-fixes-anchored, anchor-gone-splits-run-window and build-panes

## Links

### Out

- decided-by [m24-only-agent-tabs](../../../../decision/m24-only-agent-tabs.md)

### In

- verifies [headless-run-agent-tabs-verified](../../../../evidence/headless-run-agent-tabs-verified.md)
- verifies [headless-run-agent-tabs-qa-1](../../../../evidence/headless-run-agent-tabs-qa-1.md)
- verifies [headless-run-agent-tabs-reverified](../../../../evidence/headless-run-agent-tabs-reverified.md)
- verifies [headless-run-agent-tabs-qa-2](../../../../evidence/headless-run-agent-tabs-qa-2.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
