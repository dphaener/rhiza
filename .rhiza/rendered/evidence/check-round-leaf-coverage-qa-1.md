# check-round-leaf-coverage-qa-1

evidence · `01M7PZCGZQ71RB0N848T4M5Q3C` · actor [claude-code-a658b5cb](../actor/claude-code-a658b5cb.md)

## result

PASS

## statement

QA at 1bfa91e9d, scratch build (go build -trimpath) driving a real ship-change conduct step in a scratch repo with a fake claude, 4 leaves (par,a,b,c): {"findings":[]} refused exit 2 naming 'leaves par, a, b, c', 'nothing was journaled' (round count unchanged); finding on a + evidence for b only refused naming 'par, c'; finding with no leaf + b,c evidence refused naming a; unknown leaf zzz, blank line, missing leaf, extra field, evidence as object all usage exit 2; evidence:null treated as none and refused; finding a + evidence b,c,par accepted; leaf by full ID, multiple lines per leaf and stdin '-' accepted; evidence journaled on the round. go test ./internal/api -run 'RunStatus|CheckRound|RunFindings|Coverage|Evidence' and ./internal/run/... pass.

## Links

### Out

- verifies [check-round-leaf-coverage](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-leaf-coverage.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
