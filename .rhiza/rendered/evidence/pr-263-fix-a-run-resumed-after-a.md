# pr-263-fix-a-run-resumed-after-a

evidence · `01M6MES7Z8RXB5S3ENH6QCYP8W` · actor [dphaener](../actor/dphaener.md)

## head

16d5791eb3cc4a3208bfd09bb0ae127b6d3aa9c5

## pr

https://github.com/dphaener/rhiza/pull/263

## result

1 check passed: check SUCCESS

## statement

Pull request #263 (https://github.com/dphaener/rhiza/pull/263) passed CI at 16d5791

## status

verified

## Links

### Out

- verifies [run-finds-sessions-ended-unseen](../capability/rhiza/workflows/workflow-runs/run-finds-sessions-ended-unseen.md)
- verifies [lost-session-judged](../capability/rhiza/workflows/workflow-runs/lost-session-judged.md)
- verifies [rejoin-rex-unreachable](../capability/rhiza/workflows/workflow-runs/rejoin-rex-unreachable.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
