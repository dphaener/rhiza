# explicit-permission-mode

capability · `01M5TYCY4T4VXA86X91CZAXTMT` · parent [agent-adapters](../agent-adapters.md)

## acceptance

A test shows the launch command carries `--permission-mode` and a launch without one exits non-zero.

## statement

The Claude Code adapter always passes a permission mode, and a launch with none is refused.

## Links

### Out

- decided-by [claude-code-adapter](../../../../decision/claude-code-adapter.md)
- decided-by [supervised-on-the-mac](../../../../decision/supervised-on-the-mac.md)

### In

- verifies [explicit-permission-mode-test](../../../../evidence/explicit-permission-mode-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
