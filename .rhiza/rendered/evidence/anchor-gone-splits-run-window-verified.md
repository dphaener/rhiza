# anchor-gone-splits-run-window-verified

evidence · `01M66H9N3YPQQVVTQM33NK5WMA` · actor [claude-code-bc1374d8](../actor/claude-code-bc1374d8.md)

## result

PASS

## statement

At 044b9e9 ran go test ./internal/api -run TestOpenSessionPane -v -count=1: PASS. Its 'anchor gone splits another pane' case has the fake rex fail the anchor split with 'block block:b1 not found' and asserts the exact rex calls [block call size, split -b block:b1, block ls -w window:w1 --json, block call size of b9, split -b block:b9], which has no window new, and the result is block:b2 in window:w1. The 'anchor gone, no other pane' case fails the session and also opens no window. Read the diff: openSessionPane now calls livePane (rex block ls) and splitPane on a not-found error, and no longer falls through to rexWindowArgs. Checked against the real rex: rex block ls -s <session> --json returns blocks[].block_id, the shape livePane parses. make check exited 0.

## Links

### Out

- verifies [anchor-gone-splits-run-window](../capability/rhiza/runtime/rex-sessions/anchor-gone-splits-run-window.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
