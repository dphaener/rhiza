# runtime

capability · `01M5SWT6NGS1S2AS5PZN74E68P` · parent [rhiza](../rhiza.md)

## statement

Genesis runs agents: sessions in Rex, an adapter per agent, their status, and how far each machine is trusted.

## Links

### Out

- decided-by [rex-required](../../decision/rex-required.md)
- decided-by [no-tmux-fallback](../../decision/no-tmux-fallback.md)
- decided-by [rex-holds-no-state](../../decision/rex-holds-no-state.md)
- decided-by [built-in-go-adapters](../../decision/built-in-go-adapters.md)
- decided-by [supervised-on-the-mac](../../decision/supervised-on-the-mac.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
