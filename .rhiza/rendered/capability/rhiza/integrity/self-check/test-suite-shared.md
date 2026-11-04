# test-suite-shared

capability · `01M6JR8K1R9N64RNNZJD09GKMS` · parent [self-check](../self-check.md)

## acceptance

TestMakeTestTimeout runs make -n test and asserts: every go test line carries -trimpath, -timeout 20m and 9>&-; the packages other than cmd/rhiza start in the background before the lock; cmd/rhiza's scripts and then the timing budgets each run under a flock on a lock file beside GOCACHE, with a 'Waiting for another make test' line when it is held. Measured on the loaded 4-core VM: a fresh worktree of the same commit builds all 36 test binaries in 5 s with -trimpath against 143 s without (0 reused), and 22 of 35 internal packages are (cached) there; a second make test started while one holds the lock prints the waiting line and runs once it is free; where flock(1) is missing the suite runs unlocked; plan-workflow.txtar's fake rex answers block close, so the script takes 23 s instead of 509 s

## statement

make test is cheap to run from many worktrees on one machine: a second worktree at the same code reuses the first one's test builds and cached results, and only one worktree at a time runs cmd/rhiza's CPU-heavy scripts or the timing budgets, while the rest of the suite runs beside them

## Links

### Out

- decided-by [test-suite-locks-scripts-only](../../../../decision/test-suite-locks-scripts-only.md)

## History

- `01M6JJMK3M50P34PPX1QNTHYCB` The test suite is cheap to run in many worktrees at once: unchanged packages hit the shared test cache across worktrees, and concurrent full suites on one machine queue instead of thrashing
