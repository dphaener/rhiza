# no-failure-case-kind

capability · `01M7KJ4SZXSN16GPF6JZ0CEBHM` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/run/... passes with tests that run.Kinds ends with oversized-leaf then no-failure-case and KindSeverity returns major for no-failure-case; go test -trimpath ./internal/api/... passes with a round holding a confirmed no-failure-case finding failing review-check before round 3, and (failure case) a finding of an unknown kind such as no-failure still refused with an error listing the eight kinds; go test -trimpath ./internal/guide/... passes with the severities block naming no-failure-case major; go test -trimpath ./cmd/rhiza -run 'TestScripts/(describe-golden|run-review-journal)' passes after regenerating describe-golden with -update.

## statement

`rhiza run findings` accepts an eighth plan-review kind, `no-failure-case` (a leaf whose acceptance has no real failure case and no failure_waiver with a reason), listed after oversized-leaf in run.Kinds, and KindSeverity, run review-check's summary, the findings input summary, the shared severities block and the generated reference docs name it major. Files: internal/run/journal.go, internal/api/run_review.go, internal/guide/text/shared/severities.md, cmd/rhiza/testdata/script/describe-golden.txtar, cmd/rhiza/testdata/script/run-review-journal.txtar, docs/reference (regenerated with rhiza ship generate) and the tests naming the kinds (internal/run/journal_test.go, internal/guide/review_test.go, internal/api/run_round_recorded_test.go). Must not: change the guide topics' text other than severities.md, or add a leaf-field check.

## Links

### Out

- decided-by [no-failure-case-major](../../../../../decision/no-failure-case-major.md)

### In

- depends-on [plan-reviewer-leaf-contract](plan-reviewer-leaf-contract.md)
- verifies [no-failure-case-kind-verified](../../../../../evidence/no-failure-case-kind-verified.md)
- verifies [no-failure-case-kind-qa-1](../../../../../evidence/no-failure-case-kind-qa-1.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
