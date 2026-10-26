# reopen-views-verified

evidence · `01M5WT1Q1CNJW3MX4Z80365DDM`

## result

PASS under make check at commit afc1020 on branch cs/reopening-a-change-set, run by an independent verifier

## statement

Go tests fold landed, reopened and relanded change sets (including skipped and concurrent reopens), check as-is and to-be at each step, and show the land and reopen history in change show and log.

## tests

go test internal/changeset: TestFoldReopen, TestFoldReopenSkipped, TestFoldConcurrentReopen; internal/model: TestViewsReopen (landed, reopened, reopened_with_an_op, relanded); internal/api: TestChangeShowHistory, TestLogReopened, TestLogReopenedAccepted

## Links

### Out

- verifies [reopen-views](../capability/rhiza/change-sets/reopen/reopen-views.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
