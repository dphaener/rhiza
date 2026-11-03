# anchor-gone-splits-run-window-verified-3

evidence · `01M66KTKB10PR7HKBQJXW4PM4B` · actor [claude-code-d23ba6b4](../actor/claude-code-d23ba6b4.md)

## result

PASS

## statement

At ac1817d ran go test ./internal/api -run TestOpenSessionPane -v -count=1: PASS. The table rows 'anchor gone splits another pane' (fake rex: split of block:b1 fails 'not found', block ls lists b1 and b9) asserts the exact calls size, split b1, block ls -w window:w1, size b9, split b9 with no window new; 'anchor gone, no other pane opens a new tab' (ls lists only b1) and 'window gone opens a new tab' (ls fails 'window target window:w1 not found') each assert size, split, block ls, then window new. Read openSessionPane/livePane in internal/api/session.go: window new runs only when there is no anchor, or livePane reports the window not found or no other block. make check exited 0.

## Links

### Out

- verifies anchor-gone-splits-run-window (retired)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
