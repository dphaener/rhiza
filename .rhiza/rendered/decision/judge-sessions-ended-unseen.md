# judge-sessions-ended-unseen

decision · `01M6M55170XRKXHGVJWQJ9CEAQ`

## reasoning

The change set's intent asks resume to re-read each in-flight session's recorded status and treat a missing Rex block as exited. In internal/api/run.go AgentStatus judges a session only once rec.Status.Exited(), so a record with a SessionEnd hook but no exit report (the wrapper died in the reboot) polls forever; engine.go agent() rejoins without any liveness check. Checking on rejoin covers the reboot case; also checking the SessionEnd-without-exit mark each poll is a cheap file read that catches a wrapper dying mid-run, while keeping Rex calls to the rejoin path avoids load on every poll. Writing a synthetic exit into the record keeps the existing exit-is-final rule the one source of truth instead of a second judging path. The exit code is unknown, so treating it as a pass would contradict 'pass for exit 0'; but a done-waiting step's until is evidence of the work itself, so it may still pass, mirroring exitedUntil. A transient Rex outage is not evidence of a dead session; judging it ended would fail healthy runs after any daemon/Rex hiccup, so only an explicit not-found counts and the check retries.

## statement

When: a run checks an in-flight session for having ended unseen once when it rejoins it (a resume), before its first poll, and on each later poll whose record shows the SessionEnd hook mark with no exit report. Only the rejoin check asks Rex whether the session's block still exists; the per-poll check reads the session record alone. A session is found ended unseen when its record has the SessionEnd mark but no exit report, or when Rex answers that its block is not found. How it is judged: the run writes an exit into the session record (state 'lost', detail naming why: SessionEnd without exit report, or Rex block not found), so exit-is-final holds and every reader agrees. As the exit code is unknown it is never an exit 0: the step fails, except that a step with an until that waits for done passes if the agent's done mark is on the record and the until condition holds now, as for an exit 0. The pane, if any, is closed as for any finished session. Rex unreachable: an error that is not a 'not found' answer is no evidence the session ended. The run logs it, does not judge the session, keeps polling the record as now, and retries the Rex block check on later polls until Rex answers; Rex being down never fails or ends a step.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-61e610c5-1](../question/question-61e610c5-1.md)
- derived-from [resume-ended-unseen-plan](../source/resume-ended-unseen-plan.md)

### In

- decided-by [run-finds-sessions-ended-unseen](../capability/rhiza/workflows/workflow-runs/run-finds-sessions-ended-unseen.md)
- decided-by [lost-session-judged](../capability/rhiza/workflows/workflow-runs/lost-session-judged.md)
- decided-by [rejoin-rex-unreachable](../capability/rhiza/workflows/workflow-runs/rejoin-rex-unreachable.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
