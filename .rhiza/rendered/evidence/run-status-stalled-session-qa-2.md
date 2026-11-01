# run-status-stalled-session-qa-2

evidence · `01M6AVR2Q6CHMCTKT4MPGPWMTR` · actor [claude-code-2d4e2697](../actor/claude-code-2d4e2697.md)

## result

PASS

## statement

Built binary at 8276602a against this repo's own live ship-change run 01M6AREK66 at step review/review (cycle-1 gap): the working cycle-2 quality check session c98484a2 last hook report 06:08:41Z. With RHIZA_TEST_NOW set (corrected for the log offset of 5503s) to report+19m, rhiza run status --json gives no stalled and rhiza inbox --json no stalled item; at report+21m run status --json gives stalled {step review/review, session c98484a2, since 06:08:41Z}. Earlier cycle sessions (compliance f9e4d1bf done, quality cycle-1 a1d0af81) and the other run (wait-main) never show. go test -run TestScripts/(run-status-stalled|describe-golden) passes, including the stopped-run and idle cases; go test ./internal/session ./internal/api pass.

## Links

### Out

- verifies [run-status-stalled-session](../capability/rhiza/workflows/workflow-runs/run-status/run-status-stalled-session.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
