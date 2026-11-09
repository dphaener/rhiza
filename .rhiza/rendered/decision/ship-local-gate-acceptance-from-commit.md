# ship-local-gate-acceptance-from-commit

decision · `01M6YDVAKPM5P0RK7PZ730XFYY`

## reasoning

The finding and the critic are right: internal/api/daemon_settings.go keeps project settings (settingDef.project, ScopeProject) in the project's git config, which is untracked and shared by every worktree of the live repo, so a builder could only satisfy the clause by mutating the live checkout's config and a verifier's result would depend on machine state rather than the commit. Every other acceptance in this change set is a go test or a testscript, which is what a leaf's commit can prove. Moving the check into the ship-local-gate testscript (already in the acceptance, with its hermetic repo and git config per cmd/rhiza/main_test.go) keeps the intent covered: the key exists, is a project setting, defaults empty (so the gate is skipped with none set, as the statement and health-ci-stability's no-local-gate warning require) and is settable to make check. Decision ci-unstable-gate-rules (3) names make check as Rhiza's gate; that value stays the intent, but applying it to this machine's config is a deployment step like other rhiza daemon set calls, not build output.

## statement

ship-local-gate's acceptance drops the clause 'in this repo rhiza daemon settings --json lists rhiza.ship.localGate with value make check from the project settings' and checks the setting's definition instead: the ship-local-gate testscript, in its own repo, shows rhiza daemon settings --json listing rhiza.ship.localGate as a project setting with an empty default (source default), and rhiza daemon set rhiza.ship.localGate 'make check' then listing it with value make check from git config. Setting Rhiza's own gate to make check is not part of the leaf's commit: it is an operator step (rhiza daemon set rhiza.ship.localGate 'make check' in the live repo) taken once the change set lands, and the leaf's statement keeps 'make check for Rhiza itself' as the intended value, not as something the build proves.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6ycdrqqthsj4fv3sassw3v3-r3-f2](../question/review-01m6ycdrqqthsj4fv3sassw3v3-r3-f2.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
