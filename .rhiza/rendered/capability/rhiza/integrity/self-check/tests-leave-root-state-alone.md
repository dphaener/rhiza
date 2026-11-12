# tests-leave-root-state-alone

capability · `01M6T26X2D1D7TTAJ0MMF0X9FN` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./internal/rootguard/... passes with tests that, on a temp root, the helper reports a file added or changed in .rhiza-root/, a usage.limits.jsonl line appended with a hit_at outside the run span, and a usage.live.json snapshot written with a captured_at outside the run span, and reports nothing when the folder is untouched, when no root encloses the dir, when only daemon-owned files changed, when the appended line's hit_at lies inside the run span, or when usage.live.json was rewritten with snapshots captured inside the run span or left unchanged

## statement

The test runs of internal/api and cmd/rhiza fail when they write into the state folder of a root enclosing the package dir, through one shared helper in package internal/rootguard: TestMain snapshots that root's .rhiza-root/ (names, sizes, mod times, the line count of usage.limits.jsonl, and the snapshots in usage.live.json by kind and captured_at) before m.Run and, after it, exits non-zero naming each write it attributes to the test run, while writes the live root daemon, run drivers and live agents' status lines make at the same time pass. Attribution: a file that appeared or changed is reported unless the daemon or drivers own it (daemon.lock, daemon.home, daemon.state.json, requests.json, requests.lock, starts.json, starts.lock, continues.json, continues.lock, main-stale.json, usage.limits.jsonl, usage.live.json); in usage.limits.jsonl, an appended line is reported when its hit_at lies outside the wall-clock span of the test run, and in usage.live.json a new snapshot is reported when its captured_at lies outside that span (tests stamp with fixed clocks, real drivers and status lines with the current time). cmd/rhiza's TestMain uses testscript.RunMain so the check runs after the scripts

## Links

### Out

- depends-on [empty-dir-finds-no-root](../../runtime/empty-dir-finds-no-root.md)
- decided-by [empty-wd-finds-no-root](../../../../decision/empty-wd-finds-no-root.md)
- decided-by [guard-fix-recorded-here](../../../../decision/guard-fix-recorded-here.md)

### In

- verifies [tests-leave-root-state-alone-verified](../../../../evidence/tests-leave-root-state-alone-verified.md)
- verifies [tests-leave-root-state-alone-qa-1](../../../../evidence/tests-leave-root-state-alone-qa-1.md)
- verifies [pr-346-fix-tests-write-fake](../../../../evidence/pr-346-fix-tests-write-fake.md)
- verifies [pr-504-fix-the-local-ship-gate](../../../../evidence/pr-504-fix-the-local-ship-gate.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
