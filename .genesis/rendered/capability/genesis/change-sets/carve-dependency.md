# carve-dependency

capability · `01M5TS2CVJC3TH6EEMA75TNXDF` · parent [change-sets](../change-sets.md)

## statement

A carve makes the source depend on the target when the source still rests on a moved node, and refuses a carve that would close a cycle.

## Links

### Out

- decided-by [carve-adds-dependency](../../../decision/carve-adds-dependency.md)

### In

- relates-to [carve-stale-base-recovery](../../../question/carve-stale-base-recovery.md) — Open follow-up on the known limit of carve-dependency: a dependent's op with a base on a moved op gets a base conflict after a carve. Parked, outside change set 9's intent.
- verifies [carve-dependency-verified](../../../evidence/carve-dependency-verified.md)

## History

- `01M5TS2C87Z5G8HX2AWAVSGGQN` Carving that keeps IDs: genesis carve moves nodes, with their writes and links, from one unlanded change set to another under the same IDs, so milestone work can be carved out of the area bundles.
