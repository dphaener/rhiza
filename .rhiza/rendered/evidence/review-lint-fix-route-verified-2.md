# review-lint-fix-route-verified-2

evidence · `01M7BAFMBXRYEGNFGKREMDX2G1` · actor [claude-code-903034bd](../actor/claude-code-903034bd.md)

## result

PASS

## statement

Attempt 2 at 20d24dc86 (fix for quality finding f1: review-workflow.txtar expected code then review). The fix commit only changes cmd/rhiza/testdata/script/review-workflow.txtar to assert 'code ... then lints' and 'lints ... then review' (plus rendered/log records). Re-ran its Checks: go test -trimpath ./cmd/rhiza -run 'TestScripts/review-workflow$' exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (only this verifier's own actor file and index). Acceptance: go test -trimpath ./internal/api/... -run 'LintCheck' exit 0 (TestLintCheck covers no code-check entry exit 0, clean entry exit 0, two lints exit 1 with both hits in JSON and human output, later clean entry supersedes); go test -trimpath ./internal/workflow/... exit 0 (ship-change inner steps include review/lints and review/lint-fix in order); go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden$' exit 0. review.toml routes code pass->lints, lints pass->review fail->lint-fix, lint-fix agent max_visits 1 with the specified prompt, until true, pass->review. No TODOs or stubs in the diff.

## Links

### Out

- verifies [review-lint-fix-route](../capability/review-lint-fix-route.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
