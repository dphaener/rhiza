# conduct-fixtures-cover-leaves-verified

evidence · `01M7PY5AYREWX251QMMTX5130Z` · actor [claude-code-466e0ef5](../actor/claude-code-466e0ef5.md)

## result

PASS

## statement

At 4d0d6e575, go test -trimpath -p 2 -count=1 ./cmd/rhiza -run 'TestScripts/ship-change-one-ci$' (RHIZA_/SHIP_/REX_ env unset) exits 0; rhiza ship generate exits 0. The diff touches only ship-change-one-ci.txtar (plus the builder's log/actor records): the fake claude records plan, code and qa rounds each with an evidence line for leaf first, and first records a clean plan round with no evidence, which must fail; the script greps its stderr for 'leaves leaf first uncovered' and 'nothing was journaled' (matching internal/api/run_review.go:568) and asserts the round was not accepted. No Go code changed, no existing assertion dropped, no classic-check or plan-review script touched.

## Links

### Out

- verifies [conduct-fixtures-cover-leaves](../capability/rhiza/workflows/built-in-workflows/built-in-review/conduct-fixtures-cover-leaves.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
