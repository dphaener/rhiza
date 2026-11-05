# inbox-app-retry-abandon-verified

evidence · `01M6NB3XJ5S13EMJBY9TV4AGQW` · actor [claude-code-bfbdcb88](../actor/claude-code-bfbdcb88.md)

## result

PASS

## statement

At 97af3548, ran go test -trimpath -count=1 ./internal/tui/inboxapp/... (ok) and go vet (clean). With the fake Runner: TestResolveRetryCallsRunResume shows the Retry/Abandon Select on the failed item and Enter calls run.resume --retry R3 then quits; TestResolveAbandonConfirmed confirms with y and calls run.abandon R3 then quits; TestResolveAbandonNeedsConfirm submits the Confirm unconfirmed and no op beyond the inbox load is called, no quit; TestResolveEscCallsNothing cancels the form and nothing is called (rexuikit Form maps esc/ctrl+c to FormCancelledMsg, form.go:188); TestResolveOnlyForFailedKinds maps only failed/gave-up/skipped to ActionResolve, and choose() opens the resolve form only for ActionResolve. Diff read: args match RunResumeRequest (run arg 0, --retry flag) and RunAbandonRequest (run arg 0).

## Links

### Out

- verifies [inbox-app-retry-abandon](../capability/rhiza/inbox-app-retry-abandon.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
