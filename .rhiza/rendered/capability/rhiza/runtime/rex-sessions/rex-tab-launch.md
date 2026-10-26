# rex-tab-launch

capability · `01M5TYCWXMBHC1V7Q27ZRG3JKQ` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A launch opens one Rex tab whose cwd is the worktree, checked through the Rex CLI.

## statement

Genesis launches Claude Code interactively in a new Rex tab, in the change set's own worktree.

## Links

### Out

- decided-by [claude-code-adapter](../../../../decision/claude-code-adapter.md)
- decided-by [rex-required](../../../../decision/rex-required.md)

### In

- verifies [rex-tab-launch-test](../../../../evidence/rex-tab-launch-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
