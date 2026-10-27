# arbiter-answer-delivery

question · `01M5WRHH045QXSZA0X2BP6E8S0`

## statement

How does the arbiter's answer reach a session that is waiting on its question, when that session is waiting inside Claude Code's AskUserQuestion rather than at its prompt, where rhiza session message delivers? (For example: a PreToolUse hook on AskUserQuestion in autonomous runs, or the guide telling agents in autonomous runs to raise questions another way.)

## Links

### Out

- blocks [arbiter-answer-to-session](../capability/rhiza/workflows/arbiter-sessions/arbiter-answer-to-session.md)

### In

- answers [arbiter-answer-by-message](../decision/arbiter-answer-by-message.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
