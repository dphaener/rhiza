# revised-in-place-conflicts

decision · `01M5T2DHESV7R4G8513BKASNHS`

## statement

When a dependency revises an op in place and the revision changes a field that a dependent's op on the same node also writes, or drops a node or link the dependent's ops rest on, that is a conflict: check reports it, land refuses it with next, and ready gives it as the reason; no overlap, no conflict.

## wording

the orchestrator's amendment to change set 1, under the owner's authorization to amend plans

## Links

### Out

- relates-to [manual-rebase-for-now](manual-rebase-for-now.md) — applies manual-rebase-for-now to a dependency revised in place: the conflict shows only where fields overlap

### In

- decided-by [conflicts](../capability/genesis/change-sets/conflicts.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
