# test-timeout-set

capability · `01M69CDTPVNVSQSER2JM4XB0MF` · parent [self-check](../self-check.md)

## acceptance

A Go test runs make -n test at the repo root and asserts the go test line carries -timeout 20m; it skips when make is not on PATH

## statement

make test gives each package 20 minutes, so the full suite does not hit go test's 10-minute default on a slow runner

## Links

### Out

- decided-by [test-timeout-twenty-minutes](../../../../decision/test-timeout-twenty-minutes.md)

### In

- verifies [test-timeout-set-verified](../../../../evidence/test-timeout-set-verified.md)

## History

- `01M69CDDZAS7JREQN9AWKPZ7XY` make test gives each package 20 minutes instead of go test's default 10, so cmd/rhiza's growing script suite stops timing out on CI's two-core runner
