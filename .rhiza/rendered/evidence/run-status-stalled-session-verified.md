# run-status-stalled-session-verified

evidence · `01M6ASQ93QAB5V7TBF3BP9N9T1` · actor [claude-code-fb7e7161](../actor/claude-code-fb7e7161.md)

## result

PASS

## statement

At 5b2468c7 I ran go test ./internal/session -run TestStalled -v: it passes, with working at 19m and 20m not stalled, 21m stalled, and blocked, idle, usage-paused ('paused at its usage limit until ...'), exited and no-report sessions never stalled at 3h. go test ./cmd/rhiza -run 'TestScripts/(run-status-stalled|describe-golden)$' -count=3 passes. run-status-stalled.txtar hand-writes the session record's hook report at 2030-01-01T00:00:00Z: with RHIZA_TEST_NOW at 00:19 run status shows no 'Stalled' (text or --json), at 00:21 it prints '- Stalled: build session <id> has not reported since 2030-01-01T00:00:00Z' and a stalled object in --json, and an idle report shows nothing. Read the diff: session.Stalled requires the hook source, state working, a detail without the LimitPausedDetail prefix and now-At > 20m, and runInfo sets it only for an unended run's current-step session. make check: fmt, vet and all tests pass. Its only failure is a stale render of actor records written after this commit (verifier actors), which the leaf did not cause.

## Links

### Out

- verifies [run-status-stalled-session](../capability/rhiza/workflows/workflow-runs/run-status/run-status-stalled-session.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
