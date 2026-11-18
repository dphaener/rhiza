# step-session-waits-holds-verified

evidence · `01M7QHYXAE9XGC1ANFXQJSS2DJ` · actor [claude-code-eba74c98](../actor/claude-code-eba74c98.md)

## result

PASS

## statement

Verified commit 278907eb7. Re-ran its Checks (capped scopes): go vet ./internal/api ./internal/run exit 0; go test -trimpath -p 2 ./internal/run exit 0; go test -trimpath -p 2 ./internal/api -run 'StepHold|Spawn' exit 0 (TestWatchdogStepHoldNotStalled, TestRunSpawnStepHoldSameDetail PASS); rhiza ship generate exit 0 (0 files); rhiza render exit 0 (only the verifier's own actor file, restored). go test -count=1 ./internal/run -run StepHold: TestAgentStepHoldWaitsThenStarts (memory-low detail for two polls: one spawn-held entry for build with no Parent, heldStarts 0, session entry after it releases the hold, StepHeld false), TestAgentStepHoldResumedWaitsOnce, TestAgentStepHoldNotOnRejoin (no StartHold call, no spawn-held) all PASS. Diff: engine.go waitStartHold runs only when session == "" (never on rejoin or handback, agent steps only), never fails a step; stepRunner.StartHold returns spawnHoldDetail, same source as rhiza run spawn per Decision m40-spawn-gate-reads-daemon-state; journal.go releases a no-Parent hold on the step's session entry; watchdog skips a held step (needed for the not-stalled acceptance point). No stubs or TODOs; no change to daemon run picking or counting.

## Links

### Out

- verifies [step-session-waits-holds](../capability/rhiza/workflows/workflow-runs/step-session-waits-holds.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
