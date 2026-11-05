# question-61e610c5-1

question · `01M6M53H0WNY6MB88RFDBAYRPE`

## session

61e610c5-678e-4359-9e22-01a1a84f5f74

## statement

Run 01M6K32A hung because its fixer session's record has a SessionEnd hook ('done') but no exit report (the wrapper died in the reboot) and its Rex block is gone. When should a run check an in-flight session for having ended unseen? / How is a session found ended unseen judged? / If Rex cannot be reached on rejoin (not a 'not found' answer), what happens?

## Links

### In

- answers [judge-sessions-ended-unseen](../decision/judge-sessions-ended-unseen.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
