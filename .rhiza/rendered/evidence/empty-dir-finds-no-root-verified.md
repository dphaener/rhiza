# empty-dir-finds-no-root-verified

evidence · `01M6T3YAKDN8GEPP4D5PCG4VF7` · actor [claude-code-5547ffd6](../actor/claude-code-5547ffd6.md)

## result

PASS

## statement

At 2a7a01a1: go test -trimpath ./internal/rootcfg/... passes, including TestFindEmptyDirFindsNoRoot, which chdirs into a subfolder of a rhiza.toml root and asserts Find("") reports no root (and Find(".") still finds it). go test -trimpath ./internal/api/... -run Usage passes, including TestUsageLearnEmptyWdStaysInRepoFolder: Ctx.Wd empty, cwd under an enclosing root, recordLimitHit writes one hit to the repo's daemon folder and none under the root (RootStateDir not even created). In a scratch worktree with the old rootcfg.go both tests fail, so they cover the fix. Diff read: Find returns false for an empty dir; cmd/rhiza newCtx still sets Ctx.Wd from os.Getwd (absolute), and TUI projectctx callers pass os.Getwd. go vet and gofmt clean on both packages.

## Links

### Out

- verifies [empty-dir-finds-no-root](../capability/rhiza/runtime/empty-dir-finds-no-root.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
