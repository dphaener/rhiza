# run-send-back-verified

evidence · `01M6NC7DRW1YEH5QD62GEM81FC` · actor [claude-code-890f797f](../actor/claude-code-890f797f.md)

## result

PASS

## statement

Verifier at 4221f8bf: go test -trimpath ./cmd/rhiza -run TestScripts/(run-send-back|describe-golden|inbox-signoff|review-signoff.*|run-ship-signoff|tui-changes-app) passes, as do go test ./internal/api/... and ./internal/tui/inboxapp/... run-send-back.txtar drives a shiptest run with fake claude/gh/rex to signoff_waiting, refuses a blank note with exit 2 (status still at signoff, no changes asked), then send-back --note 'say more' exits 0: the fake claude log shows fixer cycle 2 then check cycle 2, status shows 'signoff (visit 1): changes ... 1 finding' and the run back at signoff, and land never ran (no order.log). Re-running the script with -testwork and rhiza run status on its repo showed the note as a finding ('alpha f1 major at owner: say more') under Sign-offs. After approve, send-back on the finished run exits 2 and land ran once only. Diff read: run.send-back op in internal/api/run_signoff.go journals an owner changes sign-off with the note as one finding and returns to the on.changes checks step, sharing the max_visits guard with request-changes; describe golden regenerated.

## Links

### Out

- verifies [run-send-back](../capability/rhiza/run-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
