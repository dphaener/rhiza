# ci-stability-state

capability · `01M6YCRKH9VEHGSH0A29RKY1ZD` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run 'TestCIStability|TestShipLocalGate|TestHealthCI' shows: an empty record stable; TestX flaked on 2 distinct heads inside the window unstable, and stable once later clean runs push the older of the two out; one flake of TestX, or a window whose only failures are infra ('scripts (4/4) / Run scripts', 'ci / UNKNOWN STEP', 'CI run 123'), stable, with a reason naming what was ignored, ship gate skipping with that reason, and rhiza health still listing them with their counts; a red dispatched main result unstable, and stable once a later ship check --main records green, whether by a tree-equal merge or a green dispatched run; rhiza.ship.ciWindow=5 shrinking the window; an invalid window refused by rhiza daemon set

## statement

Rhiza judges a repo's CI unstable while its last completed CI runs in the stability record (20 by default, the project setting rhiza.ship.ciWindow) hold a repeated named flake, one named test recorded as a flake on at least 2 distinct heads (as ci-unstable-repeated-flake states), or while main's latest result in the record is red, and stable again once neither holds. A single flake of a named test, and an infra failure (a failure naming no test: only 'job / step', 'UNKNOWN STEP' or 'CI run <id>'), stay in the record and in rhiza health with their counts but do not count toward the verdict, and the stable verdict's reason (and ship.gate's skip reason) says when one was seen and ignored. Main's latest result is the one the latest rhiza ship check --main recorded: a merge whose tree equals its tested PR head's records that head's green result for the merge commit, and a tree-differing merge records the dispatched CI run's result, so a red main clears at the next merge that settles green

## Links

### Out

- depends-on [ci-stability-record](ci-stability-record.md)
- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)
- decided-by [local-gate-trigger-identified-flakes](../../../../../decision/local-gate-trigger-identified-flakes.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../decision/local-gate-repeated-named-flakes.md)

### In

- depends-on [health-ci-stability](../../daemon/daemon-status/health-command/health-ci-stability.md)
- depends-on [ship-local-gate](ship-local-gate.md)
- depends-on [ci-unstable-repeated-flake](ci-unstable-repeated-flake.md)
- verifies [ci-stability-state-verified](../../../../../evidence/ci-stability-state-verified.md)
- verifies [ci-stability-state-qa-1](../../../../../evidence/ci-stability-state-qa-1.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)
- verifies [pr-504-fix-the-local-ship-gate](../../../../../evidence/pr-504-fix-the-local-ship-gate.md)
- verifies [pr-509-the-local-ship-gate-switches](../../../../../evidence/pr-509-the-local-ship-gate-switches.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
- `01M78RM7PFF8KV94344XM3V8DH` The local ship gate switches on only for an identified flake that repeats, and keeps its full log. Owner, 2026-11-12, after the performance report ('Well that's disappointing'), agreed to file the orchestrator's recommendation. Today one CI flake with no named test ('scripts (4/4) / UNKNOWN STEP' on #498 at ~16:24, auto-filed 01M78960; 29 of the last 30 CI runs were green) marked CI unstable and switched rhiza.ship.localGate (make check, serialized by the machine-wide test lock, ~20 min each) on for every run: about 525 run-minutes lost across 4 runs, run 01M787SV alone ~208 min over 4 local-gate visits (16:32-19:41). Wanted: (1) CI counts as unstable only for a flake that names a failing test and has failed more than once in the window (a step-level or UNKNOWN STEP failure is infra, not a flake); (2) the gate stores its full output as a file in the run dir and the journal points to it, instead of keeping only the last ~2 KB (today the failing test's name was cut off, so every gate-fix agent reran the whole suite); (3) the gate's reason names the flake and when it stops applying. Evidence: perf report ~/orchestrator/perf-2026-11-12/report-prospective.md; related bug 01M78CSH (rootguard false positive that made every gate fail).
