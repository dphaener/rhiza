# default-view-state-from-index

decision · `01M6C3Q2W4B4M5WRCH5TV2YM9J`

## statement

When a read has no view and nothing is pending in the workspace, openReader takes the effective change set's state from the index's change set table instead of folding the log; with pending writes it keeps the fold, so unflushed state is never missed; the view rule is unchanged

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m19-orchestrator-picks](../source/m19-orchestrator-picks.md)

## History

- `01M6C3PTXZA57M17Y2Y9DX6NW9` Plan milestone 19: default reads stay fast, picked by the orchestrator under the owner's standing instruction
