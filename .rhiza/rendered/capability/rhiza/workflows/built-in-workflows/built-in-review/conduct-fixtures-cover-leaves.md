# conduct-fixtures-cover-leaves

capability · `01M7PTHCV4EMFSDRTDC7CF0Z6C` · parent [built-in-review](../built-in-review.md)

## acceptance

ship-change-one-ci.txtar passes with go test -trimpath ./cmd/rhiza -run 'TestScripts/ship-change-one-ci$', its fake claude in the built-in ship-change's conduct step recording each plan, code and qa round with an evidence line for the leaf first; (failure case) the script asserts that a clean round without evidence for first, recorded by hand in the conduct step before the covered round, fails with stderr naming the uncovered leaf first.

## files

cmd/rhiza/testdata/script/ship-change-one-ci.txtar

## must_not

Change Go code; drop or weaken an existing assertion of a script; touch scripts of classic checks steps or plan review.

## statement

The cmd/rhiza testscripts whose fake conductor checks record rounds with rhiza run findings --check in a conductor step pass an evidence line for every leaf their round does not name in a finding, so they pass under the leaf-coverage refusal.

## Links

### Out

- depends-on [check-round-leaf-coverage](check-round-leaf-coverage.md)
- decided-by [shallow-round-needs-leaf-coverage](../../../../../decision/shallow-round-needs-leaf-coverage.md)

### In

- verifies [conduct-fixtures-cover-leaves-verified](../../../../../evidence/conduct-fixtures-cover-leaves-verified.md)
- verifies [conduct-fixtures-cover-leaves-qa-1](../../../../../evidence/conduct-fixtures-cover-leaves-qa-1.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
