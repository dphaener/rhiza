# spawn-visible-qa-2

evidence · `01M7BFRKVVTBMM8Z88CTZA36G4` · actor [claude-code-420fd0ba](../actor/claude-code-420fd0ba.md)

## result

PASS

## statement

QA cycle 2, scratch build: rhiza run status listed '## Spawned sessions - cond, builder, model sonnet, session <id>: ended, pass' (and 'ended, fail (exit 1)' for a failed child); --json has spawned[] with step, parent, role, model, state, outcome, detail and session; rhiza session show accepted the child's session ID. TestRunStatusSpawned, TestSpawnedSessionOpenResume and TestRunDetailSpawned pass.

## Links

### Out

- verifies [spawn-visible](../capability/rhiza/workflows/workflow-runs/spawn-visible.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
