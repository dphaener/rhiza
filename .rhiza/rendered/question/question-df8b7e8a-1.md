# question-df8b7e8a-1

question · `01M6T240AFNVXCY1N0YG48S78B`

## session

df8b7e8a-c713-40c2-ae5d-01a1b4206f96

## statement

Where should an empty working directory stop resolving the machine's root (today rootcfg.Find("") resolves the process cwd)? / How should a test run fail when something writes outside its temp dir?

## Links

### In

- answers [empty-wd-finds-no-root](../decision/empty-wd-finds-no-root.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
