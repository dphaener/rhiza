# run-status-stalled-session-qa-1

evidence · `01M6ATQ5T7MG20PKQGAQ33A4M0` · actor [claude-code-7b8ceb12](../actor/claude-code-7b8ceb12.md)

## result

FAIL

## statement

At 973ce724 I built go build -o \$TMPDIR/rhiza-qa ./cmd/rhiza and, in a scratch repo, started a real rhiza run (one agent step, fake rex and claude whose hook reported working through rhiza session hook). On the real clock: fresh report, run status shows no Stalled; hook at back-dated 2026-11-01T04:00:00Z (1h50m old) prints '- Stalled: build session 107f052b-... has not reported since 2026-11-01T04:00:00Z' and --json has stalled {step build, session, since}; turned blocked, or working with detail 'paused at its usage limit until ...', or after the run ended: no Stalled. Gap: in this repo's own ship-change run 01M6AREK66 (current step review/review), the built binary's run status --json gives stalled null and inbox --json lists no stalled item at RHIZA_TEST_NOW=2026-11-01T09:00:00Z although check session 7b8ceb12 is a live session whose hook is working and last reported 05:52Z: runStalledNow reads only st.Session[st.Current], and ship-change records its builder, verifier and check sessions as leaf-session/check-session entries, so they never stall.

## Links

### Out

- verifies [run-status-stalled-session](../capability/rhiza/workflows/workflow-runs/run-status/run-status-stalled-session.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
