# conflict-fix-topic

capability · `01M79AH36GH3GEF9Q1XQ3B68YC` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a conflict-fix test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, test-scope, command-budget and never-list blocks expanded, both parts with `rhiza ship check --json` merge_conflict and `rhiza ship check --model-conflict --json`, `rhiza set <node> --rebase`, `git add -A -- .rhiza/rendered .rhiza/log`, the Review: and Model-Conflict: trailers, the never-end-with-the-merge-in-progress rule, the one-blocking-question rule and never push, rebase or force; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide conflict-fix` gives the ship's conflict fixer its prompt on the six sections with two parts: a code merge left in progress (read both sides from `rhiza ship check --json`, resolve so both sides' behaviour holds, regenerate and render, stage .rhiza/rendered and .rhiza/log, commit the merge with the Checks: block and Change:, Actor:, Review: trailers, never end with the merge in progress) and a model conflict that refused the land (read it from `rhiza ship check --model-conflict --json`, move each stale update with `rhiza set <node> --rebase` so both edits hold, commit with the Model-Conflict: trailer); on a real clash it commits nothing and asks one blocking question naming both change sets; it never pushes, rebases or forces.

## Links

### Out

- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [conflict-fix-topic-verified](../../../../../../evidence/conflict-fix-topic-verified.md)
- verifies [conflict-fix-topic-qa-1](../../../../../../evidence/conflict-fix-topic-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
