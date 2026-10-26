# stale-hook-status-after-exit

question · `01M5WAZAAJYHCRX79TPRZVV393`

## statement

When the agent exits while its last hook report still says working or blocked, should the exit override the hook status?

## Links

### Out

- relates-to [status-precedence](../capability/rhiza/runtime/agent-status/status-precedence.md) — Status resolves hook, then Rex events, then exit, so a stale hook report outranks the exit.

### In

- answers [exit-is-final](../decision/exit-is-final.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
