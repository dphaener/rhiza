# inbox-app-signoff-verified

evidence · `01M6NATQ2BGC00QDK3DJ6Y7MNP` · actor [claude-code-b15f6e7a](../actor/claude-code-b15f6e7a.md)

## result

PASS

## statement

At cd07943e ran go test -trimpath -count=1 -v ./internal/tui/inboxapp/...: all 15 tests PASS. With the fake runner, TestSignoffDialogApprove shows the dialog naming the run (R2) and its wait, and a calls run.approve [R2] and quits; TestSignoffSendBackNote: s opens the note form, an empty note makes no call and does not quit, then 'fix the tests' calls run.send-back [R2, fix the tests] and quits; TestSignoffEscCallsNothing: esc quits with only the initial inbox load called. Diff read: openSignoff builds a ui.Dialog with the run and waits-for, approve/sendBack call the named ops and return tea.Quit, sendBack trims and refuses an empty note.

## Links

### Out

- verifies [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
