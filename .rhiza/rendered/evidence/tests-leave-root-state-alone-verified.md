# tests-leave-root-state-alone-verified

evidence · `01M6T62WCKDD0PWVD62KNX681J` · actor [claude-code-01dd94b4](../actor/claude-code-01dd94b4.md)

## result

PASS

## statement

Verified at 86414309 from a checkout under the ~/code root (rhiza.toml at /home/rhiza/code). go test -trimpath -count=1 -v ./internal/rootguard/... passes: TestFileAddedAndChanged (added and changed files reported), TestLimitLineOutsideSpan (appended usage.limits.jsonl line with hit_at outside the span reported), TestUntouched, TestNoRoot, TestDaemonOwnedFiles (all 11 owned files rewritten) and TestLimitLineInsideSpan report nothing, and TestMainFailsOnWrite shows Main turns a stray write into a non-zero exit. go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/usage passes under the new rootguard.Main + testscript.RunMain TestMain (the rhiza command keeps main's profiled/run body). go test -trimpath -count=1 -timeout 30m ./internal/api/... ran with no rootguard report, and a find snapshot of ~/code/.rhiza-root (names, sizes, mtimes) before and after both runs was identical. Its only failure was the CPU timing budget TestReadBudget (roadmap what-if median over 50ms), with load average about 6 on 4 cores and another session's internal/api suite running at the same time. Re-run alone with -p 1, TestReadBudget passed and TestReadBudgetApply hit its 200ms timing budget under the same load. The diff only wraps m.Run in a snapshot and check outside the timed code, so these are load flakes, not regressions. The budget tests run serially under make check in CI. gofmt and go vet are clean on the changed packages.

## Links

### Out

- verifies [tests-leave-root-state-alone](../capability/rhiza/integrity/self-check/tests-leave-root-state-alone.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
