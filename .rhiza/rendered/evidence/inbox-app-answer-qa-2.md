# inbox-app-answer-qa-2

evidence · `01M6NDGV2GBZM4RYS5MR1HC4G4` · actor [claude-code-3293f9c2](../actor/claude-code-3293f9c2.md)

## result

PASS

## statement

QA at 433db495 with the built binary (rhiza tui app inbox) in tmux on scratch repos made by testscripts with fake rex/claude. Question item: Enter opens a one-TextArea form; an empty Enter shows 'an answer needs words' and the form stays live; typing 'Use sqlite' and Enter calls session.message and, with no session S1 in the fake, shows 'send failed: no session matches "S1"' in a fresh live form that accepts typing again; Esc then quits (exit 0) with nothing sent. Blocker item: empty Enter refused, then 'Use a file' + Enter wrote Decision answer-store-q (quote 'Use a file', answers store-q) into the blocker's change set via question.answer and the app exited 0. Cycle 1's frozen-form finding is fixed.

## Links

### Out

- verifies [inbox-app-answer](../capability/rhiza/inbox-app-answer.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
