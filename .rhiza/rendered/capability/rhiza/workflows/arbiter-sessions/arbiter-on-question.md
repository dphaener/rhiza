# arbiter-on-question

capability · `01M5WRHG053937VJ6GEJ84EJJH` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

A testscript with a fake adapter whose agent reports blocked with kind question: autonomous, the run's change set gains a question node whose statement is the reported question, the journal has an arbiter entry naming it, and a fake arbiter is launched with its ID; a failing fake arbiter ends the run failed with the question named; normal, run status shows blocked, the change set gains no question node and no arbiter is launched.

## statement

When an agent step's session in an autonomous run reports it is blocked on a question, the engine adds that question, as the session reported it, as a question node in the run's change set, journals it, and runs the arbiter step on its ID; if the arbiter fails, the run fails naming the question. In a normal run the engine does none of this: the run is blocked on the owner, as now.

## Links

### Out

- depends-on [run-start-autonomous](../autonomous-runs/run-start-autonomous.md)
- depends-on [arbiter-step](arbiter-step.md)
- decided-by [arbiter](../../../../decision/arbiter.md)
- decided-by [autonomous-mode](../../../../decision/autonomous-mode.md)
- decided-by [arbiter-inline-step](../../../../decision/arbiter-inline-step.md)
- decided-by [engine-writes-question](../../../../decision/engine-writes-question.md)
- decided-by [question-node-shape](../../../../decision/question-node-shape.md)
- decided-by [arbiter-no-recursion](../../../../decision/arbiter-no-recursion.md)
- decided-by [resume-rejoins-arbiter](../../../../decision/resume-rejoins-arbiter.md)
- decided-by [arbiter-fail-ends-step](../../../../decision/arbiter-fail-ends-step.md)
- decided-by [one-question-per-block](../../../../decision/one-question-per-block.md)

### In

- depends-on [arbiter-answer-to-session](arbiter-answer-to-session.md)
- depends-on [autonomous-review](../autonomous-runs/autonomous-review.md)
- verifies [arbiter-on-question-verified](../../../../evidence/arbiter-on-question-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
