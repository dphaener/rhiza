# carve-entries-verified

evidence · `01M5TVXZVAQK14WB2NRTS53719`

## result

PASS under make check at commit 881d1fb on branch m2-cs9-carve

## statement

Go tests, a testscript and a corpus case show a carve written as two format-2 entries that name each other, the carve field round-tripping and rejected when malformed or in a format-1 entry, and genesis log showing what was carved from and into.

## tests

go test: TestLogCarve (internal/api), TestCarveRoundTrip, TestCarveInvalid, TestCarveNeedsFormat2 (internal/changeset), corpus case format-2/carve; testscript: carve.txtar

## Links

### Out

- verifies [carve-entries](../capability/genesis/change-sets/carve-entries.md)

## History

- `01M5TS2C87Z5G8HX2AWAVSGGQN` Carving that keeps IDs: genesis carve moves nodes, with their writes and links, from one unlanded change set to another under the same IDs, so milestone work can be carved out of the area bundles.
