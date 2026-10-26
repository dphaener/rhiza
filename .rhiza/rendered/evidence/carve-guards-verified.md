# carve-guards-verified

evidence · `01M5TVY0797CG7BPYJVP6A8C0N`

## result

PASS under make check at commit 881d1fb on branch m2-cs9-carve

## statement

Go tests and a testscript show carve refusing what it cannot move cleanly and a dry run predicting the real run's entries and IDs while writing nothing.

## tests

go test: TestCarveRefused, TestCarveDryRun (internal/api); testscript: carve-guards.txtar

## Links

### Out

- verifies [carve-guards](../capability/genesis/change-sets/carve-guards.md)

## History

- `01M5TS2C87Z5G8HX2AWAVSGGQN` Carving that keeps IDs: genesis carve moves nodes, with their writes and links, from one unlanded change set to another under the same IDs, so milestone work can be carved out of the area bundles.
