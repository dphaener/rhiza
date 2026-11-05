# rejoin-rex-unreachable-qa-1

evidence · `01M6MA2YGWT0KBK2FQYXZP3SZK` · actor [claude-code-3a61d2cb](../actor/claude-code-3a61d2cb.md)

## result

PASS

## statement

QA with the built binary: resumed a killed agent-step run while the rex wrapper failed block inspect with 'connect: connection refused' for 2s, then let the real Rex answer. The run logged 'rejoin session ...: ask Rex whether block ... exists: ... connection refused', wrote no exit, and retried block inspect on each 200ms poll (12 calls); once real Rex answered 'not found' the session got exit {state: lost, detail: 'Rex block not found'}. With Rex answering the block exists (real block inspect of a live pane), the session stayed unjudged and was polled until a later event ended it.

## Links

### Out

- verifies [rejoin-rex-unreachable](../capability/rhiza/workflows/workflow-runs/rejoin-rex-unreachable.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
