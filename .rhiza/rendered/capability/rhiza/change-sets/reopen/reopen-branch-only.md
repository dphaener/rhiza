# reopen-branch-only

capability · `01M5WPTSZ2FDZ6T77CRE763WAE` · parent [reopen](../reopen.md)

## acceptance

A testscript where the land entry is in the base gets a non-zero exit and that message; one where it is only on the branch succeeds.

## statement

Reopening is refused once the change set's land entry is on the base (the merge-base with the default branch's remote, as check uses): landed work on main is final, and the error says so and points to a new change set linked discovered-from.

## Links

### Out

- decided-by [m3-scope](../../../../decision/m3-scope.md)
- decided-by [reopen-any-land-on-base](../../../../decision/reopen-any-land-on-base.md)

### In

- verifies [reopen-branch-only-verified](../../../../evidence/reopen-branch-only-verified.md)
- verifies [pr-35-reopening-a-change-set](../../../../evidence/pr-35-reopening-a-change-set.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
