# budget-timings-in-ci-run

evidence · `01M5WKQNAC8BTMJ3M6Z2FE151M`

## result

The log of that run shows, in ms, roadmap as-is 9.94 and to-be 9.96, ready as-is 9.82 and to-be 9.94 (GOMAXPROCS=2, NumCPU=2), with min and max for every read and view. The same lines appear in attempts 2 and 3 (roadmap 14.8 and 14.8 / 13.4 and 13.9, ready 14.9 and 14.8 / 13.6 and 13.4) and in the baseline run https://github.com/dphaener/rhiza/actions/runs/40357489189.

## statement

CI prints the read-budget timings (median, min and max per read and view) on every run: make test sets RHIZA_BUDGET_LOG, and make check, which CI runs, prints a 'Read-budget timings' heading and the 15-line log after the tests, passing or failing.

## tests

Makefile test target (RHIZA_BUDGET_LOG, printed after go test, exit status kept); TestReadBudget writes the log. Run https://github.com/dphaener/rhiza/actions/runs/40366003546 (PR #28, 3f74ee3), Run checks step.

## Links

### Out

- verifies [budget-timings-in-ci](../capability/rhiza/views/fast-reads/budget-timings-in-ci.md)

## History

- `01M5WEV2063MZ6QQ2RHMPW49V4` Reads meet their budget on GitHub's runner with room to spare: roadmap and ready median under 25 ms there, and CI prints the read-budget timings on every run so the runner is measured, not guessed.
