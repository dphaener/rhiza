# gate-step-qa-2

evidence · `01M7BFRKC6N5V5NPTH4Y2YZ8GW` · actor [claude-code-420fd0ba](../actor/claude-code-420fd0ba.md)

## result

PASS

## statement

QA cycle 2 at e5f3f751e, fresh go build -trimpath ./cmd/rhiza to a scratch path, scratch root and repo, rhiza run start --foreground: the cycle-1 failure is fixed. An included gate (outer kind=workflow -> inner gate ig, commands [echo, seq 1 3000 + stderr + exit 5, pwd > file]) failed with detail 'command 2 `...`: exit 5; full output: <run dir>/gate-inner_ig-1-2.log' (3001 lines, stderr kept), the third command ran in the run's worktree; run resume --retry wrote gate-inner_ig-2-2.log beside visit 1's. A top-level failing gate kept gate-g-1-2.log and ran its later command; an all-green included gate passed ('2 commands passed'). run status --json carries the detail. Underscored step names are refused, so safeName cannot collide. Scoped TestGateStep (workflow, run) pass.

## Links

### Out

- verifies [gate-step](../capability/rhiza/workflows/workflow-definitions/gate-step.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
