# inbox-app-retry-abandon-qa-2

evidence · `01M6NDHBHHBPCGPHXWGK21BSPF` · actor [claude-code-3293f9c2](../actor/claude-code-3293f9c2.md)

## result

PASS

## statement

QA at 433db495 with the built binary in tmux: on a failed run Enter opens 'What to do with the run' with a Retry/Abandon Select; Abandon then the Confirm left at its default No re-showed the choice ('abandon needs the confirm') with the journal unchanged; Esc quit with nothing done; Abandon then Yes appended finished 'abandoned by test' and exited 0. On a real failed review run (fake claude exiting 1) Retry queued it in continues.json with retry:true and exited 0. On a skipped run (change set landed) the same Select shows and Retry surfaces run.resume's refusal ('is not resumed: its change set ... is landed') in a live form. Neither choice shows on question, signoff, blocker, accept or open items.

## Links

### Out

- verifies [inbox-app-retry-abandon](../capability/rhiza/inbox-app-retry-abandon.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
