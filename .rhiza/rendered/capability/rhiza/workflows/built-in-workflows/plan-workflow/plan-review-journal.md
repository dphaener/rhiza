# plan-review-journal

capability · `01M5WSGD28YB7NQC8BQR8MCVFM` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript with RHIZA_RUN set records two findings and a verdict confirming one and dropping the other: the journal holds both entries, and run status --json and its markdown show the round with one confirmed finding. A verdict naming a third finding, a finding of an unknown kind, or one naming a node that does not exist exits 2 with nothing written; with no RHIZA_RUN both commands exit 2. describe-golden is regenerated.

## statement

A plan review round is journaled in its run: rhiza run findings <file> records the reviewer's findings, each naming a node, one kind (missing-acceptance, missing-provenance, unreal-depends-on, outside-intent, contradiction) and its reason; rhiza run verdict <file> records the critic's verdict on each, confirmed or dropped with a reason, and refuses a finding the round does not hold, so a critic can't add one; rhiza run status shows each round's findings and verdicts.

## Links

### Out

- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)

### In

- depends-on [plan-review-guide](plan-review-guide.md)
- depends-on [plan-review-rounds](plan-review-rounds.md)
- verifies [plan-review-journal-verified](../../../../../evidence/plan-review-journal-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
