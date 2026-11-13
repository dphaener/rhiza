# fixer-topic

capability · `01M79AH2TFR25PCV7H0T92DEDC` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a fixer test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, test-scope, unstable-CI, command-budget, outside-work, never-list and severity blocks expanded, fix-only-the-named-findings, reproduce-first, `rhiza ship generate`, `rhiza render` and the Checks: block above the one final trailer paragraph, both trailer forms, the unfixed-never-dropped rule and never push, rebase, force or --no-verify; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide fixer` gives a review fixer afk's fix prompt in Rhiza's words on the six sections: fix only the confirmed findings in its prompt, each exactly as written and in place, reproduce each first, run the pre-handoff checklist (`rhiza ship generate`, `rhiza render`, a Checks: block of exit codes), make one new commit trailed Change:, Actor: and Leaf: <leaf> (leaf fixer) or Review: <run>, carry a finding it could not fix into the session's end as unfixed instead of dropping it, never push, rebase, force or use --no-verify, then session done as its own command.

## Links

### Out

- depends-on [severity-shared-block](severity-shared-block.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [fixer-topic-verified](../../../../../../evidence/fixer-topic-verified.md)
- verifies [fixer-topic-qa-1](../../../../../../evidence/fixer-topic-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
