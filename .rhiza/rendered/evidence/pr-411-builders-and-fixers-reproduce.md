# pr-411-builders-and-fixers-reproduce

evidence · `01M6ZGSF9HBFQ4CA7P1H3WJ5S5` · actor [dphaener](../actor/dphaener.md)

## head

be7cb76cba4a6218feb8a71baff4421ce2580386

## pr

https://github.com/dphaener/rhiza/pull/411

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #411 (https://github.com/dphaener/rhiza/pull/411) passed CI at be7cb76

## status

verified

## Links

### Out

- verifies [delay-repro-guidance](../capability/rhiza/integrity/self-check/delay-repro-guidance.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
