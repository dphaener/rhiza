# status-in-genesis-files

capability · `01M5TYCZ35E465M8SEAF8N9XNZ` · parent [agent-status](../agent-status.md)

## acceptance

`genesis session show` reports the last status with Rex stopped.

## statement

The status is copied into Genesis's own files, never read back from Rex as the record.

## Links

### Out

- decided-by [rex-holds-no-state](../../../../decision/rex-holds-no-state.md)

### In

- verifies [status-in-genesis-files-test](../../../../evidence/status-in-genesis-files-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
