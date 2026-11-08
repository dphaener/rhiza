# retry-revives-leaves-plan

source · `01M6WN52NAA4P2EYS63YNDRJYA`

## quotes

```json
[]
```

## raw_sha256

4789f3b0243b90223f6e461604dbe316d361dbba0c5e79e32d6fb515bd69033b

## session

a06e7f9a-f18c-4eab-837f-01a1b94e2e29

## summary

Autonomous plan run for change set 01M6WMVHAJBBHR9SFTRPZD8CHC (bug). Reading internal/run showed the cause: a retry keeps a leaves step's leaf states, and on entry reviveLeaves gives a fresh round only to a failed leaf whose own failure question has an answering Decision (in an autonomous run, only in round 1, through the arbiter). The leaf bugfix-workflow-autonomous had spent both rounds while blocked on a different question raised by its builder, so the retry found nothing to start and failed at once. The planner asked whether a retry should give every failed leaf a fresh round, or only when a blocker was answered or the branch moved; there was no owner, and the arbiter ruled every failed leaf (Decision retry-revives-every-failed-leaf). The planner linked the change set to milestone 25 (the daemon ships the backlog on its own), since a retry that cannot rebuild forces an abandon and a new run by hand. Planned one leaf under build-workflow: a retry journals a leaf-round marked as the retry's for each failed leaf, its skipped dependents reopen, the builder prompt and run status name the retry round, and a failure in it is not arbitrated again.

## title

Plan: rhiza run resume --retry rebuilds failed leaves

## Links

### In

- derived-from [retry-revives-every-failed-leaf](../decision/retry-revives-every-failed-leaf.md)
- derived-from [retry-revives-failed-leaves](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-revives-failed-leaves.md)

## History

- `01M6WMVHAJBBHR9SFTRPZD8CHC` Fix rhiza run resume --retry cannot rebuild a leaf that used up its attempts: after its blocker was fixed, retrying ship-change run 01M6W86NY11Z (change set 01M6VY7H9PYX, autonomy per workflow) re-ran build/build at 03:49Z, which failed at once ('5 of 6 leaves passed; not passed: bugfix-workflow-autonomous') with no new builder attempt, because the leaf's 6 attempts were spent while it was blocked on workflow dispatch (question bugfix-needs-workflow-dispatch, since answered and main merged in). The only way on was rhiza run abandon plus a new run. A retry of a build step should give each failed leaf a fresh attempt budget, at least when the blocking question has been answered or the branch has moved since the last attempt.
