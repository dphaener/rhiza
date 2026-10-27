# plan-gate

capability · `01M5WSGDEE854PFE4BXJM13CWS` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript ends a plan run whose record change set, landed after the interview, adds an intent served by two draft change sets: normal, both are proposed, none accepted, the landed record is left as it is, and run status --json names rhiza accept --rule owner; autonomous with a passing review, both are accepted with rule autonomous and reasons naming the run; autonomous with a question left by the review, none is accepted and the command exits non-zero naming the question; once the arbiter's Decision answers it, both are accepted. A run on an intention change set that serves an existing intent proposes that change set itself.

## statement

rhiza run gate ends a plan: it proposes each draft among the plan's change sets (the run's change set and every change set that serves an intent the run's change set adds), leaving landed or accepted ones as they are. In a normal run they wait for the owner's rhiza accept --rule owner, which run status names; in an autonomous run it accepts each with rule autonomous once the review passed and every question from the review is answered, and accepts nothing while a question from the review is open.

## Links

### Out

- depends-on [plan-review-rounds](plan-review-rounds.md)
- depends-on [autonomous-gate-rule](../../autonomous-runs/autonomous-gate-rule.md)
- decided-by [two-owner-gates](../../../../../decision/two-owner-gates.md)
- decided-by [autonomous-mode](../../../../../decision/autonomous-mode.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [plan-review-left-arbitrated](../../../../../decision/plan-review-left-arbitrated.md)

### In

- depends-on [plan-workflow-file](plan-workflow-file.md)
- relates-to [plan-record-landing](../../../../../question/plan-record-landing.md) — The gate proposes the plan's own change set with the others; who lands it after acceptance is not decided.
- relates-to [plan-review-left-autonomous](../../../../../question/plan-review-left-autonomous.md) — The gate accepts nothing while a review question is open; whether the arbiter settles those questions first is not decided.
- verifies [plan-gate-verified](../../../../../evidence/plan-gate-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
