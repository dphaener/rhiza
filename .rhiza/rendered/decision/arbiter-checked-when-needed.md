# arbiter-checked-when-needed

decision · `01M5XCKX3TZXA0K6AEZTX2CD7Y`

## statement

A run's arbiter is checked at start only when the run is autonomous and its workflow has an agent step, since only agent steps raise questions and a non-autonomous run never launches the arbiter.

## wording

the building agent's choice while building leaf arbiter-step

## Links

### In

- decided-by [arbiter-step](../capability/rhiza/workflows/arbiter-sessions/arbiter-step.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
