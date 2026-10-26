# reopen

capability · `01M5WPTSR78R8WM8WYMGWD6WJN` · parent [change-sets](../change-sets.md)

## acceptance

A testscript lands a change set on a branch, reopens it, applies one more write, lands again, and rhiza check passes; the log holds land, reopen and land entries in order.

## statement

rhiza reopen reopens a landed change set whose land is not yet on the base: it appends a reopen entry, the change set is accepted again, takes further writes, and lands again with rhiza land. No log entry is ever deleted.

## Links

### Out

- decided-by [m3-scope](../../../decision/m3-scope.md)
- decided-by [reopen-makes-current](../../../decision/reopen-makes-current.md)
- decided-by [reopen-empty-body](../../../decision/reopen-empty-body.md)

### In

- verifies [reopen-verified](../../../evidence/reopen-verified.md)
- verifies [pr-35-reopening-a-change-set](../../../evidence/pr-35-reopening-a-change-set.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
