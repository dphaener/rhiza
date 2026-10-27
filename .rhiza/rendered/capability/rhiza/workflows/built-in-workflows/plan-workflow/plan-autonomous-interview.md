# plan-autonomous-interview

capability · `01M5WSGCNQDGV9W6JFMSXNQ17J` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go guide tests find in rhiza guide plan the RHIZA_AUTONOMOUS=1 rule, the arbiter and (assumed, would have escalated). A testscript runs the plan workflow with --autonomous on fake agents whose planner raises one question: a fake arbiter answers it, the run never blocks on the owner, and the planner's step carries on to the review.

## statement

In an autonomous run the planner asks the owner nothing: the plan text tells it that when RHIZA_AUTONOMOUS=1 each question it would have asked is raised for the arbiter, which records the answering Decision as (assumed, would have escalated), and that it quotes no owner words because there are none.

## Links

### Out

- depends-on [plan-interview](plan-interview.md)
- depends-on [arbiter-answer-to-session](../../arbiter-sessions/arbiter-answer-to-session.md)
- decided-by [autonomous-mode](../../../../../decision/autonomous-mode.md)
- decided-by [arbiter](../../../../../decision/arbiter.md)

### In

- verifies [plan-autonomous-interview-verified](../../../../../evidence/plan-autonomous-interview-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
