# change-diff-op

capability · `01M6BQKGB62F75TG9YPKM13QNT` · parent [change-diff](../change-diff.md)

## acceptance

A change-diff.txtar: add then set gives statement a to b; an unset field shows old to unset; a retire shows its reason; a create shows every field with no old value; two sets of one node in a change set give one before and after; with another change set landing a different field of the node first, the diff still shows this change set's own field against its base; change diff --json and api call change.diff print the same bytes; describe-golden.txtar is regenerated

## statement

rhiza change diff [ref] lists, per node the change set writes, its action (created, updated, moved or retired), handle and kind, each changed field's old and new value taken from the node at the first op's base, and the links it adds and retires

## Links

### Out

- decided-by [diff-against-op-base](../../../../decision/diff-against-op-base.md)

### In

- verifies [change-diff-op-verified](../../../../evidence/change-diff-op-verified.md)
- verifies [change-diff-op-qa-1](../../../../evidence/change-diff-op-qa-1.md)
- verifies [pr-179-field-level-diff](../../../../evidence/pr-179-field-level-diff.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
