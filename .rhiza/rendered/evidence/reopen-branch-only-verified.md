# reopen-branch-only-verified

evidence · `01M5WT1PXSTBVVDZ7TM7G8CRV6` · actor [claude-code-b5610794](../actor/claude-code-b5610794.md)

## result

PASS under make check at commit afc1020 on branch cs/reopening-a-change-set, run by an independent verifier

## statement

A testscript and Go tests show reopen succeeding when the land entry is only on the branch, refused with a conflict when a land entry is on the base or there is no git or base, refused from every state but landed, and no writes to a landed set.

## tests

testscript: reopen-base.txtar (TestScripts/reopen-base); go test internal/api: TestReopenOnBranch, TestReopenLandedOnBase, TestReopenNoBase (NoGit, NoRef), TestReopenIllegal (draft, accepted, rejected, withdrawn), TestReopenNoWritesToLanded

## Links

### Out

- verifies [reopen-branch-only](../capability/rhiza/change-sets/reopen/reopen-branch-only.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
