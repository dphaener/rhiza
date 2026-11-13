# fixer-lint-hits-guide

capability · `01M7B0BDGMMAJBBSM9MA9YG4QX`

## acceptance

go test -trimpath ./internal/guide/... passes with fixer_test.go finding the lint-hit section naming `rhiza run lint-check --json`, the three fixes and the unfixed rule; failure case: the existing fixer rules (fix only the prompt's findings, the Checks: block, never push) are still found.

## statement

`rhiza guide fixer` gains a section for a lint-fix fixer (one whose prompt names `rhiza run lint-check`): read the hits from `rhiza run lint-check --json`; fix each hit it agrees with in place (rewrite stale statement or acceptance text with rhiza set in the change set, retire a capability whose code the branch removed, or link the node a Decision governs with decided-by), list each hit it judges wrong as unfixed with why, never reproduce a hit as a failing test, and end with one commit trailed Review: as any fixer, then session done. The rest of the fixer topic is unchanged. Files: internal/guide/text/fixer.md, internal/guide/fixer_test.go.

## Links

### Out

- depends-on [review-lint-fix-route](review-lint-fix-route.md)
- decided-by [lint-hits-feed-model-check](../decision/lint-hits-feed-model-check.md)

### In

- blocks [fixer-lint-section-over-topic-cap](../question/fixer-lint-section-over-topic-cap.md)
- blocks [leaf-failed-01m7b0sj-1](../question/leaf-failed-01m7b0sj-1.md)
- verifies [fixer-lint-hits-guide-verified](../evidence/fixer-lint-hits-guide-verified.md)
- verifies [fixer-lint-hits-guide-qa-1](../evidence/fixer-lint-hits-guide-qa-1.md)
- verifies [pr-569-model-lints-in-review](../evidence/pr-569-model-lints-in-review.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
