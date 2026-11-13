# run-spawn-qa-2

evidence · `01M7BFRKRFP72Z59HA7WPNV1ZG` · actor [claude-code-420fd0ba](../actor/claude-code-420fd0ba.md)

## result

PASS

## statement

QA cycle 2, scratch build, fake rex/claude: an agent step's session ran rhiza run spawn --role builder --brief <file> --model sonnet: the child started in the run's worktree, its spawn/<session>/brief.md ends with the contract section then '## Brief / do the child thing', and spawn printed 'Spawned builder session ... under step cond ...: pass (exit 0)'; a child exiting 1 gave 'fail (exit 1)' and spawn exit 1. A missing --brief exits 2 naming the path; outside a run it refuses naming RHIZA_RUN. TestRunSpawn and TestScripts/run-spawn pass.

## Links

### Out

- verifies [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
