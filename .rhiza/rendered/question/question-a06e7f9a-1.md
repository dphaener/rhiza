# question-a06e7f9a-1

question · `01M6WN114Z4VXB5DE9T1HFSXGM`

## session

a06e7f9a-f18c-4eab-837f-01a1b94e2e29

## statement

When `rhiza run resume --retry` restarts a failed leaves (build) step, which failed leaves should get a fresh round of MaxAttempts attempts? Today only a leaf whose own failure question has an answering Decision gets one, so a leaf that used up its attempts while blocked on another question stays failed, and the retry fails at once.

## Links

### In

- answers [retry-revives-every-failed-leaf](../decision/retry-revives-every-failed-leaf.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
