# pr-377-fix-rhiza-run-resume-retry

evidence · `01M6XBNGE6S43AKCJYC7B2SCNQ` · actor [dphaener](../actor/dphaener.md)

## head

8d7c0b2b37ca5a45e36f3f1ef9078b4876165bbc

## pr

https://github.com/dphaener/rhiza/pull/377

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #377 (https://github.com/dphaener/rhiza/pull/377) passed CI at 8d7c0b2

## status

verified

## Links

### Out

- verifies [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
