# arbiter-checks-decision

decision · `01M5XCKTAE6R38G4JHC5T13WF2`

## statement

The arbiter step looks for its Decision once its session has said done or exited: a valid Decision passes, anything else fails, and the arbiter's session is left open.

## wording

the planning agent's option, picked by the orchestrator acting as arbiter for the owner

## Links

### In

- decided-by [arbiter-step](../capability/rhiza/workflows/arbiter-sessions/arbiter-step.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
