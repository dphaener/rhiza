# resume-ended-unseen-plan

source · `01M6M568QNCQ3EJAG6XH8HJ73Z`

## quotes

```json
[]
```

## raw_sha256

ea01bccfdd7706cce221ad35957d6fe2c7274c11314751fdc678f9ed7b7f413d

## session

61e610c5-678e-4359-9e22-01a1a84f5f74

## summary

Autonomous plan of change set 01M6M4YH; no owner present, so no owner words are quoted. Investigation: run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 after the VM reboot and logged 'continuing visit 1', then nothing. Its review fixer session b0eafd65 had a SessionEnd hook report (hook state 'done', 18:27:25) but no exit report and no done mark: the session wrapper that writes the exit died in the reboot, and its Rex block is gone. stepRunner.AgentStatus (internal/api/run.go) judges a session only once its record has an exit report, and the engine rejoins a journalled session (engine.go agent, checks.go pollReviewFix, leaves.go pollLeaf) with no liveness check, so it polls forever. The planner raised one question (when to check, how to judge, what if Rex is unreachable); the arbiter answered with Decision judge-sessions-ended-unseen: check on rejoin (record and Rex block) and on each poll (record's SessionEnd mark with no exit only); write a 'lost' exit into the record naming why; never treat it as exit 0, except a done-waiting step with the done mark and its until holding passes; close the pane; a Rex error that is not 'not found' judges nothing and the block check retries. The change set serves m4-runs-hold-up, as the earlier landed-run resume fix did. Planned as three leaves under workflow-runs: detection, judgement, and Rex unreachable.

## title

Plan interview: a resumed run waits forever on a session that ended unseen

## Links

### In

- derived-from [judge-sessions-ended-unseen](../decision/judge-sessions-ended-unseen.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
