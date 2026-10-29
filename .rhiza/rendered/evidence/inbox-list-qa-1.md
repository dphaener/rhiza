# inbox-list-qa-1

evidence · `01M62KEBYMR53PG4F5XBSYA2SP` · actor [claude-code-55549e75](../actor/claude-code-55549e75.md)

## result

PASS

## statement

QA with a build of the branch head in a scratch repo: with fixture journals of a question, sign-off, failed, daemon-gave-up and done run, rhiza inbox listed exactly the first four oldest first with short IDs, in text and --json (empty repo: 'Nothing waits on you', items []). Unblocking the question, rhiza run abandon on the gave-up run, a signoff approved entry and a newer finished run on the failed run's change set emptied the inbox; the abandoned run was not listed as failed.

## Links

### Out

- verifies [inbox-list](../capability/rhiza/orchestration/inbox/inbox-list.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
