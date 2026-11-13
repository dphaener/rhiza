# review-lint-fix-route-verified

evidence · `01M7B7T0KFXF70W61KPP1Q0FHF` · actor [claude-code-a9264b25](../actor/claude-code-a9264b25.md)

## result

PASS

## statement

At a970d1bf5 re-ran the Checks block: go test -trimpath ./internal/api/... -run 'LintCheck' (TestLintCheck PASS: no code-check entry and an entry with no lints exit 0, two lints exit 1 with both hits in the response and human output, a later clean entry passes), go test -trimpath ./internal/workflow/... (ok: review order code, lints, lint-fix, review, render, signoff; code pass->lints, lints pass->review, fail->lint-fix, lint-fix max_visits 1 pass->review with the specified fixer prompt; ship-change includes review/lints and review/lint-fix), go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden$' (ok, golden has run.lint-check), rhiza ship generate (0 files changed) and rhiza render: all exit 0. Diff read: only the leaf's files plus generated docs/reference/cli.md and the actor record; no TODOs, gofmt clean.

## Links

### Out

- verifies [review-lint-fix-route](../capability/review-lint-fix-route.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
