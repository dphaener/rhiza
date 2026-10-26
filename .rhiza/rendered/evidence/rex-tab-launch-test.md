# rex-tab-launch-test

evidence · `01M5WAZ89Y96XVX6ABH7T99PF0`

## statement

cmd/rhiza/testdata/script/session-start.txtar passes, and internal/api/session_test.go TestSessionStartCreatesWorktree passes; manual check on 2026-10-26: rhiza session start --permission-mode manual in a throwaway repo opened one Rex tab titled "throwaway: Manual check" whose process cwd, read with rex block call process, was the worktree throwaway.manual-check; claude ran with the recorded --session-id and --permission-mode manual, RHIZA_CHANGE, RHIZA_ACTOR and RHIZA_FOR were in its environment, and closing the tab recorded exit 129 in the session file

## Links

### Out

- verifies [rex-tab-launch](../capability/rhiza/runtime/rex-sessions/rex-tab-launch.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
