# empty-wd-finds-no-root

decision · `01M6T25PYG4GAVPAD47AC5R6KN`

## reasoning

The leak comes from filepath.Abs("") turning an empty dir into the cwd inside rootcfg.Find; fixing it at Find closes it for every caller, not just limitsDir (daemon.go, project.go, root.go, daemon_service.go, usage_learn.go all call Find(ctx.Wd)). rootOf in daemon_settings.go already treats an empty Wd as no root, so this generalises an existing rule. cmd/rhiza/ctx.go always sets Wd from os.Getwd, and daemon_projects.go sets it explicitly, so no production path relies on the empty-dir fallback. For failing a test run, a sandbox or syscall interception is not portable or cheap in Go; the concrete damage was files appearing under ~/code/.rhiza-root, so a TestMain before/after snapshot of the enclosing root's state dir catches exactly that class in the two packages that reach root state, fails the run loudly, and is a no-op on CI where no root encloses the checkout. Setting Wd in testCtx alone would fix today's test but not stop the next one.

## statement

An empty directory never resolves a root: rootcfg.Find("") returns not found instead of falling back to the process cwd, so every caller (limitsDir and the other api ops on ctx.Wd) treats an empty Ctx.Wd as outside any root; the CLI keeps passing os.Getwd() explicitly as Ctx.Wd, so real runs are unchanged. A test run fails on writes outside its temp dir through a guard (package internal/rootguard) in TestMain of internal/api and cmd/rhiza: before m.Run it snapshots the state dir (.rhiza-root/: names, sizes, mod times, the usage.limits.jsonl line count) of any root enclosing the package dir, and after m.Run exits non-zero naming each write attributed to the run. Because the live root daemon and run drivers write there at the same time, the files they own are excluded, and an appended usage.limits.jsonl line counts only when its hit_at lies outside the run's wall-clock span. A unit test pins Find("") to not found, and a usage-learn test with an empty Wd asserts nothing is written under the enclosing root.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-df8b7e8a-1](../question/question-df8b7e8a-1.md)
- derived-from [fix-tests-write-fake-plan](../source/fix-tests-write-fake-plan.md)

### In

- decided-by [empty-dir-finds-no-root](../capability/rhiza/runtime/empty-dir-finds-no-root.md)
- decided-by [tests-leave-root-state-alone](../capability/rhiza/integrity/self-check/tests-leave-root-state-alone.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
