# spawn-visible-qa-1

evidence · `01M7BEHNAG67KCTGTMV82YDR1E` · actor [claude-code-a9fe07bb](../actor/claude-code-a9fe07bb.md)

## result

PASS

## statement

QA cycle 1, scratch build: rhiza run status listed '## Spawned sessions' with '- build, verifier on myleaf, model opus, session <id>: ended, pass', a failed child as 'ended, fail (exit 7)' and a live one as 'running'; --json spawned[] carried step, parent, session, role, leaf, model, state, outcome, detail; rhiza session show <child id> shows the child's record. The web page and pane/message were exercised only through the scoped tests TestRunDetailSpawned, TestRunStatusSpawned and TestSpawnedSessionOpenResume, which pass (no daemon or Rex in the scratch setup).

## Links

### Out

- verifies [spawn-visible](../capability/rhiza/workflows/workflow-runs/spawn-visible.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
