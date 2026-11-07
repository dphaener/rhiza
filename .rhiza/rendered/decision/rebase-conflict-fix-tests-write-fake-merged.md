# rebase-conflict-fix-tests-write-fake-merged

decision · `01M6TBBMH1Y8DFGMBS0CY88HF0`

## statement

The conflict in internal/api/usage_learn_test.go is resolved by keeping both sides' new tests (this change set's TestUsageLearnEmptyWdStaysInRepoFolder and main's TestLimitHitTeachesOnlyItsKind and TestLearnedLimitRatchets), merging origin/main into the branch; the usage tests pass on the merge and write nothing under ~/code/.rhiza-root

## wording

Stated by the orchestrator acting as arbiter: keep both sides' tests and merge main

## Links

### Out

- answers [rebase-conflict-fix-tests-write-fake](../question/rebase-conflict-fix-tests-write-fake.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
