# plan-review-rounds

capability · `01M5WSGDAB8W8BWQX6VEWR1FVY` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript in a run with journaled rounds: no confirmed finding passes; one confirmed in round 1 or 2 fails and adds nothing; one confirmed in round 3 adds a question node whose statement holds the finding and its reason, linked blocks to the capability it names, journals it with an asked entry per added question, and passes.

## statement

rhiza run review-check passes when the latest round's critic confirmed no finding and fails when it confirmed any in round 1 or 2; after a third round that still has confirmed findings it adds each as a question node (in the run's change set while it is open, else where review-questions-follow-the-node puts it), linked blocks to the node it names when that is a capability, journals them, and passes, so the plan reaches the gate with those questions open. It also journals each question it adds as an asked entry, so that in an autonomous run the arbiter decides it.

## Links

### Out

- depends-on [plan-review-journal](plan-review-journal.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [plan-review-rounds-journalled](../../../../../decision/plan-review-rounds-journalled.md)
- decided-by [plan-review-left-arbitrated](../../../../../decision/plan-review-left-arbitrated.md)
- decided-by [review-questions-follow-the-node](../../../../../decision/review-questions-follow-the-node.md)

### In

- depends-on [plan-gate](plan-gate.md)
- depends-on [review-question-handle-unique](review-question-handle-unique.md)
- verifies [plan-review-rounds-verified](../../../../../evidence/plan-review-rounds-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
