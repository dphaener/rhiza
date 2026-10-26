# known-session-id

capability · `01M5TYCYETKM4SSWXPGEJ4Z8SF` · parent [agent-adapters](../agent-adapters.md)

## acceptance

The recorded session ID equals the one on the launch command.

## statement

The adapter starts Claude Code with `--session-id` set to an ID Genesis generates and records.

## Links

### Out

- decided-by [claude-code-adapter](../../../../decision/claude-code-adapter.md)

### In

- verifies [known-session-id-test](../../../../evidence/known-session-id-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
