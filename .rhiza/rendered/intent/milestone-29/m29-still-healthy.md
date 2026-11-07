# m29-still-healthy

intent · `01M6RN37A82JPZ6VS3WSHXCR5A` · parent [milestone-29](../milestone-29.md)

## statement

make check is green, rhiza check included, read commands stay within their budget, and the default test run binds no socket.

## Links

### In

- serves change set `01M6T20EGCJZNFEY6YPJ1PSEDW` (Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
