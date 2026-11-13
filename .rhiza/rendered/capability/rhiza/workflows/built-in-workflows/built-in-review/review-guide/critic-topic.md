# critic-topic

capability · `01M79AH2PKWJWNFT52SCCKX7B0` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a critic test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, command-budget, never-list and severity blocks expanded, the isolation rule with the own-check-last-round exception, confirmed, downgraded (with a lower severity) and dropped, the never-add/raise/rewrite rule and the `rhiza run verdict --check` heredoc; a negative case asserts it carries no `rhiza run findings` example and no commit instruction.

## statement

`rhiza guide critic` gives a review check's critic afk's adversarial-critic prompt in Rhiza's words on the six sections: read only its own check's last round of this review run, re-check each blocker and major against the contract and the diff, and confirm, downgrade (the defect is real but its severity too high) or drop each with a reason, never adding a finding, raising a severity or rewriting one, recorded with `rhiza run verdict --check`.

## Links

### Out

- depends-on [severity-shared-block](severity-shared-block.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [critic-topic-verified](../../../../../../evidence/critic-topic-verified.md)
- verifies [critic-topic-qa-1](../../../../../../evidence/critic-topic-qa-1.md)
- verifies [critic-topic-verified-2](../../../../../../evidence/critic-topic-verified-2.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
