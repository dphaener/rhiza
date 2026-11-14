# gate-step

capability · `01M79KPCDGXDFMJJ3WEZ2G4H51` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

go test -trimpath ./internal/workflow -run TestGateStep covers parsing commands, rhiza workflow check refusing a gate with no commands and a gate with agent keys (prompt, agent, models); go test -trimpath ./internal/run -run TestGateStep shows a passing gate, a failing gate keeping the failing command's full output in the run directory while its later commands still run, and the journal and step detail naming the failing command. Touches internal/workflow/workflow.go, internal/workflow/check.go, internal/run/engine.go, internal/run/journal.go and internal/api/run.go (the Runner's command path, keeping output as ship_gate.go does).

## statement

A workflow step of kind gate runs its commands list, commands the daemon owns and never an agent, in order in the run's worktree; it passes only when every command exits 0, and on failure keeps each failing command's full output in a file in the run directory and names it in the step's detail for the step it routes to.

## Links

### Out

- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)

### In

- depends-on [run-handback](../workflow-runs/run-handback.md)
- verifies [gate-step-verified](../../../../evidence/gate-step-verified.md)
- verifies [gate-step-qa-1](../../../../evidence/gate-step-qa-1.md)
- verifies [gate-step-verified-2](../../../../evidence/gate-step-verified-2.md)
- verifies [gate-step-qa-2](../../../../evidence/gate-step-qa-2.md)
- verifies [pr-572-the-daemon-offers-spawn](../../../../evidence/pr-572-the-daemon-offers-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
