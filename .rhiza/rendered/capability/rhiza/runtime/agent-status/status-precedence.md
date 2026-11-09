# status-precedence

capability · `01M5TYCYW9FJ42M36KXMA9BT99` · parent [agent-status](../agent-status.md)

## acceptance

A test feeds each source alone and in conflict and gets the highest-precedence status.

## statement

A session's status comes from Claude Code's hooks first, then the status in Rex's events, then the process exiting; the one exception is a hook report blocked on a permission prompt, which a newer Rex report of idle or working outranks (permission-block-clears).

## Links

### Out

- decided-by [status-from-hooks-first](../../../../decision/status-from-hooks-first.md)
- decided-by [exit-is-final](../../../../decision/exit-is-final.md)
- decided-by [permission-block-inbox-and-clear](../../../../decision/permission-block-inbox-and-clear.md)

### In

- relates-to [stale-hook-status-after-exit](../../../../question/stale-hook-status-after-exit.md) — Status resolves hook, then Rex events, then exit, so a stale hook report outranks the exit.
- verifies [status-precedence-test](../../../../evidence/status-precedence-test.md)
- verifies [pr-418-fix-a-run-blocked-on-a](../../../../evidence/pr-418-fix-a-run-blocked-on-a.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
