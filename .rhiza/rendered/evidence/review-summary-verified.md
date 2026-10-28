# review-summary-verified

evidence · `01M601PSNVMXQG2VXTNF5HWZ4V` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran the leaf's tests at 995de72 in a detached worktree with RHIZA_* unset. review-summary.txtar passed 3 of 3 runs. Every clause is asserted: check alpha passes in cycle 1; check beta fails in cycle 1 (its critic confirms major f1 on leaf thing and drops f2), the fixer and verifier fix it in cycle 2, and beta passes on its single re-run; there are two minor-severity findings (alpha f1 at README.md:3 and beta f3 at thing.txt:2) plus one nit; `cmp stdout $WORK/summary.golden` compares the summary with a golden file, both by the latest run and by ref; at the sign-off wait, `rhiza run status` matches '^next: rhiza run summary [0-9A-Z]{26}$' (RunInfo.Human prints it beside approve and request-changes). TestReviewSummary and TestOwnerCycleFixerPrompts pass, as do review-signoff* and describe-golden, and go test ./... passes in every package.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./cmd/rhiza -run 'TestScripts/review-summary' -v -count=3; env -u ... go test ./internal/api -run 'TestReviewSummary|TestOwnerCycleFixerPrompts' -v -count=1; env -u ... go test ./cmd/rhiza -run 'TestScripts/(review-signoff|describe-golden)' -count=1; env -u ... go test ./... -count=1; commit 995de72; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-summary](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-summary.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
