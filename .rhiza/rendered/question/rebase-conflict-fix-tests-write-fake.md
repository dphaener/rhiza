# rebase-conflict-fix-tests-write-fake

question · `01M6T6Z7PGCKXDV7039ACA83G6`

## statement

Rebasing branch cs/fix-tests-write-fake onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: internal/api/usage_learn_test.go?

## Links

### Out

- blocks change set `01M6T20EGCJZNFEY6YPJ1PSEDW` (Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it)

### In

- answers [rebase-conflict-fix-tests-write-fake-merged](../decision/rebase-conflict-fix-tests-write-fake-merged.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
