# check-qa-topic

capability · `01M79AH2C9CN1F6PMN6G3SYF0Z` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a check-qa test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, command-budget, never-list and severity blocks expanded, the isolation rule (knows nothing of how the work was built; never `rhiza run status` or `rhiza run show` of the build run, its findings, run notes or session transcripts; judges only the contract from the slice, the diff from the merge-base and the model's Evidence), the clean-state rule, the scratch path outside $PATH and never-install rule, happy and failure scenarios, bounded exploration with no new test files, the broken-behaviour-only rule and the Evidence apply example with `<leaf>-qa-<cycle>` and `verifies`; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide check-qa` gives the QA check afk's hands-on QA prompt in Rhiza's words on the six sections: start from a clean state, build what changed to a scratch path outside $PATH and never install a binary, run each leaf's happy and failure scenarios with the built program (fakes alone do not count), explore in bounded pairs then edge cases without writing new test files, report only broken behaviour (never a missing test or style point), and write Evidence as itself on each leaf it exercised.

## Links

### Out

- depends-on [severity-shared-block](severity-shared-block.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [check-qa-topic-verified](../../../../../../evidence/check-qa-topic-verified.md)
- verifies [check-qa-topic-qa-1](../../../../../../evidence/check-qa-topic-qa-1.md)
- verifies [check-qa-topic-verified-2](../../../../../../evidence/check-qa-topic-verified-2.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
