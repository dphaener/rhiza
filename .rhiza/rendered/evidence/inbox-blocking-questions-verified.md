# inbox-blocking-questions-verified

evidence · `01M69VJS8PH2EC1GYKCXK86Y31` · actor [claude-code-a99f77d2](../actor/claude-code-a99f77d2.md)

## result

PASS

## statement

Verifier at dbc3266c: go test -count=1 ./cmd/rhiza -run TestScripts/(inbox-blocking-questions|describe-golden|inbox.*) passes. inbox-blocking-questions.txtar sets up an accepted change set touching store with an open model question store-q blocking it plus a run question (QUESTIONRUN) on it; rhiza inbox shows both (run question 'Which store?' and blocker 'Which-store (store-q)', 'blocks: store') and --json has kind blocker with question/handle/blocks plus kind question; after decision ssd answers store-q the blocker is gone from text and --json while the run question stays; cache-q (draft change set) and old-q (landed change set) are never listed, and the inbox is 'Nothing waits on you' before the accepted set exists. Diff read: inboxBlockers in internal/api/inbox.go reuses planCore blockers from the to-be view, skips drafts, one item per question and change set; describe-golden updated. go test ./... all ok, go vet and gofmt clean. make self-check reports only rendered actor-file drift caused by the run's uncommitted actor log entries (RHIZA_ROOT worktree), not by this leaf.

## Links

### Out

- verifies [inbox-blocking-questions](../capability/rhiza/orchestration/inbox/inbox-blocking-questions.md)

## History

- `01M686RSVTRCWEV9GPT2V16M9Y` The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares
