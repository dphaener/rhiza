# run-sessions-one-window

capability · `01M651M8TNJKY76E8G6KMV3V3Y` · parent [rex-sessions](../rex-sessions.md)

## acceptance

With a fake rex, a run that starts a review, a critic and a fixer session runs one rex window new and splits the rest against that window's block; a testscript run-panes-one-window passes

## statement

A run's agent sessions open as tiled panes of the run's own Rex window, not a window of their own

## Links

### Out

- decided-by [co-locate-run-panes](../../../../decision/co-locate-run-panes.md)

### In

- blocks [run-window-owner](../../../../question/run-window-owner.md)
- verifies [run-sessions-one-window-verified](../../../../evidence/run-sessions-one-window-verified.md)
- verifies [run-sessions-one-window-qa-1](../../../../evidence/run-sessions-one-window-qa-1.md)
- verifies [run-sessions-one-window-verified-2](../../../../evidence/run-sessions-one-window-verified-2.md)
- verifies [run-sessions-one-window-qa-2](../../../../evidence/run-sessions-one-window-qa-2.md)
- verifies [pr-117-tiled-run-sessions](../../../../evidence/pr-117-tiled-run-sessions.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
