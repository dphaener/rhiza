# reopen-views

capability · `01M5WPTT5VX3GPH672QATF03N3` · parent [reopen](../reopen.md)

## acceptance

Go tests fold a landed, reopened, relanded change set and check as-is, to-be and the history at each step.

## statement

Views treat a reopened change set as accepted until it lands again: as-is leaves it out, to-be and what-if include it, and change show and log show the land, reopen and land history.

## Links

### Out

- decided-by [m3-scope](../../../../decision/m3-scope.md)

### In

- verifies [reopen-views-verified](../../../../evidence/reopen-views-verified.md)
- verifies [pr-35-reopening-a-change-set](../../../../evidence/pr-35-reopening-a-change-set.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
