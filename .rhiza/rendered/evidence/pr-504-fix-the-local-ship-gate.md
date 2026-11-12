# pr-504-fix-the-local-ship-gate

evidence · `01M78V4GSAPK7AFAREE9NXBKQD` · actor [dphaener](../actor/dphaener.md)

## head

26308d45fea0838c3a0abd886635d76ea6e4de3a

## pr

https://github.com/dphaener/rhiza/pull/504

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #504 (https://github.com/dphaener/rhiza/pull/504) passed CI at 26308d4

## status

verified

## Links

### Out

- verifies [tests-leave-root-state-alone](../capability/rhiza/integrity/self-check/tests-leave-root-state-alone.md)
- verifies [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)
- verifies [ship-gate-full-log](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-gate-full-log.md)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
