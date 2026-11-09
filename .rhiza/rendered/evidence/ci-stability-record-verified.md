# ci-stability-record-verified

evidence · `01M6YHKHVSSPS6T9JHNK6K16S1` · actor [claude-code-3ab71de0](../actor/claude-code-3ab71de0.md)

## result

PASS

## statement

Verifier at 4f1bc3c3: go test -trimpath ./internal/api -run TestShipCheckStabilityRecord -v -count=1 passed all five subtests: a green first attempt is recorded once as a clean run (head, PR 7, not main, green, time set, no flakes, a repeat check adds nothing); a red head whose rerun passes records nothing before the rerun settles, then one green entry with flakes TestGarden and TestPond read from the first attempt's --log-failed; a head still red after its rerun records one red entry with no flake; a tree-equal ship check --main records a green main entry for merge commit fff0000aaa with no gh run call; a tree-differing one records the dispatched run's green on main. Re-run with all TestShip* tests from a clean git-archive copy outside ~/code (RHIZA_ROOT/CHANGE/RUN unset): ok. go vet and gofmt clean on internal/api. Read the diff: recordSettledCI hooks the settled paths in runShipCheck, shipCheckMainTree and shipCheckMainDispatch after shipRerunFailed (so a rerun records nothing until it settles), dedups per (head, main), skips dry runs, keeps 500 entries under the git common dir; RerunTests are captured only on a head's first rerun. In-tree, the full internal/api package hit the rootguard on .rhiza-root/usage.live.json, which the live root daemon rewrites on its own (mtime moved with no test running), and the 600s default timeout under VM load; both are environmental, not this change.

## Links

### Out

- verifies [ci-stability-record](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-record.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
