# review-lint-fix-route

capability · `01M7B0BDCY9W3A1M8XAT24DGDM`

## acceptance

go test -trimpath ./internal/api/... -run 'LintCheck' passes: a run whose latest code-check entry has two lints exits 1 with both hits in --json and the human output; failure cases: an entry with no lints and a run with no code-check entry exit 0. go test -trimpath ./internal/workflow/... passes with the review steps code, lints, lint-fix, review, render, signoff in that order, code's pass going to lints, lints' pass to review and fail to lint-fix, lint-fix's pass to review and max_visits 1, and ship-change's inner review steps including review/lints and review/lint-fix. describe-golden.txtar is regenerated with -update and go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden$' passes.

## statement

A new operation `rhiza run lint-check` (run.lint-check) exits 0 when the run's latest code-check entry has no lint hits and exits 1 listing each hit (node, lint, reference or Decision) otherwise, with the hits in its --json response. The builtin review workflow routes the code step's pass to a `lints` command step running it (on.pass review, on.fail lint-fix) and adds a `lint-fix` agent step, max_visits 1, prompted `/rhiza-fixer Fix the model lint hits that rhiza run lint-check lists for change set {change} before the checks run: follow rhiza guide fixer.`, until true, on.pass review; ship-change's inlined review gains both steps. Files: internal/api/run_lintcheck.go, internal/api/run_lintcheck_test.go, internal/workflow/builtin/review.toml, internal/workflow/workflow_test.go, cmd/rhiza/testdata/script/describe-golden.txtar.

## Links

### Out

- depends-on [code-check-runs-lints](rhiza/workflows/built-in-workflows/built-in-review/code-check-runs-lints.md)
- decided-by [lint-hits-feed-model-check](../decision/lint-hits-feed-model-check.md)

### In

- depends-on [fixer-lint-hits-guide](fixer-lint-hits-guide.md)
- verifies [review-lint-fix-route-verified](../evidence/review-lint-fix-route-verified.md)
- verifies [review-lint-fix-route-qa-1](../evidence/review-lint-fix-route-qa-1.md)
- verifies [review-lint-fix-route-verified-2](../evidence/review-lint-fix-route-verified-2.md)
- verifies [pr-569-model-lints-in-review](../evidence/pr-569-model-lints-in-review.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
