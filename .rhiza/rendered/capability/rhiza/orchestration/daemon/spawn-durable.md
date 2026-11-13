# spawn-durable

capability · `01M79KPCGPWP3YX3Q5GRCFA0B0` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./internal/run -run TestSpawnResume kills the driver while a child runs and while a hand-back is pending, restarts it, and shows a live child reattached, a child that ended while the driver was down recorded once (not twice when its spawn process also records it), the pending hand-back delivered once, and the run continuing from the same step; a child whose session record is gone is recorded failed with that reason. Touches internal/run/engine.go, internal/run/journal.go and internal/api/run.go (the resume path).

## statement

Killing the daemon while a run waits on a spawned session or a hand-back loses nothing: on restart it reattaches to live children, records children that ended while it was down, and resumes the waiting step.

## Links

### Out

- depends-on [run-spawn](../../workflows/workflow-runs/run-spawn.md)
- depends-on [run-handback](../../workflows/workflow-runs/run-handback.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)
- decided-by [spawn-starts-from-caller](../../../../decision/spawn-starts-from-caller.md)

### In

- verifies [spawn-durable-verified](../../../../evidence/spawn-durable-verified.md)
- verifies [spawn-durable-qa-1](../../../../evidence/spawn-durable-qa-1.md)
- verifies [spawn-durable-qa-2](../../../../evidence/spawn-durable-qa-2.md)
- verifies [pr-572-the-daemon-offers-spawn](../../../../evidence/pr-572-the-daemon-offers-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
