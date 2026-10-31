# question-sessions-tile-run-window

capability · `01M6857KVT0NRCTAY6EZ2DH209` · parent [rex-sessions](../rex-sessions.md)

## acceptance

Go tests TestStartArbiterSplits and TestStartAdvisorSplits build their fixtures with a run session whose tab block is block:a1 and assert the fake rex gets split -b block:a1 and never window new; a testscript arbiter-panes-one-window.txtar runs an autonomous agent step whose fake claude asks two questions and its fake rex log matches one window new then only splits

## statement

An autonomous run's arbiter sessions and the advisor sessions opened on owner questions split a pane of the run's window, so a run with many questions opens no window for them

## Links

### Out

- decided-by [m6-overnight-batch](../../../../decision/m6-overnight-batch.md)

### In

- verifies [question-sessions-tile-run-window-verified](../../../../evidence/question-sessions-tile-run-window-verified.md)
- verifies [question-sessions-tile-run-window-qa-1](../../../../evidence/question-sessions-tile-run-window-qa-1.md)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
