# ship-local-gate-verified-8

evidence · `01M6YYQ72PM6BHR485XNRBNWTN` · actor [claude-code-ffbf3f20](../actor/claude-code-ffbf3f20.md)

## result

PASS

## statement

Verifier claude-code-ffbf3f20 at 96e14104 (attempt 8): go test -trimpath ./internal/api -run TestShipLocalGate -v passes its subtests unstable_runs_the_gate_in_the_worktree, stable_skips, no_gate_skips, failing_gate_fails_with_its_output (plus TestShipLocalGateOncePerHead); go test -trimpath ./internal/workflow passes, including TestShipChangeLocalGate, which checks a gate step before every push of ship-change.toml with on.fail to a fixer, and (new in this attempt) that gate-rebase-merge and rebase-merge fail through rebase-merge-gate-failed to gate-fix-merge, whose loop render-merge -> gate-push-merge -> push-merge (ship open-pr) -> wait-rebased-merge never returns to land or any pre-land step, sets max_visits on every step, and walks repeated gate failures ending only at gate-fix-merge's cap, which addresses quality cycle-3 f1 (land went over max_visits 1). go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-local-gate|run-ship-change-rebase|run-ship-change)$' passes (ship-local-gate shows rhiza.ship.localGate as a project setting with an empty default, then make check from git config after daemon set; run-ship-change-rebase now ends the post-land conflict at rebase-merge-gate-failed). Note: the scripts fail in the run's raw env only because GOTMPDIR puts the test binary under ~/code/rhiza, so rhiza init sees that repo; with RHIZA_* unset and GOTMPDIR/TMPDIR in the session scratchpad they pass. Read the diff of ship-change.toml, workflow_test.go and the txtar.

## Links

### Out

- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
