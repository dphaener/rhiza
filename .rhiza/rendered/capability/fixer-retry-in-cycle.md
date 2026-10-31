# fixer-retry-in-cycle

capability · `01M6506EMCPMJKKCE8B3JX16RP`

## acceptance

A checks step whose change-set fixer fails its first attempt and passes its second starts no new cycle: the retry runs inside the same cycle, and that cycle's checks start only after the fixer passes

## statement

A failed change-set review fixer is retried within its cycle, and the cycle's checks start again only once its fixers pass

## Links

### Out

- decided-by [fixer-retries-within-cycle](../decision/fixer-retries-within-cycle.md)
- decided-by [fixer-retries-within-cycle-v2](../decision/fixer-retries-within-cycle-v2.md)

### In

- verifies [fixer-retry-in-cycle-verified](../evidence/fixer-retry-in-cycle-verified.md)
- verifies [fixer-retry-in-cycle-qa-1](../evidence/fixer-retry-in-cycle-qa-1.md)
- verifies [fixer-retry-in-cycle-verified-2](../evidence/fixer-retry-in-cycle-verified-2.md)
- verifies [pr-100-failed-fixer-keeps-its-cycle](../evidence/pr-100-failed-fixer-keeps-its-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
