# check-model-topic-verified-2

evidence · `01M79SBA475CG0ESWBFK3RCHJZ` · actor [claude-code-0e05ad88](../actor/claude-code-0e05ad88.md)

## result

PASS

## statement

At d73afba27 (attempt 2): go test -trimpath -count=1 ./internal/guide/... exit 0 (TestCheckModelText checks the six sections in order, the session-done, test-scope, command-budget, outside-work, never-list and severity blocks expanded, the isolation rule in Task, major/minor grading and the discovered-from draft commands; the negative test TestCheckModelNamesNoRepoCommand passes). A scratch build of `rhiza guide check-model | grep -q "only the tests covering"` exits 0. `rhiza ship generate` changed 0 files and exited 0. `rhiza render` exited 0, writing only rendered actors from this session's untracked log. `rhiza check` exits 0. The diff touches only check-model.md, its test and the builder's actor record. The trims keep every acceptance phrase, and it leaves no TODOs.

## Links

### Out

- verifies [check-model-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-model-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
