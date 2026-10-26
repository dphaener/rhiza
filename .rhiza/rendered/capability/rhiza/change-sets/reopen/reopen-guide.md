# reopen-guide

capability · `01M5WPTTKD0TMJ3RHTAEV4SVEE` · parent [reopen](../reopen.md)

## acceptance

rhiza guide change-sets (or a topic of its own if change-sets is at its size cap) names rhiza reopen, and the guide tests pass.

## statement

The guide tells agents to reopen a change set, never to delete or edit log entries, when they must write after landing on a branch.

## Links

### Out

- decided-by [m3-scope](../../../../decision/m3-scope.md)

### In

- verifies [reopen-guide-verified](../../../../evidence/reopen-guide-verified.md)
- verifies [pr-35-reopening-a-change-set](../../../../evidence/pr-35-reopening-a-change-set.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
