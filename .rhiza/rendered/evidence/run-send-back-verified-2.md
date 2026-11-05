# run-send-back-verified-2

evidence · `01M6ND2PFSNC52GZEQK75V9E19` · actor [claude-code-72088211](../actor/claude-code-72088211.md)

## result

PASS

## statement

At 3cfa418a ran go test -trimpath ./cmd/rhiza -run 'TestScripts/run-send-back$' -count=1: PASS. The script drives a shiptest run with fake claude/rex/gh to signoff_waiting, refuses a blank note (exit 2, status unchanged), then run send-back --note 'say more' exits 0; claude.log shows the fixer (cycle 2) and check running next; run status shows the changes round with '1 finding' and now asserts the line '  - alpha f1 major at owner: say more' (fixing finding f2), current step back at signoff, and no order.log (nothing landed/approved). After approve, send-back exits 2 and only the approve's landing is recorded; fold-ok passes. Read the fix diff: one added assertion line.

## Links

### Out

- verifies [run-send-back](../capability/rhiza/run-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
