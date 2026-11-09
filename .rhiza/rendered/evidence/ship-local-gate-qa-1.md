# ship-local-gate-qa-1

evidence · `01M6YR966TR2V9M6N6KYE0660H` · actor [claude-code-ac99a081](../actor/claude-code-ac99a081.md)

## result

PASS

## statement

QA claude-code-ac99a081, binary built from abb2ffe1 (go build -trimpath), run by hand in a hermetic scratch repo outside ~/code (env -i, own HOME, bare origin, fake gh on PATH): rhiza ship gate: no gate set 'Local gate skipped: no gate is set'; gate set and CI unstable (main red, then a flake) it ran sh -c in the change set's worktree repo.gate-trial (output ran-in=.../repo.gate-trial), passed and recorded the head, a second call on the same head 'it already passed on adab9e6'; a gate 'echo boom-output; exit 3' printed its output and 'Local gate failed ... exit status 3 ... fix it here, never push to see', exit 1, --json with output/error; CI stable 'Local gate skipped: CI is stable', exit 0. rhiza workflow describe ship-change shows local-gate, gate-rebase, gate-land, gate-rebase-merge before each push step, each on.fail -> gate-fix; workflow check valid. In this repo rhiza daemon settings lists rhiza.ship.localGate = make check (.rhiza/settings.json, project).

## Links

### Out

- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
