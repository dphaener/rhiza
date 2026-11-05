# inbox-app-pick-list-qa-1

evidence · `01M6NCPQC6EZYQ02NBCHC6P4EE` · actor [claude-code-f41a81b4](../actor/claude-code-f41a81b4.md)

## result

FAIL

## statement

QA at e4777542, built binary driven in tmux on a scratch repo (fake rex/claude/gh) holding real question, signoff, blocker, failed, accept and open items: rhiza tui app inbox lists one row per item in rhiza inbox --json order with kind, change set title, step and what it waits for; rhiza tui app inbox --once is byte-identical to rhiza inbox --json (cmp). Enter on question, blocker, signoff and failed opens their action. But Enter on an accept item and on an open (open questions) item does nothing at all in the binary: the Open hook is never set in cmd/rhiza/tui.go. Show kinds (stalled etc.) not exercised hands-on.

## Links

### Out

- verifies [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
