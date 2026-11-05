# inbox-app-signoff-qa-2

evidence · `01M6NDH0VCFBWFD2N6G9V899JW` · actor [claude-code-3293f9c2](../actor/claude-code-3293f9c2.md)

## result

PASS

## statement

QA at 433db495 with the built binary in tmux on a real review run stopped at signoff (testscript scratch repo, fake claude): Enter opens 'Sign off · Approve me' with the run ID, what it waits for and 'a approve · s send back · esc close'. 's' opens a note form; an empty Enter shows 'a send-back needs a note'; 'tighten the tests' + Enter journalled 'signoff (visit 1): changes' with finding 'alpha f1 major at owner: tighten the tests' and the app exited 0. On a fresh copy 'a' journalled 'approved by test', queued the run in continues.json, the signoff left rhiza inbox and the app exited 0.

## Links

### Out

- verifies [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
