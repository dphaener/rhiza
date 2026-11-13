# ci-fix-topic

capability · `01M79AH3008VB5G7BNDF0YAHJP` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a ci-fix test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, test-scope, unstable-CI, command-budget and never-list blocks expanded, `rhiza ship check --json`, reproduce-before-fix and never-push-to-see, the flake quarantine with `rhiza request --kind bug`, the Checks: block and the Review: trailer; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide ci-fix` gives the ship's CI fixer its prompt on the six sections: read the CI findings (failing tests and non-test causes with their log lines) from `rhiza ship check --json`, reproduce each locally with the repo's own scoped command and never push to see, fix it or regenerate stale files, quarantine a non-reproducing flake with a bug request named at the skip, commit once those and the touched packages pass, with the pre-handoff checklist and a `Review: <run>` trailer, then session done as its own command.

## Links

### Out

- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [ci-fix-topic-verified](../../../../../../evidence/ci-fix-topic-verified.md)
- verifies [ci-fix-topic-qa-1](../../../../../../evidence/ci-fix-topic-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
