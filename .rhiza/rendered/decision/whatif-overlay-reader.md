# whatif-overlay-reader

decision · `01M6C8EJCQZSCA7RK5SZNP0BRF`

## statement

A what-if reader answers each query from the index's as-is rows, overriding them with an overlay of the records the change set and its in-flight closure create, update, move or retire (and links they add or retire), so it never loads the whole as-is into memory; its answers equal the full model's

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- relates-to [t7-d14-read-budget](t7-d14-read-budget.md) — lifts t7-d14's what-if exemption: with the overlay reader, every budgeted read (show, ls first page, deep-subtree tree, why, log, roadmap, ready) is budgeted in what-if as well as as-is and to-be, under the same 50 ms median; the whole-tree tree stays benchmarked, not budgeted
- derived-from [m20-orchestrator-picks](../source/m20-orchestrator-picks.md)

### In

- decided-by [whatif-reads-in-budget](../capability/rhiza/views/fast-reads/whatif-reads-in-budget.md)

## History

- `01M6C8EBWCSJZC8TGSN36SZ3QR` Plan milestone 20: what-if reads meet the budget, picked by the orchestrator under the owner's standing instruction
