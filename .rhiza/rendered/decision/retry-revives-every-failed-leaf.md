# retry-revives-every-failed-leaf

decision · `01M6WN3EQ24GGSKCPVWF623BKV`

## reasoning

A retry is an explicit request from a person to try the failed step again, so a leaf that is still out of attempts makes the retry fail at once and do nothing. That is exactly the bug this change set's intent reports: the leaf bugfix-workflow-autonomous spent its attempts while blocked on bugfix-needs-workflow-dispatch, that question was answered elsewhere and main was merged in, and the retry still failed with no builder started. Today reviveLeaf only gives a fresh round when an answering Decision exists on the leaf's own failure question. The intent asks for a fresh budget 'at least' when a blocker was answered or the branch moved. Trying to work out which question really blocked a leaf, or comparing branch heads, would be brittle and would still miss cases. Making the retry itself the trigger is simple and predictable, and it costs little: each retry adds at most one bounded round of MaxAttempts per failed leaf, and each round needs a person to run another retry. The autonomous rule (arbitrate a leaf once) stays as it is, so a retry never starts arbitration loops.

## statement

When `rhiza run resume --retry` restarts a failed leaves step, every leaf of that step that failed gets one fresh round of MaxAttempts attempts, and the step's skipped leaves reopen with them. A failed leaf whose own question a Decision answers by then gets that Decision's round, as on any resume, with the Decision in its builder's prompt; every other failed leaf gets a round marked as given by the retry, not by a Decision, so no leaf gets two rounds from one retry. The retry's round does not wait on the leaf's own question, or some other blocking question, being answered, or on the branch having moved. A failure in a round the retry gave is not arbitrated again, the same as any other round after the first.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a06e7f9a-1](../question/question-a06e7f9a-1.md)
- derived-from [retry-revives-leaves-plan](../source/retry-revives-leaves-plan.md)

### In

- decided-by [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
