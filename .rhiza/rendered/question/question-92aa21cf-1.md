# question-92aa21cf-1

question · `01M78FJ01VT7SZ6CAZQHDJ4NY4`

## session

92aa21cf-cf96-4789-b4d7-01a1d0f6deed

## statement

What should switch on the full local gate (make check, ~20 min, serialized) for every run? Today one flake anywhere in the last 20 CI runs does, even an unidentified one ('scripts (4/4) / UNKNOWN STEP') while 29 of 30 were green.

## Links

### In

- answers [local-gate-trigger-identified-flakes](../decision/local-gate-trigger-identified-flakes.md)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
