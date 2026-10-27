# arbiter-answer-by-message

decision · `01M5WT5QB5QR7Y1HX1PCB5X443`

## statement

In an autonomous run, a PreToolUse hook on AskUserQuestion denies the call, records the question as a question node, and tells the session its answer will arrive as a message, so the session's turn ends at its prompt; once the arbiter's Decision answers the question, the engine delivers it with rhiza session message. Nothing is typed into a menu, and no hook waits on the arbiter.

## wording

the orchestrator's ruling as arbiter, under the owner's instruction to be the arbiter (milestone-3 conversation); not the owner's words

## Links

### Out

- answers [arbiter-answer-delivery](../question/arbiter-answer-delivery.md)

### In

- relates-to [engine-writes-question](engine-writes-question.md) — refines it: the engine, not the hook, writes the question node; the hook only denies the call and keeps the asked mark
- decided-by [arbiter-answer-to-session](../capability/rhiza/workflows/arbiter-sessions/arbiter-answer-to-session.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
