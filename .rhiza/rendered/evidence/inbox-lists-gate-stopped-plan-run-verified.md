# inbox-lists-gate-stopped-plan-run-verified

evidence · `01M71MJ7VRHFM5SE0QZFFH70V1` · actor [claude-code-062a93b9](../actor/claude-code-062a93b9.md)

## result

PASS

## statement

Verifier claude-code-062a93b9 at commit 53605755. I read the diff. inboxItems now lists a run stopped on run.OpenStop even while its lock is held, unless a process entry was journalled after the stop (a resume). inboxGateRead also looks up a gate question in the change set that created it when that set is not one of the plan's: in run 01M71E6F15HMTME3KKF96NJC5R the question review-01m71e6f15hmtme3kkf96njc5r-r3-f1 was created in record change set 01M71HBAFK92B6EQR0FGKN9H8M. The real journal (.git/rhiza/runs/01M71E6F15HMTME3KKF96NJC5R.jsonl) matches the test's shape: the plan-milestone run has owner autonomy, a gated entry naming the plan's change sets, a stop at step gate with detail 'open questions: <handle>; ...', and the resume's process entry 2.5 minutes later. I ran TestInboxListsGateStoppedPlanRun in a scratch worktree with HEAD~1's internal/api/inbox.go and it FAILED ('a stopped run whose lock is held lists 0 open items, want 1'). At HEAD, go test -trimpath ./internal/api/... passed (ok, 425s) and go test -trimpath ./cmd/rhiza -run 'TestScripts/inbox-gate-questions' passed, both run from a scratch worktree outside ~/code with GOTMPDIR overridden. Run inside the run checkout, the script failed and the root guard tripped on .rhiza-root/usage.live.json; both are environment effects, not this change. The test asserts one open item while the lock is held, none once a resume's process entry follows the stop, and the question's handle and statement through inboxGateQuestions. The root's merged inbox (mergeInbox in runInbox) reads each project through the same inboxOne, so the same item and questions reach it. Deviations: the test puts the question in a record change set, as in the real run, where the acceptance says a planned change set. It does not call the merged-inbox op directly.

## Links

### Out

- verifies [inbox-lists-gate-stopped-plan-run](../capability/rhiza/inbox-lists-gate-stopped-plan-run.md)

## History

- `01M71HJW9W123PXWWAGZTHGE6Z` Fix rhiza inbox said 'Nothing waits on you' while plan run 01M71E6F15HMTME3KKF96NJC5R (plan-milestone, owner autonomy) was stopped at step gate on open question review-01m71e6f15hmtme3kkf96njc5r-r3-f1 (left by record 01M71HBAFK92B6EQR0FGKN9H8M, plan review round 3). It was found only because rhiza run send-back refused with 'stopped at step gate: open questions'. Expected: the inbox lists the stopped run and its question, worded as what is asked. Unstuck by asking the owner, recording a Decision answering the question (record change set, commit 0ef646e in the run worktree), then rhiza run resume; the run went on to open PR #419. Daemon f7c85e91, 2026-11-10 ~01:35Z.
