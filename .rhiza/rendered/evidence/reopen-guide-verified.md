# reopen-guide-verified

evidence · `01M5WT1Q8PAAARWE2WRXXPQGAP`

## result

PASS under make check at commit afc1020 on branch cs/reopening-a-change-set, run by an independent verifier

## statement

The guide tests pass within the size caps, the guide testscripts pass, and rhiza guide reopening names rhiza reopen.

## tests

go test internal/guide: TestCaps, TestCapCatchesOversize, TestListOrder, TestEveryFileListed, TestLookup, TestTokens, TestDescriptions, TestList; testscript: guide.txtar, guide-install.txtar; rhiza guide reopening prints 'rhiza reopen --change <id> [--base <rev>] --json'

## Links

### Out

- verifies [reopen-guide](../capability/rhiza/change-sets/reopen/reopen-guide.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
