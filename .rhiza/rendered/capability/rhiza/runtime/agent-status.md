# agent-status

capability · `01M5SWT6Y6BCPQ46RHC2GGPM9J` · parent [runtime](../runtime.md)

## statement

Agent status comes from the agent's hooks, then Rex's events, then the process exiting, and is copied into Genesis's own files.

## Links

### Out

- decided-by [status-from-hooks-first](../../../decision/status-from-hooks-first.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
