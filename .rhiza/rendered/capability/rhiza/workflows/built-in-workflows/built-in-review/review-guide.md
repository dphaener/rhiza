# review-guide

capability · `01M5Z0ENY6PDKAQS7QEY0AR5RN` · parent [built-in-review](../built-in-review.md)

## acceptance

`rhiza guide ls` lists review and the nine role topics; go test -trimpath ./internal/guide/... finds the nine names in the review index and the merge-base, findings, verdict, trailer, Evidence and discovered-from rules in the role topics, each within TopicCap; `rhiza guide install` writes all ten as skills.

## statement

`rhiza guide review` is the review index: it names the check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix topics and says which role loads which (a check named scope loads check-compliance); the role topics carry reviewing from the merge-base, isolation from the build's history, the findings and verdict formats and the shared severity block, the fixers' trailers, outside work as a question or a discovered-from draft, the QA check's Evidence and the model check's missing constraints in a discovered-from draft; guide install installs them all.

## Links

### Out

- depends-on [review-findings](review-findings.md)
- decided-by [review-workflow-shape](../../../../../decision/review-workflow-shape.md)
- decided-by [review-feeds-model](../../../../../decision/review-feeds-model.md)
- decided-by [review-ruling-2](../../../../../decision/review-ruling-2.md)
- decided-by [review-ruling-25](../../../../../decision/review-ruling-25.md)
- decided-by [review-role-topics-shape](../../../../../decision/review-role-topics-shape.md)

### In

- depends-on [review-workflow-file](review-workflow-file.md)
- blocks [review-overruled-findings](../../../../../question/review-overruled-findings.md)
- verifies [review-guide-verified](../../../../../evidence/review-guide-verified.md)
- verifies [pr-45-the-review-workflow](../../../../../evidence/pr-45-the-review-workflow.md)
- verifies [review-guide-qa-1](../../../../../evidence/review-guide-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
