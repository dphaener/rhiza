# review-01m7pt83txevmbecpw6e6ee114-r3-f2

question · `01M7PV681EVKB1B7PMSEDR9K20`

## statement

Plan review round 3 left finding f2 (no-failure-case on check-round-evidence-status): The acceptance's '(failure case)' is 'a round with findings and no evidence lines prints its findings exactly as before', a regression check of the happy display path, not a failure path of the leaf. run status showing evidence has no input it refuses; the leaf needs either a real failure path (e.g. an evidence line naming a leaf no longer in the change set still prints with its ID and no handle, not a crash) or a failure_waiver giving that reason. The critic confirmed it: Real: the leaf's '(failure case)' is a round with findings and no evidence printing as before, which is the absent-data display path (already covered by its own 'no evidence key' clause), not a failure path of run status's evidence display; the leaf has no failure_waiver. A real failure path exists to test (an evidence line whose leaf no longer resolves, e.g. a leaf retired after the round, must still print with its ID rather than fail), so either that case or a waiver is needed. Major for no-failure-case stands. How should it be settled?

## Links

### Out

- blocks [check-round-evidence-status](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-evidence-status.md)

### In

- answers [evidence-status-unresolved-leaf-failure-case](../decision/evidence-status-unresolved-leaf-failure-case.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
