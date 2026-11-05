# run-finds-sessions-ended-unseen

capability · `01M6M5728P8TFP8BH3N3B7WJ16` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript resumes a review run whose journalled change-set fixer session has a SessionEnd hook report and no exit report, and another whose fake rex answers that the fixer's block is not found: in each the session record gets a 'lost' exit naming the reason, the fixer is judged failed, a fresh fixer attempt starts and the run goes on instead of waiting; an engine test shows a session whose record gains the SessionEnd mark with no exit mid-run is judged lost on the next poll with no rex call.

## statement

A run finds an in-flight agent session that ended while no run watched it: when it rejoins a journalled session (a resume), before its first poll, it reads the session record and asks Rex whether the session's block still exists, and on every later poll it reads the record alone, except that a rejoin block check Rex did not answer (an error that is not a 'not found' answer) is retried on later polls until Rex answers; a record with the SessionEnd hook mark and no exit report, or a block Rex answers is not found, gets an exit of state 'lost' written into the session record with a detail naming why, so the run stops polling it. This holds for every session a run polls: agent steps, leaf builders, verifiers and fixers, checks and critics, change-set fixers and arbiters.

## Links

### Out

- decided-by [judge-sessions-ended-unseen](../../../../decision/judge-sessions-ended-unseen.md)

### In

- depends-on [lost-session-judged](lost-session-judged.md)
- depends-on [rejoin-rex-unreachable](rejoin-rex-unreachable.md)
- verifies [run-finds-sessions-ended-unseen-verified](../../../../evidence/run-finds-sessions-ended-unseen-verified.md)
- verifies [run-finds-sessions-ended-unseen-qa-1](../../../../evidence/run-finds-sessions-ended-unseen-qa-1.md)
- verifies [run-finds-sessions-ended-unseen-verified-2](../../../../evidence/run-finds-sessions-ended-unseen-verified-2.md)
- verifies [pr-263-fix-a-run-resumed-after-a](../../../../evidence/pr-263-fix-a-run-resumed-after-a.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
