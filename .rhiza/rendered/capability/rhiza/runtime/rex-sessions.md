# rex-sessions

capability · `01M5SWT6QN894SDW5V7GJNYK3H` · parent [runtime](../runtime.md)

## statement

Agent steps run as sessions in Rex windows, interactive by default and headless when a step asks.

## Links

### In

- serves change set `01M6516H25WW66JF8N10Q83RJ4` (A run's sessions open as tiled panes in one Rex window per run, not a new window each, so no Rex session reaches its 64-window limit)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
