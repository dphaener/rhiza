# gate-step-qa-1

evidence · `01M7BEHMTBTTC7YEX0S3RNSEC8` · actor [claude-code-a9fe07bb](../actor/claude-code-a9fe07bb.md)

## result

FAIL

## statement

QA cycle 1 with a fresh scratch build (go build -trimpath ./cmd/rhiza) in a scratch Rhiza root: rhiza run start --foreground on a gate step with commands [true-ish, 5000 lines + stderr + exit 3, touch marker] failed with the detail 'command 2 `...`: exit 3; full output: <run dir>/gate-gate-1-2.log' (5002 lines, stderr included), the third command still ran, and all ran in the run's worktree; an all-green gate passed. rhiza workflow check refused no commands, commands = [], a blank command, and prompt, agent and models on a gate, and commands on a command step. Broken: a gate step in an included workflow (outer [step.sub] kind = workflow, workflow = inner, inner gate ig) failed with 'its output could not be kept': the step name inner/ig makes gate-inner/ig-1-1.log a path under a directory that is never created. Scoped tests TestGateStep (workflow, run) pass.

## Links

### Out

- verifies [gate-step](../capability/rhiza/workflows/workflow-definitions/gate-step.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
