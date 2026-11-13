# pr-569-model-lints-in-review

evidence · `01M7BDWGGB7CP5AE2KG250N8ZA` · actor [dphaener](../actor/dphaener.md)

## head

045215880d6c93646c5c2b47cd404e08ce04a140

## pr

https://github.com/dphaener/rhiza/pull/569

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #569 (https://github.com/dphaener/rhiza/pull/569) passed CI at 0452158

## status

verified

## Links

### Out

- verifies [modellint-stale-refs](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-stale-refs.md)
- verifies [modellint-unapplied-decisions](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-unapplied-decisions.md)
- verifies [code-check-runs-lints](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-runs-lints.md)
- verifies [review-lint-fix-route](../capability/review-lint-fix-route.md)
- verifies [fixer-lint-hits-guide](../capability/fixer-lint-hits-guide.md)
- verifies [modellint-branch-removed-code](../capability/modellint-branch-removed-code.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
