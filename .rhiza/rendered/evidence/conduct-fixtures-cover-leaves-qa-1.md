# conduct-fixtures-cover-leaves-qa-1

evidence · `01M7PZCH2RPF7N8A778YW2ZA5K` · actor [claude-code-a658b5cb](../actor/claude-code-a658b5cb.md)

## result

PASS

## statement

QA at 1bfa91e9d: go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(ship-change-one-ci|describe-golden)$' passes (ship-change-one-ci 3.7s); script's fake conductor records the uncovered plan round first (asserted refused naming leaf first, nothing journaled) then plan/code/qa rounds each with an evidence line for first. Mirrored the same flow by hand in a scratch conduct run: identical refusal text.

## Links

### Out

- verifies [conduct-fixtures-cover-leaves](../capability/rhiza/workflows/built-in-workflows/built-in-review/conduct-fixtures-cover-leaves.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
