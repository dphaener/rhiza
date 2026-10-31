# fixer-exhaustion-spends-cycle

capability · `01M6506FZTPP6JZJBYMMPZMQFC`

## acceptance

A checks step whose change-set fixer fails all 1 + MaxFixCycles attempts of a cycle ends that cycle's checks failed with the detail 'cycle N: fix failed', counts the cycle against MaxFixCycles, and starts no further check session in it

## statement

A review fixer that fails every retry spends its cycle, and the cycle's checks fail with the fix-failed detail

## Links

### Out

- decided-by [fixer-retries-within-cycle](../decision/fixer-retries-within-cycle.md)
- decided-by [fixer-retries-within-cycle-v2](../decision/fixer-retries-within-cycle-v2.md)

### In

- verifies [fixer-exhaustion-spends-cycle-verified](../evidence/fixer-exhaustion-spends-cycle-verified.md)
- verifies [fixer-exhaustion-spends-cycle-qa-1](../evidence/fixer-exhaustion-spends-cycle-qa-1.md)
- verifies [pr-100-failed-fixer-keeps-its-cycle](../evidence/pr-100-failed-fixer-keeps-its-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
