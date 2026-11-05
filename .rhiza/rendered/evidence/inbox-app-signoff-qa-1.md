# inbox-app-signoff-qa-1

evidence · `01M6NCQ0D7BGT7FSKAXZTJ0339` · actor [claude-code-f41a81b4](../actor/claude-code-f41a81b4.md)

## result

PASS

## statement

QA at e4777542 with the built binary in tmux on a real shiptest run stopped at signoff: Enter opens a Dialog 'Sign off · Sign it' with the run ID and what it waits for, hint 'a approve · s send back · esc close'. Esc closed it with the journal unchanged. 's' then empty Enter re-asked ('a send-back needs a note'); 'say more' + Enter journalled a signoff outcome changes with finding f1 alpha major at owner 'say more', queued the run in continues.json and the app exited 0. On a fresh copy 'a' journalled signoff outcome approved by the owner, queued the continuation and exited 0.

## Links

### Out

- verifies [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)

## History

- `01M6MYN9GRWHRJSHA31D9S838H` Inbox app: pick an inbox item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close; rebuilt on the modal shell over the inbox, question.answer, run.approve, run.resume and run.abandon operations, replacing internal/tui/inboxapp and retiring the capabilities of the old inbox app it replaces
