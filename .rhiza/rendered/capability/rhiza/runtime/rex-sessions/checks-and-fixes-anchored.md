# checks-and-fixes-anchored

capability · `01M651ME1N3KVWZ1D6M8BYFHBQ` · parent [rex-sessions](../rex-sessions.md)

## acceptance

TestStartCheck and TestStartFix with an Anchor record split the anchored session's block and open no window; both pass

## statement

Review, critic and fixer sessions split the pane of the run's latest session with a pane, so every one sits in the run's window

## Links

### Out

- decided-by [co-locate-run-panes](../../../../decision/co-locate-run-panes.md)

### In

- verifies [checks-and-fixes-anchored-verified](../../../../evidence/checks-and-fixes-anchored-verified.md)
- verifies [checks-and-fixes-anchored-qa-1](../../../../evidence/checks-and-fixes-anchored-qa-1.md)
- verifies [checks-and-fixes-anchored-qa-2](../../../../evidence/checks-and-fixes-anchored-qa-2.md)
- verifies [checks-and-fixes-anchored-verified-3](../../../../evidence/checks-and-fixes-anchored-verified-3.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
