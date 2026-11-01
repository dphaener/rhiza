# diff-against-op-base

decision · `01M6BQK4CT3NYGBWD74WAJG3Z6`

## statement

A change set's diff is built from its folded ops, one entry per node merging all its writes in the change set: old values are the node at the first op's base (replaying the node's history), new values the change set's last write; created nodes have no old values, retires show the reason, links added and retired are listed; string fields are unquoted and multi-line ones shown as a line diff from a small in-repo LCS, other values as compact JSON; change show gains a Changes section so the change-show app shows it

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m16-orchestrator-picks](../source/m16-orchestrator-picks.md)

### In

- decided-by [change-diff-op](../capability/rhiza/change-sets/change-diff/change-diff-op.md)
- decided-by [change-diff-text](../capability/rhiza/change-sets/change-diff/change-diff-text.md)
- decided-by [change-show-diff](../capability/rhiza/change-sets/change-diff/change-show-diff.md)

## History

- `01M6BQJPJPWB8CZVCKAKWWC9T7` Plan milestone 16: see what a change set changes, picked by the orchestrator under the owner's standing instruction
