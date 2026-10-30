# pr-100-failed-fixer-keeps-its-cycle

evidence · `01M65AMM3R6X4QD4GREW553J34` · actor [darin-haener](../actor/darin-haener.md)

## head

8dcbb391a80ab6ad7baed70a944fa44066dc31fe

## pr

https://github.com/dphaener/rhiza/pull/100

## result

1 check passed: check SUCCESS

## statement

Pull request #100 (https://github.com/dphaener/rhiza/pull/100) passed CI at 8dcbb39

## status

verified

## Links

### Out

- verifies [fixer-retry-in-cycle](../capability/fixer-retry-in-cycle.md)
- verifies [fixer-exhaustion-spends-cycle](../capability/fixer-exhaustion-spends-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
