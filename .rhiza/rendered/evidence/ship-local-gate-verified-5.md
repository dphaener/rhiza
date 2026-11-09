# ship-local-gate-verified-5

evidence · `01M6YXN2BGK6SZCQ4EC068VAYA` · actor [claude-code-f835634e](../actor/claude-code-f835634e.md)

## result

PASS

## statement

Verifier claude-code-f835634e at c0cba816 (attempt 5): go test -trimpath ./internal/api -run TestShipLocalGate -v passes its subtests: unstable runs the gate in the worktree, stable skips, no gate skips, a failing gate fails with its output. TestShipRebaseRunsLocalGate, TestShipGateFailed and the TestShipRebase* tests pass too. go test -trimpath ./internal/workflow passes, including TestShipChangeLocalGate: a gate step before each push step, on.fail to gate-fix, and rebase and rebase-merge failing to rebase-gate-failed (rhiza ship gate --failed, on.pass to gate-fix). go test -trimpath ./cmd/rhiza -run TestScripts/ship-local-gate passes: in its own repo, daemon settings --json lists rhiza.ship.localGate with an empty default (source default), and after daemon set 'make check' it shows make check from git config. describe-golden, run-ship-change, run-ship-change-fail and run-ship-change-rebase pass. The scripts were run with RHIZA_* unset and GOTMPDIR/TMPDIR in a scratch dir, because the run's GOTMPDIR sits inside the main repo's .git, where rhiza init finds the outer repo. The diff fixes the quality finding: a gate that fails on the rebased head leaves main merged in, unpushed, and records gate_failed so gate-fix sees the failure. .rhiza/settings.json and the project settings file layer are gone, which resolves the scope and model findings.

## Links

### Out

- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
