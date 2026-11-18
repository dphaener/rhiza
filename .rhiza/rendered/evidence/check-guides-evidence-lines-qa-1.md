# check-guides-evidence-lines-qa-1

evidence · `01M7PZCH8BFX84WTD8CKS97FG6` · actor [claude-code-a658b5cb](../actor/claude-code-a658b5cb.md)

## result

PASS

## statement

QA at 1bfa91e9d: scratch-built rhiza guide check-plan/check-code/check-qa each show the findings example with an "evidence":[{leaf,line}] list, say 'one evidence line per acceptance point of each leaf' and that an uncovered round 'is refused and not recorded'; check-qa still carries the rhiza apply Evidence node example; check-compliance/quality/model unchanged (no evidence text). go test -trimpath -p 2 ./internal/guide/... passes.

## Links

### Out

- verifies [check-guides-evidence-lines](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-guides-evidence-lines.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
