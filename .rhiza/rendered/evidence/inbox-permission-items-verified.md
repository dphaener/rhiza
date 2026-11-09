# inbox-permission-items-verified

evidence · `01M6ZSS9Z1FK24N92E9480RX6A` · actor [claude-code-5f6760f8](../actor/claude-code-5f6760f8.md)

## result

PASS

## statement

Verifier at 56d15968: read the diff (inboxItems lists live runs' permission-blocked sessions via runPermissionAll over session.Resolve, including autonomous runs; itemTab now uses the session tab for kind permission; daemonNotify passes blocked/tabOf and inboxNotice adds the detail; webInboxKinds and the web labels include permission). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-permission|daemon-notify-hook|inbox-stalled|inbox-open|describe-golden)$' -count=1: all PASS. inbox-permission.txtar covers the permission item (session, run, step, detail, tab, next=rhiza inbox open) in an owner run and an autonomous run, no item after a newer Rex idle report, no question for an autonomous run's question block, and inbox open focusing block:sess in session:s2 rather than block:root. daemon-notify-hook writes one permission line with the detail in RHIZA_INBOX_TITLE. go test -trimpath ./internal/tui/inboxapp -run 'TestRowPerItemInOrder|TestOtherKindShowsDetailAndOpensTab': PASS, including the permission row with its detail and Enter calling inbox.open k2. go test -trimpath ./internal/api -run 'TestWebInbox|TestInbox': every test PASS, including TestWebInboxViewsPermission. The package then failed only on the root-state guard, because the live root daemon rewrote ~/code/.rhiza-root/usage.live.json; its mtime also changes with no test running, so this is unrelated to the change. go test ./internal/web: PASS.

## Links

### Out

- verifies [inbox-permission-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-permission-items.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
