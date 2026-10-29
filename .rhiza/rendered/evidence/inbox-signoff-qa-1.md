# inbox-signoff-qa-1

evidence · `01M62KECD05TRNT0PMMZD6KC6S` · actor [claude-code-55549e75](../actor/claude-code-55549e75.md)

## result

PASS

## statement

QA with a build of the branch head: a real review run (fake agent, fake rex) waited at sign-off; rhiza inbox approve --dry-run journalled nothing, rhiza inbox approve <item> ended the run succeeded with 'approved by darin-haener' and an empty inbox. A second run sent back with rhiza inbox request-changes <item> <file> started fix cycle 2 (fixer saw the finding, verifier wrote Evidence) and, once resumed, waited at sign-off again under a new item. On a question item both commands exited 2 naming it a question, and the journal was unchanged; an unknown item exited 2.

## Links

### Out

- verifies [inbox-signoff](../capability/rhiza/orchestration/inbox/inbox-signoff.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
