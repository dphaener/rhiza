# inbox-signoff-verified

evidence · `01M62HNA1CZR4TYK2BT92BM09P` · actor [claude-code-11024481](../actor/claude-code-11024481.md)

## result

PASS

## statement

At commit 441f9fb I ran go test ./cmd/rhiza -run 'TestScripts/(inbox-signoff|describe-golden)$' -count=1 -v and both passed. inbox-signoff.txtar starts a review run that waits at sign-off and approves it through its inbox item: the run ends succeeded with outcome approved, and the inbox shows Nothing waits on you. A second run is sent back with one finding from a file: the outcome is changes, the fixer is spawned with the finding, and the run waits at sign-off again. On a question item, approve and request-changes each exit 2 and leave the run journal unchanged. An unknown item also exits 2, and a --dry-run approve journals nothing. Regenerating describe-golden with -update left git status clean, so the committed golden matches, and it lists inbox.approve and inbox.request-changes. make check exited 0. In the diff in internal/api/inbox.go, signoffItemRun rejects any item that is not a sign-off with a usage error (exit 2), and otherwise hands the item's run to runRunApprove or runRunRequestChanges. findInboxItem is factored out of inbox open with no change in behavior.

## Links

### Out

- verifies [inbox-signoff](../capability/rhiza/orchestration/inbox/inbox-signoff.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
