# checks-and-fixes-anchored-qa-2

evidence · `01M66KA2MV4HH6HR749DB1JTA2` · actor [claude-code-cec4c5b5](../actor/claude-code-cec4c5b5.md)

## result

PASS

## statement

At 68bf3aa, with the built binary, real Rex and a fake claude, a review run (checks alpha, beta; critic, fixer, verifier) where beta's agent stays open: the critic, the leaf fixer and verifier and the cycle-2 alpha check each ran rex split against a block of the run's one window (the fixer after its anchor's pane closed, via rex block ls, on beta's pane), and none ran rex window new; each fake agent logged a REX_BLOCK in that window. The run succeeded.

## Links

### Out

- verifies [checks-and-fixes-anchored](../capability/rhiza/runtime/rex-sessions/checks-and-fixes-anchored.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
