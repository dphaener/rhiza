# step-session-waits-holds

capability · `01M7QGA6QF8M2EWE3TNYRA3NKD` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/run and go test -trimpath ./internal/api -run 'StepHold|Spawn': an engine test with a fake Runner whose hold reports a memory-low detail for two polls sees one spawn-held entry for the step, no StartAgent call while held, then the session started and the hold released by the step's session entry; a stepRunner test sees the same detail rhiza run spawn gives for a saved tmp-low hold; (failure case) a step that rejoins its journaled session asks no hold and journals no spawn-held entry, and a held step is not reported stalled

## files

internal/run/engine.go, internal/run/journal.go, internal/run/engine_test.go, internal/api/run.go, internal/api/usage_hold.go

## must_not

Hold command steps, a rejoin or a hand-back; fail or stop a step because of a hold; change how the daemon picks or counts runs

## statement

The run engine starts a workflow step's fresh agent session only while no pre-start hold is on: while one is it waits, journaling one spawn-held entry for the step with no parent whose detail names the reading and the limit, and starts the session once the hold lifts; a step rejoining its journaled session, or handed back into it, is never held

## Links

### Out

- depends-on [spawn-waits-every-hold](spawn-waits-every-hold.md)
- decided-by [m40-spawn-gate-reads-daemon-state](../../../../decision/m40-spawn-gate-reads-daemon-state.md)

### In

- depends-on [inbox-held-spawns](../../orchestration/inbox/inbox-list/inbox-held-spawns.md)
- verifies [step-session-waits-holds-verified](../../../../evidence/step-session-waits-holds-verified.md)
- verifies [step-session-waits-holds-qa-1](../../../../evidence/step-session-waits-holds-qa-1.md)
- verifies [step-session-waits-holds-qa-2](../../../../evidence/step-session-waits-holds-qa-2.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
