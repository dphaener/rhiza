# rejoin-rex-unreachable

capability · `01M6M57BMKN39KDPQSDWK8D6M7` · parent [workflow-runs](../workflow-runs.md)

## acceptance

An engine test rejoins a live session while the fake Rex fails with a non-'not found' error: the run logs the failure, the session stays unjudged with no exit written, and once the fake Rex answers 'not found' on a later poll the session gets its 'lost' exit; with Rex answering the block exists, the session is polled as today.

## statement

When Rex cannot be reached on a rejoin's block check (any error that is not a 'not found' answer), the run logs it, judges nothing from it, keeps polling the session record as before, and retries the block check on later polls until Rex answers; Rex being down never fails or ends a step.

## Links

### Out

- depends-on [run-finds-sessions-ended-unseen](run-finds-sessions-ended-unseen.md)
- decided-by [judge-sessions-ended-unseen](../../../../decision/judge-sessions-ended-unseen.md)

### In

- verifies [rejoin-rex-unreachable-verified](../../../../evidence/rejoin-rex-unreachable-verified.md)
- verifies [rejoin-rex-unreachable-qa-1](../../../../evidence/rejoin-rex-unreachable-qa-1.md)
- verifies [pr-263-fix-a-run-resumed-after-a](../../../../evidence/pr-263-fix-a-run-resumed-after-a.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
