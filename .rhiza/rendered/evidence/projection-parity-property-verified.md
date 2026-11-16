# projection-parity-property-verified

evidence · `01M7GKQWSM12P391CFPQHEWF7J` · actor [claude-code-ec3dd358](../actor/claude-code-ec3dd358.md)

## result

PASS

## statement

In a detached checkout at c63ce1320, re-ran every Checks: command with exit 0: TestPropertyProjectionParity -v passed 100 draws in 6.6 s and logged two landed heads 47, rebuilds 76, pending 41, views 3926 plus each known case's count; TestParityReportNamesNode passed; -rapid.checks=3000 passed (162 s); internal/index TestPropertyIncremental -v passed (lands 141, accepts 167, rebuilds 128) and ./internal/index/... passed; go vet ./internal/api ./internal/index/... clean; rhiza ship generate and rhiza render left the tree clean. Read the diff: the test fails as vacuous (t.Errorf) when any of two landed heads, rebuilds or pending is zero; no draw is skipped (no Skip, the pending phase only returns early); a disagreement Fatalf prints view, node or link, field and both sides' values plus the drawn log, which rapid shrinks; TestParityReportNamesNode asserts the view, node, field and both values; every knownCases entry names its log shape; the generator and steps live in internal/index/indextest, used by both internal/index/property_test.go and internal/api/parity_test.go, which are its only importers (both tests), and indextest does not import internal/index; the only non-test Go code added is indextest; no stubs or TODOs.

## Links

### Out

- verifies [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
