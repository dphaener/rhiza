# question-session-panes-close

capability · `01M6857MT46BQ7E27CZM8H7PFP` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A Go test drives an arbiter to a journalled verdict and an advisor to a journalled recommendation with a fake rex and asserts a block close of each session's pane after the journal entry; a failed or still-running arbiter's pane is not closed

## statement

An arbiter session's pane is closed once its verdict is journalled, and an advisor session's once its recommendation is, so the run window keeps only live panes

## Links

### Out

- decided-by [question-panes-close-after-verdict](../../../../decision/question-panes-close-after-verdict.md)

### In

- verifies [question-session-panes-close-verified](../../../../evidence/question-session-panes-close-verified.md)
- verifies [question-session-panes-close-qa-1](../../../../evidence/question-session-panes-close-qa-1.md)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
