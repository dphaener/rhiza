# parity-generator-shared

decision · `01M7GGES0NJ07WCC84WMZ4X3WD`

## question

Does the parity test copy internal/index's log generator or share it?

## reasoning

genFiles draws the files the parity test needs (several change sets through the lifecycle, bases from a target's ops, corrupt and stray files), but the removals that force a rebuild are drawn in TestPropertyIncremental's step loop (property_test.go:263-273) and counted from the unexported x.stats.rebuilt (property_test.go:307-309), so moving genFiles alone would leave the parity test without the rebuild draws and count its acceptance requires: the step drawing moves too, and the parity test reads the exported index.LastCatchUp().Rebuilt. A copy would drift from the index's own property test. A _test.go file cannot be imported from another package, so sharing needs a small non-test package that only tests import; moving only the generator and the step drawing keeps the churn to property_test.go.

## statement

The log generator of internal/index/property_test.go (genFiles) and the part of TestPropertyIncremental's step loop that draws each step (which files to add, small steps half the time, and which already-written files to remove, at most two removals per draw) move to a test-only helper package (internal/index/indextest), which internal/index's property test and internal/api's parity test both import; neither is copied. indextest writes and removes the files but does not import internal/index (property_test.go is in package index); each test catches its index up after each step itself, internal/index's through its unexported catchUp and stats as today, the parity test through Refresh, counting rebuilds from index.LastCatchUp().Rebuilt. Other internal/index tests keep their own fixtures.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [parity-harness-plan](../source/parity-harness-plan.md)

### In

- decided-by [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
