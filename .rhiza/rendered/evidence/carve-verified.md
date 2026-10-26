# carve-verified

evidence · `01M5TVXZN94BWJAX70MNQXQB75`

## result

PASS under make check at commit 881d1fb on branch m2-cs9-carve

## statement

Go tests and testscripts show a carve moving nodes and their links from one change set to another under the same IDs, rewriting bases onto the new op IDs, moving children with --with-children, and leaving children behind or skipping the dependency when the target already depends on the source.

## tests

go test: TestCarve (internal/api: KeepsIDs, BasesRewritten, WithChildren, ChildLeftBehind, AlreadyDependsOnTarget); testscript: carve.txtar

## Links

### Out

- verifies [carve](../capability/rhiza/change-sets/carve.md)

## History

- `01M5TS2C87Z5G8HX2AWAVSGGQN` Carving that keeps IDs: genesis carve moves nodes, with their writes and links, from one unlanded change set to another under the same IDs, so milestone work can be carved out of the area bundles.
