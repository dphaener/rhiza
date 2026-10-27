# plan-review-guide

capability · `01M5WSGD6GSZAS9X5C75R728DG` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go guide tests: rhiza guide ls lists the plan-review topic, which names the five finding kinds, rhiza run findings and rhiza run verdict with their file formats, how the plan's change sets are found, and the rule that a critic never adds a finding; the topic stays within TopicCap.

## statement

rhiza guide plan-review tells a reviewer, given only a change set ID, to find the plan's change sets (that one and every change set serving an intent it adds), read them through rhiza slice and rhiza why, and record leaves without an acceptance check, decisions without provenance, depends-on links that aren't real, scope outside the intent and contradictions with rhiza run findings; it tells a critic to confirm or drop each finding with rhiza run verdict, never adding one, and a reviser to fix the confirmed findings in the model.

## Links

### Out

- depends-on [plan-review-journal](plan-review-journal.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)

### In

- depends-on [plan-workflow-file](plan-workflow-file.md)
- verifies [plan-review-guide-verified](../../../../../evidence/plan-review-guide-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
