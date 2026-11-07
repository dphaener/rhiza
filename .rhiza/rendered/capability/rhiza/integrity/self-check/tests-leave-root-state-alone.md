# tests-leave-root-state-alone

capability · `01M6T26X2D1D7TTAJ0MMF0X9FN` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./internal/rootguard/... passes with tests that, on a temp root, the helper reports a file added or changed in .rhiza-root/ and a usage.limits.jsonl line appended with a hit_at outside the run span, and reports nothing when the folder is untouched, when no root encloses the dir, when only daemon-owned files changed, or when the appended line's hit_at lies inside the run span; go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run TestScripts/usage pass from a checkout under a root with no write attributed to them in its .rhiza-root/

## statement

The test runs of internal/api and cmd/rhiza fail when they write into the state folder of a root enclosing the package dir, through one shared helper in a new package internal/rootguard: TestMain snapshots that root's .rhiza-root/ (names, sizes, mod times, and the line count of usage.limits.jsonl) before m.Run and, after it, exits non-zero naming each write it attributes to the test run, while writes the live root daemon and run drivers make at the same time pass. Attribution: a file that appeared or changed is reported unless the daemon or drivers own it (daemon.lock, daemon.home, daemon.state.json, requests.json, requests.lock, starts.json, starts.lock, continues.json, continues.lock, main-stale.json, usage.limits.jsonl); in usage.limits.jsonl, an appended line is reported when its hit_at lies outside the wall-clock span of the test run (tests stamp hits with fixed clocks, real drivers with the current time). cmd/rhiza's TestMain moves from testscript.Main, which exits itself, to testscript.RunMain so the check runs after the scripts

## Links

### Out

- depends-on [empty-dir-finds-no-root](../../runtime/empty-dir-finds-no-root.md)
- decided-by [empty-wd-finds-no-root](../../../../decision/empty-wd-finds-no-root.md)

### In

- verifies [tests-leave-root-state-alone-verified](../../../../evidence/tests-leave-root-state-alone-verified.md)
- verifies [tests-leave-root-state-alone-qa-1](../../../../evidence/tests-leave-root-state-alone-qa-1.md)
- verifies [pr-346-fix-tests-write-fake](../../../../evidence/pr-346-fix-tests-write-fake.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
