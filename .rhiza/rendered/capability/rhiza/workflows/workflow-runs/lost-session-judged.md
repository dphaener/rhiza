# lost-session-judged

capability · `01M6M57611NNTDR13MV8N9Z7GT` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Engine tests rejoin a lost session on a plain agent step (the step fails, detail naming the lost reason), on an until_when = done step with the done mark and the until holding (passes), and with the done mark but the until not holding (fails); in each the Runner is asked to close the session's pane.

## statement

A session with a 'lost' exit is judged as an exit whose code is unknown, never as exit 0: its step or attempt fails, except that a step with an until that waits for done passes when the record has the agent's done mark and the until holds now, and its pane, if Rex still has one, is closed as any finished session's is.

## Links

### Out

- depends-on [run-finds-sessions-ended-unseen](run-finds-sessions-ended-unseen.md)
- decided-by [judge-sessions-ended-unseen](../../../../decision/judge-sessions-ended-unseen.md)

### In

- verifies [lost-session-judged-verified](../../../../evidence/lost-session-judged-verified.md)
- verifies [lost-session-judged-qa-1](../../../../evidence/lost-session-judged-qa-1.md)
- verifies [pr-263-fix-a-run-resumed-after-a](../../../../evidence/pr-263-fix-a-run-resumed-after-a.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
