# check-model-topic

capability · `01M79AH2GKFS0RH9HHFZSPX13V` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a check-model test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, command-budget, outside-work, never-list and severity blocks expanded, the isolation rule (knows nothing of how the work was built; never `rhiza run status` or `rhiza run show` of the build run, its findings, run notes or session transcripts; judges only the contract from the slice, the diff from the merge-base and the model's Evidence), the code-contradicts-model and missing-constraint rules with their severities and the discovered-from draft commands; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide check-model` gives the model check its prompt on the six sections: the to-be model of the change set (slice, leaves, constraints, Decisions) matches the code on the branch, code that contradicts a node is a major and a missing constraint a minor, and a missing constraint or other model work outside the change set's intent goes into a discovered-from draft change set, never built in the run.

## Links

### Out

- depends-on [severity-shared-block](severity-shared-block.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [check-model-topic-verified](../../../../../../evidence/check-model-topic-verified.md)
- verifies [check-model-topic-qa-1](../../../../../../evidence/check-model-topic-qa-1.md)
- verifies [check-model-topic-verified-2](../../../../../../evidence/check-model-topic-verified-2.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
