# empty-dir-finds-no-root

capability · `01M6T26X1CPJA6YNWN7KQTP4C5` · parent [runtime](../runtime.md)

## acceptance

go test -trimpath ./internal/rootcfg/... passes with a test that, from a cwd under a folder holding rhiza.toml, asserts Find("") returns not found; and go test -trimpath ./internal/api/... -run Usage passes with a usage-learn test whose Ctx has an empty Wd under an enclosing root that asserts the hit lands in the repo's own daemon folder and nothing is written under the enclosing root

## statement

rootcfg.Find("") reports no root instead of resolving the process cwd, so every caller given an empty Ctx.Wd (limitsDir and the api ops that call Find(ctx.Wd)) treats it as outside any root, while the CLI keeps passing its absolute working directory as Ctx.Wd

## Links

### Out

- decided-by [empty-wd-finds-no-root](../../../decision/empty-wd-finds-no-root.md)

### In

- depends-on [tests-leave-root-state-alone](../integrity/self-check/tests-leave-root-state-alone.md)
- verifies [empty-dir-finds-no-root-verified](../../../evidence/empty-dir-finds-no-root-verified.md)
- verifies [empty-dir-finds-no-root-qa-1](../../../evidence/empty-dir-finds-no-root-qa-1.md)
- verifies [pr-346-fix-tests-write-fake](../../../evidence/pr-346-fix-tests-write-fake.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
