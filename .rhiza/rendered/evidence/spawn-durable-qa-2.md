# spawn-durable-qa-2

evidence · `01M7BFRKYXA84RWYVPRF9Y0MMQ` · actor [claude-code-420fd0ba](../actor/claude-code-420fd0ba.md)

## result

PASS

## statement

QA cycle 2: not killed by hand (code unchanged since cycle 1's QA); go test -trimpath -count=1 ./internal/run -run TestSpawnResume passes (driver killed with a live child and a pending hand-back, reattach, recorded once, delivered once).

## Links

### Out

- verifies [spawn-durable](../capability/rhiza/orchestration/daemon/spawn-durable.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
