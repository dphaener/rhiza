# run-send-back

capability · `01M6N099H464KP7VPM9TVD8RAZ` · parent [rhiza](../rhiza.md)

## acceptance

A testscript (run-send-back.txtar) with a fake claude drives a ship-change run to signoff, runs rhiza run send-back with a note, and shows the note in rhiza run status as a finding, the run's next step a fixer, and no signoff journaled; send-back on a run not at signoff exits non-zero and changes nothing.

## statement

rhiza run send-back <run> --note <text> sends a run waiting at a signoff step back to a fixer: the note is journaled as the next review round's findings, the run returns to its review fixer step instead of signing off, and the continuation is queued for the daemon.

## Links

### Out

- decided-by [tui-signoff-approve-or-send-back](../../decision/tui-signoff-approve-or-send-back.md)

### In

- verifies [run-send-back-verified](../../evidence/run-send-back-verified.md)
- verifies [run-send-back-qa-1](../../evidence/run-send-back-qa-1.md)
- verifies [run-send-back-verified-2](../../evidence/run-send-back-verified-2.md)
- verifies [run-send-back-qa-2](../../evidence/run-send-back-qa-2.md)
- verifies [pr-299-inbox-app-on-the-modal-shell](../../evidence/pr-299-inbox-app-on-the-modal-shell.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
