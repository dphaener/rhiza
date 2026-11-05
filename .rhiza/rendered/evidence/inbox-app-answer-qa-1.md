# inbox-app-answer-qa-1

evidence · `01M6NCPVSBXE2AD8BPK8MEM50A` · actor [claude-code-f41a81b4](../actor/claude-code-f41a81b4.md)

## result

FAIL

## statement

QA at e4777542 with the built binary in tmux: Enter on a blocker opens a one-TextArea form; typing 'Use a file' and Enter wrote the answer (quote and statement 'Use a file') into the blocker's change set via question.answer, the blocker left rhiza inbox and the app exited 0. Enter on a run question opens the form; submitting 'Use sqlite' called session.message for session S1 (failed as expected with the fake: no session matches "S1", shown in the subtitle). Esc on the form quits with nothing sent. Empty Enter sends nothing, but after it (and after the failed send) the form is frozen: input gone, keys ignored, only Esc works.

## Links

### Out

- verifies [inbox-app-answer](../capability/rhiza/inbox-app-answer.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
