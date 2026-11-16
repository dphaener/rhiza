# projection-parity-property

capability · `01M7GFXD16Z7ENTKT3FH883EFJ` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath -count=1 -p 2 ./internal/api -run '^TestPropertyProjectionParity$' -v runs at least 100 draws in under 20 s, logs how many draws had two landed heads on one target, a rebuild and pending entries (and fails as vacuous if any count is zero, as index.TestPropertyIncremental does), and reports, on a disagreement, the shrunk log, the view, the node or link, the field (version, conflict, pending, a field name, presence) and each side's value. go test -trimpath -count=1 -p 2 ./internal/index -run '^TestPropertyIncremental$' still passes on the moved generator. Failure case: TestParityReportNamesNode feeds the comparison two models that differ in one node's field and expects the report to name the view, the node and both values. It is not a timing budget: make test runs it. Reuse the generator as @parity-generator-shared says. Touches internal/api/parity_test.go (new), internal/index/indextest (new, imported only by tests) and internal/index/property_test.go. Must not change non-test code other than the new indextest package; a disagreement it finds is not fixed here: it goes into the known-cases table of @parity-known-cases-table, each entry naming the log shape, which @parity-interim-from-strict and @parity-write-index-from-strict empty.

## statement

A rapid property test in internal/api, TestPropertyProjectionParity, draws random logs and checks that, for the as-is view and every change set's write view, the interim model (newModel), the index's write model (writeIndex.model, over an index under t.TempDir()) and a fresh strict fold (changeset.FoldLog, then internal/model's AsIs, or WhatIf for an in-flight change set) give the same nodes, fields, links, versions and conflicts by the rule of @parity-version-rule; its draws include two landed heads on one target, a rebase that drops entry files after the index has seen them (so the index rebuilds), and dry-run pending entries laid over the index

## Links

### Out

- decided-by [strict-engine-is-spec](../../../../decision/strict-engine-is-spec.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [parity-version-rule](../../../../decision/parity-version-rule.md)
- decided-by [parity-known-cases-table](../../../../decision/parity-known-cases-table.md)
- decided-by [parity-generator-shared](../../../../decision/parity-generator-shared.md)

### In

- depends-on [parity-interim-from-strict](parity-interim-from-strict.md)
- depends-on [parity-write-index-from-strict](parity-write-index-from-strict.md)
- verifies [projection-parity-property-verified](../../../../evidence/projection-parity-property-verified.md)
- verifies [projection-parity-property-qa-1](../../../../evidence/projection-parity-property-qa-1.md)
- verifies [projection-parity-property-qa-2](../../../../evidence/projection-parity-property-qa-2.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
