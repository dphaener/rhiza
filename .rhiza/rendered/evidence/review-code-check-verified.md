# review-code-check-verified

evidence · `01M5ZZSAEEGZ9SK41GMWVQ3S7K` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran at 1441e92 in a detached worktree, RHIZA_* unset. cmd/rhiza/testdata/script/review-code-check.txtar PASS: (1) pass with verifier Evidence and a clean model - after alpha and beta both have Evidence by actor verifier (builder wrote the commits) and the render is clean, run code-check prints 'code-check passed: rhiza check is ok and 2 leaves have Evidence by a verifier', exit 0, and as the run's command step journals step code outcome pass and the run succeeds; (2) fail naming the cause for missing Evidence - with only beta verified it prints 'code-check failed: 1 cause' and '- leaf alpha: no Evidence verifies it', no rhiza check line, exit 1, causes JSON array names it and the journal entry is step gate outcome fail with that detail; (3) fail naming the cause for a rhiza check problem - with a stale .rhiza/rendered/index.md it prints '- rhiza check: fail render .rhiza/rendered/index.md: differs from what render writes', no leaf line, exit 1, check ok false, rule render; both causes together are named as '2 causes'. (4) review-code-check-override.txtar PASS: an override review.toml with command steps code (rhiza run code-check) -> tests (echo ran-tests > tests.log) -> checks step review; journal order is entered code, code-check code (pass), entered tests, entered review, and the fake check session saw tests=ran-tests in RHIZA_ROOT, so the command ran after code-check and before the checks step. TestCodeCheckFold, TestCodeCheckCauses (warnings ignored, fails then leaves), build-check and describe-golden scripts and go test ./... all pass. Built rhiza-vc; describe run.code-check --json shows the empty request and response with run, change, step, passed, causes, check, leaves, dry_run.

## tests

commit 1441e92; run by claude-code-0e5a9d92-verifier on 2026-10-27; each prefixed env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT: go test ./cmd/rhiza -run 'TestScripts/review-code-check' -v -count=1; go test ./internal/run -run TestCodeCheckFold -v -count=1; go test ./internal/api -run TestCodeCheckCauses -v -count=1; go test ./cmd/rhiza -run 'TestScripts/(build-check|describe-golden)' -count=1; go test ./... -count=1; go build -o $TMPDIR/rhiza-vc ./cmd/rhiza && $TMPDIR/rhiza-vc describe run.code-check --json

## Links

### Out

- verifies [review-code-check](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-code-check.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
