# test-timeout-twenty-minutes

decision · `01M69CDSGEEYK37EBFA5661FBY`

## statement

make test runs go test with -timeout 20m: cmd/rhiza's scripts took 540 s on CI's two-core runner against the 10-minute default, and arbiter-panes-one-window timed out on main at 8afcc20, failing the run's main-green

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m7-keep-picking](../source/m7-keep-picking.md)

### In

- decided-by [test-timeout-set](../capability/rhiza/integrity/self-check/test-timeout-set.md)

## History

- `01M69CDDZAS7JREQN9AWKPZ7XY` make test gives each package 20 minutes instead of go test's default 10, so cmd/rhiza's growing script suite stops timing out on CI's two-core runner
