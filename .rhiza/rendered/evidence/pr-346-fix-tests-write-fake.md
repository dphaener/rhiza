# pr-346-fix-tests-write-fake

evidence · `01M6TBQT48EBHJV148578HGM0J` · actor [dphaener](../actor/dphaener.md)

## head

870f9e8a8f99e107c5ba732a3ba59ac2d836d8bd

## pr

https://github.com/dphaener/rhiza/pull/346

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #346 (https://github.com/dphaener/rhiza/pull/346) passed CI at 870f9e8

## status

verified

## Links

### Out

- verifies [empty-dir-finds-no-root](../capability/rhiza/runtime/empty-dir-finds-no-root.md)
- verifies [tests-leave-root-state-alone](../capability/rhiza/integrity/self-check/tests-leave-root-state-alone.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
