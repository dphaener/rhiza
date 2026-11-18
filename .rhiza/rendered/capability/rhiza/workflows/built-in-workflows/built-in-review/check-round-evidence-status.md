# check-round-evidence-status

capability · `01M7PTHCQDJX30S30ZZT05GDPJ` · parent [built-in-review](../built-in-review.md)

## acceptance

go test -trimpath ./internal/api -run 'RunStatus|CheckRound' passes with a new test in internal/api/run_test.go: a conductor run whose check journaled a round with two evidence lines shows both, with leaf handles, in run status --json and in the human output under that round; a round with no evidence has no evidence key in --json and prints no evidence lines, and a round with findings and no evidence lines prints its findings exactly as before; (failure case) an evidence line whose leaf no longer resolves (a leaf retired after the round was journaled) is still shown, with the leaf's ID and an empty leaf_handle in --json and the bare ID in the human output, and run status does not fail. go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after regenerating it with -update, and make self-check passes.

## files

internal/api/run.go, internal/api/run_test.go, cmd/rhiza/testdata/script/describe-golden.txtar, docs/reference/cli.md

## must_not

Change how findings, verdicts or plan-review rounds are shown; change run findings' validation.

## statement

rhiza run status shows each review check round's evidence lines: --json gives every RunCheckRound an evidence array of {leaf, leaf_handle, line}, absent when the round has none, and the human output lists them under the round with the leaf's handle.

## Links

### Out

- depends-on [check-round-leaf-coverage](check-round-leaf-coverage.md)
- decided-by [shallow-round-needs-leaf-coverage](../../../../../decision/shallow-round-needs-leaf-coverage.md)
- decided-by [evidence-status-unresolved-leaf-failure-case](../../../../../decision/evidence-status-unresolved-leaf-failure-case.md)

### In

- blocks [review-01m7pt83txevmbecpw6e6ee114-r3-f2](../../../../../question/review-01m7pt83txevmbecpw6e6ee114-r3-f2.md)
- verifies [check-round-evidence-status-verified](../../../../../evidence/check-round-evidence-status-verified.md)
- verifies [check-round-evidence-status-qa-1](../../../../../evidence/check-round-evidence-status-qa-1.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
