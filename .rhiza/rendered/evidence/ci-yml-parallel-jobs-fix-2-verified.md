# ci-yml-parallel-jobs-fix-2-verified

evidence · `01M6QB992KMBZWY458KMMMRC1P` · actor [claude-code-6e32a8c4](../actor/claude-code-6e32a8c4.md)

## result

PASS

## statement

Verified attempt 2 (fixer commit 948ae1e1) against model finding f1. rhiza show t1-d12-ci-one-job shows a live 'superseded by ci-speed-shape' link, created by the op in log entry 01M6QB4WYQX9JAK9JMPF3CD0A1 (change set 01M6Q58T8EG0QY7YDQNWHW6NEA), so the one-job CI decision is retired and no longer contradicts the parallel-jobs shape. In a clean detached worktree at 948ae1e1, make self-check printed ok (exit 0): the model folds and rendered files are current. go test -trimpath -count=1 ./internal/ciworkflow/... passes (ok), so the leaf's ci.yml jobs, matrices, fail-fast, setup-go cache and aggregate check assertions still hold. The commit touches only .rhiza/ files. The PR CI duration point is read later from the ship run's PR CI, as the compliance check noted.

## Links

### Out

- verifies [ci-yml-parallel-jobs](../capability/rhiza/integrity/self-check/ci-yml-parallel-jobs.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
