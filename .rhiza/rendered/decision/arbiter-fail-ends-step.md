# arbiter-fail-ends-step

decision · `01M5XCKXM1AGSFCSEWF812THF5`

## statement

When the arbiter fails or errs, or the question node cannot be written, the asking step ends with an error rather than a fail, so the workflow's fail route cannot go around the arbiter.

## wording

the building agent's choice while building leaf arbiter-on-question

## Links

### In

- decided-by [arbiter-on-question](../capability/rhiza/workflows/arbiter-sessions/arbiter-on-question.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
