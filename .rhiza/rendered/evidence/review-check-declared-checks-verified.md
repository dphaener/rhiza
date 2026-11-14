# review-check-declared-checks-verified

evidence · `01M7C0QCP6WQVVY9W605QRXQNB` · actor [claude-code-9ddaccaf](../actor/claude-code-9ddaccaf.md)

## result

PASS

## statement

At a79fc39d8 re-ran the Checks block: go test -trimpath ./internal/api -run 'TestReviewCheckDeclaredChecks|TestReviewCheckRefusesStaleRound|TestReviewCheckFindingsFlag|TestReviewCheckConductor' exit 0 (all four PASS), rhiza ship generate exit 0 (0 changed), rhiza render exit 0. TestReviewCheckDeclaredChecks builds a run at step gate past an agent step conduct declaring plan, code and qa: missing rounds fail, minor-only rounds pass, a confirmed major of code fails with exit 1 and its reason in the message, a downgrade to minor and a drop each pass. TestReviewCheckRefusesStaleRound and TestReviewCheckFindingsFlag drive plan-workflow runs and still pass. The diff adds only that test (plus actor/render files); conductorReviewCheck in internal/api/run_review.go is already on main, and the answered question review-check-declared-checks-follows-main limits the build to the missing test and keeps the missing-round failure. No stubs or TODOs; fakes only.

## Links

### Out

- verifies [review-check-declared-checks](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-check-declared-checks.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
