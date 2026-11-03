# run-sessions-one-window-qa-2

evidence · `01M66KA3EQM3SJWH4RA864X2T7` · actor [claude-code-cec4c5b5](../actor/claude-code-cec4c5b5.md)

## result

PASS

## statement

At 68bf3aa, with the built binary, real Rex and a fake claude, a review run with a session left open logged exactly one rex window new and five rex splits, all in that window, then rex block close of every pane at the end; rex window ls afterwards showed no leftover run window. The run succeeded.

## Links

### Out

- verifies run-sessions-one-window (retired)

## History

- `01M6516H25WW66JF8N10Q83RJ4` A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit
