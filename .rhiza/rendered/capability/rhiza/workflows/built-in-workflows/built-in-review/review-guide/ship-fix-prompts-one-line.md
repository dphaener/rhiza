# ship-fix-prompts-one-line

capability · `01M79AH3X5BDT03Z5D6GE1DYEP` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/workflow/... passes with a test that every agent step of the builtin ship-change.toml has a one-line prompt naming its topic (ci-fix, gate-fix for gate-fix and gate-fix-merge, conflict-fix for conflict-fix, conflict-fix-landed and model-conflict-fix) with no session-done, Checks: or trailer text left in it, and a test reading .rhiza/workflows/ship-change.toml and .rhiza/workflows/docs.toml that asserts each of their fix steps' prompts equals the builtin step's; go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-ci-fix|run-ship-ci-lint-fix|run-ship-conflict-fix|run-ship-change)' passes; make self-check passes.

## statement

The builtin ship-change workflow's fix steps each have a one-line prompt naming its topic and its task variables: ci-fix `/rhiza-ci-fix Fix the CI findings of change set {change} in run {run}: follow rhiza guide ci-fix.`, gate-fix and gate-fix-merge `/rhiza-gate-fix ...: follow rhiza guide gate-fix.`, conflict-fix and conflict-fix-landed `/rhiza-conflict-fix Resolve the merge conflicts ...: follow rhiza guide conflict-fix.` and model-conflict-fix `/rhiza-conflict-fix Fix the model conflict that refused the land ...: follow rhiza guide conflict-fix.`; Rhiza's .rhiza/workflows/ship-change.toml and docs.toml use the same lines byte for byte.

## Links

### Out

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)

### In

- verifies [ship-fix-prompts-one-line-verified](../../../../../../evidence/ship-fix-prompts-one-line-verified.md)
- verifies [ship-fix-prompts-one-line-qa-1](../../../../../../evidence/ship-fix-prompts-one-line-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
