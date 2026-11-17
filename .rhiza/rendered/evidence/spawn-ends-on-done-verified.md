# spawn-ends-on-done-verified

evidence · `01M7MBM6H03HC1Q0W1WGKNWKMX` · actor [claude-code-7abf1425](../actor/claude-code-7abf1425.md)

## result

PASS

## statement

In a detached checkout at 8c10aeb1d (capped scope, RHIZA_/SHIP_/REX_ unset): go test -trimpath ./internal/api -run 'TestRunSpawn|TestStepRunnerChild' -count=1 passed (TestRunSpawnSaidDone NoExit/LaterExit/Lost/CloseFails/LiveThenExit, TestStepRunnerChild done/done-red/done-lost, TestStepRunnerChildDoneLost via markLostIfEnded); go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-spawn|describe-golden)' -count=1 passed; rhiza docs, rhiza ship generate and rhiza render left the tree clean. Diff read: childOutcome returns pass/'said done' on a done mark before checking the exit; spawn journals the end then calls closeSessionPane (rex block close -s s1 block:c1 captured via ctx.Exec), a close error only goes to stderr; live child sleeps twice then fails exit 3 with no rex call. Files stay within the leaf's list (plus .rhiza log/rendered); internal/run/engine.go and session.done untouched; no TODOs; tests use fakes only.

## Links

### Out

- verifies [spawn-ends-on-done](../capability/rhiza/workflows/workflow-runs/spawn-ends-on-done.md)

## History

- `01M7M27J4MMS8X5PT4W3EXKDTM` Fix rhiza run spawn waits for the child agent's process to exit, and with no daemon nothing closes a child that already said rhiza session done, so spawn and its run hang until someone closes the child's pane
