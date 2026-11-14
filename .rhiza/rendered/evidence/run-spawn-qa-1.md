# run-spawn-qa-1

evidence · `01M7BEHN48EFRFCEKHJYQGBQ50` · actor [claude-code-a9fe07bb](../actor/claude-code-a9fe07bb.md)

## result

PASS

## statement

QA cycle 1, scratch build: outside a run, rhiza run spawn exits 2 naming RHIZA_RUN; from an agent step's session, rhiza run spawn --role verifier --leaf myleaf --model opus --brief <file> --json started a child in the run's worktree (argv had --model opus and --permission-mode manual), wrote spawn/<sid>/brief.md as '## Contract' section then '## Brief' with the caller's file, journaled spawn then spawn-ended under step build, and printed outcome pass, detail exit 0, exit 0; a child exiting 7 gave outcome fail, detail exit 7, exit 1; a missing --brief nope.md exits 2 naming the path. With machine_hold in daemon.state.json the spawn journaled spawn-held 'machine hold: load 99 over 8', and once the state was cleared it started the child and reported held true. Scoped tests TestRunSpawn and TestScripts/run-spawn, describe-golden pass.

## Links

### Out

- verifies [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
