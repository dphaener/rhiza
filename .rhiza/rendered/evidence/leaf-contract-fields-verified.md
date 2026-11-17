# leaf-contract-fields-verified

evidence · `01M7KNFMCNXFN8F12E0EE9ZMCM` · actor [claude-code-efdcbc87](../actor/claude-code-efdcbc87.md)

## result

PASS

## statement

Verifier at 79dbbece4 re-ran the Checks: block: go test -trimpath -p 2 ./internal/op/... ./internal/run/... exit 0; go test -trimpath -p 2 ./internal/api/... exit 0; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/field-schema$' exit 0; rhiza ship generate exit 0 (0 changed); rhiza render exit 0 (only the verifier's own actor page). -run LeafContract -count=1 -v: TestLeafContractFields (capability files/must_not/failure_waiver known, constraint files unknown), TestBriefLeafContract (leaf a lists Files:/Must not:/Failure waiver: under Acceptance, leaf b none), TestSessionBriefLeafContract (add capability with files and must_not reads both back and reaches brief.md; field file refused exit 2 naming known fields acceptance, built_in, failure_waiver, files, must_not, statement) all PASS. Diff touches only fields.go, brief.go, run_brief.go and their tests, plus field-schema.txtar (its expected known-field list, required by this change) and rendered actor pages for the builder and the conductor (whose actor record f881203ca committed unrendered), both in scope. Statement/acceptance reading, existing nodes and review-check are unchanged; no stubs or TODOs.

## Links

### Out

- verifies [leaf-contract-fields](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-fields.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
