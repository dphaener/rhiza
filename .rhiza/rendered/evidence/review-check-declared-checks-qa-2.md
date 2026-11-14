# review-check-declared-checks-qa-2

evidence · `01M7DTKR23WZ3ZMWYGVRJNH419` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: go test ./internal/api -run TestReviewCheck passes (TestReviewCheckDeclaredChecks, TestReviewCheckConductor, TestReviewCheckRefusesStaleRound and TestReviewCheckFindingsFlag). The built binary's run review-check outside a run is refused naming RHIZA_RUN.

## Links

### Out

- verifies [review-check-declared-checks](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-check-declared-checks.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
