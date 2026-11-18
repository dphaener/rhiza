# question-e50b77ed-1

question · `01M7PTC9TPNZFD522P4DG5W7W7`

## session

e50b77ed-cc5e-4e7f-950d-01a1eda447fd

## statement

Change set 'Fix Shallow checks pass': how should `rhiza run findings --check` refuse a shallow empty round in a conductor step (an agent step declaring checks; classic checks steps unchanged)? Check names are free workflow config (plan/code/qa are only ship-change's names).

## Links

### In

- answers [shallow-round-needs-leaf-coverage](../decision/shallow-round-needs-leaf-coverage.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
