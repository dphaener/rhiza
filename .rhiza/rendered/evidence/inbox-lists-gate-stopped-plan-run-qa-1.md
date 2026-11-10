# inbox-lists-gate-stopped-plan-run-qa-1

evidence · `01M71MVBA82VY3X3ZB5475M5GF` · actor [claude-code-88b6eb03](../actor/claude-code-88b6eb03.md)

## result

PASS

## statement

QA claude-code-88b6eb03 at 264c0bbf, by hand with the built binary (go build -trimpath) and the merge base's build, in a scratch repo outside ~/code with RHIZA_* unset. Setup copying run 01M71E6F15HMTME3KKF96NJC5R: a record change set 'Review round' adds question review-r3-f1 ('Should the pages show plain commands?'); a plan change set 'Milestone plan' adds an intent; a hand-written journal has started (plan-milestone, plan+milestone_plan), process, entered gate, gated naming the plan set, and stopped at gate with 'open questions: review-r3-f1; ...'. With the run's lock held via flock .git/rhiza/runs/.<id>.lock, the merge base printed 'Nothing waits on you' (the bug, reproduced). With the lock free it listed the run, but as 'question:  (review-r3-f1)' with no statement. The new binary, lock held or free, lists 'open run <id>, change set <id>, step gate' plus 'question: Should the pages show plain commands? (review-r3-f1)' and the resume command; --json has kind open, change_title 'Milestone plan' and questions [{handle, statement}]. After rhiza root init and project add, the root's merged inbox shows the same: the merge base says 'Nothing waits on you', the new binary lists the item with 'project: proj' and the question's statement. With a process entry appended after the stop (a resume) and the lock held, both the root and the project print 'Nothing waits on you'. go test -trimpath ./internal/api -run TestInbox passes, TestInboxListsGateStoppedPlanRun included (one earlier attempt tripped the root guard on .rhiza-root/usage.live.json, an environment effect). go test -trimpath ./cmd/rhiza -run TestScripts/inbox-gate-questions passes. Not exercised: the Inbox app and web UI (this diff does not touch them).

## Links

### Out

- verifies [inbox-lists-gate-stopped-plan-run](../capability/rhiza/inbox-lists-gate-stopped-plan-run.md)

## History

- `01M71HJW9W123PXWWAGZTHGE6Z` Fix rhiza inbox said 'Nothing waits on you' while plan run 01M71E6F15HMTME3KKF96NJC5R (plan-milestone, owner autonomy) was stopped at step gate on open question review-01m71e6f15hmtme3kkf96njc5r-r3-f1 (left by record 01M71HBAFK92B6EQR0FGKN9H8M, plan review round 3). It was found only because rhiza run send-back refused with 'stopped at step gate: open questions'. Expected: the inbox lists the stopped run and its question, worded as what is asked. Unstuck by asking the owner, recording a Decision answering the question (record change set, commit 0ef646e in the run worktree), then rhiza run resume; the run went on to open PR #419. Daemon f7c85e91, 2026-11-10 ~01:35Z.
