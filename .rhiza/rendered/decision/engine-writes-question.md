# engine-writes-question

decision · `01M5XCKTJGJYAPA6PT90M28QJJ`

## statement

The engine, not the hook, writes the question node; the hook only denies the question and keeps an asked mark on the session record that later hooks don't overwrite, so the session reports blocked on a question until the engine clears the mark after delivering the answer.

## wording

the planning agent's option, picked by the orchestrator acting as arbiter for the owner

## Links

### Out

- relates-to [arbiter-answer-by-message](arbiter-answer-by-message.md) — refines it: the engine, not the hook, writes the question node; the hook only denies the call and keeps the asked mark

### In

- decided-by [arbiter-on-question](../capability/rhiza/workflows/arbiter-sessions/arbiter-on-question.md)
- decided-by [arbiter-answer-to-session](../capability/rhiza/workflows/arbiter-sessions/arbiter-answer-to-session.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
