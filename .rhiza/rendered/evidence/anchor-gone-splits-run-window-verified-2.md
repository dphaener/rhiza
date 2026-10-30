# anchor-gone-splits-run-window-verified-2

evidence · `01M66K3PGMD6HGHQ0THPB32D57` · actor [claude-code-d34ea5d3](../actor/claude-code-d34ea5d3.md)

## result

PASS

## statement

At 9cb629a (attempt 2), go test ./internal/api -run 'TestOpenSessionPane|TestStartCheck|TestStartFix' -count=1 -v: all PASS. TestOpenSessionPane's 'anchor gone splits another pane' case gives a fake rex whose split fails 'block block:b1 not found' and asserts calls size, split -b block:b1, block ls -w window:w1, size of block:b9, split -b block:b9, with no window new, landing in window:w1. New cases: no other pane or window gone (block ls 'not found') opens a new tab; other block ls errors fail. TestStartCheckLatestAnchor asserts paneAnchor uses the run's latest session's pane. Read the diff: livePane returns ok=false only for a gone window or no other pane, per decision anchor-gone-splits-other-pane as revised. go test ./... and gofmt pass; make self-check's only failures are render of run-local actor records outside the commit, not this change.

## Links

### Out

- verifies [anchor-gone-splits-run-window](../capability/rhiza/runtime/rex-sessions/anchor-gone-splits-run-window.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
