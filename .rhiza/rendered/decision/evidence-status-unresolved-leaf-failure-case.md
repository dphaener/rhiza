# evidence-status-unresolved-leaf-failure-case

decision · `01M7PV84S4P4F8W3T2H4ND8PQP`

## reasoning

The finding and the critic agree the leaf's stated failure case is the absent-data display path, already covered by its 'no evidence key' clause. A real failure path exists: rounds are journaled before later plan edits, so an evidence line can name a leaf retired or removed since, and the handle lookup must degrade to the ID rather than error. Testing that is cheap (one more evidence line in the same run_test.go fixture), stays inside the leaf's files and must_not, and is more useful than a waiver, so the case is chosen over a failure_waiver.

## statement

check-round-evidence-status's failure case is an evidence line whose leaf no longer resolves (e.g. a leaf retired after the round was journaled): run status still shows that line, with the leaf's ID and an empty leaf_handle in --json and the bare ID in the human output, and does not fail. The current '(failure case)' (a round with findings and no evidence printing as before) stays as a regression clause, not as the failure case; no failure_waiver.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7pt83txevmbecpw6e6ee114-r3-f2](../question/review-01m7pt83txevmbecpw6e6ee114-r3-f2.md)

### In

- decided-by [check-round-evidence-status](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-evidence-status.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
