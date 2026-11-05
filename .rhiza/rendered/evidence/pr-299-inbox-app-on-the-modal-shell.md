# pr-299-inbox-app-on-the-modal-shell

evidence · `01M6NGP5HHYYY3Q6ZCTG1Z32HD` · actor [dphaener](../actor/dphaener.md)

## head

56149251c7230f4fe7b8c3bec36bedfb30a7c323

## pr

https://github.com/dphaener/rhiza/pull/299

## result

1 check passed: check SUCCESS

## statement

Pull request #299 (https://github.com/dphaener/rhiza/pull/299) passed CI at 5614925

## status

verified

## Links

### Out

- verifies [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)
- verifies [inbox-app-answer](../capability/rhiza/inbox-app-answer.md)
- verifies [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)
- verifies [inbox-app-retry-abandon](../capability/rhiza/inbox-app-retry-abandon.md)
- verifies [run-send-back](../capability/rhiza/run-send-back.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
