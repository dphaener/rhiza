# plan-review-rounds

capability · `01M5WSGDAB8W8BWQX6VEWR1FVY` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript in a run with journaled rounds: no confirmed finding passes; one confirmed major-kind finding in round 1 or 2 fails and adds nothing; one confirmed in round 3 adds a question node whose statement holds the finding and its reason, linked blocks to the capability it names, journals it with an asked entry per added question, and passes. go test -trimpath ./internal/api -run TestPlanReviewDowngradeByKind shows a confirmed missing-provenance finding and a missing-acceptance one downgraded to minor pass review-check and gateReview, and the missing-acceptance one, once confirmed, fails review-check.

## statement

rhiza run review-check counts only the latest round's confirmed findings that are still major: a finding's own severity, or the severity a critic downgraded it to, else its kind's (missing-acceptance, contradiction and outside-intent are major; missing-provenance and unreal-depends-on are minor), a confirmed minor or nit never counting. It passes when the latest round has none and fails when it has any in round 1 or 2; after a third round that still has them it adds each as a question node (in the run's change set while it is open, else where review-questions-follow-the-node puts it), linked blocks to the node it names when that is a capability, journals them, and passes, so the plan reaches the gate with those questions open; a confirmed minor-kind finding is never added as a question. It also journals each question it adds as an asked entry, so that in an autonomous run the arbiter decides it. The gate's review check (gateReview) counts the same findings.

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
- verifies [pr-510-critic-downgrade](../../../../../evidence/pr-510-critic-downgrade.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
