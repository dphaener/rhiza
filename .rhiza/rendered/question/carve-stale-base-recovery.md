# carve-stale-base-recovery

question · `01M5TVY0D81YHQ12N5D3HCB4QP`

## statement

After a carve, how does a change set that depends on the source clear an op whose base names a moved op? Writing it again adds a second op and the old one stays; no plain command drops it.

## Links

### Out

- relates-to [carve-dependency](../capability/rhiza/change-sets/carve-dependency.md) — Open follow-up on the known limit of carve-dependency: a dependent's op with a base on a moved op gets a base conflict after a carve. Parked, outside change set 9's intent.

## History

- `01M5TS2C87Z5G8HX2AWAVSGGQN` Carving that keeps IDs: genesis carve moves nodes, with their writes and links, from one unlanded change set to another under the same IDs, so milestone work can be carved out of the area bundles.
