# question-33462c8e-1

question · `01M65048H5GPAKKSKRQ9CQ7CZF`

## session

33462c8e-2065-4b43-a930-01a18a012b7a

## statement

When a review fixer fails, how many times may it be retried inside its cycle before the cycle is spent? The code today (internal/run/checks.go nextCheckActions) spends the cycle on the first failed fix, and MaxFixCycles = 2 caps the cycles.

## Links

### In

- answers [fixer-retries-within-cycle](../decision/fixer-retries-within-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
