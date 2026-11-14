# spawn-durable-qa-1

evidence · `01M7BEHN6Z4QEKMP5FGFPR24AE` · actor [claude-code-a9fe07bb](../actor/claude-code-a9fe07bb.md)

## result

PASS

## statement

QA cycle 1, scratch build: with a spawned child blocked on a fifo, kill -9 of the foreground driver only (the spawn process lived on), then rhiza run resume --foreground: the driver reattached, waited on the step, and after the child was released the journal held exactly one spawn-ended (the spawn process and the driver did not both record it) and the run finished from step build. Scoped test TestSpawnResume passes.

## Links

### Out

- verifies [spawn-durable](../capability/rhiza/orchestration/daemon/spawn-durable.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
