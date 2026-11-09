# message-says-why-waiting

capability · `01M6ZQFFBV7TQ3R919MJNZ6QCJ` · parent [step-messages](../step-messages.md)

## acceptance

go test -trimpath ./internal/api -run TestSessionMessage: a session blocked on a permission prompt gets one stderr line naming the kind, the detail and rhiza inbox open, a new line only when the status changes to working, and a timeout error naming the last status; the session-message testscripts still pass with stdout and --json unchanged

## statement

While rhiza session message waits for a session's prompt it prints to stderr why it waits, once at the start and again whenever the resolved status changes: for a blocked session the block's kind and detail and how to answer it (rhiza inbox open), for a working or starting one that it is working, each with the timeout; the timeout error names the last status, and the JSON output and the 10m default are unchanged

## Links

### Out

- decided-by [permission-block-inbox-and-clear](../../../../decision/permission-block-inbox-and-clear.md)

### In

- verifies [message-says-why-waiting-verified](../../../../evidence/message-says-why-waiting-verified.md)
- verifies [message-says-why-waiting-qa-1](../../../../evidence/message-says-why-waiting-qa-1.md)
- verifies [pr-418-fix-a-run-blocked-on-a](../../../../evidence/pr-418-fix-a-run-blocked-on-a.md)

## History

- `01M6ZQ588DTHT2T3MHK7STGR3Q` Fix A run blocked on a permission prompt does not reach the inbox, and rhiza session message hangs on it. Run 01M6ZMRGYP9BQFAC0QNG2J9KHH (pre-publish audit), builder session 411ae041 (bypassPermissions) hit a plugin hook confirmation (cc-plugin-privacy: the command held the owner's email) at 2026-11-09T07:55:57Z; the run log shows 'step build/build: blocked on permission' and session show reports 'blocked permission', but rhiza inbox said 'Nothing waits on you' for 36+ min. After I answered the prompt (No, via rex send-key), 'rhiza session message 411ae041 <text>' hung >2 min although the session was idle at its prompt; typing via rex block send worked. Expect: a blocked-permission session is an inbox item naming the prompt and how to answer it; session message delivers to an idle session or times out with an error.
