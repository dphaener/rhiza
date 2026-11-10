# inbox-lists-gate-stopped-plan-run

capability · `01M71HSF7GVD9YRPGVAZ0F63N6` · parent [rhiza](../rhiza.md)

## acceptance

A Go test in internal/api builds the shape of plan-milestone run 01M71E6F15HMTME3KKF96NJC5R (owner autonomy, a gated entry naming proposed or draft change sets, then stopped at step gate with an 'open questions: <handle>; ...' detail, the question in a planned change set of the run worktree) under the condition the builder finds made the inbox say 'Nothing waits on you'; the test fails before the fix, and after it rhiza inbox lists one open item for the run whose questions carry the question's statement and handle, both for the project's inbox and for the merged inbox the root reads. go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/inbox-gate-questions' pass.

## statement

A plan run stopped at its gate on open questions is always listed by rhiza inbox as an open item, from the project's own inbox and from the root's merged inbox, until the run is resumed, abandoned or ended. rhiza inbox's human text and --json name the run, its change set and each open question by its full statement and handle; the Inbox app and the web UI show the same item with each open question's full text and name the run and change set by title, with no handles or IDs.

## Links

### Out

- derived-from [plan-inbox-gate-stop-missed](../../source/plan-inbox-gate-stop-missed.md)

### In

- verifies [inbox-lists-gate-stopped-plan-run-verified](../../evidence/inbox-lists-gate-stopped-plan-run-verified.md)
- verifies [inbox-lists-gate-stopped-plan-run-qa-1](../../evidence/inbox-lists-gate-stopped-plan-run-qa-1.md)
- verifies [pr-421-fix-rhiza-inbox-said-nothing](../../evidence/pr-421-fix-rhiza-inbox-said-nothing.md)

## History

- `01M71HJW9W123PXWWAGZTHGE6Z` Fix rhiza inbox said 'Nothing waits on you' while plan run 01M71E6F15HMTME3KKF96NJC5R (plan-milestone, owner autonomy) was stopped at step gate on open question review-01m71e6f15hmtme3kkf96njc5r-r3-f1 (left by record 01M71HBAFK92B6EQR0FGKN9H8M, plan review round 3). It was found only because rhiza run send-back refused with 'stopped at step gate: open questions'. Expected: the inbox lists the stopped run and its question, worded as what is asked. Unstuck by asking the owner, recording a Decision answering the question (record change set, commit 0ef646e in the run worktree), then rhiza run resume; the run went on to open PR #419. Daemon f7c85e91, 2026-11-10 ~01:35Z.
