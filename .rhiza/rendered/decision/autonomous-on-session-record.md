# autonomous-on-session-record

decision · `01M5XCKVF033W8HSQR6CAE3N54`

## statement

A session started by an autonomous run is marked autonomous on its record when it is created; its launch environment sets RHIZA_AUTONOMOUS=1 from that mark, and the hook reads it to deny questions.

## wording

the planning agent's option, picked by the orchestrator acting as arbiter for the owner

## Links

### In

- decided-by [run-start-autonomous](../capability/rhiza/workflows/autonomous-runs/run-start-autonomous.md)
- decided-by [arbiter-answer-to-session](../capability/rhiza/workflows/arbiter-sessions/arbiter-answer-to-session.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
