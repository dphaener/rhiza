# check-round-evidence-status-qa-1

evidence · `01M7PZCH55NGQGE5YAS3N42Y18` · actor [claude-code-a658b5cb](../actor/claude-code-a658b5cb.md)

## result

PASS

## statement

QA at 1bfa91e9d in a scratch conduct run: rhiza run status --json gives each round an evidence array of {leaf, leaf_handle, line} with handles; human output lists '- evidence for <handle>: <line>' under the round, after the findings line which is unchanged ('1 finding, 0 confirmed' / '- f1 minor at x:1 on a: r'); a multi-line evidence line is folded to one line. After retiring leaf b, status still exits 0, JSON keeps b's ID with leaf_handle "" and human text prints the bare ID. go test ./internal/api -run RunStatus passes (evidence-less round has no evidence key).

## Links

### Out

- verifies [check-round-evidence-status](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-evidence-status.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
