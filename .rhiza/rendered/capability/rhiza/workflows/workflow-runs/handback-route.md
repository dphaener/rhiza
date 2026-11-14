# handback-route

capability · `01M7AHW9B0Y6HMMVGBQ4QX4Z3W` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/workflow -run TestHandbackTarget covers parsing handback:<step>, handback_cap defaulting to 3, and rhiza workflow check refusing a hand-back to an unknown step, to a step that is not an agent step and a handback_cap below 1, each naming the step; go test -trimpath ./internal/run -run TestHandbackRoute shows a command step failing twice then passing re-entering the agent step each time, and a fourth failure parking the run stopped with the cause; go test -trimpath ./internal/api -run TestInboxHandbackParked shows the inbox item naming the step and the cause. Touches internal/workflow/workflow.go, internal/workflow/check.go, internal/run/engine.go, internal/run/journal.go and internal/api/inbox.go; it delivers nothing to the session (run-handback does).

## statement

A step's on.<outcome> may name handback:<agent step> of the same workflow: the run re-enters that agent step with the routing step's cause, counts hand-backs per run, and once a run has handed back handback_cap times (a workflow setting, default 3) the next red parks the run stopped with an inbox item naming the step and the cause.

## Links

### Out

- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)
- decided-by [handback-split-route-deliver](../../../../decision/handback-split-route-deliver.md)

### In

- depends-on [run-handback](run-handback.md)
- verifies [handback-route-verified](../../../../evidence/handback-route-verified.md)
- verifies [handback-route-qa-1](../../../../evidence/handback-route-qa-1.md)
- verifies [handback-route-qa-2](../../../../evidence/handback-route-qa-2.md)
- verifies [pr-572-the-daemon-offers-spawn](../../../../evidence/pr-572-the-daemon-offers-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
