# status-precedence

capability · `01M5TYCYW9FJ42M36KXMA9BT99` · parent [agent-status](../agent-status.md)

## acceptance

A test feeds each source alone and in conflict and gets the highest-precedence status.

## statement

A session's status comes from Claude Code's hooks first, then the status in Rex's events, then the process exiting.

## Links

### Out

- decided-by [status-from-hooks-first](../../../../decision/status-from-hooks-first.md)
- decided-by [exit-is-final](../../../../decision/exit-is-final.md)

### In

- relates-to [stale-hook-status-after-exit](../../../../question/stale-hook-status-after-exit.md) — Status resolves hook, then Rex events, then exit, so a stale hook report outranks the exit.
- verifies [status-precedence-test](../../../../evidence/status-precedence-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
