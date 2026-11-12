# ci-stability-state

capability · `01M6YCRKH9VEHGSH0A29RKY1ZD` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestCIStability|TestShipLocalGate|TestHealthCI' shows: an empty record stable; an identified flake ('scripts (4/4) / Run scripts' or '--- FAIL: TestX') inside the window unstable, and stable once 20 later clean runs push it out; a window whose only flakes are unidentified ('ci / UNKNOWN STEP', 'CI run 123') stable, with a reason naming the ignored flake, ship gate skipping with that reason, and rhiza health still listing the unidentified flake with its count; a red dispatched main result unstable, and stable once a later ship check --main records green, whether by a tree-equal merge or a green dispatched run; rhiza.ship.ciWindow=5 shrinking the window; an invalid window refused by rhiza daemon set

## statement

Rhiza judges a repo's CI unstable while its last completed CI runs in the stability record (20 by default, the project setting rhiza.ship.ciWindow) hold an identified flake, or while main's latest result in the record is red, and stable again once neither holds. An identified flake names its failing test ('--- FAIL:' test name) or a 'job / step' whose step is known; an unidentified flake (step read as UNKNOWN STEP, or 'CI run <id>' whose log names no job or step) stays in the record, in rhiza health with its count and in the watchdog's filed bugs, but does not count toward the verdict, and the stable verdict's reason (and ship.gate's skip reason) says when one was seen and ignored. Main's latest result is the one the latest rhiza ship check --main recorded: a merge whose tree equals its tested PR head's records that head's green result for the merge commit, and a tree-differing merge records the dispatched CI run's result, so a red main clears at the next merge that settles green

## Links

### Out

- depends-on [ci-stability-record](ci-stability-record.md)
- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)
- decided-by [local-gate-trigger-identified-flakes](../../../../../decision/local-gate-trigger-identified-flakes.md)

### In

- depends-on [health-ci-stability](../../daemon/daemon-status/health-command/health-ci-stability.md)
- depends-on [ship-local-gate](ship-local-gate.md)
- verifies [ci-stability-state-verified](../../../../../evidence/ci-stability-state-verified.md)
- verifies [ci-stability-state-qa-1](../../../../../evidence/ci-stability-state-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)
- verifies [pr-504-fix-the-local-ship-gate](../../../../../evidence/pr-504-fix-the-local-ship-gate.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
