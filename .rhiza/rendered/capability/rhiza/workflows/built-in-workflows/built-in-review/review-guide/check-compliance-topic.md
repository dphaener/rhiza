# check-compliance-topic

capability · `01M79AH1XVVY9F9AW9TFH2B12H` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a check-compliance test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, test-scope, command-budget, never-list and severity blocks expanded, the isolation rule (knows nothing of how the work was built; never `rhiza run status` or `rhiza run show` of the build run, its findings, run notes or session transcripts; judges only the contract from the slice, the diff from the merge-base and the model's Evidence), the under-built, scope-creep and unaccounted-files criteria with their severities, the one-line-with-evidence rule and the `rhiza run findings --check` heredoc with `{"findings":[]}` when clean; a negative case asserts it names no Go, make or repo-specific command and no commit or trailer instruction.

## statement

`rhiza guide check-compliance` gives the compliance check afk's compliance prompt in Rhiza's words on the six sections: each leaf's acceptance met as written, line by line, plus the scope criteria graded beside it: under-built (a leaf missing: blocker; an acceptance point unmet: major), scope creep (behaviour outside the intent and the leaves: major) and unaccounted files (a changed file no leaf or the intent accounts for: major if it changes behaviour, minor otherwise), each finding one line with its evidence so a fixer can act without re-diagnosing, recorded with `rhiza run findings --check`.

## Links

### Out

- depends-on [severity-shared-block](severity-shared-block.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [check-compliance-topic-verified](../../../../../../evidence/check-compliance-topic-verified.md)
- verifies [check-compliance-topic-qa-1](../../../../../../evidence/check-compliance-topic-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
