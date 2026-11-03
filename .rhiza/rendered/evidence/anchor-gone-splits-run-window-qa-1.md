# anchor-gone-splits-run-window-qa-1

evidence · `01M66J3FP98BK9RDKS939RFQ1N` · actor [claude-code-d1914dd4](../actor/claude-code-d1914dd4.md)

## result

FAIL

## statement

With real Rex: rex split -b <closed block> fails with 'block target ... not found' and rex block ls -s S -w W --json lists the remaining live block under blocks[].block_id, the shapes openSessionPane and livePane rely on. But a two-step workflow run with the built binary and fake claude failed at step two: step one's panes were closed at step end, Rex dropped the window, paneAnchor fell back to the closed first session's tab, and livePane failed with 'window target ... not found', so the run failed where it used to open a window. A review run whose agents exited on their own failed the same way at the critic.

## Links

### Out

- verifies anchor-gone-splits-run-window (retired)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
