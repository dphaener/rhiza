# carve-dependency-verified

evidence · `01M5TVY01AZ9PNM8JFQMD0SS7W`

## result

PASS under make check at commit dea942b on branch m2-cs9-carve

## statement

Go tests and testscripts show a carve making the source depend on the target when ops left in the source, or in an in-flight change set depending on it directly or transitively, rest on a moved node or link; refusing a carve whose needed dependency would close a cycle; and no false dropped conflict for a dependent when a moved node still exists under the same ID, including when the carve leaves the source empty. Known limit: a dependent's op whose base names a moved op still gets a base conflict.

## tests

go test: TestCarve/ChildLeftBehind, TestCarve/DependentRests, TestCarve/TransitiveDependentRests, TestCarve/DependentRestsOnNothingMoved, TestCarveRefused/Cycle, TestCarveRefused/DependentCycle (internal/api), TestRevisedRecreated (internal/changeset); testscript: carve-children.txtar, carve-dependents.txtar (incl. emptied source E with dependent G's bkid under moved p), carve-guards.txtar

## Links

### Out

- verifies [carve-dependency](../capability/rhiza/change-sets/carve-dependency.md)

## History

- `01M5TS2C87Z5G8HX2AWAVSGGQN` Carving that keeps IDs: genesis carve moves nodes, with their writes and links, from one unlanded change set to another under the same IDs, so milestone work can be carved out of the area bundles.
