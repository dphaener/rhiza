# arbiter-step

capability · `01M5WRHFM7T3P6V5ZV7XVKE55Q` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

Go workflow tests load a file with an arbiter step and refuse one that sets prompt on it. A testscript runs an arbiter step on a question with a fake adapter: the fake agent's prompt holds the question's ID and rhiza guide arbiter and no other node's handle or ID; when the fake writes the answering Decision the step passes, when it writes nothing or a Decision without reasoning it fails.

## statement

A workflow step of kind arbiter launches a fresh agent session on the run's change set whose prompt names only a question's ID and the arbiter guide topic. It passes when a live Decision answers that question with wording (assumed, would have escalated) and a non-empty reasoning field, and fails otherwise. A workflow file may not give it a prompt.

## Links

### Out

- depends-on [arbiter-guide](arbiter-guide.md)
- decided-by [arbiter](../../../../decision/arbiter.md)
- decided-by [arbiter-inline-step](../../../../decision/arbiter-inline-step.md)
- decided-by [arbiter-checks-decision](../../../../decision/arbiter-checks-decision.md)
- decided-by [arbiter-checked-when-needed](../../../../decision/arbiter-checked-when-needed.md)
- decided-by [arbiter-session-block-fails](../../../../decision/arbiter-session-block-fails.md)

### In

- depends-on [arbiter-on-question](arbiter-on-question.md)
- verifies [arbiter-step-verified](../../../../evidence/arbiter-step-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
