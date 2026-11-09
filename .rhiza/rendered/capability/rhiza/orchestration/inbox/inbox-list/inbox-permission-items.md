# inbox-permission-items

capability · `01M6ZQFF9QA1CMC83FCNZAJDXN` · parent [inbox-list](../inbox-list.md)

## acceptance

A new inbox-permission.txtar: a live run journal whose current step's session record has a hook blocked/permission report with a detail message gives one item of kind permission in rhiza inbox --json (session, run, step, detail, tab, next naming rhiza inbox open) in both an owner-run and an autonomous run, none once a newer Rex idle report is recorded, and an autonomous run's question block still gives none; in the same script, as inbox-open.txtar does, rhiza inbox open <permission item> focuses the blocked session's own tab, not the run's journalled tab (itemTab in internal/api/inbox.go); daemon-notify-hook.txtar gains a permission-item fixture and its notify command gets one line of kind permission with the detail in its RHIZA_INBOX_* text; go test -trimpath ./internal/tui/inboxapp -run 'TestRowPerItemInOrder|TestOtherKindShowsDetailAndOpensTab' covers a permission row showing its detail and Enter opening its tab; go test -trimpath ./internal/api -run TestWebInbox shows the permission item with its owner-facing label (the webInboxKinds filter in internal/api/web_inbox.go); inbox-stalled.txtar still passes

## statement

rhiza inbox lists each session of a live run that resolves as blocked on a permission prompt as an item of kind permission as soon as it blocks, in autonomous runs too, naming the run, step, session, the prompt's message and its tab, with rhiza inbox open <item> as the next step to answer it there; the notify hook, the inbox app and the web UI show it like the other kinds, and autonomous runs still list no questions

## Links

### Out

- decided-by [permission-block-inbox-and-clear](../../../../../decision/permission-block-inbox-and-clear.md)

### In

- verifies [inbox-permission-items-verified](../../../../../evidence/inbox-permission-items-verified.md)
- verifies [inbox-permission-items-qa-1](../../../../../evidence/inbox-permission-items-qa-1.md)
- verifies [pr-418-fix-a-run-blocked-on-a](../../../../../evidence/pr-418-fix-a-run-blocked-on-a.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
