# leaf-contract-fields

capability · `01M7KJ4SXVPYE91ZCAQASJ6E4C` · parent [plan-workflow](../plan-workflow.md)

## acceptance

go test -trimpath ./internal/op/... passes with FieldKnown true for capability files, must_not and failure_waiver and (failure case) false for constraint files; go test -trimpath ./internal/run/... passes with a brief test whose leaf with all three fields lists the Files:, Must not: and Failure waiver: lines under its acceptance and whose leaf with none lists none of them; go test -trimpath ./internal/api/... passes with a test that adds a capability with files and must_not and reads both back, and (failure case) a capability write with a field `file` refused with an error whose known fields name files, must_not and failure_waiver.

## statement

A capability may carry the fields `files`, `must_not` and `failure_waiver` (rhiza add, set and apply accept them on a capability and still refuse them on other kinds), and brief.md lists, under each leaf's acceptance, a `Files:` line, a `Must not:` line and a `Failure waiver:` line for each field the leaf has, leaving out a field it lacks. Files: internal/op/fields.go, internal/run/brief.go, internal/api/run_brief.go and their tests (internal/op/fields_test.go, internal/run/brief_test.go, the internal/api brief test). Must not: change how statement or acceptance read, refuse nodes already in the log, or change review-check.

## Links

### Out

- decided-by [leaf-contract-as-fields](../../../../../decision/leaf-contract-as-fields.md)

### In

- depends-on [leaf-contract-review-check](leaf-contract-review-check.md)
- depends-on [planner-leaf-contract](planner-leaf-contract.md)
- verifies [leaf-contract-fields-verified](../../../../../evidence/leaf-contract-fields-verified.md)
- verifies [leaf-contract-fields-qa-1](../../../../../evidence/leaf-contract-fields-qa-1.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
