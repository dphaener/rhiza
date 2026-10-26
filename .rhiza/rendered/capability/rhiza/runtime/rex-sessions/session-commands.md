# session-commands

capability · `01M5TYCXTZ9SNZZF3524Q9REGP` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A testscript launches against a fake Rex and reads both outputs.

## statement

`genesis session ls` and `genesis session show <id>` list the sessions Genesis launched with their change set, tab, worktree, permission mode and status.

## Links

### Out

- decided-by [m2-scope-in](../../../../decision/m2-scope-in.md)
- decided-by [rex-holds-no-state](../../../../decision/rex-holds-no-state.md)

### In

- verifies [session-commands-test](../../../../evidence/session-commands-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
