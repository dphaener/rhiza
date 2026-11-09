# review-01m6ycdrqqthsj4fv3sassw3v3-r3-f2

question · `01M6YDS4BG7SXYASYG81XCJA5K`

## statement

Plan review round 3 left finding f2 (missing-acceptance on ship-local-gate): The acceptance's last clause, 'in this repo rhiza daemon settings --json lists rhiza.ship.localGate with value make check from the project settings', cannot be produced or checked from the leaf's commit: project settings live in the repo's git config (internal/api/daemon_settings.go), which is untracked and shared by every worktree of the live repo, so a builder would have to mutate the live checkout's config and a verifier's result depends on machine state, not the code. Make it a check on the setting's definition/default (or a testscript with its own repo), and leave configuring this repo to a separate step. The critic confirmed it: Verified: project settings are read from the project's git config (internal/api/daemon_settings.go:22-23, ScopeProject), and .git/config is untracked and shared by every worktree of the live repo. The clause 'in this repo rhiza daemon settings --json lists rhiza.ship.localGate with value make check' therefore checks machine state, not the leaf's commit. A builder cannot satisfy it with a commit, and a verifier's result depends on whether someone ran git config on that machine. The acceptance should check the setting's definition (key, scope, empty default) or use a testscript with its own repo, and leave configuring Rhiza's own gate to a separate step. How should it be settled?

## Links

### Out

- blocks [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

### In

- answers [ship-local-gate-acceptance-from-commit](../decision/ship-local-gate-acceptance-from-commit.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
