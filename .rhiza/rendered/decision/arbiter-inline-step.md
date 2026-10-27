# arbiter-inline-step

decision · `01M5XCKSY9JEVYAF4GK2GGH2TK`

## statement

The arbiter runs inside the asking agent step and is journalled there; a workflow may declare at most one arbiter step, with no prompt, no routes and not the start step, to choose its agent, permission mode and interval, and without one a built-in arbiter with the run's agent and mode is used.

## wording

the planning agent's option, picked by the orchestrator acting as arbiter for the owner

## Links

### In

- decided-by [arbiter-step](../capability/rhiza/workflows/arbiter-sessions/arbiter-step.md)
- decided-by [arbiter-on-question](../capability/rhiza/workflows/arbiter-sessions/arbiter-on-question.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
