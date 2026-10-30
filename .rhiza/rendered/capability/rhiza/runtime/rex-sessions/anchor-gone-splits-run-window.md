# anchor-gone-splits-run-window

capability · `01M651MBGX79H2AV7N9YA0EC9X` · parent [rex-sessions](../rex-sessions.md)

## acceptance

With a fake rex that reports the anchor block not found, the session runs rex split on another pane of the run's window and no rex window new; it runs rex window new only when rex block ls reports the window not found or lists no other pane. TestOpenSessionPane's 'anchor gone splits another pane', 'anchor gone, no other pane opens a new tab' and 'window gone opens a new tab' cases pass

## statement

A session whose anchor pane is gone splits a live pane of its run's window, and opens a window only when the run's window is gone or has no other pane

## Links

### Out

- decided-by [co-locate-run-panes](../../../../decision/co-locate-run-panes.md)

### In

- blocks [anchor-gone-policy](../../../../question/anchor-gone-policy.md)
- verifies [anchor-gone-splits-run-window-verified](../../../../evidence/anchor-gone-splits-run-window-verified.md)
- verifies [anchor-gone-splits-run-window-qa-1](../../../../evidence/anchor-gone-splits-run-window-qa-1.md)
- verifies [anchor-gone-splits-run-window-verified-2](../../../../evidence/anchor-gone-splits-run-window-verified-2.md)
- verifies [anchor-gone-splits-run-window-qa-2](../../../../evidence/anchor-gone-splits-run-window-qa-2.md)
- verifies [anchor-gone-splits-run-window-verified-3](../../../../evidence/anchor-gone-splits-run-window-verified-3.md)
- verifies [pr-117-tiled-run-sessions](../../../../evidence/pr-117-tiled-run-sessions.md)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
