# permission-block-clears

capability · `01M6ZQFF7C7AZMBBHAW8RPC69G` · parent [agent-status](../agent-status.md)

## acceptance

go test -trimpath ./internal/session: Resolve on a hook blocked/permission report at T with a Rex idle report at T+1s gives idle from Rex, with a Rex report at T-1s or a Rex blocked report gives the hook's blocked, and a hook blocked/question at T with Rex idle at T+1s stays blocked; session-message.txtar (or a new session-message-permission.txtar) sends to a session with a hook blocked/permission report and a newer Rex idle report and the pointer is typed once, then Enter

## statement

A session blocked on a permission prompt stops resolving as blocked once a Rex report newer than that hook report says it is idle or working; question and auth blocks keep strict hook precedence, and with no newer report the block stands

## Links

### Out

- decided-by [permission-block-inbox-and-clear](../../../../decision/permission-block-inbox-and-clear.md)

### In

- verifies [permission-block-clears-verified](../../../../evidence/permission-block-clears-verified.md)
- verifies [permission-block-clears-qa-1](../../../../evidence/permission-block-clears-qa-1.md)
- verifies [pr-418-fix-a-run-blocked-on-a](../../../../evidence/pr-418-fix-a-run-blocked-on-a.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
