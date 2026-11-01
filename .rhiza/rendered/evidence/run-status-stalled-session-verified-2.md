# run-status-stalled-session-verified-2

evidence · `01M6AVC8R6HYRTR1BWTDDV9633` · actor [claude-code-6a792fd4](../actor/claude-code-6a792fd4.md)

## result

PASS

## statement

Verified at 50ea91b0 (attempt 2, fix of review findings quality f1 and qa f1). go test ./internal/session -run TestStalled passes: working at 19m and 20m not stalled, 21m stalled; blocked, idle, usage-paused, exited and no report never. go test ./internal/api -run TestRunStalledNow passes: agent step, open leaves' latest sessions, checks (closed panes skipped), unjudged fixers, and a stopped (not live) run never stalls. Testscripts run-status-stalled (hand-written hook report at 00:00, RHIZA_TEST_NOW 00:19 no Stalled, 00:21 '- Stalled: build session ...', then after killing the runner the run is stopped with no Stalled in status, --json or run ls; idle never), inbox-stalled and describe-golden pass; go test ./... passes, gofmt and vet clean. Hands-on with a scratch build on live ship-change run 01M6AREK66 at review/review: RHIZA_TEST_NOW=2026-11-01T09:00:00Z gives stalled {step review/review, session 6a792fd4 (this verifier's leaf session)} in run status --json and a stalled item in rhiza inbox --json; at the real clock stalled is absent. Diff read: runInfo now passes live into runStalledNow; stepSessions covers agent, leaf, check and fixer sessions. make check's self-check fails only on stale .rhiza/rendered actor files from this run's model writes, not on code.

## Links

### Out

- verifies [run-status-stalled-session](../capability/rhiza/workflows/workflow-runs/run-status/run-status-stalled-session.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
