# review-guide

capability · `01M5Z0ENY6PDKAQS7QEY0AR5RN` · parent [built-in-review](../built-in-review.md)

## acceptance

Go guide tests: `rhiza guide ls` lists the review topic, which names the five checks, the merge-base, the four severities, `rhiza run findings` and `rhiza run verdict` with their formats, the critic rule, both trailers, and the Evidence and discovered-from rules; the topic stays within TopicCap; `rhiza guide install` writes it into the plugin.

## statement

`rhiza guide review` tells each check, critic and fixer its part: the five checks' briefs (compliance against each leaf's acceptance, quality, hands-on QA with the built binary where fakes alone don't count, scope inside the intent and the leaves, the model), reviewing from the merge-base with the default branch, the findings format and severities, that a critic never adds a finding, the fixers' trailers, scope work outside the intent as a question or a discovered-from draft, the QA check's Evidence on each leaf it exercised and the model check's missing constraints in a discovered-from draft; guide install installs it.

## Links

### Out

- depends-on [review-findings](review-findings.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-feeds-model](../../../../../decision/review-feeds-model.md)
- decided-by [review-ruling-2](../../../../../decision/review-ruling-2.md)
- decided-by [review-ruling-25](../../../../../decision/review-ruling-25.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- blocks [review-overruled-findings](../../../../../question/review-overruled-findings.md)
- verifies [review-guide-verified](../../../../../evidence/review-guide-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
