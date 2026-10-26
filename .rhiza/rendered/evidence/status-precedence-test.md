# status-precedence-test

evidence · `01M5WAZ9PS21FRW2373AP9WF5P` · actor [claude-code-dd482246](../actor/claude-code-dd482246.md)

## statement

internal/session/status_test.go TestResolvePrecedence passes (live reports: hook over Rex; once exited, the exit over hook and Rex); internal/api/session_test.go TestRecordRexExitIsFinal passes (a Rex working report, then Rex's child_exited removal, gives the exit, and later reports don't revive it) and TestSessionRun passes (a Rex working report then the agent's exit resolves to the exit, and a later hook is not recorded); cmd/rhiza/testdata/script/session-run.txtar passes

## Links

### Out

- verifies [status-precedence](../capability/rhiza/runtime/agent-status/status-precedence.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
