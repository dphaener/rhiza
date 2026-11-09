# permission-block-clears-verified

evidence · `01M6ZT4HHZ5Q1M2KJ0QKT3FXJJ` · actor [claude-code-a7dc5b6b](../actor/claude-code-a7dc5b6b.md)

## result

PASS

## statement

At 2ce01c1d: go test -trimpath ./internal/session -run TestResolve -v passes TestResolvePermissionBlockClears' 8 cases: hook blocked/permission at T with Rex idle or working at T+1s resolves from rex; Rex idle at T-1s or at T, a Rex blocked report, or no Rex report keeps the hook's blocked; hook blocked/question and blocked/auth with Rex idle at T+1s stay blocked from the hook. go test -trimpath ./cmd/rhiza -run TestScripts/session-message (GOTMPDIR outside the repo, -count=1, three runs) passes all four session-message scripts, including the new session-message-permission.txtar: a Notification permission_prompt hook, then session run following a fake rex events that reports block:b1 idle a minute later; session show reads idle from rex, and session message logs exactly one `send -b block:b1 <pointer>` and one `send-key -b block:b1 enter` (2 sends total). In a scratch worktree with internal/session/status.go reverted to before dfdc577e, that script fails, so it covers the fix. Diff read: permissionCleared applies only to a hook blocked/permission report overtaken by a strictly newer Rex idle/working report; exit still outranks everything.

## Links

### Out

- verifies [permission-block-clears](../capability/rhiza/runtime/agent-status/permission-block-clears.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
