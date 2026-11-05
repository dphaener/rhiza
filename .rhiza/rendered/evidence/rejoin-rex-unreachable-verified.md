# rejoin-rex-unreachable-verified

evidence · `01M6M9VJ0X7MA2W8Q54V8SZ8Z6` · actor [claude-code-9274cce9](../actor/claude-code-9274cce9.md)

## result

PASS

## statement

At c8f16230 (diff: new engine test only, over run_lost.go from 7fd9c685) ran go test ./internal/run -run TestAgentRejoinRexUnreachable: a rejoin whose Rejoin fails with 'connection refused' is logged as 'rejoin session S-earlier: ... connection refused', rejoined once, polled 3 times, and the step fails only on the later poll's lost exit ('lost: Rex block not found') - PASS. Ran go test ./internal/api -run 'TestRejoin|TestPoll|TestLost': TestRejoinRexUnreachableRetried (fake rex errors 'connection refused': Rejoin returns the error, poll not done, no exit written; then fake rex 'block not found': next poll is Done+Lost with exit detail 'Rex block not found') and TestRejoinBlockExists (block exists: polled as before, no further rex calls) - all PASS. Read run_lost.go: non-'not found' errors keep blockPending and retryBlock re-asks on each poll, logging and judging nothing. make test passed (exit 0) on rerun; a first run failed only in the timing-budget phase with no test FAIL. make check's fmt/vet/test pass; self-check fails only on stale rendered actor files of this run, to be rendered at land.

## Links

### Out

- verifies [rejoin-rex-unreachable](../capability/rhiza/workflows/workflow-runs/rejoin-rex-unreachable.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
