# check-quality-topic

capability · `01M79AH24A6DZDBJ5EDMH0G21T` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a check-quality test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, test-scope, command-budget, never-list and severity blocks expanded, the isolation rule (knows nothing of how the work was built; never `rhiza run status` or `rhiza run show` of the build run, its findings, run notes or session transcripts; judges only the contract from the slice, the diff from the merge-base and the model's Evidence), discover-commands-from-the-contributor-guide, the impact radius (direct, integration, downstream), the never-fall-back-to-the-full-suite rule and the seven checklist items; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide check-quality` gives the quality check afk's quality prompt in Rhiza's words on the six sections: discover the test commands from the repo's contributor guide, map each changed file to the tests that cover it (direct, integration, downstream) as the impact radius, run those scoped commands and never silently fall back to the full suite, and walk a checklist (swallowed errors, debug output, commented-out code, unused imports, TODOs, generic names, missing error handling), each defect graded by the severity block.

## Links

### Out

- depends-on [severity-shared-block](severity-shared-block.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [check-quality-topic-verified](../../../../../../evidence/check-quality-topic-verified.md)
- verifies [check-quality-topic-qa-1](../../../../../../evidence/check-quality-topic-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
