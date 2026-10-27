# arbiter-answer-to-session

capability · `01M5WRHGM59FSCHB3SDT7C8H75` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

A testscript with a fake adapter whose agent blocks on a question in an autonomous run, and a fake arbiter that answers it: the asking fake agent receives the Decision's statement, the journal has the block and its unblock, and the step ends with the agent's outcome.

## statement

Once the arbiter has decided, the engine gives the answering Decision to the session that asked, the run is no longer blocked, and the agent step carries on to its own outcome.

## Links

### Out

- depends-on [arbiter-on-question](arbiter-on-question.md)
- decided-by [arbiter](../../../../decision/arbiter.md)
- decided-by [arbiter-answer-by-message](../../../../decision/arbiter-answer-by-message.md)
- decided-by [engine-writes-question](../../../../decision/engine-writes-question.md)
- decided-by [autonomous-on-session-record](../../../../decision/autonomous-on-session-record.md)
- decided-by [answer-message-contents](../../../../decision/answer-message-contents.md)
- decided-by [resume-rejoins-arbiter](../../../../decision/resume-rejoins-arbiter.md)
- decided-by [answer-reread-from-model](../../../../decision/answer-reread-from-model.md)
- decided-by [answer-sent-as-caller](../../../../decision/answer-sent-as-caller.md)

### In

- blocks [arbiter-answer-delivery](../../../../question/arbiter-answer-delivery.md)
- verifies [arbiter-answer-to-session-verified](../../../../evidence/arbiter-answer-to-session-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
