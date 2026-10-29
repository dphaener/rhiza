# inbox-open-verified

evidence · `01M62HT4981244CKB52DS5CNQR` · actor [claude-code-27536f53](../actor/claude-code-27536f53.md)

## result

PASS

## statement

Verified at fc55a56 (attempt 2). `go test ./cmd/rhiza -run TestScripts/inbox` passes: inbox-open.txtar, with a fake rex, opens a question item and sees `rex focus -s session:s1 block:q1`; opens a sign-off item from REX_SESSION=session:owner and sees `rex block inspect -s session:run --json block:r1` and `rex focus -s session:run block:r1` (the fake rex finds r1 only with -s, as the real one does); gets failed run FAILEDRUN001 named with `tab gone (Rex no longer has block:gone)`, tab_gone true, nothing focused; and `rhiza inbox open zzzzzzz` exits 2 with `no inbox item`. `go test ./internal/api ./internal/run` passes, including run_process_test checking the journalled process now records rex_session. Attempt 1's finding f1 is fixed: run.Process journals RexSession (set from ctx.RexSession when a block is set) and itemTab/runInboxOpen scope inspect and focus with -s. Checked against the real rex with a scratch build: a sign-off journal recording session Rhiza and this pane's block, opened from REX_SESSION=Clipboard, returned opened=true and focused the tab (an unscoped `rex block inspect` from Clipboard exits 3, the old failure); a non-existent block in that session printed `tab gone`; an unknown item exited 2. Diff read: change limited to process.go, run.go runProcess, inbox.go itemTab, the script and the unit test. make check: fmt-check, vet and all tests pass; self-check fails only on stale renders of actors claude-code-27536f53 and claude-code-c9fefacd, created after this commit by the verify sessions, not by the commit.

## Links

### Out

- verifies [inbox-open](../capability/rhiza/orchestration/inbox/inbox-open.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
