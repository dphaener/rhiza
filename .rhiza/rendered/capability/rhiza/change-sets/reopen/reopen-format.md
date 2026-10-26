# reopen-format

capability · `01M5WPTTCK57BZS99AG7M94C6X` · parent [reopen](../reopen.md)

## acceptance

The format-1 and format-2 corpus pass unchanged, and a new corpus case with a reopen entry reads and renders.

## statement

The reopen entry is a new entry type. Format 1 and 2 entries stay byte-identical; if the type needs a new format number, older readers fail with ErrNewerFormat, as carve's did.

## Links

### Out

- decided-by [m3-scope](../../../../decision/m3-scope.md)
- decided-by [reopen-format-three](../../../../decision/reopen-format-three.md)

### In

- verifies [reopen-format-verified](../../../../evidence/reopen-format-verified.md)
- verifies [pr-35-reopening-a-change-set](../../../../evidence/pr-35-reopening-a-change-set.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
