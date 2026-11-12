# local-ship-gate-plan

source · `01M78FMVWV8MPY8Q2YCKJJ5V4J`

## quotes

```json
[]
```

## raw_sha256

93cf0ccbb58c8bfe9e194d5bc0f023cecff9248885b5ec1fa7bcf276c63a4d71

## session

92aa21cf-cf96-4789-b4d7-01a1d0f6deed

## summary

Autonomous plan run (no owner) for the bug 'Fix The local ship gate'. Found the root guard fix already on main (commit e075aefa6 in PR #498 judges .rhiza-root/usage.live.json by captured_at), so this change set records it on tests-leave-root-state-alone instead of building it again. The gate's stored output was cut to the step's 2 KB / 20-line tail, which lost the failing test's name: planned a full gate log kept on disk, with the failing test lines and the log path at the end of the output the journal keeps. Asked what should switch the full local gate on; the arbiter ruled (local-gate-trigger-identified-flakes) that only main red or a flake attributed to a named test or known job/step does, while an unidentified flake is still recorded, shown and filed but ignored by the verdict. Not planned: deduplicating parallel gate fixers working the same failure.

## title

Local ship gate fix, autonomous plan

## Links

### In

- derived-from [local-gate-trigger-identified-flakes](../decision/local-gate-trigger-identified-flakes.md)
- derived-from [gate-log-kept-whole](../decision/gate-log-kept-whole.md)
- derived-from [guard-fix-recorded-here](../decision/guard-fix-recorded-here.md)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
